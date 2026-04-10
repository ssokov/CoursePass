<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../api/client'
import QuestionCard from '../components/QuestionCard.vue'
import type { IQuestion, IExamResult, ICourse } from '../api/api'

const route = useRoute()
const router = useRouter()

const courseId = Number(route.params.courseId)

const loading = ref(true)
const error = ref('')

const examId = ref(0)
const questionIds = ref<number[]>([])
const currentIndex = ref(0)
const currentQuestion = ref<IQuestion | null>(null)
const course = ref<ICourse | null>(null)
const result = ref<IExamResult | null>(null)
const submitting = ref(false)

const progress = computed(() => {
  const total = questionIds.value.length
  if (total === 0) return 0
  return Math.round((currentIndex.value / total) * 100)
})

onMounted(async () => {
  try {
    const [examStart, courseData] = await Promise.all([
      api.exam.start({ courseID: courseId }),
      api.course.byID({ courseID: courseId }),
    ])
    examId.value = examStart.examId
    questionIds.value = examStart.questionIds
    course.value = courseData

    await loadQuestion(0)
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Failed to start exam'
  } finally {
    loading.value = false
  }
})

async function loadQuestion(index: number) {
  currentIndex.value = index
  currentQuestion.value = await api.exam.getQuestion({
    examID: examId.value,
    questionID: questionIds.value[index],
  })
}

async function handleAnswer(optionIds: number[]) {
  if (submitting.value) return
  submitting.value = true

  try {
    await api.exam.answer({
      examID: examId.value,
      questionID: questionIds.value[currentIndex.value],
      optionIDs: optionIds,
    })

    const nextIndex = currentIndex.value + 1

    if (nextIndex < questionIds.value.length) {
      await loadQuestion(nextIndex)
    } else {
      result.value = await api.exam.submit({ examID: examId.value })
    }
  } finally {
    submitting.value = false
  }
}

function goToHistory() {
  router.push('/history')
}
</script>

<template>
  <section class="panel">
    <template v-if="loading">
      <p class="text-ink-soft text-sm">Starting exam…</p>
    </template>

    <template v-else-if="error">
      <p class="text-bad text-sm">{{ error }}</p>
    </template>

    <template v-else-if="result">
      <h2 class="font-heading font-bold text-lg mb-1">Exam Complete</h2>
      <p class="text-ink-soft text-sm mb-4">{{ course?.title }}</p>
      <div class="bg-white border border-line rounded-card p-6 text-center">
        <p class="text-4xl font-heading font-bold mb-2">
          {{ result.correctAnswers }}/{{ result.totalQuestions }}
        </p>
        <p class="text-ink-soft text-sm mb-4">correct answers</p>
        <span :class="result.status === 'passed' ? 'status-ok' : 'status-bad'" class="text-base px-4 py-1.5">
          {{ result.status === 'passed' ? 'Passed' : 'Failed' }}
        </span>
        <div class="mt-6">
          <button class="btn-main" @click="goToHistory">View History</button>
        </div>
      </div>
    </template>

    <template v-else-if="currentQuestion">
      <h2 class="font-heading font-bold text-lg mb-0.5">Pass Course: {{ course?.title }}</h2>
      <p class="text-ink-soft text-sm mb-3">Answer all {{ questionIds.length }} questions.</p>

      <div class="h-2 bg-[#ebe4d5] rounded-full overflow-hidden mb-4">
        <div
          class="h-full bg-gradient-to-r from-brand-light to-brand transition-all duration-300"
          :style="{ width: progress + '%' }"
        />
      </div>

      <QuestionCard
        :question="currentQuestion"
        :index="currentIndex"
        :total="questionIds.length"
        @submit="handleAnswer"
      />
    </template>
  </section>
</template>
