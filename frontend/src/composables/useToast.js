import { ref } from 'vue'

const toasts = ref([])
let _id = 0

export function useToast() {
  const show = (message, type = 'success', duration = 3000) => {
    const id = ++_id
    toasts.value.push({ id, message, type })
    setTimeout(() => {
      toasts.value = toasts.value.filter(t => t.id !== id)
    }, duration)
  }

  const error = (message, duration = 4000) => show(message, 'error', duration)

  return { toasts, show, error }
}