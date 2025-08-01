import { useWorker } from './useWorker'
import { useFileValidation } from './useFileValidation'

export const useGpmfExtractor = () => {
  const { initializeWorker, executeWorkerMethod } = useWorker()
  const { validateMp4File } = useFileValidation()

  const extractGpmf = async (file: File): Promise<Uint8Array> => {
    if (!validateMp4File(file)) {
      throw new Error('Invalid MP4 file')
    }

    initializeWorker()
    console.log('Calling exportGPMFInWorker...')
    
    // measure time
    const startTime = performance.now()
    const gpmfData = await executeWorkerMethod<ArrayBuffer>(file, 'exportGPMF')
    console.log('GPMF data received:', gpmfData)

    const endTime = performance.now()
    const duration = endTime - startTime
    console.log("%c" + "GPMF (WASM) extraction took " + duration.toFixed(2) + "ms", "color:green;font-weight:bold;");
    
    return new Uint8Array(gpmfData)
  }

  const downloadGpmf = async (file: File) => {
    try {
      const gpmf = await extractGpmf(file)
      
      const blob = new Blob([gpmf], { type: 'application/octet-stream' })
      const url = URL.createObjectURL(blob)
      
      const a = Object.assign(document.createElement('a'), {
        href: url,
        download: `${file.name.replace(/\.mp4$/i, '')}_gpmf.bin`
      })
      
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)
    } catch (err) {
      console.error('Error exporting GPMF:', err)
      throw err
    }
  }

  return {
    extractGpmf,
    downloadGpmf
  }
}