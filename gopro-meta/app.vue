<script setup lang="ts">
import { onMounted, ref, computed } from "vue";
import { useStore } from "~/store";

let go;
let wasmInstance;
const store = useStore();
const videoElement = ref<HTMLVideoElement | null>(null);
const currentTime = ref(0);

const formattedTime = computed(() => store.videoCurrentTime.toFixed(2));

const updateCurrentTime = () => {
  if (videoElement.value) {
    store.setVideoCurrentTime(videoElement.value.currentTime);
    currentTime.value = videoElement.value.currentTime;
  }
};

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

  if (videoElement.value) {
    videoElement.value.addEventListener("timeupdate", updateCurrentTime);
  }
});
</script>

<template>
  <section class="mb-8 p-4">
    <h1>TrackBack</h1>
  </section>
  <section class="mx-4 flex flex-row gap-x-4">
    <div>
      <GoProUpload class="mb-4" />
      <FaceBox v-if="store.videoUrl">
        <video ref="videoElement" :src="store.videoUrl" controls width="600" @timeupdate="updateCurrentTime"></video>
      </FaceBox>
      <p v-if="store.videoUrl">Current Playback Time: {{ formattedTime }}</p>
      <p v-if="store.videoUrl">Current GPS Data:<pre>{{ store.currentGpsData }}</pre></p>
      <p v-if="store.videoUrl">Current Acceleration Data:<pre>{{ store.currentAccelerationData }}</pre></p>
    </div>
    <div class="h-[600px] flex-1">
      <Map />
      <AccelerationVisualizer />
    </div>
  </section>
  <section class="mx-4 flex flex-row gap-x-4">
    <p v-if="store.videoUrl">Current Playback Time: {{ formattedTime }}</p>
    <p v-if="store.videoUrl">Current GPS Data:<pre>{{ store.currentGpsData }}</pre></p>
    <p v-if="store.videoUrl">Current Acceleration Data:<pre>{{ store.currentAccelerationData }}</pre></p>
    <p v-if="store.videoUrl">Current Face Data:<pre>{{ store.currentFaceData }}</pre></p>
  </section>
</template>
