package quibit

import (
	"bytes"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

const (
	MAGIC = 6667787
)

func recvHeader(conn net.Conn, log chan string) Header {
	// ret val
	var h Header
	// a buffer for decoing
	var headerBuffer bytes.Buffer
	for {
		headerSize := int(binary.Size(h))
		log <- fmt.Sprintf("Header size: %d", headerSize)
		// Byte slice for moving to buffer
		buffer := make([]byte, headerSize)
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
			log <- fmt.Sprintf("%b", headerBuffer.Bytes())
			// Check to see if we have the whole header
			if len(headerBuffer.Bytes()) == headerSize {
				h.FromBytes(headerBuffer.Bytes())
				log <- fmt.Sprintf("%d", h.Magic)
				log <- fmt.Sprintf("%d", h.Length)
				return h
			}
		}
	}
	// Should never end here
	panic("RECV HEADER")
	return h
}

func recvPayload(conn net.Conn, h Header) (Frame, error) {
	var frame Frame
	payload := make([]byte, h.Length)
	var payloadBuffer bytes.Buffer
	// Make sure we're expecting atleast one byte.
	if h.Length < 1 {
		return frame, errors.New("Length < 1")
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
