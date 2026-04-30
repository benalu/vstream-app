// tailwind.config.js
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'vstream-dark': '#0f172a',
        'vstream-card': '#1e293b',
        'vstream-accent': '#6366f1',
      }
    },
  },
  plugins: [],
}