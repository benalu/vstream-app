<script setup>
import { CheckCircle2, XCircle } from 'lucide-vue-next'
import { useToast } from '@/composables/useToast'
const { toasts } = useToast()
</script>

<template>
  <Teleport to="body">
    <div class="toast-wrap">
      <TransitionGroup name="toast">
        <div v-for="t in toasts" :key="t.id" class="toast" :class="`toast--${t.type}`">
          <CheckCircle2 v-if="t.type === 'success'" :size="15" />
          <XCircle      v-else                       :size="15" />
          {{ t.message }}
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
.toast-wrap {
  position: fixed;
  bottom: 24px;
  right: 24px;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 8px;
  pointer-events: none;
}
.toast {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 11px 16px;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 500;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(8px);
  min-width: 260px;
  max-width: 360px;
}
.toast--success {
  background: rgba(16, 185, 129, 0.15);
  border: 1px solid rgba(16, 185, 129, 0.3);
  color: #34d399;
}
.toast--error {
  background: rgba(239, 68, 68, 0.12);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #f87171;
}
.toast-enter-active {
  animation: toast-in 0.25s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.toast-leave-active {
  animation: toast-out 0.2s ease forwards;
}
@keyframes toast-in {
  from { opacity: 0; transform: translateY(12px) scale(0.95); }
}
@keyframes toast-out {
  to { opacity: 0; transform: translateY(4px) scale(0.97); }
}
</style>