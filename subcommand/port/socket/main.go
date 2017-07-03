package socket

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
	"strings"
	"strconv"

	"github.com/docktermj/mock-proxy/common/help"
	"github.com/docopt/docopt-go"
)

// Read a message from the network and respond.
func reader(reader io.Reader) {
	byteBuffer := make([]byte, 1024)
	for {
		numberOfBytesRead, err := reader.Read(byteBuffer[:])
		if err != nil {
			return
		}
		fmt.Println("<<<", string(byteBuffer[0:numberOfBytesRead]))
	}
}

// Function for the "command pattern".
func Command(argv []string) {

	usage := `
Usage:
    mock-proxy port socket [options] 

Options:
   -h, --help
   --port=<port_number>  Port to listen on 
   --socket-file=<file>  Socket file of server
   --debug               Log debugging messages
`

	// DocOpt processing.

	args, _ := docopt.Parse(usage, nil, true, "", false)

	// Test for required commandline options.

	message := ""

	if args["--port"] == nil {
		message += "Missing '--port' parameter;"
	}

	if args["--socket-file"] == nil {
		message += "Missing '--socket-file' parameter;"

	}

	if len(message) > 0 {
		help.ShowHelp(usage)
		fmt.Println(strings.Replace(message, ";", "\n", -1))
		log.Fatalln(strings.Replace(message, ";", "; ", -1))
	}

	// Get commandline options.

	port, _ := strconv.Atoi(args["--port"].(string))
	socketFile := args["--socket-file"].(string)
	isDebug := args["--debug"].(bool)

	// Listen on the Unix Domain Socket

	if isDebug {
	    fmt.Printf("Port: %d;  Socket: %s", port, socketFile)
	}

	networkConnection, err := net.Dial("unix", socketFile)
	if err != nil {
		log.Fatal("Dial error", err)
	}
	defer networkConnection.Close()

	// Start asynchronous Reader

	go reader(networkConnection)

	// Loop through Writer

	loopNumber := 1
	for {
		loopNumber += 1
		outboundMessage := fmt.Sprintf("Sending #%d", loopNumber)
		_, err := networkConnection.Write([]byte(outboundMessage))
		if err != nil {
			log.Fatal("Write error:", err)
			break
		}
		fmt.Println(">>>", outboundMessage)
		time.Sleep(1e9)
	}
}
