package quibit

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type Peer struct {
	IP   net.IP
	Port uint16
	conn *net.Conn
}

func peerFromConn(conn *net.Conn) Peer {
	var p Peer
	if conn == nil {
		return p
	}
	
	addr := (*conn).RemoteAddr()
	if addr.Network() != "tcp" {
		return p
	}
	ip, portStr, err := net.SplitHostPort(addr.String())
	port, _ := strconv.Atoi(portStr)
	if err != nil {
		return p
	}

	// Create New Peer
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
	lConn, err := net.DialTimeout("tcp", p.String(), 10 * time.Second)
	p.conn = &lConn
	if err != nil {
		p.conn = nil
		return err
	}

	// Set Keep-Alives
	(*p.conn).(*net.TCPConn).SetKeepAlive(true)
	(*p.conn).(*net.TCPConn).SetKeepAlivePeriod(time.Second)


	return nil
}

func (p *Peer) Disconnect() {
	if p.conn == nil {
		fmt.Println("Peer already disconnected!")
		return
	}
	(*p.conn).Close()
	p.conn = nil
}

func (p *Peer) sendFrame(frame Frame) error {
	var n int
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

	n, err = (*p.conn).Write(append(headerBytes, frame.Payload...))
	fmt.Printf("Wrote %d Bytes to Peer: %s\n", n, p)
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
		fmt.Println("Called recvAll with peer: ", p)
		frame, err := recvAll(*p.conn, log)
		fmt.Println("Finished recvAll with peer: ", p)

		if err != nil {
			fmt.Println("Error receiving from Peer: ", err)
			p.Disconnect()
			break
		} else {
			frame.Peer = p.String()
			fmt.Println("Sending to Recv Channel from... ", (*p.conn).RemoteAddr())
			recvChan <- frame
		}
	} // End for
} // End receive()
