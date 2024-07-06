package main

import (
	"log"
	"net"

	dns "github.com/codecrafters-io/dns-server-starter-go/lib"
)

func dnsServer() {
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

		message := buf[:size]
		log.Println("Received %d bytes from %s:\n", size, source)
		// dns.PrintMessage(message)

		response := []byte{}

		if receivedHeader, receivedQuestions, err := dns.ParseDNSMessage(message); err != nil {
			log.Fatal(err)
		} else {
			headerFlags := dns.DecodeDNSFlags(receivedHeader.Flags)
			headerFlags.QR = 1

			receivedHeader.QDCount = 1
			receivedHeader.ID = 1234
			receivedHeader.Flags = dns.EncodeDNSFlags(headerFlags)

			for i := range receivedQuestions {
				receivedQuestions[i].QType = 1
				receivedQuestions[i].QClass = 1
			}

			response = dns.EncodeDNSMessage(receivedHeader, receivedQuestions)
			// log.Println("Here is my response")
			// dns.PrintMessage(response)
		}
		// Create an empty response

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			log.Println("Failed to send response:", err)
		}
	}
}

func main() {
	dnsServer()
}
