<script setup>
import { ref, onMounted } from 'vue'
import api from '@/lib/api'  // ← ganti axios
import { Film, Users, Activity, TrendingUp, PlayCircle, Clock } from 'lucide-vue-next'

const stats = ref({ totalMovies: 0, activeUsers: 0, totalViews: 0, serverLoad: 'Normal' })
const recentMovies = ref([])
const loading = ref(true)

const fetchDashboardData = async () => {
  try {
    const res = await api.get('/admin/dashboard')  // ← baseURL sudah di api.js
    const data = res.data.data
    stats.value = {
      totalMovies: data.total_movies,
      activeUsers: data.active_users,
      totalViews: data.total_views,
      serverLoad: data.server_load
    }
    recentMovies.value = data.recent_movies || []
  } catch (err) {
    console.error("Gagal memuat data dashboard", err)
  } finally {
    loading.value = false
  }
}

const playbackLogs = ref({ logs: [], summary: [] })
const clearingLogs = ref(false)

const fetchPlaybackLogs = async () => {
  try {
    const res = await api.get('/admin/playback-logs')
    playbackLogs.value = res.data.data
  } catch (err) {
    console.error('Gagal fetch playback logs', err)
  }
}

const clearLogs = async () => {
  if (!confirm('Hapus semua log error?')) return
  clearingLogs.value = true
  try {
    await api.delete('/admin/playback-logs')
    playbackLogs.value = { logs: [], summary: [] }
  } finally {
    clearingLogs.value = false
  }
}

onMounted(() => {
  fetchDashboardData()
  fetchPlaybackLogs()
})
</script>

