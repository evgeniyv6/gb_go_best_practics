package parser

import (
	"reflect"
	"testing"
)

var parseTests = []struct {
	input string
	want  []Statement
}{
	{
		input: "select foo; select bar;",
		want: []Statement{
			SelectReq{
				SelectExpr: SelectExpr{Select: "select", Fields: []Expression{Field{Object: Identifier{Str: "foo"}}}},
			},
			SelectReq{
				SelectExpr: SelectExpr{Select: "select", Fields: []Expression{Field{Object: Identifier{Str: "bar"}}}},
			},
		},
	},
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		p, err := Parse(test.input)
		if err != nil {
			t.Errorf("unknown error %v for %v", err.Error(), test.input)
			return
		}

		if len(test.want) != len(p) {
			t.Errorf("parsed program has %d statement(s), want %d statement(s) for %q", len(p), len(test.want), test.input)
			return
		}

		for i, stmt := range p {
			want := test.want[i]

			stmtType := reflect.TypeOf(stmt).Name()
			expectType := reflect.TypeOf(want).Name()

			if stmtType != expectType {
				t.Errorf("statement type is %v, want %v for %v", stmtType, expectType, test.input)
				continue
			}

			switch stmtType {
			case "SelectQuery":
				expectStmt := want.(SelectReq)
				parsedStmt := stmt.(SelectReq)

				if !reflect.DeepEqual(parsedStmt.SelectExpr, expectStmt.SelectExpr) {
					t.Errorf("select - %v, want %v for %v", parsedStmt.SelectExpr, expectStmt.SelectExpr, test.input)
				}
				if !reflect.DeepEqual(parsedStmt.FromExpr, expectStmt.FromExpr) {
					t.Errorf("from - %v, want %v for %v", parsedStmt.FromExpr, expectStmt.FromExpr, test.input)
				}
				if !reflect.DeepEqual(parsedStmt.WhereExpr, expectStmt.WhereExpr) {
					t.Errorf("where - %v, want %v for %q", parsedStmt.WhereExpr, expectStmt.WhereExpr, test.input)
				}
				if !reflect.DeepEqual(parsedStmt.LimitExpr, expectStmt.LimitExpr) {
					t.Errorf("limit - %v, want %v for %v", parsedStmt.LimitExpr, expectStmt.LimitExpr, test.input)
				}
			}
		}
	}
}
