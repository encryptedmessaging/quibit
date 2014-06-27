package quibit

import (
  "bytes"
  "encoding/binary"
  "fmt"
)

const (
  MAGIC = 0x420
)

type Header struct {
  Magic       uint32
  Command     uint8
  Type        uint8
  Checksum    uint32
  Length      uint32
}

func (h *Header) ToBytes() ([]byte, error) {
  fmt.Println(h)
  buf := new(bytes.Buffer)
  err := binary.Write(buf, binary.LittleEndian, h)
  return buf.Bytes(), err
}

func (h *Header) FromBytes(b []byte) error {
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, h)
  return err
}
