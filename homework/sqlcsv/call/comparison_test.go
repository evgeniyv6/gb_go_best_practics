package call

import (
	"gb_go_best_practics/homework/sqlcsv/parser"
	"testing"
)

func TestComparisonResult_String(t *testing.T) {
	if GREATER.String() != "GREATER" {
		t.Errorf("%s, got %s, want %s", GREATER.String(), "GREATER", GREATER)
	}
}

var compareCombinedlyTests = []struct {
	left   parser.Prime
	right  parser.Prime
	result ComparisonResult
}{
	{
		left:   parser.NewInteger(1),
		right:  parser.NewInteger(1),
		result: EQUAL,
	},
	{
		left:   parser.NewBoolean(true),
		right:  parser.NewBoolean(false),
		result: NOT_EQUAL,
	},
	{
		left:   parser.NewInteger(1),
		right:  parser.NewInteger(2),
		result: LESS,
	},
	{
		left:   parser.NewInteger(2),
		right:  parser.NewInteger(1),
		result: GREATER,
	},
	{
		left:   parser.NewInteger(1),
		right:  parser.NewNull(),
		result: INCONSISTENT,
	},
}

func TestCompareCombinedly(t *testing.T) {
	for _, test := range compareCombinedlyTests {
		r := CompareCombinedly(test.left, test.right)
		if r != test.result {
			t.Errorf("%s, compare want %s , left %s right %s", r, test.result, test.left, test.right)
		}
	}
}

var compareTests = []struct {
	left   parser.Prime
	right  parser.Prime
	op     string
	result bool
}{
	{
		left:   parser.NewInteger(1),
		right:  parser.NewInteger(2),
		op:     "=",
		result: false,
	},
	{
		left:   parser.NewInteger(1),
		right:  parser.NewInteger(1),
		op:     "=",
		result: true,
	},
	{
		left:   parser.NewInteger(1),
		right:  parser.NewNull(),
		op:     "=",
		result: false,
	},
	{
		left:   parser.NewInteger(1),
		right:  parser.NewInteger(2),
		op:     ">",
		result: false,
	},
	{
		left:   parser.NewInteger(2),
		right:  parser.NewInteger(1),
		op:     ">",
		result: true,
	},
	{
		left:   parser.NewInteger(1),
		right:  parser.NewNull(),
		op:     ">",
		result: false,
	},
	{
		left:   parser.NewInteger(2),
		right:  parser.NewInteger(1),
		op:     "<",
		result: false,
	},
	{
		left:   parser.NewInteger(1),
		right:  parser.NewInteger(2),
		op:     "<",
		result: true,
	},
	{
		left:   parser.NewInteger(1),
		right:  parser.NewNull(),
		op:     "<",
		result: false,
	},
	{
		left:   parser.NewInteger(1),
		right:  parser.NewInteger(2),
		op:     ">=",
		result: false,
	},
	{
		left:   parser.NewInteger(2),
		right:  parser.NewInteger(2),
		op:     ">=",
		result: true,
	},
	{
		left:   parser.NewInteger(1),
		right:  parser.NewNull(),
		op:     ">=",
		result: false,
	},
	{
		left:   parser.NewInteger(2),
		right:  parser.NewInteger(1),
		op:     "<=",
		result: false,
	},
	{
		left:   parser.NewInteger(2),
		right:  parser.NewInteger(2),
		op:     "<=",
		result: true,
	},
	{
		left:   parser.NewInteger(1),
		right:  parser.NewNull(),
		op:     "<=",
		result: false,
	},
	{
		left:   parser.NewInteger(2),
		right:  parser.NewInteger(2),
		op:     "<>",
		result: false,
	},
	{
		left:   parser.NewInteger(2),
		right:  parser.NewInteger(1),
		op:     "<>",
		result: true,
	},
}

func TestCompare(t *testing.T) {
	for _, test := range compareTests {
		res := Compare(test.left, test.right, test.op)
		if res != test.result {
			t.Errorf("%v, want %v for the expression %v %v %v", res, test.result, test.left, test.op, test.right)
		}
	}
}
