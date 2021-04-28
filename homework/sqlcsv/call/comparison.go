package call

import (
	"gb_go_best_practics/homework/sqlcsv/parser"
	"strings"
)

type ComparisonResult int

const (
	EQUAL ComparisonResult = iota
	NOT_EQUAL
	LESS
	GREATER
	INCONSISTENT
)

var comparisonResultLiterals = map[ComparisonResult]string{
	EQUAL:        "EQUAL",
	NOT_EQUAL:    "NOT_EQUAL",
	LESS:         "LESS",
	GREATER:      "GREATER",
	INCONSISTENT: "INCONSISTENT",
}

func (cr ComparisonResult) String() string {
	return comparisonResultLiterals[cr]
}

func CompareCombinedly(p1 parser.Prime, p2 parser.Prime) ComparisonResult {
	if parser.IsNull(p1) || parser.IsNull(p2) {
		return INCONSISTENT
	}

	if f1 := parser.ConvPrimeToFloat(p1); !parser.IsNull(f1) {
		if f2 := parser.ConvPrimeToFloat(p2); !parser.IsNull(f2) {
			v1 := f1.(parser.Float).Value()
			v2 := f2.(parser.Float).Value()
			if v1 == v2 {
				return EQUAL
			} else if v1 < v2 {
				return LESS
			} else {
				return GREATER
			}
		}
	}

	if d1 := parser.ConvPrimeToDatetime(p1); !parser.IsNull(d1) {
		if d2 := parser.ConvPrimeToDatetime(p2); !parser.IsNull(d2) {
			v1 := d1.(parser.Datetime).Value()
			v2 := d2.(parser.Datetime).Value()
			if v1.Equal(v2) {
				return EQUAL
			} else if v1.Before(v2) {
				return LESS
			} else {
				return GREATER
			}
		}
	}

	if b1 := parser.ConvPrimeToBool(p1); !parser.IsNull(b1) {
		// `if` block ends with a `return` statement, so drop this `else` and outdent its block (golint)
		if b2 := parser.ConvPrimeToBool(p2); !parser.IsNull(b2) {
			v1 := b1.(parser.Boolean).Logic()
			v2 := b2.(parser.Boolean).Logic()
			if v1 == v2 {
				return EQUAL
			}
			return NOT_EQUAL
		}
	}

	if s1, ok := p1.(parser.String); ok {
		if s2, ok := p2.(parser.String); ok {
			v1 := strings.ToUpper(s1.String())
			v2 := strings.ToUpper(s2.String())

			if v1 == v2 {
				return EQUAL
			} else if v1 < v2 {
				return LESS
			} else {
				return GREATER
			}
		}
	}

	return INCONSISTENT
}

func EqualTo(p1 parser.Prime, p2 parser.Prime) bool {
	if r := CompareCombinedly(p1, p2); r != INCONSISTENT {
		if r == EQUAL {
			return true
		}
	}
	return false
}

func NotEqualTo(p1 parser.Prime, p2 parser.Prime) bool {
	if r := CompareCombinedly(p1, p2); r != INCONSISTENT {
		if r != EQUAL {
			return true
		}
	}
	return false
}

func LessThan(p1 parser.Prime, p2 parser.Prime) bool {
	if r := CompareCombinedly(p1, p2); r != INCONSISTENT && r != NOT_EQUAL {
		if r == LESS {
			return true
		}
	}
	return false
}

func GreaterThan(p1 parser.Prime, p2 parser.Prime) bool {
	if r := CompareCombinedly(p1, p2); r != INCONSISTENT && r != NOT_EQUAL {
		if r == GREATER {
			return true
		}

	}
	return false
}

func LessThanOrEqualTo(p1 parser.Prime, p2 parser.Prime) bool {
	if r := CompareCombinedly(p1, p2); r != INCONSISTENT && r != NOT_EQUAL {
		if r != GREATER {
			return true
		}
	}
	return false
}

func GreaterThanOrEqualTo(p1 parser.Prime, p2 parser.Prime) bool {
	if r := CompareCombinedly(p1, p2); r != INCONSISTENT && r != NOT_EQUAL {
		if r != LESS {
			return true
		}
	}
	return false
}

func Compare(p1 parser.Prime, p2 parser.Prime, operator string) bool {
	switch operator {
	case "=":
		return EqualTo(p1, p2)
	case ">":
		return GreaterThan(p1, p2)
	case "<":
		return LessThan(p1, p2)
	case ">=":
		return GreaterThanOrEqualTo(p1, p2)
	case "<=":
		return LessThanOrEqualTo(p1, p2)
	default:
		return NotEqualTo(p1, p2)
	}
}
