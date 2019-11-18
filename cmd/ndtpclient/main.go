package main

import (
	"flag"
	"github.com/egorban/ndtpClient/pkg/ndtpclient"
)

func main() {
	var serverAddress string
    var terminalID int
    var numPackets int
    flag.StringVar(&serverAddress, "s", "", "server address (e.g. 'localhost:8080')")
    flag.IntVar(&terminalID, "i", 0, "terminal ID (e.g. '1')")
    flag.IntVar(&numPackets, "n", 0, "packets number (e.g. '100')")
    flag.Parse()
    ndtpclient.Start(serverAddress, uint32(terminalID), numPackets)
}
