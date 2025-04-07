# GoPro GPS Display / GPMF Parser

Displays Gopro's GPS data on a map.

There are two sub-projects

1. Go project for the GoPro Metadata parser (GPMF Parser) that can be complied to WASM
2. Nuxt project for displaying the data

## GPMF Parser

All commands are run from the `gpmf` directory

### Run cli 
`go run ./cli <mp4 file>`

### Compile to WASM

`GOOS=js GOARCH=wasm go build -o wasm/main.wasm ./wasm`

## Nuxt

See `README.md` in the nuxt directory

## Todos

- Move GPMF Parser to own project
- Tests
- Refactoring
- Optimize performance https://goperf.dev/01-common-patterns/mem-prealloc/#why-preallocation-matters

## Ressources
- https://github.com/gopro/gpmf-parser
- https://www.trekview.org/blog/injecting-camm-gpmd-telemetry-videos-part-3-mp4-structure-telemetry-trak/
- https://developer.apple.com/documentation/quicktime-file-format/sample-to-chunk_atom/sample-to-chunk_table