<template>
  <div v-if="store.currentHueData">
    <h3>Predominant Hues</h3>
    <div class="hues-container">
      <div v-for="(hue, index) in store.currentHueData.hues" :key="index" class="hue-item">
        <div 
          :style="{
            backgroundColor: `hsl(${hue.hue}, 100%, ${hue.weight}%)`,
            width: '50px',
            height: '50px',
            border: '1px solid #000',
            marginBottom: '5px'
          }"
        ></div>
        <div>Hue: {{ hue.hue.toFixed(1) }}°</div>
        <div>Weight: {{ hue.weight.toFixed(1) }}%</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useStore } from "~/stores";

const store = useStore()
</script>

<style scoped>
.hues-container {
  display: flex;
  gap: 20px;
  margin-top: 10px;
}

.hue-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}
</style> 