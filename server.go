/**
    Copyright 2014 JARST, LLC
    
    This file is part of Quibit.

    Quibit is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    LICENSE file for details.
**/

// Package quibit provides basic Peer-To-Peer asynchronous network
// functionality and peer management.
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
			fmt.Println("Listening on... ", port)
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println(err.Error())
			}

			// Add peer to peer channel
			p := peerFromConn(&conn, true)
			if p.conn != nil {
				fmt.Println("Adding peer... ", p)
				peerChan <- p
			}

		} // End for
	}()
	return nil
}
