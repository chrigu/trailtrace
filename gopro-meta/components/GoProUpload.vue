<script setup lang="ts">
import { ref } from "vue";
import { useStore, type GpsData } from "~/store";
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const store = useStore()


const handleFile = async (event: Event) => {
  const input = event.target as HTMLInputElement;
  if (!input.files || input.files.length === 0) {
    console.error("No file selected");
    return;
  }

  const file = input.files[0];

  store.setVideoUrl(URL.createObjectURL(file));

  if (file.type !== "video/mp4" && !file.name.endsWith(".mp4")) {
    console.error("Selected file is not an MP4");
    return;
  }

  if (window.processFile) {
    try {
      const metadata = await window.processFile(file) as {gpsData: GpsData[], gyroData: any[], faceData: any[]};
      console.log("Metadata data:", metadata);
      store.setGpsData(metadata.gpsData);
      store.setGyroData(metadata.gyroData);
      store.setFaceData(metadata.faceData);
      store.setLuminanceData(metadata.lumaData);
    } catch (err) {
      console.error("Error processing file:", err);
    }
  } else {
    console.error("WASM processFile function not available");
  }
};


</script>

<template>
  <div>
    <div class="flex flex-row justify-between">
      <Label for="videofile">GoPro file</Label>
      <Input id="videofile" type="file" @change="handleFile" accept="video/mp4" />
    </div>
  </div>
</template>

