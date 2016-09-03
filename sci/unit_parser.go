package sci

import "errors"

func (p *UnitParser) PopUnit() (Unit, error) {
	if len(p.Stack) == 0 {
		return nil, errors.New("underflow")
	}

	ret := p.Stack[len(p.Stack)-1]
	p.Stack = p.Stack[:len(p.Stack)-1]
	return ret, nil
}

func (p *UnitParser) PushUnit(u Unit) {
	p.Stack = append(p.Stack, u)
}
