package lopix

import (
	"bytes"
	"image"
	"image/draw"
	"image/gif"
	"io"
)

type Animation struct {
	Delays []int // In 100ths of a second.
	Frames []*Frame
}

func Animate(w, h int, palette [16]Color, frames ...[]string) (v *Animation) {
	v = &Animation{}
	for range frames {
		v.Delays = append(v.Delays, 30)
	}
	for _, frame := range frames {
		v.Frames = append(v.Frames, F(w, h, palette, frame...))
	}
	return v
}

func (a *Animation) EncodeGIF(w io.Writer, scale int) (err error) {
	gifImages := []*image.Paletted{}
	for _, frame := range a.Frames {
		b := bytes.Buffer{}
		// TODO?: Find a more elegant solution to get an image.Paletted from a Frame
		// without having to gif.Encode/Decode.
		err = gif.Encode(&b, ScaleBy(scale, frame.Image()), &gif.Options{Drawer: draw.Over})
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
	return gif.EncodeAll(w, &gif.GIF{Image: gifImages, Delay: a.Delays})
}
