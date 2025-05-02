<template>
  <div class="h-[500px] w-full">
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

      <!-- Display a marker for the current GPS position -->
      <LMarker
        v-if="currentGpsData"
        :lat-lng="[currentGpsData.latitude, currentGpsData.longitude]"
      >
        <LPopup>
          <p>Current Position</p>
          <p>Lat: {{ currentGpsData.latitude.toFixed(5) }}</p>
          <p>Lon: {{ currentGpsData.longitude.toFixed(5) }}</p>
        </LPopup>
      </LMarker>
    </LMap>
  </div>
</template>


<script setup lang="ts">
import { computed, ref } from "vue";
import { useStore } from "@/store";

const zoom = ref(6);
const store = useStore();

// Compute center dynamically
const computedCenter = computed(() => store.center);

// Convert gpsData to a format suitable for LPolyline
const polylinePoints = computed(() =>
  store.gpsData.map((point) => [point.latitude, point.longitude])
);

// Track the current GPS data point
const currentGpsData = computed(() => store.currentGpsData);

</script>