<template>
  <div>
    <!-- Header -->
    <div class="vs-page-header">
      <div>
        <h1 class="vs-page-title">Overview</h1>
        <p class="vs-page-subtitle">Selamat datang kembali, Admin — berikut ringkasan platform V-STREAM hari ini.</p>
      </div>
      <div class="db-date-chip">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="4" width="18" height="18" rx="2"/><line x1="16" y1="2" x2="16" y2="6"/>
          <line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/>
        </svg>
        {{ new Date().toLocaleDateString('id-ID', { weekday:'long', day:'numeric', month:'long', year:'numeric' }) }}
      </div>
    </div>

    <!-- Stat cards -->
    <div class="db-stats">
      <!-- Total Film -->
      <div class="db-stat-card">
        <div class="db-stat-card__header">
          <span class="db-stat-card__label">Total Film</span>
          <div class="db-stat-card__icon db-stat-card__icon--indigo">
            <Film :size="16" />
          </div>
        </div>
        <div class="db-stat-card__value">
          <span v-if="loading" class="db-skeleton db-skeleton--value" />
          <span v-else>{{ stats.totalMovies.toLocaleString() }}</span>
        </div>
        <div class="db-stat-card__trend db-stat-card__trend--up">
          <TrendingUp :size="12" />
          +12% bulan ini
        </div>
      </div>

      <!-- Total Tayangan -->
      <div class="db-stat-card">
        <div class="db-stat-card__header">
          <span class="db-stat-card__label">Total Tayangan</span>
          <div class="db-stat-card__icon db-stat-card__icon--pink">
            <PlayCircle :size="16" />
          </div>
        </div>
        <div class="db-stat-card__value">
          <span v-if="loading" class="db-skeleton db-skeleton--value" />
          <span v-else>{{ stats.totalViews.toLocaleString() }}</span>
        </div>
        <div class="db-stat-card__trend db-stat-card__trend--up">
          <TrendingUp :size="12" />
          +5.4% minggu ini
        </div>
      </div>

      <!-- Pengguna Aktif -->
      <div class="db-stat-card">
        <div class="db-stat-card__header">
          <span class="db-stat-card__label">Pengguna Aktif</span>
          <div class="db-stat-card__icon db-stat-card__icon--blue">
            <Users :size="16" />
          </div>
        </div>
        <div class="db-stat-card__value">
          <span v-if="loading" class="db-skeleton db-skeleton--value" />
          <span v-else>{{ stats.activeUsers.toLocaleString() }}</span>
        </div>
        <div class="db-stat-card__trend db-stat-card__trend--up">
          <TrendingUp :size="12" />
          +120 hari ini
        </div>
      </div>

      <!-- Beban Server -->
      <div class="db-stat-card">
        <div class="db-stat-card__header">
          <span class="db-stat-card__label">Beban Server</span>
          <div class="db-stat-card__icon db-stat-card__icon--green">
            <Activity :size="16" />
          </div>
        </div>
        <div class="db-stat-card__value">
          <span v-if="loading" class="db-skeleton db-skeleton--value" />
          <span v-else>{{ stats.serverLoad }}</span>
        </div>
        <div class="db-stat-card__trend">
          <Activity :size="12" />
          Status operasional
        </div>
      </div>
    </div>

    <!-- Bottom grid -->
    <div class="db-grid">

      <!-- Recent movies table -->
      <div class="vs-card db-recent">
        <div class="db-section-header">
          <h2 class="db-section-title">Film Baru Ditambahkan</h2>
          <router-link to="/admin/movies" class="db-see-all">Lihat Semua →</router-link>
        </div>

        <!-- Loading -->
        <div v-if="loading" class="db-loading-list">
          <div v-for="i in 4" :key="i" class="db-movie-row">
            <div class="db-skeleton db-skeleton--poster" />
            <div class="db-skeleton-info">
              <div class="db-skeleton db-skeleton--title" />
              <div class="db-skeleton db-skeleton--sub" />
            </div>
          </div>
        </div>

        <!-- List -->
        <div v-else>
          <div
            v-for="movie in recentMovies"
            :key="movie.id"
            class="db-movie-row"
          >
            <img :src="movie.poster" class="db-movie-poster" :alt="movie.title" />
            <div class="db-movie-info">
              <h3 class="db-movie-title">{{ movie.title }}</h3>
              <div class="db-movie-meta">
                <span><Clock :size="11" /> {{ movie.duration }}</span>
                <span class="vs-badge vs-badge--slate">{{ movie.year }}</span>
              </div>
            </div>
            <span class="db-movie-date">
              {{ new Date(movie.created_at).toLocaleDateString('id-ID', { day:'numeric', month:'short' }) }}
            </span>
          </div>

          <div v-if="recentMovies.length === 0" class="db-empty">
            <Film :size="28" />
            <p>Belum ada film ditambahkan</p>
          </div>
        </div>
      </div>

      <!-- Quick actions -->
      <div class="db-sidebar-panel">
        <div class="vs-card db-quick-actions">
          <div class="db-section-header">
            <h2 class="db-section-title">Aksi Cepat</h2>
          </div>
          <div class="db-qa-list">
            <router-link to="/admin/movies" class="db-qa-item db-qa-item--accent">
              <div class="db-qa-item__icon">
                <Film :size="15" />
              </div>
              <div>
                <span class="db-qa-item__label">Tambah Film Baru</span>
                <span class="db-qa-item__desc">Import via TMDB ID</span>
              </div>
              <svg class="db-qa-item__arrow" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M5 12h14M12 5l7 7-7 7"/></svg>
            </router-link>

            <button class="db-qa-item">
              <div class="db-qa-item__icon">
                <Activity :size="15" />
              </div>
              <div>
                <span class="db-qa-item__label">Cek Status Server</span>
                <span class="db-qa-item__desc">Monitor performa realtime</span>
              </div>
              <svg class="db-qa-item__arrow" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M5 12h14M12 5l7 7-7 7"/></svg>
            </button>
          </div>
        </div>

        <!-- Mini stat -->
        <div class="vs-card db-mini-stat">
          <p class="db-mini-stat__label">Uptime Platform</p>
          <p class="db-mini-stat__value">99.9%</p>
          <div class="db-uptime-bar">
            <div v-for="i in 30" :key="i" :class="['db-uptime-tick', i > 28 ? 'db-uptime-tick--warn' : '']" />
          </div>
          <p class="db-mini-stat__sub">30 hari terakhir</p>
        </div>
      </div>
      <!-- Playback Error Monitor -->
      <div class="vs-card" style="margin-top:14px">
        <div class="db-section-header">
          <h2 class="db-section-title">Playback Error Monitor</h2>
          <button
            class="vs-btn vs-btn--ghost"
            style="padding:5px 12px;font-size:12px"
            @click="clearLogs"
            :disabled="clearingLogs"
          >
            Hapus Log
          </button>
        </div>

        <!-- Summary: film paling sering error -->
        <div v-if="playbackLogs.summary?.length" class="db-plog-summary">
          <p class="db-plog-summary-label">Paling sering error</p>
          <div
            v-for="item in playbackLogs.summary"
            :key="item.movie_id"
            class="db-plog-summary-row"
          >
            <span class="db-plog-summary-title">{{ item.movie_title }}</span>
            <span class="db-plog-badge">{{ item.error_count }}x</span>
          </div>
        </div>

        <!-- Log kosong -->
        <div v-if="!playbackLogs.logs?.length" class="db-plog-empty">
          Belum ada error dilaporkan oleh user.
        </div>

        <!-- Log list -->
        <div v-else>
          <div
            v-for="log in playbackLogs.logs"
            :key="log.id"
            class="db-plog-row"
          >
            <span class="db-plog-dot" />
            <div class="db-plog-info">
              <span class="db-plog-title">{{ log.movie_title }}</span>
              <span class="db-plog-meta">
                <span class="db-plog-server">{{ log.server.toUpperCase() }}</span>
                · {{ log.error_type }}
                · {{ new Date(log.reported_at).toLocaleString('id-ID') }}
              </span>
            </div>
            <router-link to="/admin/movies" class="db-plog-fix">Fix →</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* ── Date chip ──────────────────────────────────────────── */
