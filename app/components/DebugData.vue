<script setup lang="ts">
import { ref } from "vue";
import { useStore } from "~/stores";
import { Eye, EyeOff } from "lucide-vue-next";
import { Button } from '@/components/ui/button'

const store = useStore();

const showData = ref(false);
</script>

<template>
  <Button
    class="inline-flex items-center justify-center w-full px-5 py-3 mb-2 mr-2 text-sm font-medium text-gray-900 bg-white border border-gray-200 rounded-lg sm:w-auto focus:outline-none hover:bg-gray-100 hover:text-sky-700 focus:z-10 focus:ring-4 focus:ring-gray-200 dark:focus:ring-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:text-white dark:hover:bg-gray-700"
    @click="showData = !showData">
    <Eye v-if="!showData" />
    <EyeOff v-else />
    Show Data
  </Button>
  <div v-if="showData">
    <div class="flex flex-col gap-1 mb-4">
      <h4>Current Playback Time</h4>
      <p>{{ store.videoCurrentTime.toFixed(2) }}</p>
    </div>
    <div class="flex flex-col gap-1 mb-4">
      <h4>Video Duration</h4>
      <p>{{ store.videoDuration.toFixed(2) }} seconds</p>
    </div>
    <div class="flex flex-col gap-1 mb-4">
      <h4>Current GPS Data</h4>
      <pre>{{ store.currentGpsData }}</pre>
    </div>
    <div class="flex flex-col gap-1 mb-4">
      <h4>Current Acceleration Data</h4>
      <pre>{{ store.currentAccelerationData }}</pre>
    </div>
    <div class="flex flex-col gap-1 mb-4">
      <h4>Current Face Data</h4>
      <pre>{{ store.currentFaceData }}</pre>
    </div>
    <div class="flex flex-col gap-1 mb-4">
      <h4>Current Luminance Data</h4>
      <pre>{{ store.currentLuminanceData }}</pre>
    </div>
    <div class="flex flex-col gap-1 mb-4">
      <h4>Current Hue Data</h4>
      <pre>{{ store.currentHueData }}</pre>
    </div>
    <div class="flex flex-col gap-1 mb-4">
      <h4>Current Scene Data</h4>
      <pre>{{ store.currentSceneData }}</pre>
    </div>
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
