<script setup lang="ts">

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const emit = defineEmits(['file-selected']);

const handleFileInput = (event: Event) => {
  const input = event.target as HTMLInputElement;
  if (input.files && input.files.length > 0) {
    const file = input.files[0];
    emit('file-selected', file);
  }
};

</script>

<template>
  <div>
    <div class="flex flex-row justify-center items-center gap-4">
      <Label for="videofile" class="inline-flex items-center justify-center w-full px-5 py-3 mb-2 mr-2 text-sm font-medium rounded border border-gray-200 sm:w-auto focus:outline-none bg-sky-700 hover:bg-sky-800 text-white focus:z-10 focus:ring-4 focus:ring-gray-200 dark:focus:ring-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:text-white dark:hover:bg-gray-700">Select a different file</Label>
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

