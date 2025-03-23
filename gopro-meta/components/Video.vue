<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useStore } from "~/store";

const store = useStore()

const videoElement = ref(null);
const currentTime = ref(0);

const formattedTime = computed(() => store.videoCurrentTime.toFixed(2));

const updateCurrentTime = () => {
  if (videoElement.value) {
    store.setVideoCurrentTime(videoElement.value.currentTime);
    currentTime.value = videoElement.value.currentTime;
  }
};

onMounted(() => {
  if (videoElement.value) {
    videoElement.value.addEventListener("timeupdate", updateCurrentTime);
  }
});
</script>

<template>
  <video v-if="store.videoUrl" ref="videoElement" :src="store.videoUrl" controls width="600" @timeupdate="updateCurrentTime"></video>
  <p v-if="store.videoUrl">Current Playback Time: {{ formattedTime }}</p>
  <p v-if="store.videoUrl">Current GPS Data:<pre>{{ store.currentGpsData }}</pre></p>
  <p v-if="store.videoUrl">Current Gyro Data:<pre>{{ store.currentGyroData }}</pre></p>
</template>

