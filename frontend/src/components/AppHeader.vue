<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import { api } from '../api/client'
import type { IStudent } from '../api/api'

const router = useRouter()
const route = useRoute()
const { logout, isAuthenticated } = useAuth()

const student = ref<IStudent | null>(null)

onMounted(async () => {
  if (isAuthenticated.value) {
    try {
      student.value = await api.course.me()
    } catch {
      // ignore
    }
  }
})

function handleLogout() {
  logout()
  router.push('/login')
}
</script>

<template>
  <header class="panel flex justify-between items-center gap-3 flex-wrap">
    <h1 class="font-heading font-bold text-2xl m-0">Backend Courses</h1>
    <div class="flex items-center gap-2 flex-wrap">
      <router-link
        to="/courses"
        class="btn-ghost text-sm"
        :class="route.path === '/courses' ? 'border-brand text-brand' : ''"
      >
        Courses
      </router-link>
      <router-link
        to="/history"
        class="btn-ghost text-sm"
        :class="route.path === '/history' ? 'border-brand text-brand' : ''"
      >
        History
      </router-link>
      <router-link
        to="/profile"
        class="btn-ghost text-sm"
        :class="route.path === '/profile' ? 'border-brand text-brand' : ''"
      >
        {{ student ? `${student.firstName} ${student.lastName}` : 'Profile' }}
      </router-link>
      <button class="btn-ghost text-sm" @click="handleLogout">Logout</button>
    </div>
  </header>
</template>
