package sci

import (
	"log"
	"strconv"

	"strings"

	"github.com/pkg/errors"
)

// Exp exponentiates the top of the unit stack, pushing a MulUnit of N units of
// the present stack top.
func (p *UnitParser) pushExpUnit(unitexpstr string) error {
	pieces := strings.SplitN(unitexpstr, "^", 2)
	if len(pieces) != 2 {
		panic("invalid PushExpUnit call: input doest contain '^'")
	}

	unitstr, expstr := pieces[0], pieces[1]

	exp, err := strconv.ParseInt(expstr, 10, 64)
	if err != nil {
		return errors.Wrap(err, "parse exponent")
	}

	abs := exp
	if exp < 0 {
		abs = -exp
	}

	if abs > MaxExp {
		return &ExpToBigError{Exp: exp}
	}

	source, err := p.System.LookupUnit(unitstr)
	if err != nil {
		return errors.Wrap(err, "lookup source")
	}

	expu := make(MulUnit, abs)
	for i := range expu {
		expu[i] = source
	}

	// if the exponent is negative, push the inverse of expu
	if exp < 0 {
		p.pushUnit(&DivUnit{
			N: Nil,
			D: &expu,
		})
		return nil
	}

	// otherwise, push expu
	p.pushUnit(&expu)
	return nil
}

func (p *UnitParser) log(args ...interface{}) {
	log.Print(args...)
}

func (p *UnitParser) logf(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}

func (p *UnitParser) popUnit() (Unit, error) {
	if len(p.stack) == 0 {
		return nil, errors.New("underflow")
	}

	ret := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return ret, nil
}

func (p *UnitParser) pushUnit(u Unit) {
	p.stack = append(p.stack, u)
}
