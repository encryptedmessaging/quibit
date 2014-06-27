package quibit

import (
	"fmt"
)

var peerList map[string]*Peer

func Initialize(log chan string, recvChan, sendChan chan Frame, peerChan chan Peer, port uint16) error {
	var err error

	err = initServer(recvChan, peerChan, fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	peerList = make(map[string]*Peer)

	go mux(recvChan, sendChan, peerChan, log)
	return nil
}

func mux(recvChan, sendChan chan Frame, peerChan chan Peer, log chan string) {
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
							p.disconnect()
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
							p.disconnect()
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
		} // End select
	} // End for
} // End mux()
