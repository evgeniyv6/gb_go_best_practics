package call

import (
	"gb_go_best_practics/homework/sqlcsv/parser"
	"reflect"
	"testing"
)

func TestHeaderField_Label(t *testing.T) {
	hf := HeaderField{
		Column: "col1",
		Alias:  "al1",
	}
	expect := "al1"

	if hf.Label() != expect {
		t.Errorf("label = %s, want %s for %#v", hf.Label(), expect, hf)
	}

	hf = HeaderField{
		Column: "col1",
	}
	expect = "col1"

	if hf.Label() != expect {
		t.Errorf("label = %s, want %s for %#v", hf.Label(), expect, hf)
	}
}

func TestHeader_TableColumns(t *testing.T) {
	h := Header{
		{
			Reference: "test1",
			Column:    "c1",
			Alias:     "a1",
			FromTable: true,
		},
		{
			Reference: "test1",
			Column:    "c2",
			Alias:     "a3",
			FromTable: false,
		},
		{
			Column:    "c3",
			FromTable: true,
		},
	}
	expect := []parser.Expression{
		parser.Identifier{Str: "test1.c1"},
		parser.Identifier{Str: "c3"},
	}

	result := h.TableColumns()
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("columns = %s, want %s for %#v", result, expect, h)
	}
}

var headerContainsTests = []struct {
	Ref    string
	Column string
	Result int
	Error  string
}{
	{
		Ref:    "test2",
		Column: "col1",
		Result: 3,
	},
	{
		Column: "al2",
		Result: 1,
	},
	{
		Column: "col2",
		Result: 1,
	},
	{
		Column: "col1",
		Error:  "field absent",
	},
	{
		Column: "col22",
		Error:  "field absent",
	},
}

func TestHeader_Contains(t *testing.T) {
	h := Header{
		{
			Reference: "test1",
			Column:    "col1",
			Alias:     "al1",
			FromTable: true,
		},
		{
			Reference: "test1",
			Column:    "col2",
			Alias:     "al2",
			FromTable: false,
		},
		{
			Column:    "col3",
			FromTable: true,
		},
		{
			Reference: "test2",
			Column:    "col1",
			Alias:     "al3",
			FromTable: true,
		},
	}

	for _, test := range headerContainsTests {
		identifier := test.Column
		if 0 < len(test.Ref) {
			identifier = test.Ref + "." + identifier
		}

		result, err := h.Contains(test.Ref, test.Column)
		if err != nil {
			if len(test.Error) < 1 {
				t.Errorf("%s: unexpected error %q", identifier, err)
			} else if err.Error() != test.Error {
				t.Errorf("%s: error %q, want error %q", identifier, err, test.Error)
			}
			continue
		}
		if 0 < len(test.Error) {
			t.Errorf("%s: no error, want error %q", identifier, test.Error)
			continue
		}
		if result != test.Result {
			t.Errorf("%s: index = %d, want %d", identifier, result, test.Result)
		}
	}
}
