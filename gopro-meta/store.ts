import { defineStore } from 'pinia'

export const useStore = defineStore('metaData', {
  state: () => ({
    gpsData: [] as GpsData[],
  }),
  
  getters: {
    center(state): [number, number] {
      if (state.gpsData.length === 0) {
        return [47.21322, -1.559482]; // Default center
      }
      // Compute center as the average of all points
      const avgLat =
        state.gpsData.reduce((sum, p) => sum + p.latitude, 0) / state.gpsData.length;
      const avgLng =
        state.gpsData.reduce((sum, p) => sum + p.longitude, 0) / state.gpsData.length;
      return [avgLat, avgLng];
    },
  },

  actions: {
    updateGpsData(data: GpsData[]) {
      this.gpsData = data;
    },
  },
});

export interface GpsData {
  latitude: number;
  longitude: number;
  altitude: number;
}
