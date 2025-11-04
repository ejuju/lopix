package main

import (
	"os"

	"github.com/ejuju/lopix/pkg/lopix"
)

func main() {
	// Define dimensions.
	const width, height = 16, 16

	// Define color palette.
	palette := lopix.Palette{
		lopix.HexColor("#d6d6d6"),
		lopix.HexColor("#ff4000"),
		lopix.HexColor("#242424"),
	}

	// Define pixel grid.
	frame := lopix.F(width, height, palette,
		"0000000000000000",
		"0000000220000000",
		"0000000200000000",
		"0011111111111100",
		"0011111111111100",
		"0011111111111100",
		"0011121111211100",
		"0011121111211100",
		"0011122112211100",
		"0011111111111100",
		"0011212121121100",
		"0011222222221100",
		"0011112121211100",
		"0011111111111100",
		"0000000000000000",
		"0000000000000000",
	)

	// Encode PNG (scaling from 16x16 to 400x400 pixels).
	err := frame.EncodePNG(os.Stdout, 400/width)
	if err != nil {
		panic(err)
	}
}
