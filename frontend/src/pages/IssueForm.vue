<template>
  <div class="container mx-auto p-4 md:p-8">
    <h1 class="text-3xl font-bold text-gray-800 mb-6">
      {{ isEditMode ? '이슈 수정' : '새 이슈 생성' }}
    </h1>

    <form @submit.prevent="handleSubmit" class="bg-white p-6 rounded-lg shadow-sm">
      <!-- 제목 -->
      <div class="mb-4">
        <label for="title" class="block text-gray-700 text-sm font-bold mb-2">제목:</label>
        <input
          id="title"
          type="text"
          v-model="issue.title"
          :disabled="isDisabledAll"
          required
          class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
        />
      </div>

      <!-- 설명 -->
      <div class="mb-4">
        <label for="description" class="block text-gray-700 text-sm font-bold mb-2">설명:</label>
        <textarea
          id="description"
          v-model="issue.description"
          :disabled="isDisabledAll"
          required
          rows="5"
          class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
        ></textarea>
      </div>

      <!-- 상태 (수정 모드만 표시) -->
      <div v-if="isEditMode" class="mb-4">
        <label for="status" class="block text-gray-700 text-sm font-bold mb-2">상태:</label>
        <select
          id="status"
          v-model="issue.status"
          :disabled="!canChangeStatus"
          class="shadow border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
        >
          <option value="PENDING">대기중</option>
          <option value="IN_PROGRESS" :disabled="!hasAssignee">진행중</option>
          <option value="COMPLETED" :disabled="!hasAssignee">완료됨</option>
          <option value="CANCELLED" :disabled="!hasAssignee">취소됨</option>
        </select>
        <p v-if="!hasAssignee && canChangeStatus" class="text-sm text-red-500 mt-1">
          담당자가 지정되지 않으면 PENDING 상태만 선택할 수 있습니다.
        </p>
      </div>

      <!-- 담당자 -->
      <div class="mb-6">
        <label for="user" class="block text-gray-700 text-sm font-bold mb-2">담당자:</label>
        <select
          id="user"
          v-model="selectedUserId"
          :disabled="!canChangeAssignee"
          class="shadow border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
        >
          <option :value="0">미할당</option>
          <option v-for="user in users" :key="user.id" :value="user.id">{{ user.name }}</option>
        </select>
      </div>

      <!-- 버튼 -->
      <div class="flex items-center justify-between">
        <button
          type="submit"
          :disabled="loading || isDisabledAll"
          class="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-lg shadow-md transition duration-300 ease-in-out transform hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{ loading ? '저장 중...' : isEditMode ? '이슈 수정' : '이슈 생성' }}
        </button>
        <button
          type="button"
          @click="goBack"
          class="bg-gray-300 hover:bg-gray-400 text-gray-800 font-semibold py-2 px-4 rounded-lg shadow-md transition duration-300 ease-in-out transform hover:scale-105"
        >
          목록으로 돌아가기
        </button>
      </div>

      <!-- 메시지 -->
      <p
        v-if="message"
        :class="{ 'text-green-600': !error, 'text-red-600': error }"
        class="mt-4 text-center"
      >
        {{ message }}
      </p>
    </form>
  </div>
</template>

<script>
import axios from 'axios'
import { users as mockUsers } from '../data/mockData'

export default {
  name: 'IssueForm',
  data() {
    return {
      API_BASE_URL: 'http://localhost:8080',
      issue: {
        id: null,
        title: '',
        description: '',
        status: 'PENDING',
        user: null,
      },
      selectedUserId: 0,
      users: mockUsers,
      loading: false,
      message: '',
      error: false,
    }
  },
  computed: {
    isEditMode() {
      return this.$route.params.id !== undefined
    },
    isCompletedOrCancelled() {
      return ['COMPLETED', 'CANCELLED'].includes(this.issue.status)
    },
    hasAssignee() {
      return this.selectedUserId !== 0
    },
    isDisabledAll() {
      return this.isEditMode && this.isCompletedOrCancelled
    },
    canChangeStatus() {
      return this.isEditMode && !this.isCompletedOrCancelled
    },
    canChangeAssignee() {
      return !this.isCompletedOrCancelled
    },
  },
  watch: {
    '$route.params.id': {
      handler(newId) {
        if (newId && this.isEditMode) {
          this.fetchIssue(newId)
        } else {
          this.resetForm()
        }
      },
      immediate: true,
    },
    selectedUserId(newUserId) {
      if (this.isEditMode && newUserId === 0 && this.issue.status !== 'PENDING') {
        if (!this.isCompletedOrCancelled) {
          this.issue.status = 'PENDING'
        }
      }
    },
  },
  methods: {
    resetForm() {
      this.issue = {
        id: null,
        title: '',
        description: '',
        status: 'PENDING',
        user: null,
      }
      this.selectedUserId = 0
      this.message = ''
      this.error = false
    },
    async fetchIssue(id) {
      this.loading = true
      this.message = ''
      this.error = false
      try {
        const response = await axios.get(`${this.API_BASE_URL}/issue/${id}`)
        this.issue = response.data
        this.selectedUserId = this.issue.user ? this.issue.user.id : 0
      } catch (err) {
        console.error('이슈 상세 정보 로딩 중 오류 발생:', err)
        this.message = `이슈를 로드하는 데 실패했습니다: ${err.response?.data?.error || err.message}`
        this.error = true
        this.$router.push('/')
      } finally {
        this.loading = false
      }
    },
    async handleSubmit() {
      this.loading = true
      this.message = ''
      this.error = false

      if (this.isEditMode && this.selectedUserId === 0 && this.issue.status !== 'PENDING') {
        this.message = '담당자가 없는 이슈는 PENDING 상태만 저장할 수 있습니다.'
        this.error = true
        this.loading = false
        return
      }

      try {
        const payload = {
          title: this.issue.title,
          description: this.issue.description,
          status: this.issue.status,
          userId: this.selectedUserId === 0 ? null : this.selectedUserId,
        }

        if (this.isEditMode) {
          await axios.patch(`${this.API_BASE_URL}/issue/${this.issue.id}`, payload)
          this.message = '이슈가 성공적으로 수정되었습니다!'
          await this.fetchIssue(this.issue.id)
        } else {
          const res = await axios.post(`${this.API_BASE_URL}/issue`, payload)
          this.issue = res.data
          this.message = '이슈가 성공적으로 생성되었습니다!'
          this.$router.push('/')
        }
      } catch (err) {
        console.error('이슈 처리 중 오류:', err)
        this.message = `오류 발생: ${err.response?.data?.error || err.message}`
        this.error = true
      } finally {
        this.loading = false
      }
    },
    goBack() {
      this.$router.push('/')
    },
  },
}
</script>

<style scoped>
/* Scoped styles remain unchanged */
</style>
