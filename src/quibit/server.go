package quibit

func NewServer(recvChan chan Frame, port string) error {
  listener, err := net.Listen("tcp", port)
  if err != nil {
    return err
  }
  go func(conn net.Conn, recvChan chan Frame) {
    for {
      conn, err := listener.Accept()
      if err != nil {
        fmt.Println(err.Error())
      }
      listenForFrame(conn, revcChan)
      return nil
    }
  }(conn, recvChan)
  return nil
}
