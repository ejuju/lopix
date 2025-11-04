package lopix

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"strconv"
)

const (
	MaxWidth            = 99
	MaxHeight           = 99
	PaletteSize         = 16
	EncodedGridMaxSize  = MaxHeight * (MaxWidth + 1)
	MaxEncodedFrameSize = len("99x99\n") + (16 * len("#aabbccff\n")) + 1 + EncodedGridMaxSize
)

type Frame struct {
	w       int
	h       int
	palette Palette
	grid    []byte
}

func (a *Frame) W() int { return a.w }
func (a *Frame) H() int { return a.h }

func (a *Frame) Image() (img *image.RGBA) {
	img = image.NewRGBA(image.Rect(0, 0, a.w, a.h))
	for i, cell := range a.grid {
		img.Set(i%a.w, i/a.h, ColorRGBA(a.palette[cell]))
	}
	return img
}

func (a *Frame) EncodePNG(w io.Writer, scale int) (err error) {
	return png.Encode(w, ScaleBy(scale, a.Image()))
}

func (f *Frame) ParseFrom(r io.Reader) (err error) {
	r = io.LimitReader(r, int64(MaxEncodedFrameSize))
	p := NewParser(r)

	f.w, f.h, err = p.ParseDimensions()
	if err != nil {
		return fmt.Errorf("parse dimensions: %w", err)
	}

	f.palette, err = p.ParsePalette()
	if err != nil {
		return fmt.Errorf("parse palette: %w", err)
	}

	f.grid, err = p.ParseGrid(f.w, f.h)
	if err != nil {
		return fmt.Errorf("parse grid: %w", err)
	}

	return nil
}

// NB: In the current implementation, the frame is encoded in memory before being written.
func (f *Frame) WriteTo(w io.Writer) (n int64, err error) {
	b := &bytes.Buffer{}

	// Write dimensions line (ex: "100x150\n").
	b.WriteString(strconv.Itoa(f.w))
	b.WriteString("x")
	b.WriteString(strconv.Itoa(f.h))
	b.WriteString("\n")

	// Write palette.
	for _, color := range f.palette {
		b.WriteString("#")
		b.WriteString(u32ToHex(color))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Write grid.
	for y := range f.h {
		for x := range f.w {
			b.WriteString(string(u4ToHex(f.grid[y*f.w+(x%f.w)])))
		}
		b.WriteString("\n")
	}

	return b.WriteTo(w)
}
