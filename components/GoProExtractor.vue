<template>
  <div class="gopro-extractor">
    <div class="upload-section">
      <h2>GoPro MP4 Telemetry Extractor</h2>
      
      <div 
        class="file-drop-zone"
        :class="{ 'drag-over': isDragOver, 'processing': isProcessing }"
        @drop="handleDrop"
        @dragover.prevent="isDragOver = true"
        @dragleave="isDragOver = false"
        @click="fileInput?.click()"
      >
        <div v-if="!isProcessing" class="drop-content">
          <svg class="upload-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
          </svg>
          <p class="drop-text">
            Drop your GoPro MP4 file here or <span class="click-text">click to browse</span>
          </p>
          <p class="file-types">Supports GoPro MP4 files with telemetry data</p>
        </div>
        
        <div v-else class="processing-content">
          <div class="spinner"></div>
          <p>Processing video... {{ progress }}%</p>
        </div>
      </div>
      
      <input 
        ref="fileInput"
        type="file" 
        accept=".mp4,video/mp4" 
        style="display: none"
        @change="handleFileSelect"
      />
    </div>

    <div v-if="error" class="error-message">
      <svg class="error-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
              d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      {{ error }}
    </div>

    <div v-if="telemetryData" class="results-section">
      <div class="stats-grid">
        <div class="stat-card">
          <h3>GPS Points</h3>
          <p class="stat-number">{{ videoStats?.gpsPoints || 0 }}</p>
        </div>
        <div class="stat-card">
          <h3>Gyroscope</h3>
          <p class="stat-number">{{ videoStats?.gyroPoints || 0 }}</p>
        </div>
        <div class="stat-card">
          <h3>Accelerometer</h3>
          <p class="stat-number">{{ videoStats?.accelerometerPoints || 0 }}</p>
        </div>
        <div class="stat-card">
          <h3>Total Points</h3>
          <p class="stat-number">{{ videoStats?.totalDataPoints || 0 }}</p>
        </div>
      </div>

      <div class="action-buttons">
        <button 
          v-if="hasGpsData" 
          @click="downloadGpsAsGeoJson()" 
          class="btn btn-primary"
        >
          <svg class="btn-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
          Download GPS (GeoJSON)
        </button>
        
        <button 
          @click="downloadTelemetryAsJson()" 
          class="btn btn-secondary"
        >
          <svg class="btn-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          Download All Data (JSON)
        </button>
        
        <button @click="reset()" class="btn btn-outline">
          <svg class="btn-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Process Another File
        </button>
      </div>

      <!-- GPS Map Preview (if you want to add a map) -->
      <div v-if="hasGpsData" class="map-preview">
        <h3>GPS Track Preview</h3>
        <div class="map-placeholder">
          <p>Map would go here - integrate with your preferred mapping library</p>
          <p>Bounds: {{ gpsBounds }}</p>
          <p>Total GPS points: {{ gpsTrack.length }}</p>
        </div>
      </div>

      <!-- Data Preview -->
      <div class="data-preview">
        <h3>Data Preview</h3>
        <div class="data-tabs">
          <button 
            v-for="tab in dataTabs" 
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="['tab-button', { active: activeTab === tab.id }]"
          >
            {{ tab.label }} ({{ tab.count }})
          </button>
        </div>
        
        <div class="data-content">
          <div v-if="activeTab === 'gps' && hasGpsData" class="data-table">
            <div class="table-header">
              <span>Latitude</span>
              <span>Longitude</span>
              <span>Altitude</span>
              <span>Speed</span>
            </div>
            <div 
              v-for="(point, index) in telemetryData.gps.slice(0, 10)" 
              :key="index"
              class="table-row"
            >
              <span>{{ point.latitude.toFixed(6) }}</span>
              <span>{{ point.longitude.toFixed(6) }}</span>
              <span>{{ point.altitude.toFixed(1) }}m</span>
              <span>{{ point.speed2d.toFixed(1) }}m/s</span>
            </div>
            <p v-if="telemetryData.gps.length > 10" class="more-data">
              ... and {{ telemetryData.gps.length - 10 }} more points
            </p>
          </div>
          
          <div v-else-if="activeTab === 'gyro' && hasGyroData" class="data-table">
            <div class="table-header">
              <span>X</span>
              <span>Y</span>
              <span>Z</span>
              <span>Timestamp</span>
            </div>
            <div 
              v-for="(point, index) in telemetryData.gyro.slice(0, 10)" 
              :key="index"
              class="table-row"
            >
              <span>{{ point.x }}</span>
              <span>{{ point.y }}</span>
              <span>{{ point.z }}</span>
              <span>{{ point.timestamp }}</span>
            </div>
            <p v-if="telemetryData.gyro.length > 10" class="more-data">
              ... and {{ telemetryData.gyro.length - 10 }} more points
            </p>
          </div>
          
          <div v-else-if="activeTab === 'accel' && hasAccelerometerData" class="data-table">
            <div class="table-header">
              <span>X</span>
              <span>Y</span>
              <span>Z</span>
              <span>Timestamp</span>
            </div>
            <div 
              v-for="(point, index) in telemetryData.accelerometer.slice(0, 10)" 
              :key="index"
              class="table-row"
            >
              <span>{{ point.x }}</span>
              <span>{{ point.y }}</span>
              <span>{{ point.z }}</span>
              <span>{{ point.timestamp }}</span>
            </div>
            <p v-if="telemetryData.accelerometer.length > 10" class="more-data">
              ... and {{ telemetryData.accelerometer.length - 10 }} more points
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useGoProExtractor } from '../composables/useGoProExtractor'

