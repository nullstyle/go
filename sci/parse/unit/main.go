package unit

import "strings"

//go:generate peg -switch -inline parser.peg

type U interface {
	_ast()
}

type Div struct {
	N U
	D U
}

type Exp struct {
	U   U
	Exp int64
}

// ParseError represents the error produces when trying to operate on a
// value whose magnitude (the M field) is invalid.
type ParseError struct {
	Input        string
	FailurePhase string
}

type Ref struct {
	Name string
}

type Mul []U

type Nil struct{}

// Parse parses the provided input into the AST for a unit
func Parse(input string) (U, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return &Nil{}, nil
	}

	p := &Parser{Buffer: input}
	p.Init()

	err := p.Parse()
	if err != nil {
		return nil, err
	}

	p.Execute()
	return p.Result, nil
}

var _ U = &Ref{}
var _ U = &Nil{}
var _ U = &Div{}
var _ U = &Mul{}
var _ U = &Exp{}

func (u *Ref) _ast() {}
func (u *Nil) _ast() {}
func (u *Div) _ast() {}
func (u *Mul) _ast() {}
func (u *Exp) _ast() {}
