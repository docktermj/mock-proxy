package port

import (
	"github.com/docktermj/mock-proxy/common/runner"
	"github.com/docktermj/mock-proxy/subcommand/port/socket"
)

func Command(argv []string) {

	usage := `
Usage:
    mock-proxy port <subcommand> [<args>...]

Subcommands:
    socket    Relay to a socket
`

	functions := map[string]interface{}{
		"socket": socket.Command,
	}

	runner.Run(argv, functions, usage)
}
