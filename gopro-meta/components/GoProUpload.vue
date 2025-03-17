<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useStore, type GpsData } from "~/store";
import { Input } from '@/components/ui/input'

const store = useStore()

const fileInput = ref(null);
const videoElement = ref(null);
const videoUrl = ref("");
const currentTime = ref(0);

const formattedTime = computed(() => store.videoCurrentTime.toFixed(2));

const handleFile = async (event: Event) => {
  const input = event.target as HTMLInputElement;
  if (!input.files || input.files.length === 0) {
    console.error("No file selected");
    return;
  }

  const file = input.files[0];

  videoUrl.value = URL.createObjectURL(file);

  if (file.type !== "video/mp4" && !file.name.endsWith(".mp4")) {
    console.error("Selected file is not an MP4");
    return;
  }

  if (window.processFile) {
    try {
      const gpsData = await window.processFile(file) as GpsData[];
      store.updateGpsData(gpsData);
    } catch (err) {
      console.error("Error processing file:", err);
    }
  } else {
    console.error("WASM processFile function not available");
  }
};

const updateCurrentTime = () => {
  if (videoElement.value) {
    store.updateVideoCurrentTime(videoElement.value.currentTime);
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
  <div>
    <h1>GoPro File Upload</h1>
    <Input type="file" @change="handleFile" accept="video/mp4" />
  </div>
  <video v-if="videoUrl" ref="videoElement" :src="videoUrl" controls width="600" @timeupdate="updateCurrentTime"></video>
  <p v-if="videoUrl">Current Playback Time: {{ formattedTime }}</p>
  <p v-if="videoUrl">Current GPS Data:<pre>{{ store.currentGpsData }}</pre></p>
</template>

