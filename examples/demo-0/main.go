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
	err := frame.ParseFrom(bytes.NewReader(src))
	if err != nil {
		panic(err)
	}
	err = frame.EncodePNG(os.Stdout, 320/16)
	if err != nil {
		panic(err)
	}
}
