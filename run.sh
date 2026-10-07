#!/bin/sh
# Runs both reproductions with the go command on PATH (or $GO) and Node.js.
set -e
GO=${GO:-go}
cd "$(dirname "$0")"
"$GO" version
node --version

echo
echo '== repro1-asm (GOOS=js GOARCH=wasm)'
(cd repro1-asm && GOOS=js GOARCH=wasm "$GO" run -exec "$("$GO" env GOROOT)/lib/wasm/go_js_wasm_exec" .)

echo
echo '== repro2-wasmimport (GOOS=wasip1 GOARCH=wasm GOEXPERIMENT=simd)'
(cd repro2-wasmimport && GOEXPERIMENT=simd GOOS=wasip1 GOARCH=wasm "$GO" build -o main.wasm . && node --no-warnings run.mjs)
