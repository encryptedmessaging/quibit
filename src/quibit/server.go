package quibit

import (
	"fmt"
	"net"
)

func initServer(recvChan chan Frame, peerChan chan Peer, port string) error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	go func() {
		for {
			// Listen for new Peer
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println(err.Error())
			}

			// Add peer to peer channel
			p := peerFromConn(conn)
			if p != nil {
				peerChan <- *p
			}

		} // End for
	}()
	return nil
}
