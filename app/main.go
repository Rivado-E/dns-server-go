package main

import (
	"bytes"
	"encoding/binary"
	// "fmt"
	"log"
	"net"
)

type DnsMsg struct {
	ID      int
	Message []byte
}

func (msg *DnsMsg) toResponse(response *bytes.Buffer) {
	err := binary.Write(response, binary.BigEndian, uint16(msg.ID))
	if err != nil {
		log.Println("Unable to construct id for the header: ", err)
	}

	if _, err := response.Write(msg.Message); err != nil {
		log.Println("Unable to add the response: ", err)
	}
}

func dnsServer() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	// udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:4200")
	if err != nil {
		log.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Println("Failed to bind to address:", err)
		return
	}

	log.Println("DNS server started at ", udpAddr.String())
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Error receiving data:", err)
			break
		}

		receivedData := string(buf[:size])
		log.Printf("Received %d bytes from %s : %s\n", size, source, receivedData)

		// Create an empty response
		response := new(bytes.Buffer)
		msg := DnsMsg{1234, []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}

		msg.toResponse(response)

		_, err = udpConn.WriteToUDP(response.Bytes(), source)
		if err != nil {
			log.Println("Failed to send response:", err)
		}
	}
}

func main() {
	dnsServer()
}
