<script setup lang="ts">
import { useGpmfExtractor } from '~/composables/useGpmfExtractor'

const { downloadGpmf } = useGpmfExtractor()

async function handleFile(file: File) {
  console.log('handleFile', file)
  try {
    await downloadGpmf(file)
  } catch (err) {
    console.error('Error processing file:', err)
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
