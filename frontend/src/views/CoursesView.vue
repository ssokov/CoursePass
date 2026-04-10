<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../api/client'
import CourseCard from '../components/CourseCard.vue'
import type { ICourseSummary } from '../api/api'

const loading = ref(true)
const courses = ref<ICourseSummary[]>([])

onMounted(async () => {
  try {
    courses.value = await api.course.list({ page: 1, pageSize: 50 })
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <section class="panel">
    <h2 class="font-heading font-bold text-lg mb-1">Choose a Course</h2>
    <p class="text-ink-soft text-sm mb-4">
      Select any available course and click <strong>Pass</strong> to start questions.
    </p>

    <div v-if="loading" class="text-ink-soft text-sm py-4">Loading courses…</div>

    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <CourseCard
        v-for="course in courses"
        :key="course.courseId"
        :course="course"
      />
    </div>
  </section>
</template>
