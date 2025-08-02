import { useWorker } from './useWorker'
import { useFileValidation } from './useFileValidation'
import { useStore, type GpsData } from '~/stores'

export interface TelemetryData {
  gpsData: GpsData[]
  gyroData: any[]
  faceData: any[]
  lumaData: any[]
  hueData: any[]
  sceneData: any[]
}

export const useTelemetryProcessor = () => {
  const { initializeWorker, executeWorkerMethod } = useWorker()
  const { validateMp4File } = useFileValidation()
  const store = useStore()

  const extractTelemetry = async (file: File): Promise<TelemetryData> => {
    if (!validateMp4File(file)) {
      throw new Error('Invalid MP4 file')
    }

    initializeWorker()
    console.log('Calling exportTelemetryInWorker...')
    
    const telemetry = await executeWorkerMethod<TelemetryData>(file, 'processFile')
    console.log('Telemetry data received:', telemetry)
    
    return telemetry
  }

  const processFileWithTelemetry = async (file: File) => {
    try {
      // Set video URL first
      store.setVideoUrl(URL.createObjectURL(file))
      
      // Extract telemetry data
      const telemetry = await extractTelemetry(file)
      
      if (telemetry) {
        // Update store with all telemetry data
        store.setGpsData(telemetry.gpsData)
        store.setGyroData(telemetry.gyroData)
        store.setFaceData(telemetry.faceData)
        store.setLuminanceData(telemetry.lumaData)
        store.setHueData(telemetry.hueData)
        store.setSceneData(telemetry.sceneData)
      } else {
        console.error('No telemetry data received')
        throw new Error('No telemetry data received')
      }
    } catch (err) {
      console.error('Error processing file:', err)
      throw err
    }
  }

  return {
    extractTelemetry,
    processFileWithTelemetry
  }
}