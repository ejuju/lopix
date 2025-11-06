package lopix

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/ejuju/lopix/pkg/cli"
)

func Run(args ...string) (exitcode int) {
	return cli.Run(commands, nil, nil, args...)
}

func init() {
	commands = append([]*cli.Command{
		{
			Keyword:     "help",
			Description: "Prints available commands",
			Do:          func(args ...string) (exitcode int) { cli.PrintAvailableCommands(commands); return 0 },
		},
	}, commands...)
}

var commands = []*cli.Command{
	{
		Keyword:     "png",
		Description: "Generates a PNG from a Lopix file",
		Do:          runParseAndEncode,
	},
	{
		Keyword:     "gif",
		Description: "Generates a GIF from a Lopix file",
		Do:          runParseAndEncode,
	},
}

func runParseAndEncode(args ...string) (exitcode int) {
	if len(args) <= 1 {
		fmt.Println("missing arguments: {format} {src} {out}")
		return 1
	} else if len(args) <= 2 {
		fmt.Println("missing arguments: {src} {out}")
		return 1
	} else if len(args) <= 3 {
		fmt.Println("missing arguments: {out}")
		return 1
	}
	format := args[1]
	fpathIn := args[2]
	fpathOut := args[3]
	scaleTxt := "1"
	if len(args) >= 5 {
		scaleTxt = args[4]
	}

	// Open input file.
	inputFile, err := os.Open(fpathIn)
	if err != nil {
		fmt.Printf("open input file (%q): %s\n", fpathIn, err.Error())
		return 1
	}
	defer inputFile.Close()

	// Open output file.
	outputFile, err := os.OpenFile(fpathOut, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("open output file (%q): %s\n", fpathOut, err.Error())
		return 1
	}
	defer outputFile.Close()

	// Parse scale.
	scale, err := strconv.Atoi(scaleTxt)
	if err != nil {
		fmt.Printf("parse scale (%q): %s\n", scaleTxt, err)
		return 1
	} else if scale <= 0 && scale > 2048 {
		fmt.Printf("invalid scale: %d\n", scale)
		return 1
	}

	// Parse and encode frame/animation.
	var parseFrom func(io.Reader) error
	var encodeTo func(io.Writer, int) error
	switch format {
	default:
		fmt.Printf("unknown format: %q\n", format)
		return 1
	case "png":
		v := &Frame{}
		parseFrom = v.ParseFrom
		encodeTo = v.EncodePNG
	case "gif":
		v := &Animation{}
		parseFrom = v.ParseFrom
		encodeTo = v.EncodeGIF
	}
	err = parseFrom(inputFile)
	if err != nil {
		fmt.Printf("parse lopix: %s\n", err)
		return 1
	}
	err = encodeTo(outputFile, scale)
	if err != nil {
		fmt.Printf("encode PNG: %s\n", err)
		return 1
	}

	return 0
}
