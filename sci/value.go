package sci

import (
	"math/big"
)

// Add adds `l` and `r` together and stores the result in `v`, providing that
// `l` and `r` can be added together.
func (v *Value) Add(l, r *Value) error {
  // TODO: compare the units for compatibility
  if !l.U.Compat(r.U) {
    return 	ErrIncompatibleTypes
  }

  var vf, lf, rf big.Float{}
  _, ok := lf.SetString(l.M)

	return nil
}

// Div divides `l` and `r` together and stores the result in `v`.
func (v *Value) Div(l, r *Value) error {
	return nil
}

// Mul multiplies `l` and `r` together and stores the result in `v`.
func (v *Value) Mul(l, r *Value) error {
	return nil
}
