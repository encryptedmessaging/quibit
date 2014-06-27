package main

import (
  "net"
  "quibit"
  "fmt"
  "crypto/sha512"
)


func main() {
  conn, err := net.Dial("tcp", "localhost:1337")
  if err != nil {
    fmt.Println(err.Error())
  }
  data := []byte("Hello, World!")
  var h quibit.Header
  h.Magic = quibit.MAGIC
  h.Length = uint32(len(data))
  h.Command = quibit.HELO
  h.Checksum = sha512.Sum384(data)
  //Write the header to buffer
  buf,err := h.ToBytes()
  fmt.Println(buf)
  conn.Write(buf)
  //Write data to buffer
  conn.Write(data)

}
