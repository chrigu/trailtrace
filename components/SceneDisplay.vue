<template>
  <div class="scene-display">
    <div v-if="currentScene" class="scene-value text-lg">
      {{ sceneEmoji }} {{ currentScene }}
    </div>
    <div v-else class="no-scene">
      No scene data available
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useStore, type SceneData } from '../store';

const store = useStore();

const currentScene = computed(() => {
  const sceneData = store.currentSceneData as SceneData | null;
  if (!sceneData?.scenes?.length) return null;
  
  // Find the scene with highest probability
  const highestProbScene = sceneData.scenes.reduce((prev: { type: string; probability: number }, current: { type: string; probability: number }) => 
    (current.probability > prev.probability) ? current : prev
  );
  
  return highestProbScene.type;
});

const sceneEmoji = computed(() => {
  switch (currentScene.value) {
    case 'SNOW':
      return '❄️';
    case 'URBA':
      return '🏙️';
    case 'INDO':
      return '🏠';
    case 'WATR':
      return '🌊';
    case 'VEGE':
      return '🌿';
    case 'BEAC':
      return '🏖️';
    default:
      return '';
  }
});
</script>

<style scoped>
</style> 