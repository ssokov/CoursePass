<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../api/client'
import type { IExamSummary, ICourseSummary } from '../api/api'

const loading = ref(true)
const history = ref<IExamSummary[]>([])
const courseNames = ref<Map<number, string>>(new Map())

onMounted(async () => {
  try {
    const [exams, courses] = await Promise.all([
      api.exam.history({ page: 1, pageSize: 50 }),
      api.course.list({ page: 1, pageSize: 50 }),
    ])
    history.value = exams
    courses.forEach((c: ICourseSummary) => courseNames.value.set(c.courseId, c.title))
  } finally {
    loading.value = false
  }
})

function formatDate(iso: string) {
  return iso ? iso.slice(0, 10) : '—'
}

function courseName(courseId: number) {
  return courseNames.value.get(courseId) ?? `Course #${courseId}`
}
</script>

<template>
  <section class="panel">
    <h2 class="font-heading font-bold text-lg mb-1">Course History</h2>
    <p class="text-ink-soft text-sm mb-4">
      Result of passed courses: number of correct answers, final status and date.
    </p>

    <div v-if="loading" class="text-ink-soft text-sm py-4">Loading history…</div>

    <template v-else>
      <div v-if="history.length === 0" class="text-ink-soft text-sm py-4">
        No exams taken yet.
      </div>
      <table v-else class="w-full border-collapse text-sm">
        <thead>
          <tr>
            <th class="text-left text-xs uppercase tracking-wide text-ink-soft py-2.5 px-2 border-b border-line">
              Course
            </th>
            <th class="text-left text-xs uppercase tracking-wide text-ink-soft py-2.5 px-2 border-b border-line">
              Correct Answers
            </th>
            <th class="text-left text-xs uppercase tracking-wide text-ink-soft py-2.5 px-2 border-b border-line">
              Result
            </th>
            <th class="text-left text-xs uppercase tracking-wide text-ink-soft py-2.5 px-2 border-b border-line">
              Date
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="exam in history" :key="exam.examId">
            <td class="py-2.5 px-2 border-b border-line">{{ courseName(exam.courseId) }}</td>
            <td class="py-2.5 px-2 border-b border-line">
              {{ Math.round(exam.finalScore) }}%
            </td>
            <td class="py-2.5 px-2 border-b border-line">
              <span :class="exam.status === 'passed' ? 'status-ok' : 'status-bad'">
                {{ exam.status === 'passed' ? 'Passed' : 'Failed' }}
              </span>
            </td>
            <td class="py-2.5 px-2 border-b border-line text-ink-soft">
              {{ formatDate(exam.finishedAt) }}
            </td>
          </tr>
        </tbody>
      </table>
    </template>
  </section>
</template>
