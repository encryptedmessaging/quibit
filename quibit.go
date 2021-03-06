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
)

var peerList map[string]Peer
var quit chan bool
var incomingConnections int


//Message Types:
//
//Broadcast goes to all connected peers except for the Peer specified in the Frame.
//
//Request and Reply go to the Peer specified in the Frame.
const (
	BROADCAST = iota
	REQUEST   = iota
	REPLY     = iota
)


//Initialize the Quibit Service
//
//Frames from the network will be sent to recvChan, and includes the sending peer
//
//Frames for the network should be sent to sendChan, and include the receiving peer
//
//New Peers for connecting should be sent to peerChan.
//
//A local server will be started on the port specified by "port"
//
//If an error is returned, than neither the server or mux has been started.
func Initialize(log chan string, recvChan, sendChan chan Frame, peerChan chan Peer, port uint16) error {
	var err error

	incomingConnections = 0

	err = initServer(recvChan, peerChan, fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	peerList = make(map[string]Peer)
	quit = make(chan bool)

	go mux(recvChan, sendChan, peerChan, quit, log)
	return nil
}


//Cleanup the Quibit Service
//
//End the mux and server routines, and Disconnect from all peers.
func Cleanup() {
	quit <- true
}


//KillPeer Force Disconnects from a Peer
//
//All incoming data is dropped and the peer is removed from the Peer List
func KillPeer(p string) {
	peer, ok := peerList[p]
	if ok {
		peer.Disconnect()
		delete(peerList, p)
	}
}


//Get Peer associated with the given <IP>:<Host> string
//
//<nil> Signifies a disconnected or unknown peer.
func GetPeer(p string) *Peer {
	peer, ok := peerList[p]
	if ok {
		return &peer
	} else {
		return nil
	}
}


//Status returns the Current Connection Status
//
// Returns 0 on disconnected.  
// Returns 1 on Client Connection (Outgoing Connections Only).  
// Returns 2 On Full Connection (Incoming and Outgoing Connections).
func Status() int {
	if len(peerList) < 1 {
		return 0
	}
	if incomingConnections < 1 {
		return 1
	}
	return 2
}

func mux(recvChan, sendChan chan Frame, peerChan chan Peer, quit chan bool, log chan string) {
	var frame Frame
	var peer Peer
	var err error

	for {
		select {
		case frame = <-sendChan:
			// Received frame to send to peer(s)
			if frame.Header.Type == BROADCAST {
				// Send to all peers
				for key, p := range peerList {
					if key == frame.Peer {
						// Exclude peer in message
						continue
					}

					err = p.sendFrame(frame)
					if err != nil {
						if err.Error() != QuibitError(eHEADER).Error() {
							// Disconnect from Peer
							p.Disconnect()
							delete(peerList, key)
						}
						// Malformed header, break out of for loop
						log <- fmt.Sprintln("Error sending frame: ", err)
					}
				}
			} else {
				// Send to one peer
				if frame.Peer == "" {
					// Error, can't broadcast a non-broadcast message
					break
				}
				p, ok := peerList[frame.Peer]
				if ok {
					err = p.sendFrame(frame)
					if err != nil {
						if err.Error() != QuibitError(eHEADER).Error() {
							// Disconnect from Peer
							p.Disconnect()
							delete(peerList, frame.Peer)
						}
						// Malformed header
						log <- fmt.Sprintln("Malformed header in frame!")
					}
				} else {
					log <- fmt.Sprintln("Peer not found: ", frame.Peer)
				}
			}

		case peer = <-peerChan:
			// Received a new peer to connect to...
			err = peer.connect()
			if err == nil {
				peerList[peer.String()] = peer

				// Prevent overwriting...
				rawPeer := new(Peer)
				*rawPeer = peer
				go rawPeer.receive(recvChan, log)

			} else {
				log <- fmt.Sprintln("Error adding peer: ", err)
			}
		case <-quit:
			for _, p := range peerList {
				p.Disconnect()
			}
			return
		} // End select
	} // End for
} // End mux()
