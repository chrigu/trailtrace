<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { ref } from 'vue'

interface DemoFile {
  name: string
  url: string
  description?: string
  thumbnail?: string
}

const props = defineProps<{
  demoFiles?: DemoFile[]
  compact?: boolean
}>()

const emit = defineEmits(['file-selected'])

const isLoadingDemo = ref<string | null>(null)

const handleDemoFile = async (demoFile: DemoFile) => {
  try {
    isLoadingDemo.value = demoFile.name
    
    // Fetch the file from S3
    const response = await fetch(demoFile.url)
    if (!response.ok) {
      throw new Error(`Failed to fetch demo file: ${response.statusText}`)
    }
    
    // Convert to blob and then to File object
    const blob = await response.blob()
    const file = new File([blob], demoFile.name, { 
      type: blob.type || 'video/mp4' 
    })
    
    // Emit the same event as regular file selection
    emit('file-selected', file)
  } catch (error) {
    console.error('Error loading demo file:', error)
    // You might want to show a toast notification here
  } finally {
    isLoadingDemo.value = null
  }
}
</script>

<template>
  <div v-if="demoFiles && demoFiles.length > 0" class="space-y-4">
    <div class="text-center" v-if="!compact">
      <h3 class="text-lg font-medium text-gray-900 mb-2">Demo Files</h3>
      <p class="text-sm text-gray-500 mb-4">Try a demo file:</p>
    </div>
    
    <div :class="[
      compact 
        ? 'mx-auto' 
        : 'grid grid-cols-1 md:grid-cols-1 lg:grid-cols-1 gap-4'
    ]">
      <div 
        v-for="demoFile in demoFiles" 
        :key="demoFile.name"
        :class="[
          'border border-gray-200 rounded-lg hover:border-sky-400 transition-colors',
          compact ? 'p-2' : 'p-4'
        ]"
      >
        <!-- Thumbnail if available and not compact -->
        <div v-if="demoFile.thumbnail && !compact" class="mb-3">
          <img 
            :src="demoFile.thumbnail" 
            :alt="demoFile.name"
            class="w-full h-64 object-cover rounded"
          />
        </div>
        
        <!-- File info -->
        <div :class="compact ? 'space-y-1 text-center' : 'space-y-2'">
          <h4 :class="[
            'font-medium text-gray-900',
            compact ? 'text-sm' : 'text-base'
          ]">
            {{ demoFile.name }}
          </h4>
          
          <p v-if="demoFile.description && !compact" class="text-sm text-gray-500">
            {{ demoFile.description }}
          </p>
          
          <!-- Load button -->
          <Button 
            @click="handleDemoFile(demoFile)"
            :disabled="isLoadingDemo === demoFile.name"
            :class="compact ? 'w-auto' : 'w-full'"
            :size="compact ? 'sm' : 'sm'"
          >
            <span v-if="isLoadingDemo === demoFile.name">Loading...</span>
            <span v-else>{{ compact ? 'Load' : 'Load Demo' }}</span>
          </Button>
        </div>
      </div>
    </div>
  </div>
</template> 