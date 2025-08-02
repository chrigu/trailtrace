# Trailtrace

Test it online https://trailtrace.video/

Trailtrace let's you explore the follwing GoPro metadata from your mp4-files directly in your browser:
- 🌍 GPS
- 🚀 Acceleration
- ☀️ Luminance
- 🎨 Hue
- 😀 Face
- 🎬 Scene

Your footage is processed locally — no need to upload anything to a server.

Other features:
- 💾 Extract GPMF data
- 🛰️ Extract GPX data

## GPMF Parser

The GPMF (GoPro Metadata Format) parser is written in Go and complied to WASM. 

The mp4-file is never read as a whole, this allows to process files > 2GB, which is the memory limit for WASM in most browsers.

You can find the Go implementation at https://github.com/chrigu/go-gpmf

## Nuxt

### Setup

Make sure to install dependencies:

```bash
# npm
npm install

# pnpm
pnpm install

# yarn
yarn install

# bun
bun install
```

### Development Server

Start the development server on `http://localhost:3000`:

```bash
# npm
npm run dev

# pnpm
pnpm dev

# yarn
yarn dev

# bun
bun run dev
```

### Production

Build the application for production:

```bash
# npm
npm run build

# pnpm
pnpm build

# yarn
yarn build

# bun
bun run build
```

Locally preview production build:

```bash
# npm
npm run preview

# pnpm
pnpm preview

# yarn
yarn preview

# bun
bun run preview
```

## Todo
- Rename main.wasm to something meaningful
- Update README