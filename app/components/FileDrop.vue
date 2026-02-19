<script setup lang="ts">

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { ref } from 'vue'

const emit = defineEmits(['file-selected'])

const isDragging = ref(false)


const handleFileInput = (event: Event) => {
  const input = event.target as HTMLInputElement;
  if (input.files && input.files.length > 0) {
    const file = input.files[0];
    emit('file-selected', file);
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
    emit('file-selected', file);
  }
};
</script>

<template>
  <div>
    <div
      class="border-2 border-dashed rounded-xl p-6 text-center transition-all bg-white/60 backdrop-blur-sm"
      :class="[
        isDragging ? 'border-primary bg-blue-50/80 shadow-lg scale-105' : 'border-border hover:border-primary/50 hover:shadow-md',
        'cursor-pointer',
      ]"
      @dragover="handleDragOver"
      @dragleave="handleDragLeave"
      @drop="handleDrop"
    >
      <div class="space-y-4">
        <div class="text-foreground">
          <p class="md:text-lg lg:text-xl font-medium">Drag and drop your GoPro file here</p>
        </div>
        <p class="space-y-4 text-md text-muted-foreground">or</p>
        <div class="flex flex-row justify-center items-center gap-4">
          <Label for="videofile" class="cursor-pointer bg-primary hover:bg-primary/90 text-primary-foreground rounded px-6 py-2 font-medium transition-all">Select a file</Label>
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
