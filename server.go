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
    message := <- Log
    fmt.Println(message)
  }
}

