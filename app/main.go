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

		if receivedHeader, receivedQuestions, _, err := dns.ParseDNSMessage(message); err != nil {
			log.Fatal(err)
		} else {
			headerFlags := dns.DecodeDNSFlags(receivedHeader.Flags)
			headerFlags.QR = 1
			headerFlags.AA = 0
			headerFlags.TC = 0
			headerFlags.RA = 0
			headerFlags.Z = 0
			headerFlags.RCODE = 4

			receivedHeader.QDCount = 1
			receivedHeader.ANCount = 1
			// receivedHeader.ID = 1234
			receivedHeader.Flags = dns.EncodeDNSFlags(headerFlags)
			answers := []dns.DNSRecord{}

			for i, question := range receivedQuestions {
				receivedQuestions[i].QType = 1
				receivedQuestions[i].QClass = 1
				data, err := dns.IPAddressStringToBytes("8.8.8.8")

				if err != nil {
					log.Fatal("Error trying to convert ip to bytes", err)
				}
				record := dns.DNSRecord{
					Name:     question.QName,
					Type:     receivedQuestions[i].QType, // A record
					Class:    receivedQuestions[i].QType, // IN
					TTL:      60,
					RDLength: uint16(len(data)),
					// RData:    []byte{8, 8, 8, 8}, // 192.0.2.1
					RData: data, // 192.0.2.1
				}
				answers = append(answers, record)
			}

			response = dns.EncodeDNSMessage(receivedHeader, receivedQuestions, answers)
		}

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			log.Println("Failed to send response:", err)
		}
	}
}

func main() {
	dnsServer()
}
