<script setup lang="ts">
import { useRouter } from 'vue-router'
import type { ICourseSummary } from '../api/api'

const props = defineProps<{ course: ICourseSummary; questionCount?: number }>()
const router = useRouter()

function pass() {
  router.push(`/exam/${props.course.courseId}`)
}
</script>

<template>
  <article class="bg-white border border-line rounded-card p-4 flex flex-col gap-3">
    <h3 class="font-heading font-bold text-lg m-0">{{ course.title }}</h3>
    <p v-if="questionCount !== undefined" class="text-ink-soft text-sm m-0">
      Questions: {{ questionCount }}
    </p>
    <button class="btn-main w-full mt-auto" @click="pass">Pass</button>
  </article>
</template>
