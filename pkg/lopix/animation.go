package lopix

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"io"
	"strconv"
	"strings"
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

func (a *Animation) ReadFrom(r io.Reader) (n int64, err error) {
	// Limit reads to maximum possible size (to prevent malicious large allocation attempts).
	r = io.LimitReader(r, int64(MaxEncodedAnimationSize))

	// Read first line (where number of frames and delays are declared).
	bufr := bufio.NewReader(r)
	lineN := 1 // Line number (index+1).
	line, err := bufr.ReadString('\n')
	n += int64(len(line))
	if err != nil {
		return n, fmt.Errorf("read first line: %w", err)
	}
	lineN++
	line = strings.TrimSpace(line)
	parts := strings.Split(line, "*")
	if len(parts) != 2 {
		return n, fmt.Errorf("invalid first line: %q", line)
	}
	numFrames, err := strconv.Atoi(parts[0])
	if err != nil {
		return n, fmt.Errorf("invalid number of frames: %w", err)
	} else if numFrames <= 0 || numFrames > MaxAnimationFrames {
		return n, fmt.Errorf("invalid number of frames: %d", numFrames)
	}
	frameDelay, err := strconv.Atoi(parts[1])
	if err != nil {
		return n, fmt.Errorf("invalid frame delay: %w", err)
	} else if frameDelay <= 0 || frameDelay > MaxFrameDelay {
		return n, fmt.Errorf("invalid frame delay: %d", frameDelay)
	}
	for range numFrames {
		a.delays = append(a.delays, frameDelay)
	}

	// Read dimensions line (where width and height are declared).
	line, err = bufr.ReadString('\n')
	n += int64(len(line))
	if err != nil {
		return n, fmt.Errorf("read first line: %w", err)
	}
	lineN++
	line = strings.TrimSpace(line)
	parts = strings.Split(line, "x")
	if len(parts) != 2 {
		return n, fmt.Errorf("invalid dimensions line: %q", line)
	}
	a.w, err = strconv.Atoi(parts[0])
	if err != nil {
		return n, fmt.Errorf("invalid width: %w", err)
	} else if a.w <= 0 || a.w > MaxWidth {
		return n, fmt.Errorf("invalid width: %d", a.w)
	}
	a.h, err = strconv.Atoi(parts[1])
	if err != nil {
		return n, fmt.Errorf("invalid height: %w", err)
	} else if a.h <= 0 || a.h > MaxHeight {
		return n, fmt.Errorf("invalid height: %d", a.h)
	}

	// Read palette of colors (NB: 16 colors maximum).
	a.palette = Palette{}
	for i := range 16 {
		line, err = bufr.ReadString('\n')
		n += int64(len(line))
		if err != nil {
			return n, fmt.Errorf("read palette (at line %d): %w", lineN, err)
		}
		lineN++
		line = strings.TrimSpace(line)
		if line == "" {
			break // Reached end of palette (= blank line).
		}
		a.palette[i], err = HexColor(line)
		if err != nil {
			return n, fmt.Errorf("parse palette color (at line %d): %w", lineN, err)
		}
	}

	// Read grids (grid by grid and row by row, with blank lines as seperators).
	for frameI := range numFrames {
		grid := make([]byte, a.w*a.h)
		for y := range a.h {
			// Read one row.
			line, err = bufr.ReadString('\n')
			n += int64(len(line))
			if err != nil {
				return n, fmt.Errorf("read grid row (at line %d): %w", lineN, err)
			}
			lineN++
			line = strings.TrimSpace(line)

			// Ensure that we have one character for each cell of the expected width.
			if line == "" {
				return n, fmt.Errorf("at line %d: unexpected blank line in grid: %w", lineN, err)
			} else if len(line) != a.w {
				return n, fmt.Errorf("at line %d: len(row) == %d bytes but width == %d", lineN, len(line), a.w)
			}

			// Parse row.
			for x := range a.w {
				cell := hexToU4(line[x])
				if cell > 0x0F {
					return n, fmt.Errorf("at line %d: invalid reserved cell value %d", lineN, cell)
				}
				grid[y*a.w+(x%a.w)] = cell
			}
		}

		// Consume separator (blank) line.
		if frameI != numFrames-1 {
			line, err = bufr.ReadString('\n')
			n += int64(len(line))
			if err != nil {
				return n, fmt.Errorf("read grid row (at line %d): %w", lineN, err)
			} else if line != "\n" {
				return n, fmt.Errorf("unexpected non-blank line separator (at line %d): %q", lineN, line)
			}
			lineN++
		}

		a.grids = append(a.grids, grid)
	}

	return n, nil
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
