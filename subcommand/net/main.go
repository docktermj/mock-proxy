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

// Read a message from the network and respond.
func outboundTraffic(reader io.Reader) {
	byteBuffer := make([]byte, 1024)
	for {
		numberOfBytesRead, err := reader.Read(byteBuffer[:])
		if err != nil {
			return
		}
		fmt.Println("<<<", string(byteBuffer[0:numberOfBytesRead]))
	}
}


// Read a message from the network and respond.
func inboundTraffic(inboundNetworkConnection net.Conn, outboundNetworkConnection net.Conn) {
	for {
		byteBuffer := make([]byte, 512)

		// Read the inbound network connection.

		numberOfBytesRead, err := inboundNetworkConnection.Read(byteBuffer)
		if err != nil {
			return
		}

		// Print what was received and sent.
		
		inboundMessage := byteBuffer[0:numberOfBytesRead]
		fmt.Println(">>>", string(inboundMessage))
		outboundMessage := fmt.Sprintf("Relay: \"%s\"", inboundMessage)
		fmt.Println("<<<", outboundMessage)
		
		// Write to outbound network connection.
		_, err = outboundNetworkConnection.Write([]byte(outboundMessage))
		if err != nil {
			log.Fatal("Writing client error: ", err)
		}
	}
}

// Function for the "command pattern".
func Command(argv []string) {

	usage := `
Usage:
    mock-relay net [options]

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

	// ...

	if isDebug {
		log.Printf("Listening on '%s' network with address '%s'", inboundNetwork, inboundAddress)
		log.Printf("Sending   to '%s' network with address '%s'", outboundNetwork, outboundAddress)
		
	}

	// Inbound network connection.  net.Listen creates a server.
	
	inboundConnection, err := net.Listen(inboundNetwork, inboundAddress)
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
	}(inboundConnection, sigc)
	
    // Outbound network connection.  net.Dial creates a client.
    
	outboundConnection, err := net.Dial(outboundConnection, outboundAddress)
	if err != nil {
		log.Fatal("Dial error", err)
	}
	defer networkConnection.Close()



	// As a server, Read and Echo loop.

	for {
		inboundNetworkConnection, err := inboundConnection.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}
		go inboundTraffic(inboundNetworkConnection, outboundConnection)
	}	
	
	
	// *****************************************************


	
	// Start asynchronous Reader.

	go reader(networkConnection)

	// Configure listener to exit when program ends.

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(listener net.Listener, c chan os.Signal) {
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)
		listener.Close()
		os.Exit(0)
	}(listener, sigc)

	// Read and Echo loop.

	for {
		networkConnection, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}
		go inboundTraffic(networkConnection)
	}
}
