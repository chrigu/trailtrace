<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useStore, type GpsData } from "~/store";

const store = useStore()

const fileInput = ref(null);
const videoElement = ref(null);
const videoUrl = ref("");
const currentTime = ref(0);
let go;
let wasmInstance;


const formattedTime = computed(() => store.videoCurrentTime.toFixed(2));

const loadWasmExec = async () => {
  return new Promise((resolve, reject) => {
    const script = document.createElement("script");
    script.src = "/wasm_exec.js"; // Load from public folder
    script.onload = resolve;
    script.onerror = reject;
    document.body.appendChild(script);
  });
};

// Load wasm_exec.js and WebAssembly
onMounted(async () => {
  try {
    await loadWasmExec();
    go = new Go();
    const wasmResponse = await fetch("/main.wasm");
    const wasmBytes = await wasmResponse.arrayBuffer();
    const { instance } = await WebAssembly.instantiate(wasmBytes, go.importObject);
    go.run(instance);
    wasmInstance = instance;
  } catch (error) {
    console.error("Error loading WebAssembly:", error);
  }
});

const handleFile = async () => {
  
  if (!fileInput.value || !fileInput.value.files.length) {
    console.error("No file selected");
    return;
  }

  const file = fileInput.value.files[0];

  videoUrl.value = URL.createObjectURL(file);

  if (file.type !== "video/mp4" && !file.name.endsWith(".mp4")) {
    console.error("Selected file is not an MP4");
    return;
  }

  if (window.processFile) {
    try {
      const gpsData = await window.processFile(file) as GpsData[];
      store.updateGpsData(gpsData);
      // gpsData.forEach((point, index) => {
      //   console.log(`Point ${index}: Lat ${point.latitude}, Lon ${point.longitude}, Alt ${point.altitude}`);
      // });
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
    <h1>WASM File Upload</h1>
    <input type="file" ref="fileInput" />
    <button @click="handleFile">Process File</button>
  </div>
  <video v-if="videoUrl" ref="videoElement" :src="videoUrl" controls width="600" @timeupdate="updateCurrentTime"></video>
  <p v-if="videoUrl">Current Playback Time: {{ formattedTime }}</p>
  <p v-if="videoUrl">Current GPS Data:<pre>{{ store.currentGpsData }}</pre></p>
</template>

