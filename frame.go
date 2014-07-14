package quibit

type Frame struct {
	Peer    string // Peer who sent or will receive this frame
	Header  Header // Header associated with frame
	Payload []byte // Serialize frame payload
}

// Configure Frame f with a proper header for
// payload "data" interpreted as "command"
func (f *Frame) Configure(data []byte, command, t uint8) {
	var h Header
	h.Configure(data)
	h.Command = command
	h.Type = t
	f.Header = h
	f.Payload = data
}
