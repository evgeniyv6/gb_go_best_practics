package reader

import (
	"gb_go_best_practics/homework/sqlcsv/parser"
	"reflect"
	"strings"
	"testing"
)

// go test -v -cover ./...

var tests = []struct {
	input string
	want  [][]parser.Prime
	err   string
}{
	{
		input: "hello,world\nfoo,bar",
		want: [][]parser.Prime{
			{parser.NewString("hello"), parser.NewString("world")},
			{parser.NewString("foo"), parser.NewString("bar")},
		},
	},
	{
		input: "hello,\"world\",\"foo\nbar\"\nbaz,qux,",
		want: [][]parser.Prime{
			{parser.NewString("hello"), parser.NewString("world"), parser.NewString("foo\nbar")},
			{parser.NewString("baz"), parser.NewString("qux"), parser.NewString("")},
		},
	},
	{
		input: "foo,bar,baz\nqux,ose",
		err:   "line 2, col 8, err: wrong number of fields in line",
	},
	{
		input: "foo,\"bar\",\"baz\nqux\nose,",
		err:   "line 3, col 5, err: bad symbol \" in field",
	},
}

var headerTest = struct {
	input string
	want  [][]parser.Prime
}{
	input: "col1,col2\n1,2\nfoo,bar",
	want: [][]parser.Prime{
		{parser.NewString("1"), parser.NewString("2")},
		{parser.NewString("foo"), parser.NewString("bar")},
	},
}

func TestReader_ReadAll(t *testing.T) {
	for _, test := range tests {
		nr := NewReader(strings.NewReader(test.input))

		rec, err := nr.ReadAll()

		if err != nil {
			if test.err == "" {
				t.Errorf(": unexpected error %q", err.Error())
			} else if test.err != err.Error() {
				t.Errorf(": error %q, want error %q", err.Error(), test.err)
			}
			continue
		}

		if !reflect.DeepEqual(rec, test.want) {
			t.Errorf(": records = %q, want %q", rec, test.want)
		}

	}
}

func TestReader_ReadHeader(t *testing.T) {
	nr := NewReader(strings.NewReader(headerTest.input))
	h, err := nr.ReadHeader()

	if err != nil {
		t.Errorf("unknown error %v", err.Error())
	}
	if !reflect.DeepEqual(h, []string{"col1", "col2"}) {
		t.Errorf("header - %q, want %q", h, []string{"col1", "col2"})
	}

	rec, err := nr.ReadAll()
	if err != nil {
		t.Errorf("unexpected error %q", err.Error())
	}
	if !reflect.DeepEqual(rec, headerTest.want) {
		t.Errorf("records = %q, want %q", rec, headerTest.want)
	}
}
