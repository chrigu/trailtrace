/* eslint-disable no-restricted-globals */
importScripts("/wasm_exec.js");          // <- the helper shipped by Go
console.log('worker_wasm.js imported');
// 1. boot Go + Wasm once -------------------------------------------------
const go = new self.Go();
const ready = (async () => {
  const { instance } = await WebAssembly.instantiateStreaming(
    fetch("/main.wasm"), go.importObject);
  go.run(instance);                      // boots package main (exports appear on globalThis)
})();

// 2. helper the Go shim calls:  (offset,len) → Uint8Array  ---------------
self.getBuf = (offset, length) => {
  // File object is stored in module-scoped variable `file`
  const blob  = file.slice(Number(offset), Number(offset) + Number(length));
  const abuf  = new FileReaderSync().readAsArrayBuffer(blob);   // sync, 0-copy in JS heap
  return new Uint8Array(abuf);            // Uint8Array is what js.CopyBytesToGo expects
};

let file;                                 // current file (set per request)

self.onmessage = async ({ data }) => {
  const { id, method, file: f } = data;
  file = f; // needed for getBuf()

  try {
    await ready;

    let result;

    switch (method) {
      case "exportGPMF":
        result = await globalThis.exportGPMF(file); // returns Uint8Array
        postMessage({ id, ok: true, payload: result.buffer }, [result.buffer]);
        break;

      case "processFile":
        result = await globalThis.processFile(file); // returns JS object
        postMessage({ id, ok: true, payload: result }); // no transfer needed
        break;

      default:
        throw new Error(`Unknown method: ${method}`);
    }
  } catch (err) {
    console.error(`Worker error in ${method}:`, err);
    postMessage({ id, ok: false, error: String(err) });
  }
};