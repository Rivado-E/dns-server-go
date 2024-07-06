package parser
package main

import (
	"encoding/binary"
	"fmt"
	"log"
)

// DNSHeader represents the DNS header structure
type DNSHeader struct {
	ID      uint16
	Flags   uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

// DNSQuestion represents a DNS question structure
type DNSQuestion struct {
	QName  string
	QType  uint16
	QClass uint16
}

// DNSRecord represents a DNS resource record structure
type DNSRecord struct {
	Name     string
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLength uint16
	RData    []byte
}

func main() {
	// Example DNS message bytes (replace with your actual byte slice)
	message := []byte{
		0x12, 0x34, // ID
		0x01, 0x00, // Flags
		0x00, 0x01, // QDCount
		0x00, 0x02, // ANCount
		0x00, 0x00, // NSCount
		0x00, 0x01, // ARCount

		// Question section (example with one question)
		0x03, 'w', 'w', 'w', 0x05, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 0x03, 'c', 'o', 'm', 0x00, // QName
		0x00, 0x01, // QType (A record)
		0x00, 0x01, // QClass (IN)

		// Answer section (example with one answer)
		0xc0, 0x0c, // Name pointer (compression)
		0x00, 0x01, // Type (A record)
		0x00, 0x01, // Class (IN)
		0x00, 0x00, 0x00, 0x0a, // TTL
		0x00, 0x04, // RDLength
		0x7f, 0x00, 0x00, 0x01, // RData (127.0.0.1)

		// Additional record section (example with one additional record)
		0xc0, 0x0c, // Name pointer (compression)
		0x00, 0x01, // Type (A record)
		0x00, 0x01, // Class (IN)
		0x00, 0x00, 0x00, 0x0a, // TTL
		0x00, 0x04, // RDLength
		0x7f, 0x00, 0x00, 0x01, // RData (127.0.0.1)
	}

	// Parse DNS message
	header, questions, answers, authorities, additionals, err := ParseDNSMessage(message)
	if err != nil {
		log.Fatalf("Error parsing DNS message: %v", err)
	}

	// Print parsed data
	fmt.Printf("DNS Header:\n")
	fmt.Printf("  ID: %d\n", header.ID)
	fmt.Printf("  Flags: %d\n", header.Flags)
	fmt.Printf("  QDCount: %d\n", header.QDCount)
	fmt.Printf("  ANCount: %d\n", header.ANCount)
	fmt.Printf("  NSCount: %d\n", header.NSCount)
	fmt.Printf("  ARCount: %d\n", header.ARCount)

	fmt.Printf("\nDNS Questions:\n")
	for _, q := range questions {
		fmt.Printf("  QName: %s\n", q.QName)
		fmt.Printf("  QType: %d\n", q.QType)
		fmt.Printf("  QClass: %d\n", q.QClass)
	}

	fmt.Printf("\nDNS Answers:\n")
	for _, a := range answers {
		fmt.Printf("  Name: %s\n", a.Name)
		fmt.Printf("  Type: %d\n", a.Type)
		fmt.Printf("  Class: %d\n", a.Class)
		fmt.Printf("  TTL: %d\n", a.TTL)
		fmt.Printf("  RDLength: %d\n", a.RDLength)
		fmt.Printf("  RData: %v\n", a.RData)
	}

	fmt.Printf("\nDNS Authorities:\n")
	for _, auth := range authorities {
		fmt.Printf("  Name: %s\n", auth.Name)
		fmt.Printf("  Type: %d\n", auth.Type)
		fmt.Printf("  Class: %d\n", auth.Class)
		fmt.Printf("  TTL: %d\n", auth.TTL)
		fmt.Printf("  RDLength: %d\n", auth.RDLength)
		fmt.Printf("  RData: %v\n", auth.RData)
	}

	fmt.Printf("\nDNS Additionals:\n")
	for _, add := range additionals {
		fmt.Printf("  Name: %s\n", add.Name)
		fmt.Printf("  Type: %d\n", add.Type)
		fmt.Printf("  Class: %d\n", add.Class)
		fmt.Printf("  TTL: %d\n", add.TTL)
		fmt.Printf("  RDLength: %d\n", add.RDLength)
		fmt.Printf("  RData: %v\n", add.RData)
	}
}

// ParseDNSMessage parses a DNS message byte slice into structured DNS components
func ParseDNSMessage(message []byte) (header DNSHeader, questions []DNSQuestion, answers []DNSRecord, authorities []DNSRecord, additionals []DNSRecord, err error) {
	// Parse DNS header
	header.ID = binary.BigEndian.Uint16(message[0:2])
	header.Flags = binary.BigEndian.Uint16(message[2:4])
	header.QDCount = binary.BigEndian.Uint16(message[4:6])
	header.ANCount = binary.BigEndian.Uint16(message[6:8])
	header.NSCount = binary.BigEndian.Uint16(message[8:10])
	header.ARCount = binary.BigEndian.Uint16(message[10:12])

	// Offset to start of questions section
	offset := 12

	// Parse questions
	for i := 0; i < int(header.QDCount); i++ {
		qname, qnameLen := parseDomainName(message, offset)
		offset += qnameLen

		qtype := binary.BigEndian.Uint16(message[offset : offset+2])
		offset += 2

		qclass := binary.BigEndian.Uint16(message[offset : offset+2])
		offset += 2

		question := DNSQuestion{
			QName:  qname,
			QType:  qtype,
			QClass: qclass,
		}

		questions = append(questions, question)
	}

	// Parse answers, authorities, and additionals (similar structure)
	answers, offset = parseDNSRecords(message, offset, int(header.ANCount))
	authorities, offset = parseDNSRecords(message, offset, int(header.NSCount))
	additionals, offset = parseDNSRecords(message, offset, int(header.ARCount))

	return header, questions, answers, authorities, additionals, nil
}

// parseDNSRecords parses DNS resource records from the message
func parseDNSRecords(message []byte, offset int, count int) (records []DNSRecord, newOffset int) {
	for i := 0; i < count; i++ {
		name, nameLen := parseDomainName(message, offset)
		offset += nameLen

		typ := binary.BigEndian.Uint16(message[offset : offset+2])
		offset += 2

		class := binary.BigEndian.Uint16(message[offset : offset+2])
		offset += 2

		ttl := binary.BigEndian.Uint32(message[offset : offset+4])
		offset += 4

		rdlength := binary.BigEndian.Uint16(message[offset : offset+2])
		offset += 2

		rdata := message[offset : offset+int(rdlength)]
		offset += int(rdlength)

		record := DNSRecord{
			Name:     name,
			Type:     typ,
			Class:    class,
			TTL:      ttl,
			RDLength: rdlength,
			RData:    rdata,
		}

		records = append(records, record)
	}

	return records, offset
}

// parseDomainName parses a DNS domain name from the message
func parseDomainName(message []byte, offset int) (name string, length int) {
	var parts []string
	for {
		length := int(message[offset])
		if length == 0 {
			offset++
			break
		}

		offset++
		part := string(message[offset : offset+length])
		parts = append(parts, part)
		offset += length
	}

	name = ""
	for i, part := range parts {
		if i > 0 {
			name += "."
		}
		name += part
	}

	return name, offset - offset
}

