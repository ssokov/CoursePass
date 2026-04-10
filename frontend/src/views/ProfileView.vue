<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../api/client'
import type { IStudent } from '../api/api'

const student = ref<IStudent | null>(null)
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    student.value = await api.course.me()
  } catch {
    error.value = 'Failed to load profile'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <section class="panel max-w-lg">
    <h2 class="font-heading font-bold text-lg mb-4">Profile</h2>

    <div v-if="loading" class="text-ink-soft text-sm py-4">Loading…</div>

    <div v-else-if="error" class="p-3 rounded-btn bg-bad-bg text-bad text-sm">{{ error }}</div>

    <dl v-else-if="student" class="flex flex-col gap-3">
      <div class="flex flex-col gap-0.5">
        <dt class="text-xs text-ink-soft uppercase tracking-wide">First Name</dt>
        <dd class="text-sm font-medium">{{ student.firstName }}</dd>
      </div>
      <div class="flex flex-col gap-0.5">
        <dt class="text-xs text-ink-soft uppercase tracking-wide">Last Name</dt>
        <dd class="text-sm font-medium">{{ student.lastName }}</dd>
      </div>
      <div class="flex flex-col gap-0.5">
        <dt class="text-xs text-ink-soft uppercase tracking-wide">Login</dt>
        <dd class="text-sm font-medium">{{ student.login }}</dd>
      </div>
      <div class="flex flex-col gap-0.5">
        <dt class="text-xs text-ink-soft uppercase tracking-wide">Email</dt>
        <dd class="text-sm font-medium">{{ student.email }}</dd>
      </div>
    </dl>
  </section>
</template>
