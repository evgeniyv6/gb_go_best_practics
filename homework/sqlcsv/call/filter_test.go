package call

import (
	"gb_go_best_practics/homework/sqlcsv/parser"
	"log"
	"os"
	"path"
	"reflect"
	"testing"
)

var filterTests = []struct {
	filter Filter
	expr   parser.Expression
	want   parser.Prime
	err    string
}{
	{
		expr: parser.NewString("str"),
		want: parser.NewString("str"),
	},
	{
		filter: []FilterRecord{
			{
				View: &View{
					Header: NewHeader("table1", []string{"col1", "col2"}),
					Records: []Record{
						NewRecord([]parser.Prime{
							parser.NewInteger(1),
							parser.NewString("newstr1"),
						}),
						NewRecord([]parser.Prime{
							parser.NewInteger(2),
							parser.NewString("newstr2"),
						}),
					},
				},
				RecordIndex: 1,
			},
		},
		expr: parser.Identifier{Str: "col2"},
		want: parser.NewString("newstr2"),
	},
	{
		filter: []FilterRecord{
			{
				View: &View{
					Header: NewHeader("table1", []string{"col1", "col2"}),
					Records: []Record{
						NewRecord([]parser.Prime{
							parser.NewInteger(1),
							parser.NewString("newstr1"),
						}),
						NewRecord([]parser.Prime{
							parser.NewInteger(2),
							parser.NewString("newstr2"),
						}),
					},
				},
				RecordIndex: 1,
			},
		},
		expr: parser.Identifier{Str: "col3"},
		err:  "rec is unrecognized",
	},
	{
		expr: parser.Logic{
			LHS:      parser.NewLogicOp(true),
			RHS:      parser.NewLogicOp(false),
			Operator: parser.Token{Token: parser.OR, Str: "or"},
		},
		want: parser.NewLogicOp(true),
	},
	{
		expr: parser.Logic{
			LHS:      parser.Identifier{Str: "absent"},
			RHS:      parser.NewLogicOp(false),
			Operator: parser.Token{Token: parser.AND, Str: "and"},
		},
		err: "rec not exist",
	},
	{
		expr: parser.Logic{
			LHS:      parser.NewLogicOp(false),
			RHS:      parser.Identifier{Str: "absent"},
			Operator: parser.Token{Token: parser.AND, Str: "and"},
		},
		err: "rec not exist",
	},
}

func TestFilter_Evaluate(t *testing.T) {
	curDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chdir(path.Join(curDir, "..", "testfiles"))
	if err != nil {
		log.Fatal(err)
	}

	for _, test := range filterTests {
		res, err := test.filter.Evaluate(test.expr)
		if !reflect.DeepEqual(res, test.want) {
			t.Errorf("got - %v, expected - %v", res.Logic(), test.want.Logic())
		}

		if err != nil {
			if len(test.err) < 1 {
				t.Errorf("unknown error %v", err)
			} else if err.Error() != test.err {
				t.Errorf("got - %v, expected - %v", err, test.err)
			}
			continue
		}
		if 0 < len(test.err) {
			t.Errorf("no error, expected - %v", test.err)
			continue
		}
	}
}
