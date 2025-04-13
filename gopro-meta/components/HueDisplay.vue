<template>
  <div v-if="store.currentHueData">
    <!-- Show color box -->
    <div
      :style="{
        backgroundColor: colorStyle,
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

// Compute the CSS color style from the RGB values
const colorStyle = computed(() => {
  if (!store.currentHueData) return 'transparent'
  
  const { red, green, blue } = store.currentHueData
  return `rgb(${red}, ${green}, ${blue})`
})
</script> 