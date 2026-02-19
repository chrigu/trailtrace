<script setup lang="ts">
import DataPrivacy from '~/components/DataPrivacy.vue';

const route = useRoute();

useHead({
  bodyAttrs: {
    class: "bg-background"
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
  <div class="min-h-screen flex flex-col bg-gradient-to-br from-background via-background to-blue-50">
    <header class="mb-0 py-4 px-10 flex flex-col md:flex-row justify-between border-b border-border bg-white/50 backdrop-blur-sm">
      <h1 class="text-primary font-roboto-title text-2xl font-bold">TrailTrace</h1>
    <div class="flex flex-row gap-6 items-center">
      <NuxtLink
        to="/"
        :class="[
          'transition-all duration-200 font-medium',
          route.path === '/'
            ? 'text-primary font-semibold border-b-2 border-primary'
            : 'text-muted-foreground hover:text-primary'
        ]"
      >
        Explorer
      </NuxtLink>
      <NuxtLink
        to="/gpmf-extractor"
        :class="[
          'transition-all duration-200 font-medium',
          route.path === '/gpmf-extractor'
            ? 'text-primary font-semibold border-b-2 border-primary'
            : 'text-muted-foreground hover:text-primary'
        ]"
      >
        Extractor
      </NuxtLink>
      <NuxtLink
        to="/about"
        :class="[
          'transition-all duration-200 font-medium',
          route.path === '/about'
            ? 'text-primary font-semibold border-b-2 border-primary'
            : 'text-muted-foreground hover:text-primary'
        ]"
      >
        About
      </NuxtLink>
    </div>
  </header>
  <main class="px-10 flex-grow">
    <NuxtPage />
  </main>
  <DataPrivacy />
  </div>
</template>
