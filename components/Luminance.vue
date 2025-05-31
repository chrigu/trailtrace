<template>
  <div v-if="store.currentLuminanceData">
    <!-- Show numeric value -->
    <p>Luminance: {{ store.currentLuminanceData.luminance }}</p>

    <!-- Show gray box representing luminance -->
    <div
      :style="{
        backgroundColor: luminanceToGray(store.currentLuminanceData.luma),
        width: '100px',
        height: '100px',
        border: '1px solid #000',
      }"
    ></div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useStore } from "~/store";

const store = useStore()

// Converts a luminance value (0 to 1) to a gray CSS color string
const luminanceToGray = (luminance) => {
  return `rgb(${luminance}, ${luminance}, ${luminance})`
}
</script>