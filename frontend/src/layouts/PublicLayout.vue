<template>
  <div class="min-h-screen flex flex-col bg-slate-950">
    <nav class="bg-black/50 backdrop-blur-md border-b border-slate-800 p-4 sticky top-0 z-50">
      <div class="container mx-auto flex justify-between items-center">
        <h1 class="text-2xl font-bold text-white"><span class="text-indigo-500">V-</span>STREAM</h1>
        <div>
          <router-link to="/" class="mr-4 hover:text-indigo-400">Home</router-link>
          <router-link to="/admin" class="text-sm bg-indigo-600 px-4 py-2 rounded-lg text-white">Admin</router-link>
        </div>
      </div>
    </nav>

    <main class="flex-1 container mx-auto p-4">
      <router-view />
    </main>
  </div>
</template>