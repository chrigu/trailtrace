<script setup lang="ts">
import { computed } from 'vue';
import { useStore } from '~/store';
import type { FaceData } from '~/store';

const store = useStore();

const currentFaceData = computed(() => store.currentFaceData as FaceData | null);

// Calculate the position and size of the face box relative to the video dimensions
const faceBoxStyle = computed(() => {
  if (!currentFaceData.value) return {};
  
  return {
    position: 'absolute' as const,
    left: `${100 - currentFaceData.value.x * 100}%`,
    top: `${100 - currentFaceData.value.y * 100}%`,
    width: `${currentFaceData.value.w * 100}%`,
    height: `${currentFaceData.value.h * 100}%`,
    border: '2px solid red',
    boxSizing: 'border-box' as const,
  };
});

// Display face metrics
const faceMetrics = computed(() => {
  if (!currentFaceData.value) return null;
  
  return {
    confidence: `${(currentFaceData.value.confidence * 100).toFixed(1)}%`,
    smile: `${(currentFaceData.value.smile * 100).toFixed(1)}%`,
    blink: `${(currentFaceData.value.blink * 100).toFixed(1)}%`,
  };
});
</script>

<template>
  <div class="relative">
    <slot></slot>
    <div v-if="currentFaceData" :style="faceBoxStyle">
      <!-- <div class="absolute -top-6 left-0 bg-black/75 text-white text-xs p-1 rounded">
        <div>ID: {{ currentFaceData.id }}</div>
        <div>Confidence: {{ faceMetrics?.confidence }}</div>
        <div>Smile: {{ faceMetrics?.smile }}</div>
        <div>Blink: {{ faceMetrics?.blink }}</div>
      </div> -->
    </div>
  </div>
</template> 