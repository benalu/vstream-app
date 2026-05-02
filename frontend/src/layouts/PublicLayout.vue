<script setup>
import TheNavbar from '@/components/public/TheNavbar.vue'
import TheFooter from '@/components/public/TheFooter.vue'
</script>

<template>
  <div class="pub-layout">
    <TheNavbar />
    <router-view />
    <TheFooter />
  </div>
</template>

<style scoped>
.pub-layout {
  min-height: 100vh;
  background: #07070e;
  color: #e2e8f0;
  font-family: 'DM Sans', system-ui, sans-serif;
}
</style>