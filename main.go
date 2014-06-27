package main

import (
  "fmt"
  "bytes"
  "net"
  "reflect"
  "quibit"
  "crypto/sha512"
)


func main() {
  listener, err := net.Listen("tcp", ":1337")
  if err != nil {
    fmt.Println(err.Error())
  }
  log := make(chan string)
  go func() {
    for {
      conn, err := listener.Accept()
      if err != nil {
        log <- err.Error()
      }
      log <- "Connection Accepted!"
      go func() {
        // So now we have a connection.  Let's shake hands.
        log <- fmt.Sprintf("%v", conn)
        var h quibit.Header
        var headerBuffer bytes.Buffer
        for {
          headerSize := int(reflect.TypeOf(h).Size())
          log <- fmt.Sprintf("Header size: %d", headerSize)
          buffer := make([]byte, headerSize)
          n, err := conn.Read(buffer)
          if err != nil {
            if err.Error() == "EOF" {
              break
            }
            log <- err.Error()
          }
          if n > 0 {
            fmt.Println(buffer)
            headerBuffer.Write(buffer)
            log <- fmt.Sprintf("%b", headerBuffer.Bytes())
            if len(headerBuffer.Bytes()) == headerSize {
              h.FromBytes(headerBuffer.Bytes())
              log <- fmt.Sprintf("%d", h.Magic)
              log <- fmt.Sprintf("%d", h.Length)
              var payloadBuffer bytes.Buffer
              if h.Length > 0 {
                for {
                  log <- fmt.Sprintf("Payload size: %d", h.Length)
                  buffer := make([]byte, h.Length)
                  n, err := conn.Read(buffer)
                  if err != nil {
                    log <- err.Error()
                  }
                  if n > 0 {
                    payloadBuffer.Write(buffer)
                    if len(payloadBuffer.Bytes()) == int(h.Length) {
                      //Checksum
                      if h.Checksum != sha512.Sum384(payloadBuffer.Bytes()) {
                        log <- "Incorrect Checksum!"
                      }
                      log <- fmt.Sprintf("Payload: %s", payloadBuffer.Bytes())
                      break
                    }
                  }
                }
              }
            }
          }
        }
      }()
    }
  }()
  for {
    message := <- log
    fmt.Println(message)
  }
}


