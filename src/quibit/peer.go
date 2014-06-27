package quibit

import (
	"fmt"
	"net"
	"strconv"
)

type Peer struct {
	IP   net.IP
	Port uint16
	conn net.Conn
}

func peerFromConn(conn net.Conn) *Peer {
	addr := conn.RemoteAddr()
	if addr.Network() != "tcp" {
		return nil
	}
	ip, portStr, err := net.SplitHostPort(addr.String())
	port, _ := strconv.Atoi(portStr)
	if err != nil {
		return nil
	}

	// Create New Peer
	p := new(Peer)
	p.IP = net.ParseIP(ip)
	p.Port = uint16(port)

	p.conn = conn
	return p
}

func (p *Peer) String() string {
	if p == nil {
		return ""
	}
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}

func (p *Peer) IsConnected() bool {
	if p == nil {
		return false
	}
	return (p.conn == nil)
}

func (p *Peer) connect() error {
	// Check for sane peer object
	if p == nil {
		return QuibitError(eNILOBJ)
	}
	if p.conn != nil {
		return nil
	}

	var err error
	p.conn, err = net.Dial("tcp", p.String())
	if err != nil {
		p.conn.Close()
		p.conn = nil
		return err
	}
	return nil
}

func (p *Peer) disconnect() {
	if p.conn == nil {
		fmt.Println("Peer already disconnected!")
		return
	}
	p.conn.Close()
	p.conn = nil
}

func (p *Peer) sendFrame(frame Frame) error {
	if p == nil {
		return QuibitError(eNILOBJ)
	}
	if p.conn == nil {
		return QuibitError(eNILOBJ)
	}

	var err error

	headerBytes, err := frame.Header.ToBytes()
	if err != nil {
		return QuibitError(eHEADER)
	}

	_, err = p.conn.Write(headerBytes)
	if err != nil {
		return err
	}

	_, err = p.conn.Write(frame.Payload)
	if err != nil {
		return err
	}

	return nil
}

func (p *Peer) receive(recvChan chan Frame, log chan string) {
	if p.conn == nil {
		return
	}
	for {
		// So now we have a connection.  Let's shake hands.
		header := recvHeader(p.conn, log)
		frame, err := recvPayload(p.conn, header)
		if err != nil {
			fmt.Println("Error receiving from Peer: ", err)
		} else {
			recvChan <- frame
		}
	} // End for
} // End receive()
