<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { useStore } from "~/store";

const store = useStore();

const handleFile = async (file: File) => {
  if (file.type !== "video/mp4" && !file.name.endsWith(".mp4")) {
    console.error("Selected file is not an MP4");
    return;
  }

  if (window.exportGPMF) {
    try {
      const gpmfData = await window.exportGPMF(file) as Uint8Array;
      
      // Create a blob from the Uint8Array
      const blob = new Blob([gpmfData], { type: 'application/octet-stream' });
      
      // Create a download link
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `${file.name.replace('.mp4', '')}_gpmf.bin`;
      document.body.appendChild(a);
      a.click();
      
      // Cleanup
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    } catch (err) {
      console.error("Error exporting GPMF:", err);
    }
  } else {
    console.error("WASM exportGPMF function not available");
  }
};


</script>

<template>
  <section class="p-4 flex flex-col items-stretch h-full gap-x-8">
    <div class="flex-1 flex flex-col mx-auto">
      <h2 class="max-w-2xl mb-4 text-4xl font-roboto-title font-bold leading-none tracking-tight md:text-5xl">GPMF Export</h2>
      <p class="max-w-2xl mb-6 font-light text-gray-500 lg:mb-8 md:text-lg lg:text-xl dark:text-gray-400">
        This tool allows you to export GPMF data from a GoPro video.
      </p>
      <p class="max-w-2xl mb-6 font-light text-gray-500 lg:mb-8 md:text-lg lg:text-xl dark:text-gray-400">
        And everything happens in your browser. Your data stays local.
      </p>
      <FileDrop class="h-full" @file-selected="handleFile" />
    </div>
  </section>
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

