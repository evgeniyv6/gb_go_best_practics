package call

import (
	"errors"
	"fmt"
	"gb_go_best_practics/homework/sqlcsv/parser"
)

type FilterRecord struct {
	View        *View
	RecordIndex int
}

type Filter []FilterRecord

func (f Filter) Evaluate(expr parser.Expression) (parser.Prime, error) {
	var primary parser.Prime
	var err error

	if parser.IsPrime(expr) {
		primary = expr.(parser.Prime)
	} else {
		switch expr.(type) {
		case parser.Identifier:
			primary, err = f.evalIdentifier(expr.(parser.Identifier))
		case parser.Comparison:
			primary, err = f.evalComparison(expr.(parser.Comparison))
		case parser.Logic:
			primary, err = f.evalLogic(expr.(parser.Logic))
		default:
			//  should replace `errors.New(fmt.Sprintf(...))` with `fmt.Errorf(...)` (golint)
			return nil, fmt.Errorf("syntax error: %s", expr)
		}
	}

	return primary, err
}

func (f Filter) evalIdentifier(expr parser.Identifier) (parser.Prime, error) {
	ref, field, err := expr.Field()
	if err != nil {
		return nil, err
	}
	var p parser.Prime
	for _, v := range f {
		idx, err := v.View.Header.Contains(ref, field)
		if err != nil {
			return nil, errors.New("rec is unrecognized")
		}
		if p != nil {
			return nil, errors.New("rec is unrecognized")
		}
		p = v.View.Records[v.RecordIndex][idx].Zero()
	}
	if p == nil {
		return nil, errors.New("rec not exist")
	}
	return p, nil
}

func (f Filter) evalComparison(expr parser.Comparison) (parser.Prime, error) {
	lhs, err := f.Evaluate(expr.LHS)
	if err != nil {
		return nil, err
	}
	rhs, err := f.Evaluate(expr.RHS)
	if err != nil {
		return nil, err
	}

	t := Compare(lhs, rhs, expr.Operator.Str)
	return parser.NewLogicOp(t), nil
}

func (f Filter) evalLogic(expr parser.Logic) (parser.Prime, error) {
	lhs, err := f.Evaluate(expr.LHS)
	if err != nil {
		return nil, err
	}
	rhs, err := f.Evaluate(expr.RHS)
	if err != nil {
		return nil, err
	}

	var t bool
	switch expr.Operator.Token {
	case parser.AND:
		switch {
		case lhs.Logic() == false || rhs.Logic() == false:
			t = false
		default:
			t = true
		}
	case parser.OR:
		switch {
		case lhs.Logic() == true || rhs.Logic() == true:
			t = true
		default:
			t = false
		}
	}
	return parser.NewLogicOp(t), nil
}
