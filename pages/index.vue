<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import { useStore } from "~/store";
import VideoControls from '../components/VideoControls.vue'
import DebugData from '../components/DebugData.vue'
import DemoFiles from '../components/DemoFiles.vue'
import { useProcessGoproFile } from '../composables'

const store = useStore();
const videoElement = ref<HTMLVideoElement | null>(null);
const videoControls = ref<InstanceType<typeof VideoControls> | null>(null);
const currentTime = ref(0);

const { processFile } = useProcessGoproFile()

// Demo files configuration - centralized
const demoFiles = [
  {
    name: "Evenes Approach Time-lapse",
    url: "https://gopro-meta.s3.eu-west-1.amazonaws.com/evenes_approach_4k.MP4",
    description: "Evenes Approach Time-lapse",
    thumbnail: "https://gopro-meta.s3.eu-west-1.amazonaws.com/evenes_thumb.jpg"
  },
  // {
  //   name: "Nauders Bike", 
  //   url: "https://gopro-meta.s3.eu-west-1.amazonaws.com/nauders_bike.mp4",
  //   description: "Nauders Bike",
  //   thumbnail: "https://gopro-meta.s3.eu-west-1.amazonaws.com/nauders_bike.mp4"
  // },
  // {
  //   name: "Skiing Mürren",
  //   url: "https://gopro-meta.s3.eu-west-1.amazonaws.com/ski_murren.MP4", 
  //   description: "Skiing Mürren",
  //   thumbnail: "https://gopro-meta.s3.eu-west-1.amazonaws.com/ski_murren.MP4"
  // }
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
      <h2 class="max-w-2xl mb-4 text-4xl font-roboto-title font-bold leading-none tracking-tight md:text-5xl">Visualizer - What is this?</h2>
      <p class="max-w-2xl mb-6 font-light text-gray-500 lg:mb-8 md:text-lg lg:text-xl dark:text-gray-400">
        GoPro Metadata is a tool that allows you to extract and display metadata from GoPro videos.
        It can extract GPMF data in a binary format and export GPS data as GPX.
      </p>
      <p class="max-w-2xl mb-6 font-light text-gray-500 lg:mb-8 md:text-lg lg:text-xl dark:text-gray-400">
        And everything happens in your browser. Your data stays local.
      </p>
      
      <div class="space-y-8">
        <!-- File Drop -->
        <FileDrop @file-selected="processFile" />
        
        <!-- Demo Files -->
        <DemoFiles :demo-files="demoFiles" @file-selected="processFile" />
      </div>
    </div>
  </section>
  
  <!-- Video loaded state -->
  <Transition name="fade">
    <section class="mx-4 flex flex-col lg:flex-row gap-x-4 h-[calc(100vh-200px)]" v-if="store.videoUrl">
      <!-- Main content area -->
      <div class="flex-1 flex items-center justify-center">
        <div>
          <FaceBox class="mb-4">
            <video ref="videoElement" :src="store.videoUrl" controls @timeupdate="updateCurrentTime"></video>
          </FaceBox>
          <UploadButton />
          <DemoFiles :demo-files="demoFiles" @file-selected="processFile" compact />
        </div>
      </div>
      
      <!-- Sidebar with metadata and demo files -->
      <div class="flex-1 h-full">
        <div class="h-full overflow-y-auto pr-4 space-y-6">
          
          <!-- Metadata sections -->
          <section v-if="(store as any).showMap" class="mb-8">
            <h2 class="font-roboto-title font-bold">GPS</h2>
            <Map />
          </section>
          <section class="mb-8">
            <h2 class="font-roboto-title font-bold">Scene</h2>
            <SceneDisplay />
          </section>
          <section class="mb-8">
            <h2 class="font-roboto-title font-bold">Acceleration</h2>
            <AccelerationVisualizer />
          </section>
          <section class="mb-8">
            <HueDisplay />
          </section>
          <section class="mb-8">
            <Luminance class="mb-4"/>
          </section>
          <section class="mb-4">
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

