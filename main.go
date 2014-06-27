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
      // Second go routine so we can accept more connections
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
  // ret val
  var h quibit.Header
  // a buffer for decoing
  var headerBuffer bytes.Buffer
  for {
    headerSize := int(reflect.TypeOf(h).Size())
    Log <- fmt.Sprintf("Header size: %d", headerSize)
    // Byte slice for moving to buffer
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
      // Add to header buffer
      headerBuffer.Write(buffer)
      Log <- fmt.Sprintf("%b", headerBuffer.Bytes())
      // Check to see if we have the whole header
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
  // Make sure we're expecting atleast one byte.
  if h.Length < 1 {
    return payload, errors.New("Length < 1")
  }
  for {
    // store in byte array
    n, err := conn.Read(payload)
    if err != nil {
      return payload, err
    }
    if n > 1 {
      // write to buffer
      payloadBuffer.Write(payload)
      // Check to see if we have whole payload
      if len(payloadBuffer.Bytes()) == int(h.Length) {
        // Verify checksum
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
