# GoPro Metadata Display / GPMF Parser

Displays Gopro's  metadata like:
- 🌍 GPS
- 🚀 Acceleration
- ☀️ Luminance
- 🎨 Hue
- 😀 Face
- 🎬 Scene

There are two sub-projects

1. Go project for the GoPro Metadata parser (GPMF Parser) that can be complied to WASM
2. Nuxt project for displaying the data

## GPMF Parser

Moded to https://github.com/chrigu/go-gpmf

Use generated wasm file in this project.

## Nuxt

See `README.md` in the nuxt directory

## Todos

- Design
- Fix Accleration display
- Refactor timed data, extraction
- Test older GoPros
- Select metadata to export
- Tests
- Refactoring
- Optimize performance https://goperf.dev/01-common-patterns/mem-prealloc/#why-preallocation-matters

## Resources
- https://github.com/gopro/gpmf-parser
- https://www.trekview.org/blog/injecting-camm-gpmd-telemetry-videos-part-3-mp4-structure-telemetry-trak/
- https://developer.apple.com/documentation/quicktime-file-format/sample-to-chunk_atom/sample-to-chunk_table