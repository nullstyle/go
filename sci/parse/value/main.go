package value

import (
	"strings"

	"github.com/pkg/errors"
)

//go:generate peg -switch -inline parser.peg

var (
	// ErrBlankValue is returned when attempting to parse a blank string into a
	// value.
	ErrBlankValue = errors.New("blank value string")
)

// ParseError represents the error produces when trying to operate on a
// value whose magnitude (the M field) is invalid.
type ParseError struct {
	Input        string
	FailurePhase string
}

// Value represents the parsed output of this parser
type V struct {
	M string
	U string
}

// Parse parses the provided input into the AST for a unit
func Parse(input string) (V, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return V{}, ErrBlankValue
	}

	p := &Parser{Buffer: input}
	p.Init()

	err := p.Parse()
	if err != nil {
		return V{}, err
	}

	p.Execute()
	return p.Result, nil
}
