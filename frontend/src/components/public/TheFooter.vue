<template>
  <footer class="vs-footer">
    <div class="vs-footer__brand">
      <svg width="16" height="16" viewBox="0 0 20 20" fill="none">
        <path d="M4 4h5v5H4zM11 4h5v5h-5zM4 11h5v5H4zM11 11h5v5h-5z" fill="currentColor" opacity=".35"/>
        <path d="M7 2L2 7l5 5 5-5z" fill="currentColor"/>
      </svg>
      <span>V-STREAM</span>
    </div>
    <p class="vs-footer__copy">© {{ year }} V-STREAM. Platform streaming pribadi.</p>
  </footer>
</template>

<script setup>
const year = new Date().getFullYear()
</script>

<style scoped>
.vs-footer {
  border-top: 1px solid rgba(255,255,255,0.07);
  padding: 24px 48px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}
.vs-footer__brand {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: 'Bebas Neue', sans-serif;
  font-size: 16px;
  letter-spacing: 0.1em;
  color: #64748b;
}
.vs-footer__copy {
  font-size: 12px;
  color: #64748b;
  opacity: 0.6;
}

@media (max-width: 768px) {
  .vs-footer {
    padding: 20px;
    flex-direction: column;
    text-align: center;
  }
}
</style>