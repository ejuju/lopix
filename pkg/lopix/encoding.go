package lopix

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

const HeaderSize = (1 + 1) + (16 * 4) // Width + Height + Palette
const MaxEncodedFrameSize = HeaderSize + (math.MaxUint8 * math.MaxUint8)

func (f *Frame) ReadFrom(r io.Reader) (n int64, err error) {
	b := [HeaderSize]byte{}
	nRead, err := io.ReadFull(r, b[:])
	n += int64(nRead)
	if err != nil {
		return n, fmt.Errorf("read header: %w", err)
	}
	// NB: Since a value of 0 would be invalid, we assume offset value by one.
	// This allows a maximum of 256x256 cells instead of 255x255.
	f.w, f.h = int(b[0])+1, int(b[1])+1
	f.palette = [16]Color{}
	for i := range 16 {
		f.palette[i] = Color(binary.BigEndian.Uint32(b[2+(i*4) : 2+(i*4)+4]))
	}

	f.grid = make([]byte, int(f.w)*int(f.h))
	nRead, err = io.ReadFull(r, f.grid)
	n += int64(nRead)
	if err != nil {
		return n, fmt.Errorf("read grid: %w", err)
	}

	return n, err
}
