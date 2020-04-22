package main

import (
	"flag"
	"github.com/egorban/ndtpClient/pkg/ndtpclient"
)

func main() {
	var serverAddress string
	var terminalID int
	flag.StringVar(&serverAddress, "s", "localhost:9000", "server address (e.g. 'localhost:9000')")
	flag.IntVar(&terminalID, "i", 1, "terminal ID (e.g. '1')")
	flag.Parse()
	ndtpclient.Start(serverAddress, terminalID)
}
