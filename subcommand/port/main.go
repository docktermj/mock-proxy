package port

import (
	"github.com/docktermj/mock-relay/common/runner"
	"github.com/docktermj/mock-relay/subcommand/port/socket"
)

func Command(argv []string) {

	usage := `
Usage:
    mock-relay port <subcommand> [<args>...]

Subcommands:
    socket    Relay to a socket
`

	functions := map[string]interface{}{
		"socket": socket.Command,
	}

	runner.Run(argv, functions, usage)
}
