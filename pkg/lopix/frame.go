package lopix

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"math"
	"strconv"
	"strings"
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

func F(w, h int, p Palette, rows ...string) (f *Frame) {
	if w > math.MaxUint8 || w < 0 {
		panic(fmt.Errorf("invalid width: %d", w))
	} else if h > math.MaxUint8 || h < 0 {
		panic(fmt.Errorf("invalid height: %d", h))
	}
	grid, err := ParseGrid(rows, w, h)
	if err != nil {
		panic(fmt.Errorf("parse grid: %w", err))
	}
	return &Frame{w, h, p, grid}
}

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

func (f *Frame) ReadFrom(r io.Reader) (n int64, err error) {
	// Limit reads to maximum possible grid size (to prevent malicious large allocation attempts).
	r = io.LimitReader(r, int64(MaxEncodedFrameSize))

	// Read dimensions line (where width and height are declared).
	bufr := bufio.NewReader(r)
	lineI := 0 // Line index.
	line, err := bufr.ReadString('\n')
	n += int64(len(line))
	if err != nil {
		return n, fmt.Errorf("read first line: %w", err)
	}
	lineI++
	line = strings.TrimSpace(line)
	parts := strings.Split(line, "x")
	if len(parts) != 2 {
		return n, fmt.Errorf("invalid first line: %q", line)
	}
	f.w, err = strconv.Atoi(parts[0])
	if err != nil {
		return n, fmt.Errorf("invalid width: %w", err)
	} else if f.w <= 0 || f.w > MaxWidth {
		return n, fmt.Errorf("invalid width: %d", f.w)
	}
	f.h, err = strconv.Atoi(parts[1])
	if err != nil {
		return n, fmt.Errorf("invalid height: %w", err)
	} else if f.h <= 0 || f.h > MaxHeight {
		return n, fmt.Errorf("invalid height: %d", f.h)
	}

	// Read palette of colors (NB: 16 colors maximum).
	f.palette = Palette{}
	for i := range 16 {
		line, err = bufr.ReadString('\n')
		n += int64(len(line))
		if err != nil {
			return n, fmt.Errorf("read palette (at line %d): %w", lineI, err)
		}
		lineI++
		line = strings.TrimSpace(line)
		if line == "" {
			break // Reached end of palette (= blank line).
		}
		f.palette[i], err = HexColor(line)
		if err != nil {
			return n, fmt.Errorf("parse palette color (at line %d): %w", lineI, err)
		}
	}

	// Read grid (row by row).
	f.grid = make([]byte, f.w*f.h)
	for y := range f.h {
		// Read one row.
		line, err = bufr.ReadString('\n')
		n += int64(len(line))
		if err != nil {
			return n, fmt.Errorf("read grid row (at line %d): %w", lineI, err)
		}
		lineI++
		line = strings.TrimSpace(line)

		// Ensure that we have one character for each cell of the expected width.
		if line == "" {
			return n, fmt.Errorf("at line %d: unexpected blank line in grid: %w", lineI, err)
		} else if len(line) != f.w {
			return n, fmt.Errorf("at line %d: len(row) == %d bytes but width == %d", lineI, len(line), f.w)
		}

		// Parse row.
		for x := range f.w {
			cell := hexToU4(line[x])
			if cell > 0x0F {
				return n, fmt.Errorf("at line %d: invalid reserved cell value %d", lineI, cell)
			}
			f.grid[y*f.w+(x%f.w)] = cell
		}
	}

	return n, nil
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

func ParseGrid(rows []string, w, h int) (grid []byte, err error) {
	grid = make([]byte, w*h)
	for y, row := range rows {
		if len(row) != w {
			return nil, fmt.Errorf("at row %d: %d bytes doesn't match declared width %d", y, len(row), w)
		}
		for x := range w {
			cell := hexToU4(row[x])
			if cell > 0x0F {
				return nil, fmt.Errorf("at row %d: invalid reserved cell value %d", y, cell)
			}
			grid[y*w+(x%w)] = cell
		}
	}
	return grid, nil
}
