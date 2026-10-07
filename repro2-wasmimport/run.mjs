// Usage: GOEXPERIMENT=simd GOOS=wasip1 GOARCH=wasm go build -o main.wasm . && node run.mjs
import { readFile } from 'node:fs/promises';
import { WASI } from 'node:wasi';

// A module exporting vec: () -> v128, which returns v128.const i32x4 1 2 3 4.
const vecModule = new Uint8Array([
  0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00, // header
  0x01, 0x05, 0x01, 0x60, 0x00, 0x01, 0x7b,       // type: () -> v128
  0x03, 0x02, 0x01, 0x00,                         // func 0: type 0
  0x07, 0x07, 0x01, 0x03, 0x76, 0x65, 0x63, 0x00, 0x00, // export "vec"
  0x0a, 0x16, 0x01, 0x14, 0x00, 0xfd, 0x0c,       // code: v128.const
  0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00,
  0x03, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00,
  0x0b,                                           // end
]);
const vec = await WebAssembly.instantiate(vecModule);

const wasi = new WASI({ version: 'preview1' });
const { instance } = await WebAssembly.instantiate(await readFile('main.wasm'), {
  wasi_snapshot_preview1: wasi.wasiImport,
  test: { vec: vec.instance.exports.vec },
});
wasi.start(instance);
