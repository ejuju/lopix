package cli

import (
	"fmt"
)

type HandlerFunc func(args ...string) (exitcode int)

type Command struct {
	Keyword     string
	Description string
	Do          HandlerFunc
}

func Run(commands []*Command, fallback, handleNotFound HandlerFunc, args ...string) (exitcode int) {
	if fallback == nil {
		fallback = CommandRequiredHandler(commands)
	}
	if handleNotFound == nil {
		handleNotFound = NotFoundHandler(commands)
	}

	if len(args) < 2 {
		return fallback(args...)
	}
	for _, cmd := range commands {
		if cmd.Keyword == args[1] {
			return cmd.Do(args...)
		}
	}
	return handleNotFound(args...)
}

func PrintAvailableCommands(commands []*Command) {
	fmt.Println("Available commands:")
	for _, cmd := range commands {
		fmt.Printf("%q: %s\n", cmd.Keyword, cmd.Description)
	}
}

func NotFoundHandler(commands []*Command) HandlerFunc {
	return func(args ...string) (exitcode int) {
		fmt.Printf("Command not found: %q\n\n", args[1])
		PrintAvailableCommands(commands)
		return 1
	}
}

func CommandRequiredHandler(commands []*Command) HandlerFunc {
	return func(args ...string) (exitcode int) {
		fmt.Printf("A command is required!\n\n")
		PrintAvailableCommands(commands)
		return 1
	}
}
