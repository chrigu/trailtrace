/* eslint-disable no-restricted-globals */
importScripts("/wasm_exec.js");
console.log('worker_wasm.js imported');
const go = new self.Go();
const ready = (async () => {
  const { instance } = await WebAssembly.instantiateStreaming(
    fetch("/main.wasm"), go.importObject);
  go.run(instance);
})();

self.getBuf = (offset, length) => {
  // File object is stored in module-scoped variable `file`
  const blob  = file.slice(Number(offset), Number(offset) + Number(length));
  const abuf  = new FileReaderSync().readAsArrayBuffer(blob);
  return new Uint8Array(abuf);
};

let file;

self.onmessage = async ({ data }) => {
  const { id, method, file: f } = data;
  file = f;

  try {
    await ready;

    let result;

    switch (method) {
      case "exportGPMF":
        result = await globalThis.exportGPMF(file);
        postMessage({ id, ok: true, payload: result.buffer }, [result.buffer]);
        break;

      case "processFile":
        result = await globalThis.processFile(file);
        postMessage({ id, ok: true, payload: result });
        break;

      default:
        throw new Error(`Unknown method: ${method}`);
    }
  } catch (err) {
    console.error(`Worker error in ${method}:`, err);
    postMessage({ id, ok: false, error: String(err) });
  }
};