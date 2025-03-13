<template>
  <div style="height: 100vh; width: 100vw">
    <LMap
      ref="map"
      :zoom="zoom"
      :center="computedCenter"
      :use-global-leaflet="false"
    >
      <LTileLayer
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        attribution="&amp;copy; <a href=&quot;https://www.openstreetmap.org/&quot;>OpenStreetMap</a> contributors"
        layer-type="base"
        name="OpenStreetMap"
      />
      <!-- Draw polyline if there are at least two points -->
      <LPolyline
        v-if="store.gpsData.length > 1"
        :lat-lngs="polylinePoints"
        color="blue"
        :weight="4"
      />
    </LMap>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import { useStore } from "@/store"; // Adjust the import path based on your project structure

const zoom = ref(6);
const store = useStore();

// Compute center dynamically
const computedCenter = computed(() => store.center);

// Convert gpsData to a format suitable for LPolyline
const polylinePoints = computed(() =>
  store.gpsData.map((point) => [point.latitude, point.longitude])
);
</script>
