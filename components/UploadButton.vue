<script setup lang="ts">

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useProcessGoproFile } from '../composables'

const { processFile } = useProcessGoproFile()

const handleFileInput = (event: Event) => {
  const input = event.target as HTMLInputElement;
  if (input.files && input.files.length > 0) {
    processFile(input.files[0]);
  }
};

</script>

<template>
  <div>
    <div class="flex flex-row justify-center items-center gap-4">
      <Label for="videofile" class="inline-flex items-center justify-center w-full px-5 py-3 mb-2 mr-2 text-sm font-medium text-gray-900 bg-white border border-gray-200 rounded-lg sm:w-auto focus:outline-none hover:bg-gray-100 hover:text-sky-700 focus:z-10 focus:ring-4 focus:ring-gray-200 dark:focus:ring-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:text-white dark:hover:bg-gray-700">Select a different file</Label>
      <Input 
        id="videofile" 
        type="file" 
        @change="handleFileInput" 
        accept="video/mp4" 
        class="hidden" 
      />
    </div>
  </div>
</template>

