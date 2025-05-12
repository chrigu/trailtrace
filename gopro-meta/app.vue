<script setup lang="ts">
import { onMounted, ref, computed, watch } from "vue";
import { useStore } from "~/store";
import VideoControls from './components/VideoControls.vue'

let go;
let wasmInstance;
const store = useStore();
const videoElement = ref<HTMLVideoElement | null>(null);
const videoControls = ref<InstanceType<typeof VideoControls> | null>(null);
const currentTime = ref(0);

const formattedTime = computed(() => store.videoCurrentTime.toFixed(2));
const videoDuration = computed(() => store.videoDuration.toFixed(2));

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
    videoElement.value.addEventListener("loadedmetadata", () => {
      if (videoElement.value) {
        store.setVideoDuration(videoElement.value.duration);
      }
    });
  }
});

watch(
  () => [videoElement.value, videoControls.value],
  ([videoEl, controls]) => {
    if (videoEl && controls) {
      controls.videoElement = videoEl;
    }
  },
  { immediate: true }
);

</script>

<template>
  <section class="mb-0 p-4 bg-gray-700 flex flex-col md:flex-row justify-between">
    <h1 class="text-sky-500 font-sans font-bold text-4xl">GoGoPro</h1>
    <UploadButton v-if="store.videoUrl" />
  </section>
  <section class="mb-8">
    <VideoControls ref="videoControls" v-if="store.videoUrl" />
  </section>
  <section v-if="!store.videoUrl">
    <FileDrop class="mb-4" />
  </section>
  <section class="mx-4 flex flex-col lg:flex-row gap-x-4 h-[calc(100vh-200px)]">
    <div class="flex-1">
      <div v-show="store.videoUrl">
      <section>
        <SceneDisplay />
      </section>
      <FaceBox>
        <video ref="videoElement" :src="store.videoUrl" controls @timeupdate="updateCurrentTime"></video>
      </FaceBox>
        <section>
          <Luminance />
        </section>
        <section>
          <HueDisplay />
        </section>
      </div>
    </div>
    <div class="flex-1" v-if="store.videoUrl">
      <div>
        <section v-if="store.showMap">
          <h2>GPS</h2>
          <Map />
          <GoProExport />
        </section>
        <section>
          <h2>Acceleration</h2>
          <AccelerationVisualizer />
        </section>
      </div>
    </div>
  </section>
  <section class="mx-4" v-if="store.videoUrl">
    <p>Current Playback Time: {{ formattedTime }}</p>
    <p>Video Duration: {{ videoDuration }} seconds</p>
    <p>Current GPS Data:<pre>{{ store.currentGpsData }}</pre></p>
    <p>Current Acceleration Data:<pre>{{ store.currentAccelerationData }}</pre></p>
    <p>Current Face Data:<pre>{{ store.currentFaceData }}</pre></p>
    <p>Current Luminance Data:<pre>{{ store.currentLuminanceData }}</pre></p>
    <p>Current Hue Data:<pre>{{ store.currentHueData }}</pre></p>
    <p>Current Scene Data:<pre>{{ store.currentSceneData }}</pre></p>
  </section>
</template>

