// Package m implements functions for working with magnitudes and multitudes
package m

import (
	"math/big"
)

var (
	// One is 1.0
	One big.Float

	// Zero is 0.0
	Zero big.Float
)

// M represents either a magnitude or a multitude, encoded as a string of
// decimal digits.
type M string

func Singular(m M) {
	var mf big.Float
	_, ok := mf.SetString(m)
}

func init() {
	_, ok := One.SetString("1.0")
	if !ok {
		panic("couldn't initialize One")
	}

	_, ok = Zero.SetString("0.0")
	if !ok {
		panic("couldn't initialize Zero")
	}
}
