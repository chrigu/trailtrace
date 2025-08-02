<script setup lang="ts">
import { computed } from 'vue';
import { useStore } from '~/stores';
import type { FaceData } from '~/stores';

const store = useStore();

const currentFaceData = computed(() => store.currentFaceData as FaceData | null);

// Calculate the position and size of the face box relative to the video dimensions
const faceBoxStyle = computed(() => {
  if (!currentFaceData.value) return {};
  
  return {
    position: 'absolute' as const,
    left: `${currentFaceData.value.x * 100}%`,
    top: `${currentFaceData.value.y * 100}%`,
    width: `${currentFaceData.value.w * 100}%`,
    height: `${currentFaceData.value.h * 100}%`,
    border: '2px solid red',
    boxSizing: 'border-box' as const,
  };
});

// Display face metrics
const faceMetrics = computed(() => {
  if (!currentFaceData.value) return null;
  let smile = "😶"
  let blink = "😳"

  if (currentFaceData.value.smile > 0.65) {
    smile = "😊"
  } else if (currentFaceData.value.smile > 0.35) {
    smile = "😐"
  }

  if (currentFaceData.value.blink > 0.5) {
    blink = "😵"
  }

  
  return {
    confidence: currentFaceData.value.confidence,
    smile: `${smile} ${currentFaceData.value.smile}`,
    blink: `${blink} ${currentFaceData.value.blink}`,
  };
});
</script>

<template>
  <div class="relative">
    <slot></slot>
    <div v-if="currentFaceData" :style="faceBoxStyle">
      <div class="absolute left-full top-0 ml-2 bg-black/75 text-white p-1 rounded">
        <div>ID: {{ currentFaceData.id }}</div>
        <div>Confidence: {{ faceMetrics?.confidence }}</div>
        <div>Smile: {{ faceMetrics?.smile }}</div>
        <div>Blink: {{ faceMetrics?.blink }}</div>
      </div>
    </div>
  </div>
</template> 