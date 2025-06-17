<!-- src/App.vue -->
<template>
  <div id="app" class="min-h-screen bg-gray-100">
    <header
      class="fixed top-0 w-full z-50 bg-gradient-to-r from-indigo-600 to-purple-600 text-white shadow-md"
    >
      <div class="container mx-auto flex justify-between items-center px-6 py-4">
        <router-link
          to="/"
          class="text-2xl font-bold tracking-tight hover:text-gray-200 transition"
        >
          이슈 관리 시스템
        </router-link>
        <nav>
          <ul class="flex space-x-3">
            <li>
              <router-link
                to="/"
                class="inline-flex items-center px-4 py-2 rounded-lg bg-white text-indigo-600 font-semibold shadow hover:bg-indigo-100 transition"
              >
                <i class="fas fa-list-alt mr-2"></i> 이슈 목록
              </router-link>
            </li>
            <li>
              <router-link
                to="/issue/new"
                class="inline-flex items-center px-4 py-2 rounded-lg bg-indigo-500 text-white font-semibold shadow hover:bg-indigo-600 transition"
              >
                <i class="fas fa-plus-circle mr-2"></i> 새 이슈
              </router-link>
            </li>
          </ul>
        </nav>
      </div>
    </header>

    <main class="pt-28 px-4 container mx-auto">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import '@fortawesome/fontawesome-free/css/all.css'
</script>

<style>
#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
