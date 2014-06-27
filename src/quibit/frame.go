package quibit


type Frame struct {
  Header Header
  Payload []byte
}

func (f *Frame) Configure(data []byte, command uint8) {
  var h Header
  h.Configure(data)
  h.Command = command
  f.Header = h
  f.Payload = data
}
