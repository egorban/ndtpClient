package main

import (
	"flag"
	"github.com/egorban/ndtpClient/pkg/ndtpclient"
)

func main() {
	var serverAddress string
	var terminalID int
	var numPackets int
	var numControlPackets int
	flag.StringVar(&serverAddress, "s", "localhost:9000", "server address (e.g. 'localhost:8080')")
	flag.IntVar(&terminalID, "i", 1, "terminal ID (e.g. '1')")
	flag.IntVar(&numPackets, "n", 100, "packets number (e.g. '100')")
	flag.IntVar(&numControlPackets, "c", 0, "control packets number (e.g. '1')")
	flag.Parse()
	ndtpclient.Start(serverAddress, terminalID, numPackets, numControlPackets)
}
