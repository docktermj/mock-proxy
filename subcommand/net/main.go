package net

// Inspirations:
//  - https://gist.github.com/hakobe/6f70d69b8c5243117787fd488ae7fbf2

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/docktermj/mock-server/common/help"
	"github.com/docopt/docopt-go"
)

func proxy(inbound net.Conn, outbound net.Conn) {
	byteBuffer := make([]byte, 1024)

	for {

		// Read the inbound network connection.

		numberOfBytesRead, err := inbound.Read(byteBuffer)
		if err != nil {
			log.Println("Proxy Read return")
			return
		}

		// Print the message handled.

		message := byteBuffer[0:numberOfBytesRead]
		fmt.Println(">>>", string(message))

		// Write to outbound network connection.

		_, err = outbound.Write([]byte(message))
		if err != nil {
			log.Println("Proxy Write return")
			return
		}
	}
}

// Function for the "command pattern".
func Command(argv []string) {

	usage := `
Usage:
    mock-proxy net [options]

Options:
   -h, --help
   --inbound-network=<network_type>   Type of network used for communication
   --inbound-address=<address>        Address for network_type.
   --outbound-network=<network_type>  Type of network used for communication
   --outbound-address=<address>       Address for network_type.   
   --debug                            Log debugging messages

Where:
   network_type   Examples: 'unix', 'tcp'
   address        Examples: '/tmp/test.sock', '127.0.0.1:12345'
`

	// DocOpt processing.

	args, _ := docopt.Parse(usage, nil, true, "", false)

	// Test for required commandline options.

	message := ""

	if args["--inbound-network"] == nil {
		message += "Missing '--inbound-network' parameter;"
	}

	if args["--inbound-address"] == nil {
		message += "Missing '--inbound-address' parameter;"
	}

	if args["--outbound-network"] == nil {
		message += "Missing '--outbound-network' parameter;"
	}

	if args["--outbound-address"] == nil {
		message += "Missing '--outbound-address' parameter;"
	}

	if len(message) > 0 {
		help.ShowHelp(usage)
		fmt.Println(strings.Replace(message, ";", "\n", -1))
		log.Fatalln(strings.Replace(message, ";", "; ", -1))
	}

	// Get commandline options.

	inboundNetwork := args["--inbound-network"].(string)
	inboundAddress := args["--inbound-address"].(string)
	outboundNetwork := args["--outbound-network"].(string)
	outboundAddress := args["--outbound-address"].(string)
	isDebug := args["--debug"].(bool)

	// Debugging information.

	if isDebug {
		log.Printf("Listening on '%s' network with address '%s'", inboundNetwork, inboundAddress)
		log.Printf("Sending   to '%s' network with address '%s'", outboundNetwork, outboundAddress)
	}

	// Inbound listener.  net.Listen creates a server.

	inboundListener, err := net.Listen(inboundNetwork, inboundAddress)
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	// Configure listener to exit when program ends.

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(listener net.Listener, c chan os.Signal) {
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)
		listener.Close()
		os.Exit(0)
	}(inboundListener, sigc)

	// As a server, Read and Echo loop.

	for {

		// As a server, listen for a connection request.

		inboundConnection, err := inboundListener.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}
		if isDebug {
			log.Println("Accepted inbound connection.")
		}

		// Create a new outbound network connection.  net.Dial creates a client.

		outboundConnection, err := net.Dial(outboundNetwork, outboundAddress)
		if err != nil {
			log.Fatal("Dial error", err)
		}
		defer outboundConnection.Close()

		// Asynchronously handle bi-directional traffic.

		go proxy(inboundConnection, outboundConnection)
		go proxy(outboundConnection, inboundConnection)

	}
}
