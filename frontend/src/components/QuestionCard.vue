<script setup lang="ts">
import { ref, watch } from 'vue'
import type { IQuestion } from '../api/api'

const props = defineProps<{
  question: IQuestion
  index: number
  total: number
}>()

const emit = defineEmits<{
  (e: 'submit', optionIds: number[]): void
}>()

const selected = ref<number[]>([])

watch(() => props.question, () => {
  selected.value = []
})

const isMultiple = () => props.question.questionType === 'multiple'

function toggleOption(optionId: number) {
  if (isMultiple()) {
    const idx = selected.value.indexOf(optionId)
    if (idx === -1) {
      selected.value = [...selected.value, optionId]
    } else {
      selected.value = selected.value.filter((id) => id !== optionId)
    }
  } else {
    selected.value = [optionId]
  }
}

function isSelected(optionId: number) {
  return selected.value.includes(optionId)
}

function handleSubmit() {
  if (selected.value.length === 0) return
  emit('submit', [...selected.value])
}

const isLast = () => props.index === props.total - 1
</script>

<template>
  <article class="bg-white border border-line rounded-card p-4">
    <p class="text-ink-soft text-sm m-0">Question {{ index + 1 }} of {{ total }}</p>
    <p class="text-ink-soft text-sm mt-0.5 mb-2">
      {{ isMultiple() ? 'Multiple choice (checkboxes)' : 'Single choice (radio)' }}
    </p>
    <h3 class="font-heading font-bold text-lg m-0 mb-3">{{ question.questionText }}</h3>

    <img
      v-if="question.photoUrl"
      :src="question.photoUrl"
      alt="Question illustration"
      class="w-full max-h-72 object-cover rounded-card border border-line mb-4 bg-[#e5f7f5]"
    />

    <div class="flex flex-col gap-2 mt-2">
      <label
        v-for="option in question.options"
        :key="option.optionId"
        class="flex items-center gap-3 border rounded-btn px-3 py-2.5 cursor-pointer transition-colors select-none"
        :class="isSelected(option.optionId) ? 'border-brand bg-[#f0fdf9]' : 'border-line hover:bg-gray-50'"
        @click="toggleOption(option.optionId)"
      >
        <input
          :type="isMultiple() ? 'checkbox' : 'radio'"
          :checked="isSelected(option.optionId)"
          class="accent-brand"
          @change.prevent
        />
        <span class="text-sm">{{ option.optionText }}</span>
      </label>
    </div>

    <div class="mt-4">
      <button
        class="btn-main"
        :disabled="selected.length === 0"
        @click="handleSubmit"
      >
        {{ isLast() ? 'Finish Course' : 'Submit Answer' }}
      </button>
    </div>
  </article>
</template>
