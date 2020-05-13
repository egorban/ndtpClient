package main

import (
	"flag"

	"github.com/egorban/ndtpClient/pkg/ndtpclient"
)

func main() {
	var serverAddress string
	var terminalID int
	var time int
	flag.StringVar(&serverAddress, "s", "localhost:9000", "server address (e.g. 'localhost:9000')")
	flag.IntVar(&terminalID, "i", 1, "terminal ID (e.g. '1')")
	flag.IntVar(&time, "t", 20, "time period seconds (e.g. '20')")
	flag.Parse()
	ndtpclient.Start(serverAddress, terminalID, time)
}
