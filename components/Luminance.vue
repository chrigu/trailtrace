<template>
  <div v-if="store.currentLuminanceData">
    <h3>Luminance</h3>
    <div class="luminance-container">
      <div class="luminance-item">
        <div
          :style="{
            backgroundColor: luminanceToGray(store.currentLuminanceData.luma),
            width: '50px',
            height: '50px',
            border: '1px solid #000',
            marginBottom: '5px'
          }"
        ></div>
        <div>Value: {{ store.currentLuminanceData.luma.toFixed(1) }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useStore } from "~/store";

const store = useStore()

// Converts a luminance value (0 to 1) to a gray CSS color string
const luminanceToGray = (luminance) => {
  return `rgb(${luminance}, ${luminance}, ${luminance})`
}
</script>

<style scoped>
.luminance-container {
  display: flex;
  gap: 20px;
  margin-top: 10px;
}

.luminance-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}
</style>