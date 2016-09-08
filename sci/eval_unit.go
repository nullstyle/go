package sci

import (
	"github.com/nullstyle/go/sci/parse/unit"
	"github.com/pkg/errors"
)

// eval is the simple stack machine that evaluates the AST for a unit into a
// fully-fledged unit.
type evalUnit struct {
	System *System
	stack  []Unit
}

func (eu *evalUnit) eval(input unit.U) (Unit, error) {
	eu.stack = nil
	err := eu.visit(input)
	if err != nil {
		return nil, errors.Wrap(err, "visit failed")
	}

	if len(eu.stack) != 1 {
		return nil, errors.New("unconsumed stack")
	}

	return eu.stack[0], nil
}

func (eu *evalUnit) pop() (Unit, error) {
	if len(eu.stack) == 0 {
		return nil, errors.New("underflow")
	}

	ret := eu.stack[len(eu.stack)-1]
	eu.stack = eu.stack[:len(eu.stack)-1]
	return ret, nil
}

func (eu *evalUnit) push(u Unit) {
	eu.stack = append(eu.stack, u)
}

func (eu *evalUnit) visit(cur unit.U) error {
	switch cur := cur.(type) {
	case *unit.Nil:
		eu.push(eu.System.Nil())
	case *unit.Ref:
		found, err := eu.System.LookupUnit(cur.Name)
		if err != nil {
			return errors.Wrap(err, "lookup failed")
		}
		eu.push(found)
	case *unit.Exp:
		err := eu.visit(cur.U)
		if err != nil {
			return errors.Wrap(err, "visit child failed")
		}

		source, err := eu.pop()
		if err != nil {
			return errors.Wrap(err, "pop failed")
		}

		if cur.Abs() > MaxExp {
			return &ExpToBigError{Exp: cur.Exp}
		}

		mul := make(MulUnit, cur.Abs())
		for i := range mul {
			mul[i] = source
		}

		var ret Unit
		if cur.Exp < 0 {
			ret = &DivUnit{N: eu.System.Nil(), D: &mul}
		} else {
			ret = &mul
		}

		eu.push(ret)
	case *unit.Mul:
		for _, child := range *cur {
			err := eu.visit(child)
			if err != nil {
				return errors.Wrap(err, "visit child failed")
			}
		}

		ret := make(MulUnit, len(*cur))
		for i := len(ret) - 1; i >= 0; i-- {
			u, err := eu.pop()
			if err != nil {
				return errors.Wrap(err, "pop failed")
			}
			ret[i] = u
		}
		eu.push(&ret)
	case *unit.Div:
		err := eu.visit(cur.N)
		if err != nil {
			return errors.Wrap(err, "visit child N failed")
		}

		err = eu.visit(cur.D)
		if err != nil {
			return errors.Wrap(err, "visit child D failed")
		}

		d, err := eu.pop()
		if err != nil {
			return errors.Wrap(err, "pop D failed")
		}

		n, err := eu.pop()
		if err != nil {
			return errors.Wrap(err, "pop N failed")
		}

		eu.push(&DivUnit{N: n, D: d})
	default:
		return errors.Errorf("Unknown unit: %s", cur)
	}

	return nil
}
