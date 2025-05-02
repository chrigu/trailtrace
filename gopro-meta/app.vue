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

useHead({
  bodyAttrs: {
    class: "bg-stone-100"
  }
})

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
  <section class="mb-8 p-4 bg-gray-900 flex flex-row justify-between">
    <h1 class="text-pink-500 font-display text-4xl">GoGoPro</h1>
    <GoProUpload class="mb-4" />
  </section>
  <section class="mx-4 flex flex-row gap-x-4">
    <div class="flex-1">
      <FaceBox v-if="store.videoUrl">
        <video ref="videoElement" :src="store.videoUrl" controls @timeupdate="updateCurrentTime"></video>
      </FaceBox>
    </div>
    <div class="flex-1">
      <div v-if="store.videoUrl">
        <section v-if="store.showMap">
          <h2>GPS</h2>
          <Map />
          <GoProExport />
        </section>
        <section>
          <h2>Acceleration</h2>
          <AccelerationVisualizer />
        </section>
        <section>
          <h2>Luminance</h2>
          <Luminance />
        </section>
        <section>
          <h2>Hue</h2>
          <HueDisplay />
        </section>
        <section>
          <h2>Scene</h2>
          <SceneDisplay />
        </section>
      </div>
    </div>
  </section>
  <section class="mx-4" v-if="store.videoUrl">
    <p>Current Playback Time: {{ formattedTime }}</p>
    <p>Current GPS Data:<pre>{{ store.currentGpsData }}</pre></p>
    <p>Current Acceleration Data:<pre>{{ store.currentAccelerationData }}</pre></p>
    <p>Current Face Data:<pre>{{ store.currentFaceData }}</pre></p>
    <p>Current Luminance Data:<pre>{{ store.currentLuminanceData }}</pre></p>
    <p>Current Hue Data:<pre>{{ store.currentHueData }}</pre></p>
    <p>Current Scene Data:<pre>{{ store.currentSceneData }}</pre></p>
    <p>Current Playback Time: {{ formattedTime }}</p>
    <p>Current GPS Data:<pre>{{ store.currentGpsData }}</pre></p>
    <p>Current Acceleration Data:<pre>{{ store.currentAccelerationData }}</pre></p>
  </section>
</template>
