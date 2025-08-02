<script setup>
import { Line } from 'vue-chartjs'
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend } from 'chart.js'
import { useStore } from '~/stores'
import { computed } from 'vue'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend)

const store = useStore()

const chartData = computed(() => ({
  labels: store.accelerationData.map(d => (d.timestamp - store.accelerationData[0].timestamp)/1000),
  datasets: [
    {
      label: 'X',
      data: store.accelerationData.map(d => d.x),
      borderColor: 'rgb(59, 130, 246)',
      tension: 0.1,
      pointRadius: store.accelerationData.map(d => d.timestamp === store.currentAccelerationData?.timestamp ? 5 : 0)
    },
    {
      label: 'Y',
      data: store.accelerationData.map(d => d.y),
      borderColor: 'rgb(236, 72, 153)',
      tension: 0.1,
      pointRadius: store.accelerationData.map(d => d.timestamp === store.currentAccelerationData?.timestamp ? 5 : 0)
    },
    {
      label: 'Z',
      data: store.accelerationData.map(d => d.z),
      borderColor: 'rgb(249, 115, 22)',
      tension: 0.1,
      pointRadius: store.accelerationData.map(d => d.timestamp === store.currentAccelerationData?.timestamp ? 5 : 0)
    }
  ]
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    y: {
      title: {
        display: true,
        text: 'm/s^2'
      },
      ticks: {
        callback: (value) => `${value.toFixed(2)}`
      }
    },
    x: {
      title: {
        display: true,
        text: 'Time',
      }
    }
  }
}
</script>

<template>
  <div class="w-full h-[400px]">
    <Line
      v-if="chartData.labels.length > 0"
      :data="chartData"
      :options="chartOptions"
    />
  </div>
</template> 