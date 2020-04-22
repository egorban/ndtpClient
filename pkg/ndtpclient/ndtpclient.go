package ndtpclient

import (
	"github.com/ashirko/navprot/pkg/ndtp"
	"log"
	"net"
	"time"
)

var packetNav = []byte{126, 126, 74, 0, 2, 0, 107, 210, 2, 0, 0, 0, 0, 0, 0, 1, 0, 101, 0, 1, 0, 171,
	20, 0, 0, 0, 0, 36, 141, 198, 90, 87, 110, 119, 22, 201, 186, 64, 33, 224, 203, 0, 0, 0, 0, 83, 1, 0,
	0, 220, 0, 4, 0, 2, 0, 22, 0, 67, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 167, 97, 0, 0, 31, 6, 0, 0, 8,
	0, 2, 0, 0, 0, 0, 0}

var packetAuth = []byte{126, 126, 59, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 1, 0, 0, 0, 0,
	0, 6, 0, 2, 0, 2, 3, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 51, 53, 53, 48, 57, 52, 48, 52, 51, 49, 56,
	56, 51, 49, 49, 50, 53, 48, 48, 49, 54, 53, 48, 53, 56, 49, 53, 53, 51, 55, 0}

const (
	defaultBufferSize     = 1024
	writeTimeout          = 10 * time.Second
	readTimeout           = 180 * time.Second
	sendPeriod            = 20 * time.Second
	nphSrvGenericControls = 0
	nphSrvNavdata         = 1
	nphResult             = 0
)

func Start(addr string, terminalID int) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("got error: %v", err)
		return
	}
	defer conn.Close()
	log.Printf("NDTP client was started. Server address: %v; Terminal ID: %v", addr, terminalID)
	err = setConnection(conn, terminalID)
	if err != nil {
		log.Printf("got error: %v", err)
	}
	go sendData(conn)
	receiveReply(conn)
}

func setConnection(conn net.Conn, terminalID int) (err error) {
	packetAuth = formAuthPacket(terminalID)
	err = send(conn, packetAuth)
	if err != nil {
		return
	}
	log.Printf("send first packet: %v", packetAuth)
	var b [defaultBufferSize]byte
	_, err = conn.Read(b[:])
	if err != nil {
		return
	}
	parsedPacket := new(ndtp.Packet)
	_, err = parsedPacket.Parse(b[:])
	if err != nil {
		return
	}
	return
}

func sendData(conn net.Conn) {
    i:=0
	for _ = range time.Tick(sendPeriod) {
		err := sendNewMessage(conn, i)
		if err != nil {
			log.Printf("got error: %v", err)
		}
		i++
	}
}

func receiveReply(conn net.Conn) {
	var restBuf []byte
	for {
		var b [defaultBufferSize]byte
		err := conn.SetReadDeadline(time.Now().Add(readTimeout))
		if err != nil {
			log.Printf("got error: %v", err)
		}
		n, err := conn.Read(b[:])
		if err != nil {
			log.Printf("got error: %v", err)
			break
		}
		log.Printf("read n byte: %d", n)
		restBuf = append(restBuf, b[:n]...)
		for len(restBuf) != 0 {
			parsedPacket := new(ndtp.Packet)
			restBuf, err = parsedPacket.Parse(restBuf)
			if err != nil {
				log.Printf("error while parsing NDTP: %v", err)
				break
			}
			if parsedPacket.Nph.ServiceID == nphSrvGenericControls {
				log.Printf("receive control packet: %v", parsedPacket.String())
			} else if parsedPacket.Nph.ServiceID == nphSrvNavdata && parsedPacket.Nph.PacketType == nphResult {
				log.Printf("receive confirm: %v", parsedPacket.String())
			} else {
				log.Printf("receive other reply: %v", parsedPacket.String())
			}
		}
	}
}

func formAuthPacket(terminalID int) []byte {
	changes := map[string]int{ndtp.PeerAddress: terminalID}
	packetAuth = ndtp.Change(packetAuth, changes)
	return packetAuth
}

func sendNewMessage(conn net.Conn, i int) (err error) {
	changes := map[string]int{ndtp.NphReqID: i}
	packetNav = ndtp.Change(packetNav, changes)
	parsedPacket := new(ndtp.Packet)
	_, err = parsedPacket.Parse(packetNav)
	if err != nil {
		return
	}
	err = send(conn, packetNav)
	if err != nil {
		return
	}
	log.Printf("send packet: %v", parsedPacket.String())
	return
}

func send(conn net.Conn, packet []byte) (err error) {
	err = conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	if err != nil {
		return
	}
	_, err = conn.Write(packet)
	return
}
