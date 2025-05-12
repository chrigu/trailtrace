<script setup lang="ts">
import { ref, computed } from 'vue'
import { useStore } from '~/store'
import { Button } from '@/components/ui/button'
import { Play, Pause } from 'lucide-vue-next'

const store = useStore()
const isPlaying = ref(false)
const videoElement = ref<HTMLVideoElement | null>(null)
const progressBar = ref<HTMLElement | null>(null)

const progressPercentage = computed(() => {
  if (store.videoDuration === 0) return 0
  return (store.videoCurrentTime / store.videoDuration) * 100
})

const playPause = () => {
  if (videoElement.value) {
    if (videoElement.value.paused) {
      videoElement.value.play()
      isPlaying.value = true
    } else {
      videoElement.value.pause()
      isPlaying.value = false
    }
  }
}

const handleProgressClick = (event: MouseEvent) => {
  console.log(videoElement.value)
  if (!videoElement.value) return
  
  const rect = progressBar.value?.getBoundingClientRect()
  if (!rect) return
  const x = event.clientX - rect.left
  const percentage = x / rect.width
  const newTime = percentage * store.videoDuration
  
  videoElement.value.currentTime = newTime
  store.setVideoCurrentTime(newTime)
}

const formatTime = (seconds: number) => {
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = Math.floor(seconds % 60)
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`
}

// Expose the video element ref to parent
defineExpose({
  videoElement
})
</script>

<template>
  <div>
    <div ref="progressBar" class="relative h-5 bg-gray-300 my-0 cursor-pointer" @click="handleProgressClick">
      <div 
        class="absolute h-full bg-sky-500 transition-[width] duration-100 ease-linear"
        :style="{ width: progressPercentage + '%' }"
      ></div>
      <div class="absolute right-4 top-1/2 -translate-y-1/2">
        {{ formatTime(store.videoCurrentTime) }} / {{ formatTime(store.videoDuration) }}
      </div>
    </div>
    <div class="mt-4 flex justify-center">
      <Button variant="outline" size="icon" @click="playPause">
        <Play v-if="!isPlaying" />
        <Pause v-else />
      </Button>
    </div>
  </div>
</template>
