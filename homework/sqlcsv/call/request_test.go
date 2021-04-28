package call

import (
	"gb_go_best_practics/homework/sqlcsv/parser"
	"log"
	"os"
	"path"
	"reflect"
	"testing"
)

// go test -v -cover ./...

var reqTests = []struct {
	input string
	want  []Result
	err   string
}{
	{
		input: "select col1 from test1.csv limit 1",
		want: []Result{
			{
				View: &View{
					Header: []HeaderField{
						{
							Reference:  "test1",
							Column:     "col1",
							FromTable:  true,
							IsGroupKey: false,
						},
					},
					Records: []Record{
						{
							NewCell(parser.NewString("1")),
						},
					},
					Filename: "test1.csv",
				},
			},
		},
	},
	{
		input: "select col2 from test2.csv where col1 > '2020-02-28'",
		want: []Result{
			{
				View: &View{
					Header: []HeaderField{
						{
							Reference:  "test2",
							Column:     "col2",
							FromTable:  true,
							IsGroupKey: false,
						},
					},
					Records: []Record{
						{
							NewCell(parser.NewString("bar")),
						},
					},
					Filename: "test2.csv",
				},
			},
		},
	},
	{
		input: "select from absent",
		err:   "syntax error",
	},
	{
		input: "select col1 from null",
		err:   "file null absent",
	},
	{
		input: "select col1 from test1.csv where empty = 'foo'",
		err:   "rec is unrecognized",
	},
}

func TestRequest(t *testing.T) {
	curDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chdir(path.Join(curDir, "..", "testfiles"))
	if err != nil {
		log.Fatal(err)
	}

	for _, test := range reqTests {
		res, err := Execute(test.input)
		if !reflect.DeepEqual(res, test.want) {
			t.Errorf("%v: got - %v, expected - %v", test.input, res, test.want)
		}

		if err != nil {
			if len(test.err) < 1 {
				t.Errorf("%v: unknown error %v", test.input, err)
			} else if err.Error() != test.err {
				t.Errorf("%v: got - %v, expected - %v", test.input, err, test.err)
			}
			continue
		}
		if 0 < len(test.err) {
			t.Errorf("no error, expected - %v for %v", test.err, test.input)
			continue
		}
	}
}
