package main

import (
  "fmt"
  "bytes"
  "net"
  "reflect"
  "quibit"
  "crypto/sha512"
  "errors"
)


var Log = make(chan string)

func main() {
  listener, err := net.Listen("tcp", ":1337")
  if err != nil {
    fmt.Println(err.Error())
  }
  go func() {
    for {
      conn, err := listener.Accept()
      if err != nil {
        Log <- err.Error()
      }
      Log <- "Connection Accepted!"
      go func() {
        // So now we have a connection.  Let's shake hands.
        Log <- fmt.Sprintf("%v", conn)
        header := RecvHeader(conn)
        payload, err := RecvPayload(conn, header)
        if err != nil {
          Log <- fmt.Sprint("%s", err.Error())
        }
        Log <- fmt.Sprintf("%s", payload)
      }()
    }
  }()
  for {
    message := <- Log
    fmt.Println(message)
  }
}

func RecvHeader(conn net.Conn) quibit.Header {
  var h quibit.Header
  var headerBuffer bytes.Buffer
  for {
    headerSize := int(reflect.TypeOf(h).Size())
    Log <- fmt.Sprintf("Header size: %d", headerSize)
    buffer := make([]byte, headerSize)
    n, err := conn.Read(buffer)
    if err != nil {
      if err.Error() == "EOF" {
        break
      }
      Log <- err.Error()
    }
    if n > 0 {
      fmt.Println(buffer)
      headerBuffer.Write(buffer)
      Log <- fmt.Sprintf("%b", headerBuffer.Bytes())
      if len(headerBuffer.Bytes()) == headerSize {
        h.FromBytes(headerBuffer.Bytes())
        Log <- fmt.Sprintf("%d", h.Magic)
        Log <- fmt.Sprintf("%d", h.Length)
        return h
      }
    }
  }
  // Should never end here
  panic("RECV HEADER")
  return h
}

func RecvPayload(conn net.Conn, h quibit.Header) ([]byte, error) {
  payload := make([]byte, h.Length)
  var payloadBuffer bytes.Buffer
  if h.Length < 1 {
    return payload, errors.New("Length < 1")
  }
  for {
    n, err := conn.Read(payload)
    if err != nil {
      return payload, err
    }
    if n > 1 {
      payloadBuffer.Write(payload)
      if len(payloadBuffer.Bytes()) == int(h.Length) {
        if h.Checksum != sha512.Sum384(payloadBuffer.Bytes()) {
          return payloadBuffer.Bytes(), errors.New("Incorrect Checksum")
        }
        return payloadBuffer.Bytes(), nil
      }
    }
  }
  //Should never end here
  panic("RECV PAYLOAD")
  return payload, nil
}
