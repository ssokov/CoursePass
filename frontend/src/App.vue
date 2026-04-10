<script setup lang="ts">
import { useRoute } from 'vue-router'
import AppHeader from './components/AppHeader.vue'
import AppToast from './components/AppToast.vue'
import { setUnauthorizedHandler } from './api/client'
import { useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

setUnauthorizedHandler(() => {
  router.push('/login')
})
</script>

<template>
  <div class="max-w-[1100px] mx-auto px-5 py-5 flex flex-col gap-3.5">
    <AppHeader v-if="route.path !== '/login'" />
    <router-view />
  </div>
  <AppToast />
</template>
