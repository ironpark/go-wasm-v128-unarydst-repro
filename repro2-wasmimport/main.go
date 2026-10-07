package main

import (
	"fmt"
	"simd/archsimd"
)

// vec is provided by the host and returns the v128 (1, 2, 3, 4).
//
//go:wasmimport test vec
func vec() archsimd.Int32x4

func main() {
	var a [4]int32
	vec().StoreArray(&a)
	fmt.Println("wasmimport v128 result:", a, "(want [1 2 3 4])")
}
