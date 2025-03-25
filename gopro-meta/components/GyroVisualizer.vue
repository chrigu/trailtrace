<template>
  <div ref="container" class="w-full h-screen relative">
  </div>
</template>

<script setup>
import * as THREE from 'three'
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useStore } from "@/store";

const container = ref(null)
let scene, camera, renderer, cube = 0
let animationId = null

const store = useStore();

const setupScene = () => {
  scene = new THREE.Scene()
  camera = new THREE.PerspectiveCamera(75, container.value.clientWidth / container.value.clientHeight, 0.1, 1000)
  renderer = new THREE.WebGLRenderer({ antialias: true })
  renderer.setSize(container.value.clientWidth, container.value.clientHeight)
  container.value.appendChild(renderer.domElement)

  // Cube
  const geometry = new THREE.BoxGeometry()
  const material = new THREE.MeshNormalMaterial()
  cube = new THREE.Mesh(geometry, material)
  scene.add(cube)
  camera.position.z = 5
}

watch(() => store.currentGyroData, (newGyroData) => {
  const { x, y, z } = newGyroData

  cube.rotation.x = THREE.MathUtils.degToRad(x)
  cube.rotation.y = THREE.MathUtils.degToRad(y)
  cube.rotation.z = THREE.MathUtils.degToRad(z)

  renderer.render(scene, camera)
})

onMounted(() => {
  setupScene()
})

onUnmounted(() => {
  cancelAnimationFrame(animationId)
  renderer.dispose()
})
</script>

<style scoped>
canvas {
  display: block;
}
</style>
