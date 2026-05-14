<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/lib/api'
import { Loader2, Eye, EyeOff } from 'lucide-vue-next'
import { confirmSessionCache } from '@/router'

const router = useRouter()

const key = ref('')
const loading = ref(false)
const error = ref('')
const showKey = ref(false)

const login = async () => {
  if (!key.value.trim()) return
  loading.value = true
  error.value = ''

  try {
    await api.post('/auth/login', { key: key.value })
    confirmSessionCache()
    router.push('/admin')
  } catch (err) {
    error.value = err.response?.data?.error || 'Terjadi kesalahan, coba lagi'
    // Shake animation
    const input = document.getElementById('admin-key-input')
    input?.classList.add('lg-shake')
    setTimeout(() => input?.classList.remove('lg-shake'), 500)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="lg-root">
    <!-- Background grain + glow -->
    <div class="lg-bg">
      <div class="lg-glow" />
    </div>

    <div class="lg-card">
      <!-- Brand -->
      <div class="lg-brand">
        <div class="lg-brand-icon">
          <svg width="22" height="22" viewBox="0 0 20 20" fill="none">
            <path d="M4 4h5v5H4zM11 4h5v5h-5zM4 11h5v5H4zM11 11h5v5h-5z" fill="currentColor" opacity=".3"/>
            <path d="M7 2L2 7l5 5 5-5z" fill="currentColor"/>
          </svg>
        </div>
        <div>
          <h1 class="lg-brand-name">V-STREAM <span class="lg-brand-tag">PRO</span></h1>
          <p class="lg-brand-sub">Admin Dashboard</p>
        </div>
      </div>

      <!-- Divider -->
      <div class="lg-divider" />

      <!-- Form -->
      <div class="lg-form">
        <div class="lg-form-header">
          <h2 class="lg-title">Masuk ke Dashboard</h2>
          <p class="lg-subtitle">Masukkan admin key untuk melanjutkan</p>
        </div>

        <form @submit.prevent="login">
          <!-- Error message -->
          <Transition name="lg-error">
            <div v-if="error" class="lg-error-box">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
              </svg>
              {{ error }}
            </div>
          </Transition>

          <!-- Key input -->
          <div class="lg-field">
            <label class="lg-label">Admin Key</label>
            <div class="lg-input-wrap">
              <input
                id="admin-key-input"
                v-model="key"
                :type="showKey ? 'text' : 'password'"
                placeholder="Masukkan admin key..."
                autocomplete="current-password"
                required
                class="lg-input"
              />
              <button
                type="button"
                class="lg-eye-btn"
                @click="showKey = !showKey"
                tabindex="-1"
              >
                <EyeOff v-if="showKey" :size="15" />
                <Eye v-else :size="15" />
              </button>
            </div>
          </div>

          <!-- Submit -->
          <button
            type="submit"
            class="lg-submit"
            :disabled="loading || !key.trim()"
          >
            <Loader2 v-if="loading" :size="15" class="lg-spin" />
            <span>{{ loading ? 'Memverifikasi...' : 'Masuk' }}</span>
          </button>
        </form>
      </div>

      <!-- Footer note -->
      <p class="lg-note">
        Akses terbatas. Hanya administrator yang berwenang.
      </p>
    </div>
  </div>
</template>

<style scoped>

/* ── Root ───────────────────────────────────────────────────*/
.lg-root {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  background: #09090f;
  font-family: 'DM Sans', system-ui, sans-serif;
  position: relative;
  overflow: hidden;
}

/* ── Background ─────────────────────────────────────────────*/
.lg-bg {
  position: fixed;
  inset: 0;
  pointer-events: none;
}
.lg-glow {
  position: absolute;
  top: -20%;
  left: 50%;
  transform: translateX(-50%);
  width: 600px;
  height: 400px;
  border-radius: 50%;
  background: radial-gradient(ellipse, rgba(99,102,241,0.12) 0%, transparent 70%);
}

/* ── Card ───────────────────────────────────────────────────*/
.lg-card {
  width: 100%;
  max-width: 380px;
  background: #0f0f1a;
  border: 1px solid rgba(255,255,255,0.07);
  border-radius: 20px;
  overflow: hidden;
  box-shadow:
    0 0 0 1px rgba(255,255,255,0.03),
    0 40px 80px rgba(0,0,0,0.4);
  position: relative;
  z-index: 1;
  animation: card-in 0.4s cubic-bezier(0.34,1.2,0.64,1);
}
@keyframes card-in {
  from { opacity: 0; transform: translateY(20px) scale(0.97); }
}

/* ── Brand ──────────────────────────────────────────────────*/
.lg-brand {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 24px 24px 20px;
}
.lg-brand-icon {
  width: 44px; height: 44px;
  border-radius: 12px;
  background: rgba(99,102,241,0.12);
  color: #818cf8;
  display: grid; place-items: center;
  flex-shrink: 0;
  border: 1px solid rgba(99,102,241,0.2);
}
.lg-brand-name {
  font-size: 15px;
  font-weight: 700;
  letter-spacing: 0.05em;
  color: #fff;
  display: flex;
  align-items: center;
  gap: 6px;
}
.lg-brand-tag {
  font-size: 9px;
  font-weight: 600;
  letter-spacing: 0.08em;
  color: #818cf8;
  background: rgba(99,102,241,0.12);
  border: 1px solid rgba(99,102,241,0.25);
  padding: 2px 6px;
  border-radius: 4px;
}
.lg-brand-sub {
  font-size: 12px;
  color: #64748b;
  margin-top: 2px;
}

.lg-divider {
  height: 1px;
  background: rgba(255,255,255,0.06);
  margin: 0 24px;
}

/* ── Form ───────────────────────────────────────────────────*/
.lg-form { padding: 24px; }

.lg-form-header { margin-bottom: 20px; }
.lg-title {
  font-size: 16px;
  font-weight: 600;
  color: #fff;
  margin-bottom: 4px;
}
.lg-subtitle { font-size: 12.5px; color: #64748b; }

/* ── Error ──────────────────────────────────────────────────*/
.lg-error-box {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  border-radius: 8px;
  background: rgba(239,68,68,0.08);
  border: 1px solid rgba(239,68,68,0.2);
  color: #f87171;
  font-size: 12.5px;
  margin-bottom: 16px;
}
.lg-error-enter-active { animation: error-in 0.2s ease; }
.lg-error-leave-active { animation: error-in 0.15s ease reverse; }
@keyframes error-in { from { opacity:0; transform: translateY(-4px); } }

/* ── Field ──────────────────────────────────────────────────*/
.lg-field { display: flex; flex-direction: column; gap: 6px; margin-bottom: 16px; }
.lg-label { font-size: 12px; font-weight: 500; color: #94a3b8; }

.lg-input-wrap { position: relative; }
.lg-input {
  width: 100%;
  background: #14141f;
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 10px;
  padding: 11px 42px 11px 14px;
  font-size: 13.5px;
  color: #e2e8f0;
  font-family: inherit;
  outline: none;
  transition: all 0.15s;
  box-sizing: border-box;
}
.lg-input:focus {
  border-color: rgba(99,102,241,0.5);
  box-shadow: 0 0 0 3px rgba(99,102,241,0.08);
}
.lg-input::placeholder { color: #475569; }

.lg-eye-btn {
  position: absolute;
  right: 12px; top: 50%;
  transform: translateY(-50%);
  background: none; border: none;
  color: #475569;
  cursor: pointer;
  display: grid; place-items: center;
  padding: 4px;
  border-radius: 4px;
  transition: color 0.15s;
}
.lg-eye-btn:hover { color: #94a3b8; }

/* Shake animation untuk input salah */
@keyframes shake {
  0%, 100% { transform: translateX(0); }
  20%, 60% { transform: translateX(-6px); }
  40%, 80% { transform: translateX(6px); }
}
.lg-shake { animation: shake 0.4s ease; }

/* ── Submit button ──────────────────────────────────────────*/
.lg-submit {
  width: 100%;
  padding: 12px;
  border-radius: 10px;
  border: none;
  background: #6366f1;
  color: #fff;
  font-size: 13.5px;
  font-weight: 600;
  font-family: inherit;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: all 0.15s;
  box-shadow: 0 0 20px rgba(99,102,241,0.3);
}
.lg-submit:hover:not(:disabled) {
  background: #818cf8;
  box-shadow: 0 0 28px rgba(99,102,241,0.45);
  transform: translateY(-1px);
}
.lg-submit:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.lg-spin { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

/* ── Footer note ────────────────────────────────────────────*/
.lg-note {
  text-align: center;
  font-size: 11px;
  color: #334155;
  padding: 0 24px 20px;
}
</style>