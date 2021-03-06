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
	"bytes"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"net"
	"time"
	"io"
)

const (
	MAGIC = 6667787
)

func recvAll(conn net.Conn, log chan string) (Frame, error) {
	// ret val
	var h Header
	var t time.Time
	var frame Frame
	// a buffer for decoing
	var headerBuffer bytes.Buffer
	for {
		headerSize := int(binary.Size(h))
		// Byte slice for moving to buffer
		buffer := make([]byte, headerSize)
		if conn == nil {
			return frame, errors.New("Nil connection")
		}
		conn.SetReadDeadline(t)
		n, err := io.ReadFull(conn, buffer)
		if err != nil {
			log <- err.Error()
			return frame, err
		}
		if n > 0 {
			// Add to header buffer
			headerBuffer.Write(buffer)
			// Check to see if we have the whole header
			if len(headerBuffer.Bytes()) != headerSize {
				return frame, errors.New("Incorrect header size...")
			}
			h.FromBytes(headerBuffer.Bytes())
			if h.Magic != MAGIC {
				return frame, errors.New("Incorrect Magic Number!")
			}
			frame.Header = h
			break
		}
	}

	payload := make([]byte, h.Length)
	var payloadBuffer bytes.Buffer
	if h.Length < 1 {
		frame.Payload = nil
		frame.Header = h
		return frame, nil
	}
	for {
		// store in byte array
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		n, err := io.ReadFull(conn, payload)
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
