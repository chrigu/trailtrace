<script setup lang="ts">

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { ref } from 'vue'
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
      <Label for="videofile" class="cursor-pointer text-sky-500 text-xl">Select a file</Label>
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

