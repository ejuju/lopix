package main

import (
	"bytes"
	_ "embed"
	"os"

	"github.com/ejuju/lopix/pkg/lopix"
)

//go:embed demo.lopix
var src []byte

func main() {
	animation := &lopix.Animation{}
	_, err := animation.ReadFrom(bytes.NewReader(src))
	if err != nil {
		panic(err)
	}
	err = animation.EncodeGIF(os.Stdout, 400/animation.W())
	if err != nil {
		panic(err)
	}
}
