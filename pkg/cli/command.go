package cli

type HandlerFunc func(args ...string) (exitcode int)

type Command struct {
	Keyword     string
	Description string
	Do          HandlerFunc
}

func Run(commands []*Command, handleDefault, handleNotFound HandlerFunc, args ...string) (exitcode int) {
	if len(args) < 2 {
		return handleDefault(args...)
	}
	for _, cmd := range commands {
		if cmd.Keyword == args[1] {
			return cmd.Do(args...)
		}
	}
	return handleNotFound(args...)
}
