package lopix

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"io"
	"strconv"
)

const (
	MaxAnimationFrames      = 99
	MaxFrameDelay           = 99
	MaxEncodedAnimationSize = len("99*99\n") + len("99x99\n") + (16 * len("#aabbccff\n")) + 1 + (MaxAnimationFrames * (EncodedGridMaxSize + 1))
)

type Animation struct {
	delays  []int // In 100ths of a second.
	w       int
	h       int
	palette Palette
	grids   [][]byte
}

func (a *Animation) W() int { return a.w }
func (a *Animation) H() int { return a.h }

func (a *Animation) EncodeGIF(w io.Writer, scale int) (err error) {
	if len(a.delays) == 0 {
		for range a.grids {
			a.delays = append(a.delays, 30)
		}
	}
	gifImages := []*image.Paletted{}
	for _, grid := range a.grids {
		b := bytes.Buffer{}
		// TODO?: Find a more elegant solution to get an image.Paletted from a Frame
		// without having to gif.Encode/Decode.
		err = gif.Encode(&b, ScaleBy(scale, (&Frame{a.w, a.h, a.palette, grid}).Image()), &gif.Options{Drawer: draw.Over})
		if err != nil {
			return err
		}
		decodedImg, err := gif.Decode(&b)
		if err != nil {
			return err
		}
		palettedImg, ok := decodedImg.(*image.Paletted)
		if !ok {
			panic("unreachable")
		}
		gifImages = append(gifImages, palettedImg)
	}
	return gif.EncodeAll(w, &gif.GIF{Image: gifImages, Delay: a.delays})
}

func (a *Animation) ParseFrom(r io.Reader) (err error) {
	r = io.LimitReader(r, int64(MaxEncodedAnimationSize))
	p := NewParser(r)

	numFrames, frameDelay, err := p.ParseAnimationInfo()
	if err != nil {
		return fmt.Errorf("parse animation info: %w", err)
	}
	for range numFrames {
		a.delays = append(a.delays, frameDelay)
	}

	a.w, a.h, err = p.ParseDimensions()
	if err != nil {
		return fmt.Errorf("parse dimensions: %w", err)
	}

	a.palette, err = p.ParsePalette()
	if err != nil {
		return fmt.Errorf("parse palette: %w", err)
	}

	for frameI := range numFrames {
		grid, err := p.ParseGrid(a.w, a.h)
		if err != nil {
			return fmt.Errorf("parse frame %d: %w", frameI, err)
		}

		if frameI != numFrames-1 {
			err := p.ReadBlankLine()
			if err != nil {
				return fmt.Errorf("expect separator between animation grids: %w", err)
			}
		}

		a.grids = append(a.grids, grid)
	}

	return nil
}

// NB: In the current implementation, we encode the whole animation in memory before being written.
func (a *Animation) WriteTo(w io.Writer) (n int64, err error) {
	b := &bytes.Buffer{}

	// Write dimensions line (ex: "100x150\n").
	b.WriteString(strconv.Itoa(a.w))
	b.WriteString("x")
	b.WriteString(strconv.Itoa(a.h))
	b.WriteString("\n")

	// Write palette.
	for _, color := range a.palette {
		b.WriteString("#")
		b.WriteString(u32ToHex(color))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Write grid.
	for i, grid := range a.grids {
		if i > 0 {
			b.WriteString("\n")
		}
		for y := range a.h {
			for x := range a.w {
				b.WriteString(string(u4ToHex(grid[y*a.w+(x%a.w)])))
			}
			b.WriteString("\n")
		}
	}

	return b.WriteTo(w)
}
