<script setup lang="ts">

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { ref } from 'vue'

const emit = defineEmits(['file-selected'])

const isDragging = ref(false)
const isFileTooLarge = ref(false)

function checkFileSize(file: File) {
  return file.size > 2 * 1024 * 1024 * 1024;
}

const handleFileInput = (event: Event) => {
  const input = event.target as HTMLInputElement;
  if (input.files && input.files.length > 0) {
    const file = input.files[0];
    if (checkFileSize(file)) {
      isFileTooLarge.value = true;
    } else {
      isFileTooLarge.value = false;
      emit('file-selected', file);
    }
  }
};

const handleDragOver = (event: DragEvent) => {
  event.preventDefault();
  isDragging.value = true;
};

const handleDragLeave = (event: DragEvent) => {
  event.preventDefault();
  isDragging.value = false;
};

const handleDrop = (event: DragEvent) => {
  event.preventDefault();
  isDragging.value = false;
  
  if (event.dataTransfer?.files && event.dataTransfer.files.length > 0) {
    const file = event.dataTransfer.files[0];
    if (checkFileSize(file)) {
      isFileTooLarge.value = true;
    } else {
      isFileTooLarge.value = false;
      emit('file-selected', file);
    }
  }
};
</script>

<template>
  <div>
    <div 
      class="border-2 border-dashed rounded-xl p-6 text-center transition-colors bg-white"
      :class="[
        isDragging ? 'border-sky-700 bg-sky-50' : 'border-gray-300 hover:border-sky-400',
        'cursor-pointer',
        isFileTooLarge ? 'border-red-500 bg-red-50' : ''
      ]"
      @dragover="handleDragOver"
      @dragleave="handleDragLeave"
      @drop="handleDrop"
    >
      <div class="space-y-4">
        <div v-if="isFileTooLarge">
          <p class="text-red-500">File is too large</p>
        </div>
        <div class="text-gray-900">
          <p class="md:text-lg lg:text-xl font-medium">Drag and drop your GoPro file here</p>
        </div>
        <div class="text-gray-900">
          <p class="text-sm">⚠️ Unfortunately, most browsers only support files up to 2GB ⚠️</p>
        </div>
        <p class="space-y-4 text-md">or</p>
        <div class="flex flex-row justify-center items-center gap-4">
          <Label for="videofile" class="cursor-pointer bg-sky-700 hover:bg-sky-800 text-white rounded px-4 py-2">Select a file</Label>
          <Input 
            id="videofile" 
            type="file" 
            @change="handleFileInput" 
            accept="video/mp4" 
            class="hidden" 
          />
        </div>
      </div>
    </div>
  </div>
</template>

