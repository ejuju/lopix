package main

import (
	"os"

	"github.com/ejuju/lopix/pkg/lopix"
)

func main() {
	os.Exit(lopix.Run(os.Args...))
}
