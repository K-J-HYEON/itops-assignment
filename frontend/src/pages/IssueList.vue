<template>
  <div class="max-w-4xl mx-auto">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold text-gray-800">📝 이슈 목록</h1>
      <router-link
        to="/issue/new"
        class="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 transition"
      >
        + 새 이슈
      </router-link>
    </div>

    <div class="mb-4">
      <div class="flex flex-wrap gap-2">
        <button
          v-for="s in statuses"
          :key="s"
          @click="selectedStatus = s"
          :class="[
            'px-4 py-1 rounded-full text-sm font-medium',
            selectedStatus === s
              ? 'bg-indigo-600 text-white'
              : 'bg-gray-200 text-gray-800 hover:bg-gray-300',
          ]"
        >
          {{ s }}
        </button>
      </div>
    </div>

    <div v-if="loading" class="text-gray-500 py-10 flex justify-center items-center">
      <i class="fas fa-spinner fa-spin mr-2"></i> 불러오는 중...
    </div>

    <div v-else-if="filteredIssues.length === 0" class="text-gray-500 mt-10 text-center">
      이슈가 없습니다.
    </div>

    <div v-else class="space-y-4">
      <div
        v-for="issue in filteredIssues"
        :key="issue.id"
        class="bg-white p-4 rounded-lg shadow hover:shadow-md transition"
      >
        <div class="text-lg font-semibold text-gray-800 mb-1">{{ issue.title }}</div>
        <div class="text-sm text-gray-600 mb-2">
          <span class="mr-4"
            >상태: <strong>{{ issue.status }}</strong></span
          >
          <span
            >담당자: <strong>{{ issue.assignee || '미지정' }}</strong></span
          >
        </div>
        <router-link
          :to="`/issue/edit/${issue.id}`"
          class="text-sm text-indigo-600 hover:underline"
        >
          수정하기 →
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const issues = ref([])
const loading = ref(true)
const selectedStatus = ref('전체')
const statuses = ['전체', '대기중', '진행중', '완료됨', '취소됨']

onMounted(() => {
  setTimeout(() => {
    issues.value = [
      { id: 1, title: '로그인 오류 수정', status: '진행중', assignee: '홍길동' },
      { id: 2, title: 'UI 리팩터링', status: '대기중', assignee: null },
      { id: 3, title: '백엔드 API 수정', status: '완료됨', assignee: '김개발' },
    ]
    loading.value = false
  }, 1000)
})

const filteredIssues = computed(() => {
  if (selectedStatus.value === '전체') return issues.value
  return issues.value.filter((issue) => issue.status === selectedStatus.value)
})
</script>