.db-date-chip {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--muted);
  padding: 6px 12px;
  border-radius: 20px;
  background: var(--surface);
  border: 1px solid var(--border);
  white-space: nowrap;
  align-self: flex-start;
  margin-top: 4px;
}

/* ── Stat cards grid ────────────────────────────────────── */
.db-stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px;
  margin-bottom: 20px;
}
@media (max-width: 1100px) { .db-stats { grid-template-columns: repeat(2,1fr); } }
@media (max-width: 540px)  { .db-stats { grid-template-columns: 1fr; } }

.db-stat-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 18px 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  position: relative;
  overflow: hidden;
  transition: border-color 0.2s, transform 0.2s;
}
.db-stat-card:hover {
  border-color: rgba(255,255,255,0.12);
  transform: translateY(-1px);
}
.db-stat-card::before {
  content: '';
  position: absolute;
  top: 0; left: 0; right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.06), transparent);
}

.db-stat-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.db-stat-card__label {
  font-size: 12px;
  font-weight: 500;
  color: var(--muted);
  letter-spacing: 0.02em;
}
.db-stat-card__icon {
  width: 32px; height: 32px;
  border-radius: 8px;
  display: grid; place-items: center;
}
.db-stat-card__icon--indigo { background: rgba(99,102,241,0.12); color: #818cf8; }
.db-stat-card__icon--pink   { background: rgba(236,72,153,0.12);  color: #f472b6; }
.db-stat-card__icon--blue   { background: rgba(59,130,246,0.12);  color: #60a5fa; }
.db-stat-card__icon--green  { background: rgba(16,185,129,0.12);  color: #34d399; }

.db-stat-card__value {
  font-size: 28px;
  font-weight: 700;
  color: #fff;
  letter-spacing: -0.03em;
  min-height: 34px;
}
.db-stat-card__trend {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 11.5px;
  color: var(--muted);
}
.db-stat-card__trend--up { color: var(--success); }

/* ── Bottom grid ────────────────────────────────────────── */
.db-grid {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: 14px;
  align-items: start;
}
@media (max-width: 900px) {
  .db-grid { grid-template-columns: 1fr; }
}

.db-section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 18px;
  border-bottom: 1px solid var(--border);
}
.db-section-title {
  font-size: 13.5px;
  font-weight: 600;
  color: #fff;
}
.db-see-all {
  font-size: 12px;
  color: var(--accent-hi);
  text-decoration: none;
  opacity: 0.8;
  transition: opacity 0.15s;
}
.db-see-all:hover { opacity: 1; }

/* ── Movie rows ─────────────────────────────────────────── */
.db-movie-row {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 12px 18px;
  border-bottom: 1px solid var(--border2);
  transition: background 0.15s;
}
.db-movie-row:last-child { border-bottom: none; }
.db-movie-row:hover { background: rgba(255,255,255,0.02); }

.db-movie-poster {
  width: 38px; height: 54px;
  border-radius: 6px;
  object-fit: cover;
  flex-shrink: 0;
  background: var(--surface2);
}
.db-movie-info { flex: 1; min-width: 0; }
.db-movie-title {
  font-size: 13px;
  font-weight: 500;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.db-movie-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 11px;
  color: var(--muted);
  margin-top: 4px;
}
.db-movie-meta span { display: flex; align-items: center; gap: 3px; }
.db-movie-date {
  font-size: 11px;
  color: var(--muted2);
  white-space: nowrap;
}

/* ── Quick actions ──────────────────────────────────────── */
.db-sidebar-panel { display: flex; flex-direction: column; gap: 14px; }

.db-qa-list { padding: 10px; display: flex; flex-direction: column; gap: 4px; }
.db-qa-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  border: none;
  background: transparent;
  color: var(--muted);
  font-size: 13px;
  font-family: var(--font);
  text-decoration: none;
  transition: all 0.15s;
  width: 100%;
  text-align: left;
}
.db-qa-item:hover { background: var(--border2); color: var(--text); }
.db-qa-item--accent:hover { background: var(--accent-lo); color: var(--accent-hi); }
.db-qa-item--accent:hover .db-qa-item__icon { background: rgba(99,102,241,0.2); color: var(--accent-hi); }

.db-qa-item__icon {
  width: 32px; height: 32px;
  border-radius: 8px;
  background: var(--border2);
  display: grid; place-items: center;
  flex-shrink: 0;
  transition: all 0.15s;
}
.db-qa-item__label {
  display: block;
  font-size: 12.5px;
  font-weight: 600;
  color: var(--text);
}
.db-qa-item__desc {
  display: block;
  font-size: 11px;
  color: var(--muted);
}
.db-qa-item__arrow { margin-left: auto; opacity: 0.4; flex-shrink: 0; }
.db-qa-item:hover .db-qa-item__arrow { opacity: 0.8; }

/* ── Mini stat ──────────────────────────────────────────── */
.db-mini-stat { padding: 16px 18px; }
.db-mini-stat__label { font-size: 11px; color: var(--muted); text-transform: uppercase; letter-spacing: 0.08em; font-weight: 600; }
.db-mini-stat__value { font-size: 32px; font-weight: 700; color: var(--success); letter-spacing: -0.03em; margin: 6px 0; }
.db-uptime-bar { display: flex; gap: 3px; margin: 8px 0; }
.db-uptime-tick { flex: 1; height: 24px; border-radius: 3px; background: rgba(16,185,129,0.3); }
.db-uptime-tick--warn { background: rgba(245,158,11,0.4); }
.db-mini-stat__sub { font-size: 11px; color: var(--muted2); }

/* ── Skeleton ───────────────────────────────────────────── */
.db-skeleton {
  background: linear-gradient(90deg, var(--surface2) 25%, var(--border) 50%, var(--surface2) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.4s infinite;
  border-radius: 6px;
  display: block;
}
@keyframes shimmer { to { background-position: -200% 0; } }
.db-skeleton--value { width: 80px; height: 34px; }
.db-skeleton--poster { width: 38px; height: 54px; border-radius: 6px; flex-shrink: 0; }
.db-skeleton--title { width: 60%; height: 14px; margin-bottom: 8px; }
.db-skeleton--sub { width: 40%; height: 11px; }

.db-loading-list { display: flex; flex-direction: column; }
.db-skeleton-info { flex: 1; padding: 4px 0; }

/* ── Playback Log ───────────────────────────────────────── */
.db-plog-summary {
  padding: 12px 18px;
  border-bottom: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.db-plog-summary-label {
  font-size: 10.5px; font-weight: 600;
  text-transform: uppercase; letter-spacing: 0.07em;
  color: var(--muted); margin-bottom: 2px;
}
.db-plog-summary-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12.5px;
}
.db-plog-summary-title { color: var(--text); }
.db-plog-badge {
  font-size: 11px; font-weight: 700; color: var(--danger);
  background: rgba(239,68,68,0.1);
  padding: 1px 7px; border-radius: 4px;
}

.db-plog-empty {
  padding: 32px; text-align: center;
  font-size: 13px; color: var(--muted2);
}

.db-plog-row {
  display: flex; align-items: center; gap: 12px;
  padding: 10px 18px;
  border-bottom: 1px solid var(--border2);
  background: rgba(239,68,68,0.02);
  transition: background 0.15s;
}
.db-plog-row:last-child { border-bottom: none; }
.db-plog-row:hover { background: rgba(239,68,68,0.05); }

.db-plog-dot {
  width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0;
  background: var(--danger);
  box-shadow: 0 0 6px var(--danger);
}
.db-plog-info { flex: 1; min-width: 0; }
.db-plog-title {
  display: block; font-size: 13px; font-weight: 500; color: #fff;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}
.db-plog-meta {
  display: flex; align-items: center; gap: 6px;
  font-size: 11px; color: var(--muted); margin-top: 2px;
}
.db-plog-server {
  font-weight: 700; color: var(--accent-hi);
  background: var(--accent-lo);
  padding: 1px 5px; border-radius: 3px;
}
.db-plog-fix {
  font-size: 11.5px; color: var(--accent-hi);
  text-decoration: none; flex-shrink: 0;
  opacity: 0.7; transition: opacity 0.15s;
}
.db-plog-fix:hover { opacity: 1; }

/* ── Empty ──────────────────────────────────────────────── */
.db-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 40px;
  color: var(--muted2);
  font-size: 13px;
}
</style>