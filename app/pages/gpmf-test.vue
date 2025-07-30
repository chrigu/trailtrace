<template>
  <div class="container mx-auto px-4 py-8 max-w-4xl">
    <h1 class="text-3xl font-bold text-center mb-8">GPMF Binary Extractor Test</h1>
    
    <!-- File Input Section -->
    <div class="mb-8">
      <div 
        class="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center transition-colors duration-200"
        :class="{ 'border-blue-500 bg-blue-50': isDragOver, 'border-red-500 bg-red-50': error }"
        @drop="handleDrop"
        @dragover.prevent="isDragOver = true"
        @dragleave="isDragOver = false"
        @click="fileInput?.click()"
      >
        <div v-if="!isProcessing" class="cursor-pointer">
          <svg class="mx-auto h-12 w-12 text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
          </svg>
          <p class="text-lg text-gray-600 mb-2">
            Drop your GoPro MP4 file here or <span class="text-blue-600 underline">click to browse</span>
          </p>
          <p class="text-sm text-gray-500">Supports GoPro MP4 files (Hero5 and later)</p>
        </div>
        
        <div v-else class="flex flex-col items-center">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mb-4"></div>
          <p class="text-lg text-gray-600">Processing... {{ Math.round(progress) }}%</p>
        </div>
      </div>
      
      <input 
        ref="fileInput"
        type="file" 
        accept=".mp4,video/mp4" 
        class="hidden"
        @change="handleFileSelect"
      />
    </div>

    <!-- Error Display -->
    <div v-if="error" class="mb-6 p-4 bg-red-100 border border-red-400 text-red-700 rounded-lg">
      <div class="flex">
        <svg class="h-5 w-5 mr-2 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
        </svg>
        <div>
          <h3 class="text-sm font-medium">Error</h3>
          <p class="text-sm">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- Success Message -->
    <div v-if="downloadComplete" class="mb-6 p-4 bg-green-100 border border-green-400 text-green-700 rounded-lg">
      <div class="flex">
        <svg class="h-5 w-5 mr-2 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
        </svg>
        <div>
          <h3 class="text-sm font-medium">Success!</h3>
          <p class="text-sm">GPMF binary data extracted and downloaded successfully.</p>
          <p class="text-xs mt-1">File: {{ lastDownloadedFile }}</p>
        </div>
      </div>
      
      <div class="mt-4">
        <button 
          @click="reset"
          class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
        >
          Process Another File
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

// Reactive state
const isDragOver = ref(false)
const isProcessing = ref(false)
const progress = ref(0)
const error = ref('')
const downloadComplete = ref(false)
const fileName = ref('')
const fileSize = ref(0)
const lastDownloadedFile = ref('')
const fileInput = ref<HTMLInputElement>()

// Methods
const handleDrop = (event: DragEvent) => {
  event.preventDefault()
  isDragOver.value = false
  
  const files = event.dataTransfer?.files
  if (files && files.length > 0) {
    processFile(files[0])
  }
}

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  const files = target.files
  if (files && files.length > 0) {
    processFile(files[0])
  }
}

const processFile = async (file: File) => {
  if (!file) return
  
  // Reset state
  error.value = ''
  downloadComplete.value = false
  isProcessing.value = true
  progress.value = 0
  fileName.value = file.name
  fileSize.value = file.size
  
  try {
    // Validate file type
    if (!file.name.toLowerCase().endsWith('.mp4') && file.type !== 'video/mp4') {
      throw new Error('Please select an MP4 file')
    }
    
    // Dynamic import of gpmf-extract for browser compatibility
    const { default: gpmfExtract } = await import('gpmf-extract')
    
    // Extract GPMF data with progress tracking
    const result = await gpmfExtract(file, {
      browserMode: true,
      progress: (percent: number) => {
        progress.value = percent
      }
    })
    
    if (!result) {
      throw new Error('No GPMF data found in this file')
    }
    
    console.log('GPMF extraction completed:', result)
    
    // Automatically download the binary data
    const blob = new Blob([result.rawData], { type: 'application/octet-stream' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    const downloadFileName = `${file.name.replace('.mp4', '')}_gpmf_binary.bin`
    a.download = downloadFileName
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
    
    // Show success message
    lastDownloadedFile.value = downloadFileName
    downloadComplete.value = true
    
  } catch (err) {
    console.error('Error extracting GPMF data:', err)
    error.value = err instanceof Error ? err.message : 'An unknown error occurred'
  } finally {
    isProcessing.value = false
    progress.value = 0
  }
}

const reset = () => {
  downloadComplete.value = false
  error.value = ''
  fileName.value = ''
  fileSize.value = 0
  lastDownloadedFile.value = ''
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

// Set page title
useHead({
  title: 'GPMF Binary Extractor Test'
})
</script> 