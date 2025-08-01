# Trailtrace

See it online https://trailtrace.video/

Displays Gopro's  metadata like:
- 🌍 GPS
- 🚀 Acceleration
- ☀️ Luminance
- 🎨 Hue
- 😀 Face
- 🎬 Scene

Other features:
- 💾 Extract GPMF data
- 🛰️ Extract GPX data

## GPMF Parser

Uses wasm file from https://github.com/chrigu/go-gpmf

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

Check out the [deployment documentation](https://nuxt.com/docs/getting-started/deployment) for more information.

## Todo
- Rename main.wasm to something meaningful
- Update README
