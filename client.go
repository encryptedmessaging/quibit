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
  data := []byte("Hello, World!")
  var frame quibit.Frame
  frame.Configure(data, quibit.HELO)
  //Write the header to buffer
  buf,err := frame.Header.ToBytes()
  fmt.Println(buf)
  conn.Write(buf)
  //Write data to buffer
  conn.Write(frame.Payload)

}