const {
  isProcessing,
  telemetryData,
  error,
  progress,
  hasGpsData,
  hasGyroData,
  hasAccelerometerData,
  processFile,
  downloadGpsAsGeoJson,
  downloadTelemetryAsJson,
  getGpsTrack,
  getGpsBounds,
  getVideoStatistics,
  reset
} = useGoProExtractor()

const isDragOver = ref(false)
const activeTab = ref('gps')
const fileInput = ref<HTMLInputElement>()

const videoStats = computed(() => getVideoStatistics())
const gpsTrack = computed(() => getGpsTrack())
const gpsBounds = computed(() => getGpsBounds())

const dataTabs = computed(() => [
  { id: 'gps', label: 'GPS', count: telemetryData.value?.gps?.length || 0 },
  { id: 'gyro', label: 'Gyroscope', count: telemetryData.value?.gyro?.length || 0 },
  { id: 'accel', label: 'Accelerometer', count: telemetryData.value?.accelerometer?.length || 0 }
].filter(tab => tab.count > 0))

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
</script>

<style scoped>
.gopro-extractor {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.upload-section h2 {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
}

.file-drop-zone {
  border: 2px dashed #ddd;
  border-radius: 12px;
  padding: 60px 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #fafafa;
}

.file-drop-zone:hover,
.file-drop-zone.drag-over {
  border-color: #007bff;
  background: #f0f8ff;
}

.file-drop-zone.processing {
  cursor: not-allowed;
  opacity: 0.7;
}

.upload-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 20px;
  color: #666;
}

.drop-text {
  font-size: 18px;
  margin-bottom: 10px;
  color: #333;
}

.click-text {
  color: #007bff;
  text-decoration: underline;
}

.file-types {
  color: #666;
  font-size: 14px;
}

.processing-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #007bff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-message {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 15px;
  background: #fee;
  border: 1px solid #fcc;
  border-radius: 8px;
  color: #c33;
  margin: 20px 0;
}

.error-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.results-section {
  margin-top: 40px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  text-align: center;
}

.stat-card h3 {
  margin: 0 0 10px;
  color: #666;
  font-size: 14px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.stat-number {
  margin: 0;
  font-size: 24px;
  font-weight: bold;
  color: #333;
}

.action-buttons {
  display: flex;
  gap: 15px;
  margin-bottom: 30px;
  flex-wrap: wrap;
}

.btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  border-radius: 8px;
  border: none;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
  text-decoration: none;
}

.btn-icon {
  width: 16px;
  height: 16px;
}

.btn-primary {
  background: #007bff;
  color: white;
}

.btn-primary:hover {
  background: #0056b3;
}

.btn-secondary {
  background: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background: #545b62;
}

.btn-outline {
  background: transparent;
  color: #6c757d;
  border: 1px solid #6c757d;
}

.btn-outline:hover {
  background: #6c757d;
  color: white;
}

.map-preview {
  margin: 30px 0;
}

.map-placeholder {
  padding: 40px;
  background: #f8f9fa;
  border-radius: 8px;
  text-align: center;
  color: #666;
}

.data-preview {
  margin-top: 30px;
}

.data-tabs {
  display: flex;
  gap: 0;
  margin-bottom: 20px;
  border-bottom: 1px solid #ddd;
}

.tab-button {
  padding: 12px 20px;
  border: none;
  background: none;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: all 0.3s ease;
}

.tab-button:hover {
  background: #f8f9fa;
}

.tab-button.active {
  border-bottom-color: #007bff;
  color: #007bff;
  font-weight: 500;
}

.data-table {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.table-header {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  padding: 15px 20px;
  background: #f8f9fa;
  font-weight: 600;
  color: #333;
  border-bottom: 1px solid #ddd;
}

.table-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  padding: 12px 20px;
  border-bottom: 1px solid #eee;
  font-family: monospace;
  font-size: 13px;
}

.table-row:last-child {
  border-bottom: none;
}

.more-data {
  padding: 15px 20px;
  color: #666;
  font-style: italic;
  text-align: center;
  background: #f8f9fa;
}

@media (max-width: 768px) {
  .action-buttons {
    flex-direction: column;
  }
  
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .table-header,
  .table-row {
    grid-template-columns: repeat(2, 1fr);
    font-size: 12px;
  }
}
</style> 