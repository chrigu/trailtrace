<script setup lang="ts">
import { ref } from "vue";
import { useStore } from "~/stores";
import { Eye, EyeOff } from "lucide-vue-next";
import { Button } from '@/components/ui/button'

const store = useStore();

const showData = ref(false);
</script>

<template>
  <Button class="inline-flex items-center justify-center w-full px-5 py-3 mb-2 mr-2 text-sm font-medium text-gray-900 bg-white border border-gray-200 rounded-lg sm:w-auto focus:outline-none hover:bg-gray-100 hover:text-sky-700 focus:z-10 focus:ring-4 focus:ring-gray-200 dark:focus:ring-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:text-white dark:hover:bg-gray-700"
  @click="showData = !showData">
    <Eye v-if="!showData" />
    <EyeOff v-else />
    Show Data
  </Button>
  <div v-if="showData">
    <p>Current Playback Time: {{ store.videoCurrentTime.toFixed(2) }}</p>
    <p>Video Duration: {{ store.videoDuration.toFixed(2) }} seconds</p>
    <p>Current GPS Data:<pre>{{ store.currentGpsData }}</pre></p>
    <p>Current Acceleration Data:<pre>{{ store.currentAccelerationData }}</pre></p>
    <p>Current Face Data:<pre>{{ store.currentFaceData }}</pre></p>
    <p>Current Luminance Data:<pre>{{ store.currentLuminanceData }}</pre></p>
    <p>Current Hue Data:<pre>{{ store.currentHueData }}</pre></p>
    <p>Current Scene Data:<pre>{{ store.currentSceneData }}</pre></p>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.4s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>

