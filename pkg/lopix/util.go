package lopix

import (
	"fmt"
	"image"
	"image/color"
	"strings"
)

type Color = uint32
type Palette = [16]Color

func Uint32RGBA(rgba Color) color.RGBA {
	return color.RGBA{
		byte((rgba >> 24) & 0xFF),
		byte((rgba >> 16) & 0xFF),
		byte((rgba >> 8) & 0xFF),
		byte((rgba >> 0) & 0xFF),
	}
}

func HexColor(v string) (c Color) {
	v = strings.TrimPrefix(v, "#")
	if !(len(v) == 6 || len(v) == 8) {
		panic("invalid hexadecimal RGB(A) color")
	}
	c |= uint32(uint8FromHex2(v[0:2])) << 24
	c |= uint32(uint8FromHex2(v[2:4])) << 16
	c |= uint32(uint8FromHex2(v[4:6])) << 8
	if len(v) == 8 {
		c |= uint32(uint8FromHex2(v[6:8])) << 0
	} else {
		c |= 0xff
	}
	return c
}

func ScaleBy(factor int, img image.Image) image.Image {
	scaled := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx()*factor, img.Bounds().Dy()*factor))

	for y := range img.Bounds().Dy() {
		for x := range img.Bounds().Dx() {
			for y2 := range factor {
				for x2 := range factor {
					scaled.Set((x*factor)+x2, (y*factor)+y2, img.At(x, y))
				}
			}
		}
	}

	return scaled
}

func uint8FromHex(v byte) uint8 {
	switch {
	default:
		panic(fmt.Errorf("invalid hexadecimal character: %q", v))
	case v >= '0' && v <= '9':
		return v - '0'
	case v >= 'a' && v <= 'f':
		return v - 'a' + 10
	case v >= 'A' && v <= 'F':
		return v - 'A' + 10
	}
}

func uint8FromHex2(v string) uint8 {
	return uint8FromHex(v[0])<<4 | uint8FromHex(v[1])
}
