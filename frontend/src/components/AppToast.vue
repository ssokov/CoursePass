<script setup lang="ts">
import { useToast } from '../composables/useToast'
const { toasts } = useToast()
</script>

<template>
  <Teleport to="body">
    <div class="fixed bottom-5 right-5 flex flex-col gap-2 z-50 pointer-events-none">
      <TransitionGroup name="toast">
        <div
          v-for="t in toasts"
          :key="t.id"
          class="px-4 py-3 rounded-btn text-sm font-medium shadow-panel max-w-xs"
          :class="t.type === 'ok' ? 'bg-ok-bg text-ok' : 'bg-bad-bg text-bad'"
        >
          {{ t.message }}
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.25s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateX(20px);
}
</style>
