package quibit

import (
	"fmt"
)

var peerList map[string]*Peer
var quit chan bool

// Initialize the Quibit Service
// Frames from the network will be sent to recvChan, and includes the sending peer
// Frames for the network should be sent to sendChan, and include the receiving peer
// New Peers for connecting should be sent to peerChan.
// A local server will be started on the port specified by "port"
// If an error is returned, than neither the server or mux has been started.
func Initialize(log chan string, recvChan, sendChan chan Frame, peerChan chan Peer, port uint16) error {
	var err error

	err = initServer(recvChan, peerChan, fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	peerList = make(map[string]*Peer)
	quit = make(chan bool)

	go mux(recvChan, sendChan, peerChan, quit, log)
	return nil
}

// Cleanup the Quibit Service
// End the mux and server routines, and Disconnect from all peers.
func Cleanup() {
	quit <- true
}

func KillPeer(p string) {
	peer, ok := peerList[p]
	if ok {
		peer.Disconnect()
		delete(peerList, p)
	}
}

func GetPeer(p string) *Peer {
	peer, ok := peerList[p]
	if ok {
		return peer
	} else {
		return nil
	}
}

func mux(recvChan, sendChan chan Frame, peerChan chan Peer, quit chan bool, log chan string) {
	var frame Frame
	var peer Peer
	var err error

	for {
		select {
		case frame = <-sendChan:
			// Received frame to send to peer(s)
			if frame.Peer == "" {
				// Send to all peers
				for key, p := range peerList {
					err = p.sendFrame(frame)
					if err != nil {
						if err.Error() != QuibitError(eHEADER).Error() {
							// Disconnect from Peer
							p.Disconnect()
							delete(peerList, key)
						}
						// Malformed header, break out of for loop
						fmt.Println("Malformed header in frame!")
						break
					}
				}
			} else {
				// Send to one peer
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
						fmt.Println("Malformed header in frame!")
					}
				}
			}

		case peer = <-peerChan:
			// Received a new peer to connect to...
			err = peer.connect()
			if err == nil {
				go peer.receive(recvChan, log)
				peerList[peer.String()] = &peer

			} else {
				fmt.Println("Error adding peer: ", err)
			}
		case <-quit:
			for _, p := range peerList {
				p.Disconnect()
			}
			return
		} // End select
	} // End for
} // End mux()
