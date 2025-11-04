package lopix

import (
	"fmt"
	"image"
)

// NB: This only works for square images for now...
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

func hexToU4(v byte) uint8 {
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

func hexToU8(v string) uint8 {
	return hexToU4(v[0])<<4 | hexToU4(v[1])
}

func u4ToHex(v uint8) byte {
	v &= 0x0F // Ensure we are handling a uint4.
	if v > 9 {
		return 'a' + v
	}
	return '0' + v
}

func u32ToHex(v uint32) (txt string) {
	txt += string(u4ToHex(uint8((v >> 28) & 0xF)))
	txt += string(u4ToHex(uint8((v >> 24) & 0xF)))
	txt += string(u4ToHex(uint8((v >> 20) & 0xF)))
	txt += string(u4ToHex(uint8((v >> 16) & 0xF)))
	txt += string(u4ToHex(uint8((v >> 12) & 0xF)))
	txt += string(u4ToHex(uint8((v >> 8) & 0xF)))
	txt += string(u4ToHex(uint8((v >> 4) & 0xF)))
	txt += string(u4ToHex(uint8((v >> 0) & 0xF)))
	return txt
}
