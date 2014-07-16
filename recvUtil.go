package quibit

import (
	"bytes"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"net"
	"fmt"
	"time"
	"io"
)

const (
	MAGIC = 6667787
)

func recvAll(conn net.Conn, log chan string) (Frame, error) {
	fmt.Println("Called recvAll with connection: ", conn.RemoteAddr())
	fmt.Printf("Called recvAll with connection pointer: %p\n", conn)
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
		fmt.Println("Reading header...")
		n, err := io.ReadFull(conn, buffer)
		fmt.Println("Received header with connection: ", conn.RemoteAddr())
		fmt.Printf("Received header with connection pointer: %p\n", conn)
		if err != nil {
			if err.Error() == "EOF" || err.Error() == "use of closed network connection" {
				return frame, err
			}
			log <- err.Error()
			continue
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

	fmt.Println("Payload Length: ", h.Length)
	fmt.Println("Payload from: ", conn.RemoteAddr())
	payload := make([]byte, h.Length)
	var payloadBuffer bytes.Buffer
	if h.Length < 1 {
		frame.Payload = nil
		frame.Header = h
		return frame, nil
	}
	for {
		fmt.Println("In For Loop for...", conn.RemoteAddr())
		// store in byte array
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		n, err := io.ReadFull(conn, payload)
		fmt.Println("Done reading from...", conn.RemoteAddr())
		if err != nil {
			return frame, err
		}
		if n > 1 {
			fmt.Println("Writing payload to buffer from...", conn.RemoteAddr())
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
