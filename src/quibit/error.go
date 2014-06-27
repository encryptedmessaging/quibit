package quibit

type QuibitError int

const (
	eNILOBJ = iota	
	eHEADER = iota
)

func (e QuibitError) Error() string {
	switch int(e) {
	case eNILOBJ:
		return "Received unexpected nil object."
	case eHEADER:
		return "Malformed header, could not serialize"
	default:
		return "Unknown Quibit Error!"
	}
}