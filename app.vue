<script setup lang="ts">
import { onMounted} from "vue";
import { useStore } from "~/store";

let go;
const store = useStore();

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
    go = new Go();
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
    class: "bg-stone-50"
  }
})
</script>

<template>
  <section class="mb-0 p-4 flex flex-col md:flex-row justify-between">
    <h1 class="text-sky-700 font-roboto-title font-bold text-2xl">GoPro Metadata</h1>
    <NuxtLink to="/gpmf-export">GPMF Export</NuxtLink>
    <Transition name="fade">
      <UploadButton v-if="store.videoUrl"/>
    </Transition>
  </section>
  <NuxtPage />
</template>
