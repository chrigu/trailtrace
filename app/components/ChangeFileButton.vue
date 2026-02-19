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
      <Label for="videofile" class="inline-flex items-center justify-center w-full px-5 py-3 mb-2 mr-2 text-sm font-medium rounded border border-primary sm:w-auto focus:outline-none bg-primary hover:bg-primary/90 text-primary-foreground focus:z-10 focus:ring-4 focus:ring-primary/20 transition-all">Select a different file</Label>
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
