package lopix

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"math"
)

type Frame struct {
	w       int
	h       int
	palette Palette
	grid    []byte
}

func F(w, h int, p Palette, rows ...string) (f *Frame) {
	if w > math.MaxUint8 || w < 0 {
		panic(fmt.Errorf("invalid width: %d", w))
	} else if h > math.MaxUint8 || h < 0 {
		panic(fmt.Errorf("invalid height: %d", h))
	}
	f = &Frame{w, h, p, make([]byte, w*h)}
	for y, row := range rows {
		if len(row) != w {
			panic(fmt.Errorf("at row %d: %d bytes doesn't match declared width %d", y, len(row), w))
		}
		for x := range w {
			cell := uint8FromHex(row[x])
			if cell > 0x0F {
				panic(fmt.Errorf("at row %d: invalid reserved cell value %d", y, cell))
			}
			f.grid[y*w+(x%w)] = cell
		}
	}
	return f
}

func (f *Frame) Image() (img *image.RGBA) {
	img = image.NewRGBA(image.Rect(0, 0, f.w, f.h))
	for i, cell := range f.grid {
		img.Set(i%f.w, i/f.h, Uint32RGBA(f.palette[cell]))
	}
	return img
}

func (f *Frame) EncodePNG(w io.Writer, scale int) (err error) {
	return png.Encode(w, ScaleBy(scale, f.Image()))
}
