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
  Magic       int32
  PayloadLen  int32
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
  fmt.Sprint(h)
  return err
}
