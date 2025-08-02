<script setup lang="ts">
import { Button } from '~/components/ui/button'

</script>

<template>
  <section class="p-4 flex flex-col items-stretch h-full gap-x-8">
    <div class="flex-1 flex flex-col mx-auto">
      <section>
      <h2 class="max-w-2xl mb-6 text-4xl font-roboto-title font-bold leading-none tracking-tight md:text-5xl">ℹ️ About</h2>
      <h3 class="max-w-2xl mb-4 text-2xl font-roboto-title font-bold leading-none tracking-tight md:text-3xl">TrailTrace – What is this?</h3>
      <p class="max-w-2xl mb-6 font-light text-gray-500 lg:mb-8 md:text-lg lg:text-xl dark:text-gray-400">
        As a GoPro action camera owner, I was toying for quite some time with the idea of displaying a map alongside a video. While digging deeper, I discovered that GoPro cameras actually record much more than just GPS data — including acceleration, gyroscope readings, and even face detection
(see: <a href="https://github.com/gopro/gpmf-parser" class="text-sky-700" target="_blank">gopro/gpmf-parser</a>).
      </p>
      <p class="max-w-2xl mb-6 font-light text-gray-500 lg:mb-8 md:text-lg lg:text-xl dark:text-gray-400">
        Thanks to modern web technologies, all of this can now be extracted directly in your browser — no need to upload anything to a server.
      </p>
      <p class="max-w-2xl mb-6 font-light text-gray-500 lg:mb-8 md:text-lg lg:text-xl dark:text-gray-400">
        Note: Only tested with GoPro Hero 11 and 13 so far — other models might behave differently.
      </p>
      </section>
      <section>
      <h3 class="max-w-2xl mb-4 text-2xl font-roboto-title font-bold leading-none tracking-tight md:text-3xl">Tech Stuff</h3>
      <p class="max-w-2xl mb-6 font-light text-gray-500 lg:mb-8 md:text-lg lg:text-xl dark:text-gray-400">
        The core of the app is a <a href="https://github.com/chrigu/go-gpmf" class="text-sky-700" target="_blank">Go-based GPMF parser</a>, compiled to WebAssembly (WASM). This enables fast, local parsing of GoPro metadata directly in your browser.
      </p>
      <p class="max-w-2xl mb-6 font-light text-gray-500 lg:mb-8 md:text-lg lg:text-xl dark:text-gray-400">
        The visual interface is powered by a <a href="https://github.com/chrigu/trailtrace" class="text-sky-700" target="_blank">Nuxt.js frontend</a>, which lets you explore and inspect the extracted metadata — from maps to timestamps and raw sensor values.
      </p>
      </section>
      <section>
      <h4 class="text-xl md:text-2xl font-roboto-title font-bold text-gray-800 dark:text-white">
        Want to dig deeper?
      </h4>
      <p class="text-gray-600 dark:text-gray-300 font-light text-base md:text-lg">
        Both the parser and the frontend are open source. Explore the code, suggest improvements, or build your own tool!
      </p>
      <div class="flex flex-col sm:flex-row gap-4 mt-4">
        <Button variant="outline" size="lg" as-child>
          <a href="https://github.com/chrigu/go-gpmf" target="_blank" rel="noopener" class="w-full sm:w-auto">
            Go-based GPMF Parser ↗
          </a>
        </Button>
        <Button variant="outline" size="lg" as-child>
          <a href="https://github.com/chrigu/trailtrace" target="_blank" rel="noopener" class="w-full sm:w-auto">
            Nuxt Frontend ↗
          </a>
        </Button>
      </div>
      </section>
    </div>
  </section>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.4s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>

