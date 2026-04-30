// src/lib/api.js
import axios from 'axios'
import { clearSessionCache } from '@/router'

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  withCredentials: true,
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Reset cache agar guard re-validasi di navigasi berikutnya
      clearSessionCache()

      // Jangan redirect kalau sudah di halaman login — mencegah infinite loop
      if (!window.location.pathname.includes('/admin/login')) {
        window.location.href = '/admin/login'
      }
    }
    return Promise.reject(error)
  }
)

export default api