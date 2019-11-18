package ndtpclient

import (
	"encoding/binary"
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

var (
	numSend    = 0
	numConfirm = 0
	numControl = 0
)

const (
	defaultBufferSize        = 1024
	writeTimeout             = 10 * time.Second
	readTimeout              = 180 * time.Second
	NPH_SRV_GENERIC_CONTROLS = 0
	NPH_SRV_NAVDATA          = 1
	NPH_RESULT               = 0
)

func Start(addr string, terminalID uint32, num int) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("got error: %v", err)
	}
	defer conn.Close()
	packetAuth = formAuthPacket(terminalID)
	err = sendAndReceive(conn, packetAuth)
	if err != nil {
		log.Printf("got error: %v", err)
	}
	for i := 0; i < num; i++ {
		err = sendNewMessage(conn, i)
		if err != nil {
			log.Printf("got error: %v", err)
		}
	}
	log.Printf("Results:\nNum sent packets: %d\nNum received confirm packets: %d\nNum received control packets: %d",
		numSend, numConfirm, numControl)
}

func formAuthPacket(terminalID uint32) []byte {
	binary.LittleEndian.PutUint32(packetAuth[31:35], terminalID)
	crc := crc16(packetAuth[15:])
	binary.BigEndian.PutUint16(packetAuth[6:8], crc)
	return packetAuth
}

func sendNewMessage(conn net.Conn, i int) error {
	changes := map[string]int{ndtp.NphReqID: i}
	packetNav = ndtp.Change(packetNav, changes)
	return sendAndReceive(conn, packetNav)
}

func sendAndReceive(c net.Conn, packet []byte) (err error) {
	err = send(c, packet)
	if err != nil {
		return
	}
	parsedPacket := new(ndtp.Packet)
	parsedPacket.Parse(packet)
	log.Printf("send: %v", parsedPacket.String())
	numSend++
	var b [defaultBufferSize]byte
	_, err = c.Read(b[:])
	if err != nil {
		return
	}
	parsedPacket.Parse(b[:])
	log.Printf("receive: %v", parsedPacket.String())
	if parsedPacket.Nph.ServiceID == NPH_SRV_GENERIC_CONTROLS {
		numControl++
	}
	if parsedPacket.Nph.ServiceID == NPH_SRV_NAVDATA && parsedPacket.Nph.PacketType == NPH_RESULT {
		numConfirm++
	}
	return
}

func send(conn net.Conn, packet []byte) error {
	err := conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	if err != nil {
		return err
	}
	_, err = conn.Write(packet)
	return err
}

func parsePacket(packet []byte) string {
	ndtpPacket := new(ndtp.Packet)
	ndtpPacket.Parse(packet)
	return ndtpPacket.String()
}

func crc16(bs []byte) (crc uint16) {
	l := len(bs)
	crc = 0xFFFF
	for i := 0; i < l; i++ {
		crc = (crc >> 8) ^ crc16tab[(crc&0xff)^uint16(bs[i])]
	}
	return
}

var crc16tab = [256]uint16{
	0x0000, 0xC0C1, 0xC181, 0x0140, 0xC301, 0x03C0, 0x0280, 0xC241,
	0xC601, 0x06C0, 0x0780, 0xC741, 0x0500, 0xC5C1, 0xC481, 0x0440,
	0xCC01, 0x0CC0, 0x0D80, 0xCD41, 0x0F00, 0xCFC1, 0xCE81, 0x0E40,
	0x0A00, 0xCAC1, 0xCB81, 0x0B40, 0xC901, 0x09C0, 0x0880, 0xC841,
	0xD801, 0x18C0, 0x1980, 0xD941, 0x1B00, 0xDBC1, 0xDA81, 0x1A40,
	0x1E00, 0xDEC1, 0xDF81, 0x1F40, 0xDD01, 0x1DC0, 0x1C80, 0xDC41,
	0x1400, 0xD4C1, 0xD581, 0x1540, 0xD701, 0x17C0, 0x1680, 0xD641,
	0xD201, 0x12C0, 0x1380, 0xD341, 0x1100, 0xD1C1, 0xD081, 0x1040,
	0xF001, 0x30C0, 0x3180, 0xF141, 0x3300, 0xF3C1, 0xF281, 0x3240,
	0x3600, 0xF6C1, 0xF781, 0x3740, 0xF501, 0x35C0, 0x3480, 0xF441,
	0x3C00, 0xFCC1, 0xFD81, 0x3D40, 0xFF01, 0x3FC0, 0x3E80, 0xFE41,
	0xFA01, 0x3AC0, 0x3B80, 0xFB41, 0x3900, 0xF9C1, 0xF881, 0x3840,
	0x2800, 0xE8C1, 0xE981, 0x2940, 0xEB01, 0x2BC0, 0x2A80, 0xEA41,
	0xEE01, 0x2EC0, 0x2F80, 0xEF41, 0x2D00, 0xEDC1, 0xEC81, 0x2C40,
	0xE401, 0x24C0, 0x2580, 0xE541, 0x2700, 0xE7C1, 0xE681, 0x2640,
	0x2200, 0xE2C1, 0xE381, 0x2340, 0xE101, 0x21C0, 0x2080, 0xE041,
	0xA001, 0x60C0, 0x6180, 0xA141, 0x6300, 0xA3C1, 0xA281, 0x6240,
	0x6600, 0xA6C1, 0xA781, 0x6740, 0xA501, 0x65C0, 0x6480, 0xA441,
	0x6C00, 0xACC1, 0xAD81, 0x6D40, 0xAF01, 0x6FC0, 0x6E80, 0xAE41,
	0xAA01, 0x6AC0, 0x6B80, 0xAB41, 0x6900, 0xA9C1, 0xA881, 0x6840,
	0x7800, 0xB8C1, 0xB981, 0x7940, 0xBB01, 0x7BC0, 0x7A80, 0xBA41,
	0xBE01, 0x7EC0, 0x7F80, 0xBF41, 0x7D00, 0xBDC1, 0xBC81, 0x7C40,
	0xB401, 0x74C0, 0x7580, 0xB541, 0x7700, 0xB7C1, 0xB681, 0x7640,
	0x7200, 0xB2C1, 0xB381, 0x7340, 0xB101, 0x71C0, 0x7080, 0xB041,
	0x5000, 0x90C1, 0x9181, 0x5140, 0x9301, 0x53C0, 0x5280, 0x9241,
	0x9601, 0x56C0, 0x5780, 0x9741, 0x5500, 0x95C1, 0x9481, 0x5440,
	0x9C01, 0x5CC0, 0x5D80, 0x9D41, 0x5F00, 0x9FC1, 0x9E81, 0x5E40,
	0x5A00, 0x9AC1, 0x9B81, 0x5B40, 0x9901, 0x59C0, 0x5880, 0x9841,
	0x8801, 0x48C0, 0x4980, 0x8941, 0x4B00, 0x8BC1, 0x8A81, 0x4A40,
	0x4E00, 0x8EC1, 0x8F81, 0x4F40, 0x8D01, 0x4DC0, 0x4C80, 0x8C41,
	0x4400, 0x84C1, 0x8581, 0x4540, 0x8701, 0x47C0, 0x4680, 0x8641,
	0x8201, 0x42C0, 0x4380, 0x8341, 0x4100, 0x81C1, 0x8081, 0x4040}
