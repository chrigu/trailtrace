<script setup lang="ts">

let worker: Worker;
let req = 0;                      // simple request-id counter

onMounted(() => {
  // Spawn the module worker (Vite 5+, Nuxt 3, etc. can URL-import; otherwise use raw path)
  worker = new Worker('/worker_wasm.js');
  console.log('Worker spawned', worker);

});

function exportGPMFInWorker(file: File): Promise<Uint8Array> {
  return new Promise((resolve, reject) => {
    const id = ++req;
    const listener = ({ data }: MessageEvent) => {
      if (data.id !== id) return;
      worker.removeEventListener('message', listener);
      data.ok ? resolve(new Uint8Array(data.payload)) : reject(data.error);
    };
    worker.addEventListener('message', listener);
    console.log('posting message', id, file.size);
    worker.postMessage({ id, file, method: 'exportGPMF' });
  });
}

async function handleFile(file: File) {
  console.log('handleFile', file);
  if (file.type !== "video/mp4" && !file.name.toLowerCase().endsWith(".mp4")) {
    console.error('Selected file is not an MP4');
    return;
  }
  try {
    // const gpmf = await exportGPMFInWorker(file);     // **no main-thread Wasm call**
    console.log('Calling exportGPMFInWorker...');
    const gpmf = await exportGPMFInWorker(file);
    console.log('GPMF data received:', gpmf);
    const blob = new Blob([gpmf], { type: 'application/octet-stream' });
    const url = URL.createObjectURL(blob);
    const a = Object.assign(document.createElement('a'), {
      href: url,
      download: `${file.name.replace(/\.mp4$/i, '')}_gpmf.bin`
    });
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  } catch (err) {
    console.error('Error exporting GPMF:', err);
  }
}

</script>

<template>
  <section class="p-4 flex flex-col items-stretch h-full gap-x-8">
    <div class="flex-1 flex flex-col mx-auto">
      <h2 class="max-w-2xl mb-4 text-4xl font-roboto-title font-bold leading-none tracking-tight md:text-5xl">🛠️
        Extractor</h2>
      <p class="max-w-2xl mb-6 font-light text-gray-500 lg:mb-8 md:text-lg lg:text-xl dark:text-gray-400">
        This tool allows you to export GPMF data from a GoPro video.
      </p>
      <p class="max-w-2xl mb-6 font-light text-gray-500 lg:mb-8 md:text-lg lg:text-xl dark:text-gray-400">
        All processing is done locally in your browser — no files are uploaded.
      </p>

      <div class="space-y-8">
        <!-- File Drop -->
        <FileDrop @file-selected="handleFile" />
      </div>
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
