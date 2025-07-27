<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import { useStore } from "~/store";
import VideoControls from '../components/VideoControls.vue'
import DebugData from '../components/DebugData.vue'
import DemoFiles from '../components/DemoFiles.vue'
import { Button } from '~/components/ui/button'
import { useProcessGoproFile } from '../composables'

const store = useStore();
const videoElement = ref<HTMLVideoElement | null>(null);
const videoControls = ref<InstanceType<typeof VideoControls> | null>(null);
const currentTime = ref(0);

const { processFile } = useProcessGoproFile()

const exportGpx = () => {
  store.exportGpx();
}

// Demo files configuration - centralized
const demoFiles = [
  {
    title: "Evenes Approach Time-lapse",
    name: "evenes_approach_4k.MP4",
    url: "/videos/evenes_approach_4k.MP4",
    description: "Evenes Approach Time-lapse",
    thumbnail: "/thumbnails/evenes_thumb.jpg"
  },
  {
    title: "Nauders Bike", 
    name: "nauders_bike.MP4",
    url: "/videos/nauders_bike.MP4",
    description: "Nauders Bike",
    thumbnail: "/thumbnails/nauders_thumb.jpg"
  }
]

const updateCurrentTime = () => {
  if (videoElement.value) {
    store.setVideoCurrentTime(videoElement.value.currentTime);
    currentTime.value = videoElement.value.currentTime;
  }
};

// Load wasm_exec.js and WebAssembly
onMounted(async () => {
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
      (controls as any).videoElement = videoEl;
    }
  },
  { immediate: true }
);

</script>

<template>
  <!-- <section class="mb-8">
    <Transition name="fade">
      <VideoControls ref="videoControls" v-if="store.videoUrl" />
    </Transition>
  </section> -->
  
  <!-- Initial state: no video loaded -->
  <section v-if="!store.videoUrl" class="p-4 flex flex-col items-stretch h-full gap-x-8">
    <div class="flex-1 flex flex-col mx-auto">
      <h2 class="max-w-2xl mb-4 text-4xl font-roboto-title font-bold leading-none tracking-tight md:text-5xl">🧭 Explorer</h2>
      <p class="max-w-2xl mb-6 font-light text-gray-500 lg:mb-8 md:text-lg lg:text-xl dark:text-gray-400">
        The Explorer is a tool that lets you visualize metadata from GoPro videos — including GPS paths, sensor data, and more.
      </p>
      <p class="max-w-2xl mb-6 font-light text-gray-500 lg:mb-8 md:text-lg lg:text-xl dark:text-gray-400">
        Everything happens in your browser. Your data stays local.
      </p>
      <div class="space-y-8">
        <FileDrop @file-selected="processFile" />
        <DemoFiles :demo-files="demoFiles" @file-selected="processFile" />
      </div>
    </div>
  </section>
  
  <!-- Video loaded state -->
  <Transition name="fade">
    <section class="flex flex-col lg:flex-row gap-x-4 h-[calc(100vh-50px)]" v-if="store.videoUrl">
      <!-- Main content area -->
      <div class="flex-1 flex justify-center mt-6">
        <div>
          <FaceBox class="mb-4">
            <video ref="videoElement" :src="store.videoUrl" controls @timeupdate="updateCurrentTime"></video>
          </FaceBox>
          <UploadButton />
        </div>
      </div>
      
      <!-- Sidebar with metadata and demo files -->
      <div class="flex-1 h-full">
        <div class="h-full overflow-y-auto space-y-6">
          
          <!-- Metadata sections -->
          <section v-if="(store as any).showMap" class="mb-8 bg-white p-4">
            <div class="flex justify-between items-center mb-4">
              <h2 class="font-roboto-title text-lg text-gray-900">GPS</h2>
              <Button 
                variant="outline" 
                size="sm" 
                @click="exportGpx"
                :disabled="store.gpsData.length === 0"
                class="text-sm"
              >
                Export GPX
              </Button>
            </div>
            <Map />
          </section>
          <section class="mb-8 bg-white p-4">
            <h2 class="font-roboto-title text-lg text-gray-900">Scene</h2>
            <SceneDisplay />
          </section>
          <section class="mb-8 bg-white p-4">
            <h2 class="font-roboto-title text-lg text-gray-900">Acceleration</h2>
            <AccelerationVisualizer />
          </section>
          <section class="mb-8 bg-white p-4">
            <h2 class="font-roboto-title text-lg text-gray-900">Color and Luminance</h2>
            <div class="flex flex-row justify-between gap-4">
              <HueDisplay />
              <Luminance />
            </div>
          </section>
          <section class="mb-4 bg-white p-4">
            <DebugData />
          </section>
        </div>
      </div>
    </section>
  </Transition>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.4s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>

