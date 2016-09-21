package calc

import (
	"context"

	"github.com/nullstyle/go/influx"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// ---- operators -----

// Add is an action that sets the queued operator to addition
type Add struct{}

func (o Add) operate(l, r decimal.Decimal) decimal.Decimal {
	return l.Add(r)
}

// Sub is an action that sets the queued operator to subtraction
type Sub struct{}

func (o Sub) operate(l, r decimal.Decimal) decimal.Decimal {
	return l.Sub(r)
}

// Mul is an action that sets the queued operator to multiplication
type Mul struct{}

func (o Mul) operate(l, r decimal.Decimal) decimal.Decimal {
	return l.Mul(r)
}

// Div is an action that sets the queued operator to division
type Div struct{}

func (o Div) operate(l, r decimal.Decimal) decimal.Decimal {
	return l.Div(r)
}

var _ Operator = Add{}
var _ Operator = Sub{}
var _ Operator = Div{}
var _ Operator = Mul{}

// ---- end operators -----

// ---- digits ----

// Digit represents a digit button on the calculator. For the purposes of this
// example application, it also shows how you can use vars (see the declarations
// below) to provide enum-like actions.
type Digit struct {
	digit string
}

var (
	// Zero represents the action of pressing the 0 key
	Zero = Digit{"0"}
	// One represents the action of pressing the 1 key
	One = Digit{"1"}
	// Two represents the action of pressing the 2 key
	Two = Digit{"2"}
	// Three represents the action of pressing the 3 key
	Three = Digit{"3"}
	// Four represents the action of pressing the 4 key
	Four = Digit{"4"}
	// Five represents the action of pressing the 5 key
	Five = Digit{"5"}
	// Six represents the action of pressing the 6 key
	Six = Digit{"6"}
	// Seven represents the action of pressing the 7 key
	Seven = Digit{"7"}
	// Eight represents the action of pressing the 8 key
	Eight = Digit{"8"}
	// Nine represents the action of pressing the 9 key
	Nine = Digit{"9"}
)

// ---- end digits -----

// Clear is an action that clears the current number
type Clear struct{}

// Equals is an action that causes the calculator to calculate an answer by
// taking evaluating the saved operation using the current and the last saved
// number as operands.
type Equals struct{}

// PlusMinus is an action that causes the current number's sign to flip.
type PlusMinus struct{}

// Calculator represents the state of a handheld calculator that has 0-9 digit
// keys, +-*/ operators, a +/- key, a clear key, and an = key.  It is the root
// state of this example app.
type Calculator struct {
	QueuedOperator Operator
	CurrentNumber  string
	SavedNumber    string
	ShowingResult  bool
}

// HandleAction implements influx.Handler.
func (c *Calculator) HandleAction(
	ctx context.Context,
	action influx.Action,
) error {

	switch action := action.(type) {
	case Clear:
		c.CurrentNumber = ""
		c.ShowingResult = false
	case Digit:
		if c.ShowingResult {
			c.ShowingResult = false
			c.SavedNumber = c.CurrentNumber
			c.CurrentNumber = ""
		}

		c.CurrentNumber = c.CurrentNumber + action.digit
	case Equals:
		if c.QueuedOperator == nil {
			return nil
		}

		r, err := decimal.NewFromString(c.CurrentNumber)
		if err != nil {
			return errors.Wrap(err, "parse right failed")
		}

		l, err := decimal.NewFromString(c.SavedNumber)
		if err != nil {
			return errors.Wrap(err, "parse left failed")
		}
		ret := c.QueuedOperator.operate(l, r)

		c.CurrentNumber = ret.String()
		c.ShowingResult = true
	case Operator:
		c.QueuedOperator = action
		c.SavedNumber = c.CurrentNumber
		c.CurrentNumber = ""
		c.ShowingResult = false
	case PlusMinus:
		cur, err := decimal.NewFromString(c.CurrentNumber)
		if err != nil {
			return errors.Wrap(err, "parse current failed")
		}

		result := cur.Mul(negOne)
		c.CurrentNumber = result.String()
		c.ShowingResult = false
	}

	return nil
}

// Operator represents a type that can perform arithmetic given two values.
type Operator interface {
	operate(l decimal.Decimal, r decimal.Decimal) decimal.Decimal
}

var negOne = decimal.NewFromFloat(-1.0)
