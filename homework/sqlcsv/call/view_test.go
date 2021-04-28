package call

import (
	"gb_go_best_practics/homework/sqlcsv/parser"
	"log"
	"os"
	"path"
	"reflect"
	"testing"
)

var newViewTests = []struct {
	name   string
	from   parser.FromExpr
	filter Filter
	want   *View
	err    string
}{
	{
		name: "load file",
		from: parser.FromExpr{
			Tables: []parser.Expression{
				parser.Table{
					Object: parser.Identifier{Str: "test1.csv"},
				},
			},
		},
		want: &View{
			Header: NewHeader("test1", []string{"col1", "col2", "col3"}),
			Records: []Record{
				NewRecord([]parser.Prime{
					parser.NewString("1"),
					parser.NewString("2"),
					parser.NewString("3"),
				}),
				NewRecord([]parser.Prime{
					parser.NewString("hello"),
					parser.NewString("foo"),
					parser.NewString("bar"),
				}),
			},
			Filename: "test1.csv",
		},
	},
	{
		name: "load file error",
		from: parser.FromExpr{
			Tables: []parser.Expression{
				parser.Table{
					Object: parser.Identifier{Str: "null"},
				},
			},
		},
		err: "file null absent",
	},
}

func TestNewView(t *testing.T) {
	curDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chdir(path.Join(curDir, "..", "testfiles"))
	if err != nil {
		log.Fatal(err)
	}

	for _, test := range newViewTests {
		res, err := NewView(test.from, test.filter)
		if !reflect.DeepEqual(res, test.want) {
			t.Errorf("test name %s: got - %v, expected - %v", test.name, res, test.want)
		}

		if err != nil {
			if len(test.err) < 1 {
				t.Errorf("test name %s: unknown error %v", test.name, err)
			} else if err.Error() != test.err {
				t.Errorf("test name %s:: got - %v, expected - %v", test.name, err.Error(), test.err)
			}
			continue
		}
		if 0 < len(test.err) {
			t.Errorf("test name: %s, no error, expected - %v", test.name, test.err)
			continue
		}
	}
}

var viewWhereTests = []struct {
	name  string
	view  *View
	where parser.WhereExpr
	want  []int
	err   string
}{
	{
		name: "where",
		view: &View{
			Header: NewHeader("table1", []string{"col1", "col2", "col3"}),
			Records: []Record{
				NewRecord([]parser.Prime{
					parser.NewString("1"),
					parser.NewString("2"),
					parser.NewString("3"),
				}),
				NewRecord([]parser.Prime{
					parser.NewString("hello"),
					parser.NewString("foo"),
					parser.NewString("bar"),
				}),
			},
		},
		where: parser.WhereExpr{
			Filter: parser.Comparison{
				LHS:      parser.Identifier{Str: "col1"},
				RHS:      parser.NewInteger(1),
				Operator: parser.Token{Token: parser.COMPARISON_OP, Str: "="},
			},
		},
		want: []int{0},
	},
	{
		name: "where filter error",
		view: &View{
			Header: NewHeader("test1", []string{"col1", "col2", "col3"}),
			Records: []Record{
				NewRecord([]parser.Prime{
					parser.NewString("1"),
					parser.NewString("2"),
					parser.NewString("3"),
				}),
				NewRecord([]parser.Prime{
					parser.NewString("hello"),
					parser.NewString("foo"),
					parser.NewString("bar"),
				}),
			},
		},
		where: parser.WhereExpr{
			Filter: parser.Comparison{
				LHS:      parser.Identifier{Str: "null"},
				RHS:      parser.NewInteger(2),
				Operator: parser.Token{Token: parser.COMPARISON_OP, Str: "="},
			},
		},
		err: "rec is unrecognized",
	},
}

func TestView_Where(t *testing.T) {
	for _, test := range viewWhereTests {
		err := test.view.Where(test.where)
		if err != nil {
			if len(test.err) < 1 {
				t.Errorf("test name %s: unexpected error %v", test.name, err)
			} else if err.Error() != test.err {
				t.Errorf("test name %s: error %v, want error %v", test.name, err.Error(), test.err)
			}
			continue
		}
		if 0 < len(test.err) {
			t.Errorf("%s: no error, want error %q", test.name, test.err)
			continue
		}
		if !reflect.DeepEqual(test.view.filteredIndices, test.want) {
			t.Errorf("test name %s: result = %v, want %v", test.name, test.view.filteredIndices, test.want)
		}
	}
}

var selectTests = []struct {
	name      string
	view      *View
	selectExp parser.SelectExpr
	want      *View
	err       string
}{
	{
		name: "select",
		view: &View{
			Header: []HeaderField{
				{Reference: "test1", Column: "col1", FromTable: true},
				{Reference: "test1", Column: "col2", FromTable: true},
				{Reference: "test1", Column: "col3", FromTable: true},
			},
			Records: []Record{
				NewRecord([]parser.Prime{
					parser.NewString("1"),
					parser.NewString("2"),
					parser.NewString("3"),
				}),
				NewRecord([]parser.Prime{
					parser.NewString("hello"),
					parser.NewString("foo"),
					parser.NewString("bar"),
				}),
			},
		},
		selectExp: parser.SelectExpr{
			Fields: []parser.Expression{
				parser.Field{Object: parser.Identifier{Str: "col3"}},
				parser.Field{Object: parser.Asterisk{}},
				parser.Field{Object: parser.NewInteger(100500), Alias: parser.Identifier{Str: "s"}},
			},
		},
		want: &View{
			Header: []HeaderField{
				{Reference: "test1", Column: "col1", FromTable: true},
				{Reference: "test1", Column: "col2", FromTable: true},
				{Reference: "test1", Column: "col3", FromTable: true},
				{Alias: "s"},
			},
			Records: []Record{
				NewRecord([]parser.Prime{
					parser.NewString("1"),
					parser.NewString("2"),
					parser.NewString("3"),
					parser.NewInteger(100500),
				}),
				NewRecord([]parser.Prime{
					parser.NewString("hello"),
					parser.NewString("foo"),
					parser.NewString("bar"),
					parser.NewInteger(100500),
				}),
			},
			selectFields: []int{2, 0, 1, 2, 3},
		},
	},
}

func TestView_Select(t *testing.T) {
	for _, test := range selectTests {
		err := test.view.Select(test.selectExp)
		if err != nil {
			if len(test.err) < 1 {
				t.Errorf("test name %s: unexpected error %v", test.name, err)
			} else if err.Error() != test.err {
				t.Errorf("test name %s: error %q, want error %v", test.name, err.Error(), test.err)
			}
			continue
		}
		if 0 < len(test.err) {
			t.Errorf("test name %s: no error, want error %v", test.name, test.err)
			continue
		}
		if !reflect.DeepEqual(test.view.Header, test.want.Header) {
			t.Errorf("test name %s: header = %v, want %v", test.name, test.view.Header, test.want.Header)
		}
		if !reflect.DeepEqual(test.view.Records, test.want.Records) {
			t.Errorf("test name %s: records = %v, want %v", test.name, test.view.Records, test.want.Records)
		}
		if !reflect.DeepEqual(test.view.selectFields, test.want.selectFields) {
			t.Errorf("test name %s: select indices = %v, want %v", test.name, test.view.selectFields, test.want.selectFields)
		}
	}
}
