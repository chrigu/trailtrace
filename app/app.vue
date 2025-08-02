<script setup lang="ts">

const route = useRoute();

useHead({
  bodyAttrs: {
    class: "bg-stone-100"
  },
  title: 'TrailTrace – GoPro GPS Visualizer',
  meta: [
    { name: 'description', content: 'Free browser-based tool to extract and visualize GoPro telemetry data like GPS, acceleration and more.' },
    { property: 'og:title', content: 'TrailTrace – GoPro GPS Visualizer' },
    { property: 'og:description', content: 'Free browser-based tool to extract and visualize GoPro telemetry data like GPS, acceleration and more.' },
    { property: 'og:image', content: 'https://trailtrace.video/og-image.jpg' },
    { property: 'twitter:card', content: 'summary_large_image' },
    { property: 'twitter:title', content: 'TrailTrace – GoPro GPS Visualizer' },
    { property: 'twitter:description', content: 'Free browser-based tool to extract and visualize GoPro telemetry data like GPS, acceleration and more.' },
    { property: 'twitter:image', content: 'https://trailtrace.video/og-image.jpg' },
  ],
})

</script>

<template>
  <header class="mb-0 py-4 px-10 flex flex-col md:flex-row justify-between">
    <h1 class="text-sky-700 font-roboto-title text-2xl">TrailTrace</h1>
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
        Explorer
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
      <NuxtLink 
        to="/about" 
        :class="[
          'transition-colors duration-200',
          route.path === '/about' 
            ? 'text-sky-700 font-semibold border-b-2 border-sky-700' 
            : 'text-gray-600 hover:text-sky-600'
        ]"
      >
        About
      </NuxtLink>
    </div>
  </header>
  <main class="px-10">
    <NuxtPage />
  </main>
</template>