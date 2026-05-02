<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const countdown = ref(5)
let timer = null

onMounted(() => {
  timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) router.push('/')
  }, 1000)
})

onUnmounted(() => clearInterval(timer))
</script>

<template>
  <div class="nf-root">
    <div class="nf-card">
      <div class="nf-code">404</div>
      <h1 class="nf-title">Halaman Tidak Ditemukan</h1>
      <p class="nf-desc">
        URL yang kamu akses tidak ada. Kamu akan diarahkan ke halaman utama dalam
        <span class="nf-count">{{ countdown }}</span> detik.
      </p>
      <button class="nf-btn" @click="router.push('/')">Kembali ke Home</button>
    </div>
  </div>
</template>

<style scoped>
.nf-root {
  min-height: 100vh;
  background: #07070e;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: 'DM Sans', system-ui, sans-serif;
  padding: 24px;
}

.nf-card {
  text-align: center;
  max-width: 420px;
}

.nf-code {
  font-family: 'Bebas Neue', sans-serif;
  font-size: clamp(100px, 20vw, 160px);
  line-height: 1;
  color: transparent;
  -webkit-text-stroke: 2px rgba(99, 102, 241, 0.4);
  letter-spacing: 0.05em;
  margin-bottom: 8px;
  animation: nf-pulse 3s ease-in-out infinite;
}
@keyframes nf-pulse {
  0%, 100% { -webkit-text-stroke-color: rgba(99,102,241,0.4); }
  50%       { -webkit-text-stroke-color: rgba(99,102,241,0.8); }
}

.nf-title {
  font-size: 22px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 12px;
  letter-spacing: -0.02em;
}

.nf-desc {
  font-size: 14px;
  color: #64748b;
  line-height: 1.6;
  margin-bottom: 28px;
}

.nf-count {
  color: #818cf8;
  font-weight: 700;
  font-size: 16px;
}

.nf-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 11px 28px;
  border-radius: 40px;
  border: none;
  background: #6366f1;
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  font-family: inherit;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.3);
}
.nf-btn:hover {
  background: #818cf8;
  transform: translateY(-2px);
  box-shadow: 0 12px 32px rgba(99, 102, 241, 0.4);
}
</style>