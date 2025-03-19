import { defineStore } from 'pinia'


export interface GpsData {
  latitude: number;
  longitude: number;
  altitude: number;
  timestamp: number;
  videoUrl: string;
}


export const useStore = defineStore('metaData', {
  state: () => ({
    gpsData: [] as GpsData[],
    videoCurrentTime: 0,
    videoUrl: '',
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
    currentGpsData(state) {
      const starTime = state.gpsData[0]?.timestamp;
      return findClosestObject(state.gpsData, state.videoCurrentTime, starTime);
    },
  },

  actions: {
    setGpsData(data: GpsData[]) {
      this.gpsData = data;
    },
    setVideoCurrentTime(time: number) {
      this.videoCurrentTime = time;
    },
    setVideoUrl(url: string) {
      this.videoUrl = url;
    },
  },
});

const findClosestObject = (arr: { timestamp: number }[], targetTimestamp: number, starTime: number) => {
  if (arr.length === 0) return null;

  let left = 0;
  let right = arr.length - 1;

  while (left < right) {
    const mid = Math.floor((left + right) / 2);

    if ((arr[mid].timestamp - starTime)/1000 === targetTimestamp) {
      return arr[mid]; // Exact match
    } else if ((arr[mid].timestamp - starTime)/1000  < targetTimestamp) {
      left = mid + 1;
    } else {
      right = mid;
    }
  }

  // After the loop, 'left' is the closest index or the one just after the target
  if (left === 0) return arr[0]; // Target is before the first element
  if (left >= arr.length) return arr[arr.length - 1]; // Target is after the last element

  // Compare the two closest candidates
  const prev = arr[left - 1];
  const next = arr[left];

  return Math.abs((prev.timestamp - starTime)/1000 - targetTimestamp) <= Math.abs((next.timestamp - starTime)/1000 - targetTimestamp)
    ? prev
    : next;
};
