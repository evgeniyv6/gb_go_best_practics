package crypto

import (
	"gb_go_best_practics/homework/sqlcsv/call"
	"gb_go_best_practics/homework/sqlcsv/parser"
	"testing"
)

// go test -v -cover ./...

var tests = []struct {
	view *call.View
	want string
	err  string
}{
	{
		view: &call.View{
			Header: call.NewHeader("test", []string{"col1", "col2\ncol3", "c4"}),
			Records: []call.Record{
				call.NewRecord([]parser.Prime{parser.NewInteger(123), parser.NewBoolean(true)}),
				call.NewRecord([]parser.Prime{parser.NewFloat(3.14), parser.NewString("foobar")}),
				call.NewRecord([]parser.Prime{parser.NewInteger(123456789), parser.NewString(" helloworld\nbaz\"qux\n"), parser.NewNull()}),
			},
		},
		want: "col1,col2\ncol3,c4\n123 , true , \n3.14 , foobar , \n123456789 ,  helloworld\nbaz\"qux\n , ",
	},
}

func TestEncodeCSV(t *testing.T) {
	for _, test := range tests {
		res := call.Result{
			View: test.view,
		}
		s := EncodeCSV(res)
		if s != test.want {
			t.Errorf("result = %q, want %q", s, test.want)
		}
	}
}
