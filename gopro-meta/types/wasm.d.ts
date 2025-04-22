interface Window {
  processFile: (file: File) => Promise<{
    gpsData: any[];
    gyroData: any[];
    faceData: any[];
    lumaData: any[];
    hueData: any[];
    sceneData: any[];
  }>;
  exportGPMF: (file: File) => Promise<Uint8Array>;
} 