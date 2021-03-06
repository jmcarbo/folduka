// this file was generated by gomacro command: import _b "math/bits"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package imports

import (
	. "reflect"
	bits "math/bits"
)

// reflection: allow interpreted code to import "math/bits"
func init() {
	Packages["math/bits"] = Package{
	Name: "bits",
	Binds: map[string]Value{
		"Add":	ValueOf(bits.Add),
		"Add32":	ValueOf(bits.Add32),
		"Add64":	ValueOf(bits.Add64),
		"Div":	ValueOf(bits.Div),
		"Div32":	ValueOf(bits.Div32),
		"Div64":	ValueOf(bits.Div64),
		"LeadingZeros":	ValueOf(bits.LeadingZeros),
		"LeadingZeros16":	ValueOf(bits.LeadingZeros16),
		"LeadingZeros32":	ValueOf(bits.LeadingZeros32),
		"LeadingZeros64":	ValueOf(bits.LeadingZeros64),
		"LeadingZeros8":	ValueOf(bits.LeadingZeros8),
		"Len":	ValueOf(bits.Len),
		"Len16":	ValueOf(bits.Len16),
		"Len32":	ValueOf(bits.Len32),
		"Len64":	ValueOf(bits.Len64),
		"Len8":	ValueOf(bits.Len8),
		"Mul":	ValueOf(bits.Mul),
		"Mul32":	ValueOf(bits.Mul32),
		"Mul64":	ValueOf(bits.Mul64),
		"OnesCount":	ValueOf(bits.OnesCount),
		"OnesCount16":	ValueOf(bits.OnesCount16),
		"OnesCount32":	ValueOf(bits.OnesCount32),
		"OnesCount64":	ValueOf(bits.OnesCount64),
		"OnesCount8":	ValueOf(bits.OnesCount8),
		"Reverse":	ValueOf(bits.Reverse),
		"Reverse16":	ValueOf(bits.Reverse16),
		"Reverse32":	ValueOf(bits.Reverse32),
		"Reverse64":	ValueOf(bits.Reverse64),
		"Reverse8":	ValueOf(bits.Reverse8),
		"ReverseBytes":	ValueOf(bits.ReverseBytes),
		"ReverseBytes16":	ValueOf(bits.ReverseBytes16),
		"ReverseBytes32":	ValueOf(bits.ReverseBytes32),
		"ReverseBytes64":	ValueOf(bits.ReverseBytes64),
		"RotateLeft":	ValueOf(bits.RotateLeft),
		"RotateLeft16":	ValueOf(bits.RotateLeft16),
		"RotateLeft32":	ValueOf(bits.RotateLeft32),
		"RotateLeft64":	ValueOf(bits.RotateLeft64),
		"RotateLeft8":	ValueOf(bits.RotateLeft8),
		"Sub":	ValueOf(bits.Sub),
		"Sub32":	ValueOf(bits.Sub32),
		"Sub64":	ValueOf(bits.Sub64),
		"TrailingZeros":	ValueOf(bits.TrailingZeros),
		"TrailingZeros16":	ValueOf(bits.TrailingZeros16),
		"TrailingZeros32":	ValueOf(bits.TrailingZeros32),
		"TrailingZeros64":	ValueOf(bits.TrailingZeros64),
		"TrailingZeros8":	ValueOf(bits.TrailingZeros8),
		"UintSize":	ValueOf(bits.UintSize),
	}, Untypeds: map[string]string{
		"UintSize":	"int:64",
	}, 
	}
}
