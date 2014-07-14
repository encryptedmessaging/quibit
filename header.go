package quibit

import (
	"bytes"
	"crypto/sha512"
	"encoding/binary"
)

// Used as a fixed-length description of a frame
type Header struct {
	Magic    uint32   // Known magic number
	Command  uint8    // How to interpret payload
	Type     uint8    // How to interpret payload
	Checksum [48]byte // SHA-384 Checksum of Payload
	Length   uint32   // Length of Payload
}

// Configure a new header given the frame payload
func (h *Header) Configure(data []byte) {
	h.Magic = MAGIC
	h.Checksum = sha512.Sum384(data)
	h.Length = uint32(len(data))
}

// Serialize Header for sending over the wire
func (h *Header) ToBytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, h)
	return buf.Bytes(), err
}

// Unserialize header from the wire
func (h *Header) FromBytes(b []byte) error {
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, h)
	return err
}
