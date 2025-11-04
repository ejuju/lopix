package lopix

import (
	"fmt"
	"image/color"
	"strings"
)

type Color = uint32
type Palette = [16]Color

// Parses a hexadecimal RGB(A) color (ex: "#aabbccff" or "#aabbcc").
func HexColor(v string) (c Color, err error) {
	v = strings.TrimPrefix(v, "#")
	if !(len(v) == 6 || len(v) == 8) {
		return 0, fmt.Errorf("invalid hexadecimal RGB(A) color: %q", v)
	}
	c |= uint32(hexToU8(v[0:2])) << 24
	c |= uint32(hexToU8(v[2:4])) << 16
	c |= uint32(hexToU8(v[4:6])) << 8
	if len(v) == 8 {
		c |= uint32(hexToU8(v[6:8])) << 0
	} else {
		c |= 0xff
	}
	return c, nil
}

// Like HexColor, but panics on invalid input.
func C(v string) (c Color) {
	c, err := HexColor(v)
	if err != nil {
		panic(err)
	}
	return c
}

func ColorRGBA(rgba Color) color.RGBA {
	return color.RGBA{
		byte((rgba >> 24) & 0xFF),
		byte((rgba >> 16) & 0xFF),
		byte((rgba >> 8) & 0xFF),
		byte((rgba >> 0) & 0xFF),
	}
}
