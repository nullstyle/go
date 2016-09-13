package sci

import (
	"fmt"
	"math/big"

	"github.com/pkg/errors"
)

// Add adds `l` and `r` together and stores the result in `v`, providing that
// `l` and `r` can be added together.
func (v *Value) Add(l, r *Value) error {
	// TODO: compare the units for compatibility

	var vf, lf, rf big.Float

	_, ok := lf.SetString(l.M)
	if !ok {
		return &MagnitudeError{M: l.M}
	}

	_, ok = rf.SetString(r.M)
	if !ok {
		return &MagnitudeError{M: r.M}
	}

	vf.Add(&lf, &rf)

	// BUG(scott): we don't have unit conversions yet, so the following adoption
	// of l.U for the return value is broken.
	v.M = vf.String()
	v.U = l.U
	return nil
}

// Convert takes the source value, converts it to the target unit and sets the
// result on the method receiver
func (v *Value) Convert(source *Value, target Unit) error {
	sys := source.U.System()
	sc, err := sys.getConverter(source.U)
	if err != nil {
		return errors.Wrap(err, "failed to get source converter")
	}

	tc, err := sys.getConverter(target)
	if err != nil {
		return errors.Wrap(err, "failed to get target converter")
	}

	nm := sc.ToBase(source.M)

	var result Value
	result.U = target
	result.M = tc.FromBase(nm)
	*v = result
	return nil
}

// Div divides `l` and `r` together and stores the result in `v`.
func (v *Value) Div(l, r *Value) error {
	return nil
}

// Eq checks v and other for equality
func (v *Value) Eq(other *Value) bool {
	return false
}

// Mul multiplies `l` and `r` together and stores the result in `v`.
func (v *Value) Mul(l, r *Value) error {
	return nil
}

// String implements fmt.Stringer
func (v *Value) String() string {
  var (
    one big.Float
    m big.Float
  )

  
  


	return fmt.Sprintf("%s %s", v.M, v.U)
}
