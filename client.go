package main

import (
  "net"
  "quibit"
  "fmt"
)


func main() {
  conn, err := net.Dial("tcp", "localhost:1337")
  if err != nil {
    fmt.Println(err.Error())
  }
  var h quibit.Header
  h.Magic = 420
  h.PayloadLen = 0
  buf,err := h.ToBytes()
  fmt.Println(buf)
  conn.Write(buf)
}
