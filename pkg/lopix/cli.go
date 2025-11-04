package lopix

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func Run(args ...string) (exitcode int) {
	if len(args) <= 4 {
		log.Println("missing argument(s)")
		return 1
	}
	format := args[1]
	fpathIn := args[2]
	fpathOut := args[3]
	scaleTxt := args[4]

	// Open input file.
	fIn, err := os.Open(fpathIn)
	if err != nil {
		log.Printf("open input file (%q): %s", fpathIn, err.Error())
		return 1
	}
	defer fIn.Close()

	// Open output file.
	fOut, err := os.OpenFile(fpathOut, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Printf("open output file (%q): %s", fpathOut, err.Error())
		return 1
	}
	defer fOut.Close()

	// Parse scale.
	scale, err := strconv.Atoi(scaleTxt)
	if err != nil {
		log.Printf("parse scale (%q): %s", scaleTxt, err)
		return 1
	} else if scale <= 0 && scale > 2048 {
		log.Printf("invalid scale: %d", scale)
		return 1
	}

	// Parse and encode frame/animation.
	var encode func(w io.Writer, scale int) error
	switch format {
	default:
		log.Printf("invalid format: %q", format)
		return 1
	case "", "png":
		v := &Frame{}
		_, err = v.ReadFrom(fIn)
		if err != nil {
			log.Printf("parse frame: %s", err)
			return 1
		}
		encode = v.EncodePNG

	case "gif":
		v := &Animation{}
		_, err = v.ReadFrom(fIn)
		if err != nil {
			log.Printf("parse animation: %s", err)
			return 1
		}
		encode = v.EncodeGIF
	}
	err = encode(fOut, scale)
	if err != nil {
		log.Printf("encode %s: %s", strings.ToUpper(format), err)
		return 1
	}

	return 0
}
