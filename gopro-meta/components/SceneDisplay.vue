<template>
  <div class="scene-display">
    <h3>Current Scene</h3>
    <div v-if="currentScene" class="scene-value">
      {{ currentScene }}
    </div>
    <div v-else class="no-scene">
      No scene data available
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useStore } from '../store';

const store = useStore();
const currentScene = computed(() => store.currentSceneData?.scene || null);
</script>

<style scoped>
</style> 