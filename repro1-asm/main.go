package main

import "fmt"

// store16 stores the 16 bytes at src to dst+16 with V128Store $16.
//
//go:noescape
func store16(dst, src *byte)

// lane3 returns the lane 3 of the four uint32 at p with I32x4ExtractLane $3.
//
//go:noescape
func lane3(p *uint32) uint32

// load16 loads the 16 bytes at src+16 with V128Load $16 and stores them to dst.
//
//go:noescape
func load16(dst, src *byte)

var lanes = [4]uint32{10, 11, 12, 13}

func main() {
	var dst [32]byte
	var src [16]byte
	for i := range src {
		src[i] = byte(i + 1)
	}
	store16(&dst[0], &src[0])
	fmt.Println("V128Store $16:      ", dst, "(want the bytes at 16..31)")

	fmt.Println("I32x4ExtractLane $3:", lane3(&lanes[0]), "(want 13)")

	var src2 [32]byte
	for i := range src2 {
		src2[i] = byte(i)
	}
	var dst2 [16]byte
	load16(&dst2[0], &src2[0])
	fmt.Println("V128Load $16:       ", dst2, "(want 16..31; control)")
}
