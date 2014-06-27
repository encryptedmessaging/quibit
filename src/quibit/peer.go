package quibit

import (
	"net"
)

type Peer struct {
	IP net.IP
	Port uint16
	conn net.Conn
}

func peerFromConn(conn net.Conn) *Peer {
	addr := conn.RemoteAddr()
	if addr.Network() != "tcp" {
		return nil
	}
	tcpAddr := net.TCPAddr(addr)
	p := new(Peer)
	p.IP = tcpAddr.IP
	p.Port = uint6(tcpAddr.Port)
	p.conn = conn
	return p
}

func (p *Peer) IsConnected() bool {
	if p == nil {
		return false
	}
	return (p.conn == nil)
}

func (p *Peer) Connect() error {
	if p == nil {
		return errors.New("")
	}
}