package main

import (
	"bytes"
	"log"
	"net"
	"testing"
)

func TestDnsMsg_toResponse(t *testing.T) {
	msg := DnsMsg{ID: 1234, Message: []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}
	response := new(bytes.Buffer)

	msg.toResponse(response)

	expected := []byte{0x04, 0xD2, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	if !bytes.Equal(response.Bytes(), expected) {
		t.Errorf("Expected response %v, got %v", expected, response.Bytes())
	}
}

func TestUdpServer(t *testing.T) {
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	go dnsServer()

	clientAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:6969")
	if err != nil {
		t.Fatalf("Failed to resolve client address: %v", err)
	}

	clientConn, err := net.DialUDP("udp", clientAddr, serverAddr)
	if err != nil {
		t.Fatalf("Failed to dial server: %v", err)
	}
	log.Println("Test server at addr: ", serverAddr.String())
	defer clientConn.Close()

	// request := []byte{0x01, 0x02, 0x03, 0x04}
	request := []byte("Hi dns server, this is from test")
	_, err = clientConn.Write(request)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	buf := make([]byte, 512)
	size, _, err := clientConn.ReadFromUDP(buf)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	expectedResponse := []byte{0x04, 0xD2, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	if !bytes.Equal(buf[:size], expectedResponse) {
		t.Errorf("Expected response %v, got %v", expectedResponse, buf[:size])
	}
}
