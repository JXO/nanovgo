// +build darwin linux
// +build arm arm64

package vg

import (
	"log"
	"unsafe"
)

type Float float32

var shaderHeader string = `
#define DAVINCI_GL2 1
#define UNIFORMARRAY_SIZE 11
`
//var shaderHeader = `
//#define DAVINCI_GL3 1
//#define USE_UNIFORMBUFFER 1
//`

func prepareTextureBuffer(data []byte, w, h, bpp int) []byte {
	return data
}

func castFloat32ToByte(vertexes []float32) []byte {
	// Convert []float32 list to []byte without copy
	var b []byte
	b = (*(*[1 << 20]byte)(unsafe.Pointer(&vertexes[0])))[:len(vertexes)*4]
	return b
}

func dumpLog(fmt string, values ...interface{}) {
	log.Printf(fmt, values...)
}
