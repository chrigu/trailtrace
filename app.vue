<script setup lang="ts">
import { onMounted, ref} from "vue";

let go: any;
const route = useRoute();

const loadWasmExec = async () => {
  return new Promise((resolve, reject) => {
    const script = document.createElement("script");
    script.src = "/wasm_exec.js"; // Load from public folder
    script.onload = resolve;
    script.onerror = reject;
    document.body.appendChild(script);
  });
};

onMounted(async () => {
  try {
    await loadWasmExec();
    go = new (window as any).Go();
    const wasmResponse = await fetch("/main.wasm");
    const wasmBytes = await wasmResponse.arrayBuffer();
    const { instance } = await WebAssembly.instantiate(wasmBytes, go.importObject);
    go.run(instance);
  } catch (error) {
    console.error("Error loading WebAssembly:", error);
  }
});

useHead({
  bodyAttrs: {
    class: "bg-stone-100"
  }
})

</script>

<template>
  <header class="mb-0 py-4 px-10 flex flex-col md:flex-row justify-between">
    <h1 class="text-sky-700 font-roboto-title font-bold text-2xl">GPMF Explorer</h1>
    <div class="flex flex-row gap-4 items-center">
      <NuxtLink 
        to="/" 
        :class="[
          'transition-colors duration-200',
          route.path === '/' 
            ? 'text-sky-700 font-semibold border-b-2 border-sky-700' 
            : 'text-gray-600 hover:text-sky-600'
        ]"
      >
        Visualizer
      </NuxtLink>
      <NuxtLink 
        to="/gpmf-extractor" 
        :class="[
          'transition-colors duration-200',
          route.path === '/gpmf-extractor' 
            ? 'text-sky-700 font-semibold border-b-2 border-sky-700' 
            : 'text-gray-600 hover:text-sky-600'
        ]"
      >
        Extractor
      </NuxtLink>
    </div>
  </header>
  <main class="px-10">
    <NuxtPage />
  </main>
</template>