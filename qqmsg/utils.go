package qqmsg

import (
	"encoding/binary"
	"math/big"
	"reflect"
	"unicode/utf16"
	"unsafe"
)

func DecodeUtf16(b []byte) string {
	return string(utf16.Decode(b2u16(b)))
}

func b2u16(b []byte) []uint16 {
	b = b[:len(b)/2]
	return *(*[]uint16)(unsafe.Pointer(&b))
}

type Buffer struct {
	buf []byte
	off int
}

func NewBuffer(buf []byte) *Buffer {
	return &Buffer{
		buf: buf,
		off: 0,
	}
}

func (b *Buffer) empty() bool { return len(b.buf) <= b.off }

func (b *Buffer) Read(n int) []byte {
	if b.empty() {
		return nil
	}
	r := make([]byte, n)
	copy(r, b.buf[b.off:])
	b.off += n
	return r
}

func (b *Buffer) Skip(n int) {
	b.off += n
}

func (b *Buffer) Uint32() (u uint32) {
	u = binary.LittleEndian.Uint32(b.buf[b.off:])
	b.off += 4
	return
}

func (b *Buffer) Uint16() (u uint16) {
	u = binary.LittleEndian.Uint16(b.buf[b.off:])
	b.off += 2
	return
}

func (b *Buffer) Byte() (by byte) {
	by = b.buf[b.off]
	b.off += 1
	return
}

func (b *Buffer) T() uint8 {
	return b.Byte()
}

func (b *Buffer) L() uint16 {
	return b.Uint16()
}

func (b *Buffer) TLV() (t uint8, l uint16, v []byte) {
	t = b.T()
	l = b.L()
	v = b.Read(int(l))
	return
}

const b48Dict = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ~@$%`(){}[]_"

var bigZero = big.NewInt(0)
var bigRadix = big.NewInt(48)

func encodeB48(b []byte) string {
	x := new(big.Int)
	y := (*Int)(unsafe.Pointer(x))
	yh := (*reflect.SliceHeader)(unsafe.Pointer(&y.abs))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	yh.Data = bh.Data
	yh.Cap = bh.Cap / int(unsafe.Sizeof(uint(0)))
	yh.Len = bh.Len / int(unsafe.Sizeof(uint(0)))

	ans := make([]byte, 0, 23)
	for x.CmpAbs(bigZero) > 0 {
		mod := new(big.Int)
		x.DivMod(x, bigRadix, mod)
		ans = append(ans, b48Dict[mod.Int64()])
	}
	return string(ans)
}

type Int struct {
	neg bool   // sign
	abs []uint // absolute value of the integer
}
