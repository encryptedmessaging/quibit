/**
    Copyright 2014 JARST, LLC
    
    This file is part of Quibit.

    Quibit is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    LICENSE file for details.
**/

package quibit

// Error type for Quibit-Specific Errors
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
