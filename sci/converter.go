package sci

import "math/big"

func (c *converter) ToBase(mag string) string {
	var (
		in  big.Float
		out big.Float
	)
	_, ok := in.SetString(mag)
	if !ok {
		panic("invalid magnitude: " + mag)
	}

	out.Mul(&in, &c.factor)
	return out.String()
}

func (c *converter) FromBase(mag string) string {
	var (
		in  big.Float
		out big.Float
	)
	_, ok := in.SetString(mag)
	if !ok {
		panic("invalid magnitude: " + mag)
	}

	out.Quo(&in, &c.factor)
	return out.String()
}
