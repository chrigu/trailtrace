// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },
  modules: ['@nuxtjs/leaflet', '@pinia/nuxt', '@nuxtjs/tailwindcss', 'shadcn-nuxt', 'nuxt-umami'],
  app: {
    head: {
      link: [
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Inter:wght@400;700&display=swap',
        },
      ],
    },
  },
  pinia: {
    storesDirs: ['./stores/**',],
  },
  shadcn: {
    prefix: "",
    componentDir: "./app/components/ui",
  },
  nitro: {
    compressPublicAssets: { brotli: true, gzip: true } // includes application/wasm
  },
  umami: {
    id: 'c63f819a-6bb0-4bef-9980-7325ce80f2df',
    host: 'https://cloud.umami.is',
    autoTrack: true,
    domains: ['trailtrace.video']
  },
})

