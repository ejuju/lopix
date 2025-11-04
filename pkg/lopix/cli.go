package lopix

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/ejuju/lopix/pkg/cli"
)

func Run(args ...string) (exitcode int) {
	return cli.Run(commands, onNoCommand, onCLICommandNotFound, args...)
}

func printAvailableCommands() {
	fmt.Println("Available commands:")
	for _, cmd := range commands {
		fmt.Printf("%q: %s\n", cmd.Keyword, cmd.Description)
	}
}

func onCLICommandNotFound(args ...string) (exitcode int) {
	fmt.Printf("Command not found: %q\n\n", args[1])
	printAvailableCommands()
	return 1
}

func onNoCommand(args ...string) (exitcode int) {
	fmt.Printf("A command is required!\n\n")
	printAvailableCommands()
	return 1
}

func init() {
	commands = append([]*cli.Command{
		{
			Keyword:     "help",
			Description: "Prints available commands",
			Do:          func(args ...string) (exitcode int) { printAvailableCommands(); return 0 },
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
	if len(args) <= 4 {
		fmt.Println("missing argument(s)")
		return 1
	}
	format := args[1]
	fpathIn := args[2]
	fpathOut := args[3]
	scaleTxt := args[4]

	// Open input file.
	fIn, err := os.Open(fpathIn)
	if err != nil {
		fmt.Printf("open input file (%q): %s\n", fpathIn, err.Error())
		return 1
	}
	defer fIn.Close()

	// Open output file.
	fOut, err := os.OpenFile(fpathOut, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("open output file (%q): %s\n", fpathOut, err.Error())
		return 1
	}
	defer fOut.Close()

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
	var encode func(io.Writer, int) error
	switch format {
	default:
		fmt.Printf("unknown format: %q\n", format)
		return 1
	case "png":
		v := &Frame{}
		err = v.ParseFrom(fIn)
		if err != nil {
			fmt.Printf("parse frame: %s\n", err)
			return 1
		}
		encode = v.EncodePNG
	case "gif":
		v := &Animation{}
		err = v.ParseFrom(fIn)
		if err != nil {
			fmt.Printf("parse frame: %s\n", err)
			return 1
		}
		encode = v.EncodeGIF
	}
	err = encode(fOut, scale)
	if err != nil {
		fmt.Printf("encode PNG: %s\n", err)
		return 1
	}

	return 0
}
