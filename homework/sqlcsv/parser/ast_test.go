package parser

import (
	"testing"
	"time"
)

func TestIsPrime(t *testing.T) {
	var expr Expression

	expr = Identifier{Str: "foo"}
	if IsPrime(expr) {
		t.Errorf("test name ast is prime: expr %v not prime", expr)
	}

	expr = Integer{str: "100500"}
	if !IsPrime(expr) {
		t.Errorf("test name ast is prime:expr %v is prime, but do not", expr)
	}

	if IsPrime(nil) {
		t.Error("test name ast is prime:nil is prime, want empty string")
	}
}

func TestIsNull(t *testing.T) {
	var p Prime

	p = NewInteger(1)
	if IsNull(p) {
		t.Errorf("test name ast is prime: err - value %v not null", p)
	}

	p = NewNull()
	if !IsNull(p) {
		t.Errorf("test name ast is prime: err - value %v is null", p)
	}
}

func TestString_String(t *testing.T) {
	s := "foobar"
	p := NewString(s)
	want := s
	if p.String() != want {
		t.Errorf("string = %v, want %v for %v", p.String(), want, p)
	}
}

func TestString_Value(t *testing.T) {
	s := "foobar"
	p := NewString(s)
	if p.String() != s {
		t.Errorf("value = %v, want %v for %v", p.String(), s, p)
	}
}

func TestString_Logic(t *testing.T) {
	s := "true"
	p := NewString(s)
	if p.Logic() != true {
		t.Errorf("bool = %v, want %v for %v", p.Logic(), true, p)
	}

	s = "false"
	p = NewString(s)
	if p.Logic() != false {
		t.Errorf("bool = %v, want %v for %v", p.Logic(), false, p)
	}

	s = "error"
	p = NewString(s)
	if p.Logic() != false {
		t.Errorf("bool = %v, want %v for %v", p.Logic(), false, p)
	}
}

func TestDatetime_Logic(t *testing.T) {
	d := NewDatetime(time.Time{})

	if d.Logic() != false {
		t.Errorf("bool = %v, want %v for %v", d.Logic(), false, d)
	}
}

func TestNull_Bool(t *testing.T) {
	p := NewNull()
	if p.Logic() != false {
		t.Errorf("bool = %v, want %v for %v", p.Logic(), false, p)
	}
}

func TestSelectQuery_String(t *testing.T) {
	expr := SelectReq{
		SelectExpr: SelectExpr{
			Select: "select",
			Fields: []Expression{Field{Object: Identifier{Str: "col"}}},
		},
		FromExpr: FromExpr{
			From:   "from",
			Tables: []Expression{Table{Object: Identifier{Str: "test.csv"}}},
		},
		WhereExpr: WhereExpr{
			Where: "where",
			Filter: Comparison{
				LHS:      Identifier{Str: "col"},
				Operator: Token{Token: COMPARISON_OP, Str: ">"},
				RHS:      Integer{str: "1"},
			},
		},
		LimitExpr: LimitExpr{
			Limit:  "limit",
			Number: 100500,
		},
	}
	want := "select col from test.csv where col > 1 limit 100500"
	if expr.String() != want {
		t.Errorf("string = %v, want %v for %v", expr.String(), want, expr)
	}
}

func TestFromExpr_String(t *testing.T) {
	expr := FromExpr{
		From: "from",
		Tables: []Expression{
			Table{Object: Identifier{Str: "test1.csv"}},
			Table{Object: Identifier{Str: "test2.csv"}},
		},
	}
	want := "from test1.csv, test2.csv"
	if expr.String() != want {
		t.Errorf("string = %v, want %v for %v", expr.String(), want, expr)
	}
}

func TestWhereExpr_String(t *testing.T) {
	expr := WhereExpr{
		Where: "where",
		Filter: Comparison{
			LHS:      Identifier{Str: "col1"},
			Operator: Token{Token: COMPARISON_OP, Str: ">"},
			RHS:      Integer{str: "1"},
		},
	}
	want := "where col1 > 1"
	if expr.String() != want {
		t.Errorf("string = %v, want %v for %v", expr.String(), want, expr)
	}
}

func TestLimitClause_String(t *testing.T) {
	expr := LimitExpr{Limit: "limit", Number: 100500}
	want := "limit 100500"
	if expr.String() != want {
		t.Errorf("string = %q, want %q for %#v", expr.String(), want, expr)
	}
}
