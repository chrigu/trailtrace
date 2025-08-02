import { ref } from 'vue'

export interface WorkerRequest {
  id: number
  file: File
  method: string
}

export interface WorkerResponse {
  id: number
  ok: boolean
  payload?: any
  error?: string
}

export const useWorker = () => {
  const worker = ref<Worker | null>(null)
  const requestCounter = ref(0)

  const initializeWorker = () => {
    if (!worker.value) {
      worker.value = new Worker('/worker_wasm.js')
      console.log('Worker spawned', worker.value)
    }
    return worker.value
  }

  const executeWorkerMethod = <T = any>(file: File, method: string): Promise<T> => {
    return new Promise((resolve, reject) => {
      if (!worker.value) {
        reject(new Error('Worker not initialized'))
        return
      }

      const id = ++requestCounter.value
      
      const listener = ({ data }: MessageEvent<WorkerResponse>) => {
        if (data.id !== id) return
        
        worker.value?.removeEventListener('message', listener)
        
        if (data.ok) {
          resolve(data.payload)
        } else {
          reject(new Error(data.error || 'Worker execution failed'))
        }
      }

      worker.value.addEventListener('message', listener)
      console.log('📤 posting message', id, method, file.size)
      worker.value.postMessage({ id, file, method } as WorkerRequest)
    })
  }

  const destroyWorker = () => {
    if (worker.value) {
      worker.value.terminate()
      worker.value = null
      requestCounter.value = 0
    }
  }

  return {
    worker: readonly(worker),
    initializeWorker,
    executeWorkerMethod,
    destroyWorker
  }
}