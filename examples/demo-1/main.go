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
	err := animation.ParseFrom(bytes.NewReader(src))
	if err != nil {
		panic(err)
	}
	err = animation.EncodeGIF(os.Stdout, 320/16)
	if err != nil {
		panic(err)
	}
}
