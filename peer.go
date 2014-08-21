/**
    Copyright 2014 JARST, LLC
    
    This file is part of Quibit.

    Quibit is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    LICENSE file for details.
**/

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
	external bool
}

func peerFromConn(conn *net.Conn, external bool) Peer {
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
	p.external = external
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
	return (p.conn != nil)
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

	if p.external {
		incomingConnections++
	}


	return nil
}

func (p *Peer) Disconnect() {
	if p.conn == nil {
		return
	}
	(*p.conn).Close()
	p.conn = nil
	if p.external {
		incomingConnections--
	}
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

	_, err = (*p.conn).Write(append(headerBytes, frame.Payload...))
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
		// So now we have a connection.  Let's start receiving
		frame, err := recvAll(*p.conn, log)

		if err != nil {
			log <- fmt.Sprintln("Error receiving from Peer: ", err)
			KillPeer(p.String())
			break
		} else {
			frame.Peer = p.String()
			recvChan <- frame
		}
	} // End for
} // End receive()
