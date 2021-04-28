package printer

import (
	"testing"
)

// go test -v -cover ./...

var test = struct {
	input string
	err   string
}{
	input: "select col1 from null",
	err:   "file null absent",
}

func TestPrint(t *testing.T) {
	err := Print(test.input)
	if err != nil {
		if len(test.err) < 1 {
			t.Errorf("%v: unknown error %v", test.input, err)
		} else if err.Error() != test.err {
			t.Errorf("%v: got - %v, expected - %v", test.input, err, test.err)
		}
	} else {
		if 0 < len(test.err) {
			t.Errorf("no error, expected - %v for %v", test.err, test.input)
		}
	}
}
