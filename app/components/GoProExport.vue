<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { useStore } from "~/stores";

const store = useStore();

const handleFile = async (event: Event) => {
  const input = event.target as HTMLInputElement;
  if (!input.files || input.files.length === 0) {
    console.error("No file selected");
    return;
  }

  const file = input.files[0];

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
  <div>
    <div class="flex flex-row justify-between items-center gap-4">
      <Label for="gpmffile">Export GPMF Data</Label>
      <div class="flex gap-2">
        <Input id="gpmffile" type="file" @change="handleFile" accept="video/mp4" />
      </div>
    </div>
    <div class="mt-4">
      <Button @click="store.exportGpx" :disabled="!store.gpsData.length">
        Export GPX Track
      </Button>
    </div>
  </div>
</template> 