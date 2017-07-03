package main

import (
	"fmt"
	"log"

	"github.com/docktermj/mock-proxy/common/runner"
	"github.com/docktermj/mock-proxy/subcommand/net"
	"github.com/docopt/docopt-go"
)

// Values updated via "go install -ldflags" parameters.

var programName string = "unknown"
var buildVersion string = "0.0.0"
var buildIteration string = "0"

// TODO: Add logging.

func main() {
	usage := `
Usage:
    mock-proxy [--version] [--help] <command> [<args>...]

Options:
    -h, --help

The mock-proxy commands are:
    net    Relay through different types of networks

See 'mock-proxy <command> --help' for more information on a specific command.
`
	// DocOpt processing.

	commandVersion := fmt.Sprintf("%s %s-%s", programName, buildVersion, buildIteration)
	args, _ := docopt.Parse(usage, nil, true, commandVersion, true)

	// Configure output log.

	log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds | log.LUTC)

	// Construct 'argv'.

	argv := make([]string, 1)
	argv[0] = args["<command>"].(string)
	argv = append(argv, args["<args>"].([]string)...)

	// Reference: http://stackoverflow.com/questions/6769020/go-map-of-functions

	functions := map[string]interface{}{
		"net": net.Command,
	}

	runner.Run(argv, functions, usage)
}
