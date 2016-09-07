package unit

import (
	"log"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func (p *Parser) div() {
	d, err := p.popUnit()
	if err != nil {
		p.Err = err
		return
	}

	n, err := p.popUnit()
	if err != nil {
		p.Err = err
		return
	}

	p.pushUnit(&Div{N: n, D: d})
}

func (p *Parser) expParens(expstr string) {
	exp, err := p.parseExp(expstr)
	if err != nil {
		p.Err = errors.Wrap(err, "parse exponent")
		return
	}

	source, err := p.popUnit()
	if err != nil {
		p.Err = errors.Wrap(err, "pop source")
		return
	}

	p.pushUnit(&Exp{U: source, Exp: exp})
}

func (p *Parser) expUnit(unitexpstr string) {
	pieces := strings.SplitN(unitexpstr, "^", 2)
	if len(pieces) != 2 {
		panic("invalid PushExpUnit call: input doest contain '^'")
	}

	unitstr, expstr := pieces[0], pieces[1]
	source := &Ref{Name: unitstr}

	exp, err := p.parseExp(expstr)
	if err != nil {
		p.Err = errors.Wrap(err, "parse exponent")
		return
	}

	p.pushUnit(&Exp{U: source, Exp: exp})
}

func (p *Parser) finish() {
	if len(p.stack) != 1 {
		p.Err = &ParseError{Input: p.Buffer, FailurePhase: "stack drain"}
		return
	}

	p.Result = p.stack[0]
}

func (p *Parser) invertUnit() {
	d, err := p.popUnit()
	if err != nil {
		p.Err = err
		return
	}

	p.pushUnit(&Div{N: &Nil{}, D: d})
}

func (p *Parser) mul() {
	last, err := p.popUnit()
	if err != nil {
		p.Err = err
		return
	}

	first, err := p.popUnit()
	if err != nil {
		p.Err = err
		return
	}

	p.pushUnit(&Mul{first, last})
}

func (p *Parser) log(args ...interface{}) {
	log.Print(args...)
}

func (p *Parser) logf(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}
func (p *Parser) parseExp(expstr string) (int64, error) {
	exp, err := strconv.ParseInt(expstr, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "parse exponent")
	}

	return exp, nil
}

func (p *Parser) popUnit() (U, error) {
	if len(p.stack) == 0 {
		return nil, errors.New("underflow")
	}

	ret := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return ret, nil
}

func (p *Parser) pushUnit(u U) {
	p.stack = append(p.stack, u)
}

func (p *Parser) ref(unitstr string) {
	source := &Ref{Name: unitstr}
	p.pushUnit(source)
}
