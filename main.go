package main

import (
  "quibit"
  "fmt"
)

var FrameChannel chan quibit.Frame

var TChan chan string
var LogChannel chan string

func main() {
  go func() {
    msg := <- TChan
    fmt.Println("A")
    fmt.Println(msg)
  }()
  go func() {
    msg := <- TChan
     fmt.Println("B")
    fmt.Println(msg)
  }()
  go func() {
    msg := <- TChan
    fmt.Println("C")
    fmt.Println(msg)
  }()
  TChan <- "KNC Sucks"
  for {

  }
}
