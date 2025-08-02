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
import { computed, ref, watch, nextTick } from "vue";
import { useStore } from "~/stores";

const zoom = ref(6);
const store = useStore();
const map = ref();

// Compute center dynamically
const computedCenter = computed(() => store.center);

// Convert gpsData to a format suitable for LPolyline
const polylinePoints = computed(() =>
  store.gpsData.map((point) => [point.latitude, point.longitude])
);

// Track the current GPS data point
const currentGpsData = computed(() => store.currentGpsData);

// Calculate bounding box of the track
const trackBounds = computed(() => {
  if (store.gpsData.length === 0) {
    return null;
  }

  const lats = store.gpsData.map(point => point.latitude);
  const lngs = store.gpsData.map(point => point.longitude);

  return {
    minLat: Math.min(...lats),
    maxLat: Math.max(...lats),
    minLng: Math.min(...lngs),
    maxLng: Math.max(...lngs),
  };
});

// Calculate optimal zoom level based on track bounds
const calculateOptimalZoom = (bounds: any, mapWidth: number = 800, mapHeight: number = 500) => {
  if (!bounds) return 6;

  const latDiff = bounds.maxLat - bounds.minLat;
  const lngDiff = bounds.maxLng - bounds.minLng;

  // Calculate zoom level based on the larger dimension
  // These constants are empirically determined for good visibility
  const latZoom = Math.log2(180 / latDiff) - 1;
  const lngZoom = Math.log2(360 / lngDiff) - 1;

  // Use the smaller zoom (wider view) to ensure the entire track fits
  const calculatedZoom = Math.min(latZoom, lngZoom);

  // Clamp zoom between reasonable bounds (3-18)
  return Math.max(3, Math.min(18, Math.floor(calculatedZoom)));
};

// Auto-fit the map to show the entire track
const fitTrackToView = async () => {
  if (!trackBounds.value || !map.value) {
    return;
  }

  await nextTick();
  
  // Option 1: Use Leaflet's fitBounds (recommended)
  const leafletMap = map.value.leafletObject;
  if (leafletMap && leafletMap.fitBounds) {
    const bounds = [
      [trackBounds.value.minLat, trackBounds.value.minLng],
      [trackBounds.value.maxLat, trackBounds.value.maxLng]
    ];
    
    // Add padding around the track (10% on each side)
    leafletMap.fitBounds(bounds, { 
      padding: [20, 20],
      maxZoom: 16 // Prevent zooming in too much for small tracks
    });
  } else {
    // Option 2: Manual zoom calculation (fallback)
    const optimalZoom = calculateOptimalZoom(trackBounds.value);
    zoom.value = optimalZoom;
  }
};

// Watch for changes in GPS data and auto-fit the map
watch(
  () => store.gpsData.length,
  (newLength) => {
    if (newLength > 1) {
      // Use a small delay to ensure the map is fully rendered
      setTimeout(fitTrackToView, 100);
    }
  },
  { immediate: true }
);

// Also expose the method for manual triggering
defineExpose({
  fitTrackToView,
  trackBounds,
  calculateOptimalZoom
});
</script>