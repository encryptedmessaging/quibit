package quibit

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestAcceptance(t *testing.T) {
	log := make(chan string, 100)
	recvChan := make(chan Frame, 10)
	sendChan := make(chan Frame, 10)
	peerChan := make(chan Peer)
	port := uint16(4444)

	// Initialize Quibit
	err := Initialize(log, recvChan, sendChan, peerChan, port)
	if err != nil {
		fmt.Println("ERROR INITIALIZING! ", err)
		t.FailNow()
	}

	// Test 1: Manual Connection, look for receive
	conn, err := net.Dial("tcp", "127.0.0.1:4444")
	if err != nil {
		fmt.Println("Error connecting: ", err)
		t.FailNow()
	}

	time.Sleep(time.Millisecond)
	if len(peerList) == 0 {
		fmt.Println("Not in peer list!")
		t.FailNow()
	}

	data := []byte{'a', 'b', 'c', 'd'}
	frame := new(Frame)
	frame.Configure(data, 1, 1)

	buf, _ := frame.Header.ToBytes()
	_, err = conn.Write(buf)
	if err != nil {
		fmt.Println("Error writing header: ", err)
		t.FailNow()
	}

	_, err = conn.Write(frame.Payload)
	if err != nil {
		fmt.Println("Error writing payload: ", err)
		t.FailNow()
	}

	frame2 := <-recvChan
	if string(frame2.Payload) != string(data) {
		fmt.Println("Bad frame! ", frame2)
		t.FailNow()
	}

	if frame2.Peer != conn.LocalAddr().String() {
		fmt.Println("Peer doesn't match! ", frame2.Peer, conn.LocalAddr().String())
		t.FailNow()
	}

	// Test 2: Send, look for manual receive
	sendChan <- frame2
	time.Sleep(time.Millisecond)

	// So now we have a connection.  Let's shake hands.
	header3, _ := recvHeader(conn, log)
	frame3, err := recvPayload(conn, header3)

	if err != nil {
		fmt.Println("Error Receiving Frame 3... ", err)
		t.FailNow()
	}

	if string(frame3.Payload) != string(data) {
		fmt.Println("Bad frame! ", frame3)
		t.FailNow()
	}

	conn.Close()
	Cleanup()
}
