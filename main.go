package main

import (
  "fmt"
  "net"
  "quibit"
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
        header := quibit.RecvHeader(conn, Log)
        payload, err := quibit.RecvPayload(conn, header)
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

