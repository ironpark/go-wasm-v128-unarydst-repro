# Go wasm: V128Store offsets and SIMD lane indices are encoded as 0

Reproductions for a bug in Go's wasm backend (`cmd/internal/obj/wasm`): the immediate of `V128Store` (the memory offset) and of the SIMD lane ops (`*ExtractLane*`, `*ReplaceLane`, the lane index) are silently encoded as 0.

It affects hand-written assembly and, with `GOEXPERIMENT=simd`, the code the compiler generates for `//go:wasmimport` and `//go:wasmexport` functions that pass SIMD vectors.

Reproduced with go1.27.1 and go1.28-devel_77fb1abe (master as of 2026-10-05) on darwin/arm64, with Node.js v26.1.0.

## Run

Requires Go 1.27 or later and Node.js (with SIMD and `node:wasi`, e.g. v20 or later).

```sh
./run.sh                       # with the go command on PATH
GO=gotip GOTOOLCHAIN=local ./run.sh  # with another toolchain
```

## Repro 1: assembly (`repro1-asm`)

`GOOS=js GOARCH=wasm`. [`main_wasm.s`](repro1-asm/main_wasm.s) uses `V128Store $16`, `I32x4ExtractLane $3`, and, as a control, `V128Load $16`.

```
V128Store $16:       [1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] (want the bytes at 16..31)
I32x4ExtractLane $3: 10 (want 13)
V128Load $16:        [16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31] (want 16..31; control)
```

## Repro 2: `//go:wasmimport` returning a SIMD vector (`repro2-wasmimport`)

`GOOS=wasip1 GOARCH=wasm GOEXPERIMENT=simd`, no assembly. [`main.go`](repro2-wasmimport/main.go) imports `vec() archsimd.Int32x4`, which [`run.mjs`](repro2-wasmimport/run.mjs) provides from a tiny wasm module returning `v128.const i32x4 1 2 3 4` (a v128 cannot cross the JS boundary, so the host function must be wasm).

```
wasmimport v128 result: [3 4 112 0] (want [1 2 3 4])
```

The import wrapper stores the result at SP+0 instead of its result slot, overwriting another part of the frame. The garbage differs between builds.

## Cause

The encoder in `cmd/internal/obj/wasm` reads these immediates from `p.To.Offset`, as the compiler's SSA backend sets them (`cmd/compile/internal/wasm/simdssa.go`). But the ops are missing from `unaryDst`, so both `cmd/asm` and the `appendp` helper in `wasmobj.go` (used for the `//go:wasmimport` results and the `//go:wasmexport` parameters) put the single operand in `p.From`, leaving `p.To.Offset` 0. The scalar stores (`I32Store`, `F64Store`, …) are in `unaryDst` and work.

## Fix

[`fix.patch`](fix.patch) adds `AV128Store` and the 14 lane ops to `unaryDst`. Applied to master (`77fb1abe`), both repros print the expected values, and `go test cmd/asm/internal/asm` and `GOEXPERIMENT=simd GOOS=js GOARCH=wasm go test simd/...` pass.

```sh
cd "$(gotip env GOROOT)" && git apply /path/to/fix.patch
cd src && GOTOOLCHAIN=local ../bin/go install cmd/asm cmd/compile cmd/link
```

Expected output with the fix:

```
V128Store $16:       [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16] (want the bytes at 16..31)
I32x4ExtractLane $3: 13 (want 13)
V128Load $16:        [16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31] (want 16..31; control)
wasmimport v128 result: [1 2 3 4] (want [1 2 3 4])
```
