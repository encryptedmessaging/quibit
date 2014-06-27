package quibit

import (
	"bytes"
	"crypto/sha512"
	"encoding/binary"
)

type Header struct {
	Magic    uint32
	Command  uint8
	Type     uint8
	Checksum [48]byte
	Length   uint32
}

func (h *Header) Configure(data []byte) {
	h.Magic = MAGIC
	h.Checksum = sha512.Sum384(data)
	h.Length = uint32(len(data))
}

func (h *Header) ToBytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, h)
	return buf.Bytes(), err
}

func (h *Header) FromBytes(b []byte) error {
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, h)
	return err
}
