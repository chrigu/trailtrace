import { useStore, type GpsData } from "~/store";


export const useProcessGoproFile = () => {

const store = useStore()

  const processFile = async (file: File) => {
    if (!file) {
      console.error("No file selected");
      return;
    }

    store.setVideoUrl(URL.createObjectURL(file));
    if (file.type !== "video/mp4" && !file.name.toLowerCase().endsWith(".mp4")) {
      console.error("Selected file is not an MP4");
      return;
    }

    if (window.processFile) {
      try {
        const metadata = await window.processFile(file) as {gpsData: GpsData[], gyroData: any[], faceData: any[], lumaData: any[], hueData: any[], sceneData: any[]};
        console.log("Metadata data:", metadata);
        store.setGpsData(metadata.gpsData);
        store.setGyroData(metadata.gyroData);
        store.setFaceData(metadata.faceData);
        store.setLuminanceData(metadata.lumaData);
        store.setHueData(metadata.hueData);
        store.setSceneData(metadata.sceneData);
      } catch (err) {
        console.error("Error processing file:", err);
      }
    } else {
      console.error("WASM processFile function not available");
    }
  };

return { processFile }
}