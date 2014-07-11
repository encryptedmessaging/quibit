package quibit

import (
	"bytes"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"net"
	"fmt"
)

const (
	MAGIC = 6667787
)

func recvHeader(conn net.Conn, log chan string) (Header, error) {
	// ret val
	var h Header
	// a buffer for decoing
	var headerBuffer bytes.Buffer
	for {
		headerSize := int(binary.Size(h))
		// Byte slice for moving to buffer
		buffer := make([]byte, headerSize)
		if conn == nil {
			return h, errors.New("Nil connection")
		}
		n, err := conn.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log <- err.Error()
		}
		if n > 0 {
			// Add to header buffer
			headerBuffer.Write(buffer)
			// Check to see if we have the whole header
			if len(headerBuffer.Bytes()) == headerSize {
				h.FromBytes(headerBuffer.Bytes())
				return h, nil
			}
		}
	}

	return h, errors.New("EOF")
}

func recvPayload(conn net.Conn, h Header) (Frame, error) {
	var frame Frame
	fmt.Println("Payload Length: ", h.Length)
	payload := make([]byte, h.Length)
	var payloadBuffer bytes.Buffer
	if h.Length < 1 {
		frame.Payload = nil
		frame.Header = h
		return frame, nil
	}
	for {
		// store in byte array
		n, err := conn.Read(payload)
		if err != nil {
			return frame, err
		}
		if n > 1 {
			// write to buffer
			payloadBuffer.Write(payload)
			// Check to see if we have whole payload
			if len(payloadBuffer.Bytes()) == int(h.Length) {
				// Verify checksum
				frame.Payload = payloadBuffer.Bytes()
				frame.Header = h
				if h.Checksum != sha512.Sum384(payloadBuffer.Bytes()) {
					return frame, errors.New("Incorrect Checksum")
				}
				return frame, nil
			}
		}
	}
	//Should never end here
	panic("RECV PAYLOAD")
	return frame, nil
}
