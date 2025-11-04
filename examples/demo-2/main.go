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
	frame := &lopix.Frame{}
	_, err := frame.ReadFrom(bytes.NewReader(src))
	if err != nil {
		panic(err)
	}
	err = frame.EncodePNG(os.Stdout, 400/frame.W())
	if err != nil {
		panic(err)
	}
}
