package parser

import "testing"

func TestPrimaryToDatetime(t *testing.T) {
	var a Prime
	var b Prime

	a = NewInteger(2121181845)
	b = ConvPrimeToDatetime(a)
	if _, ok := b.(Datetime); !ok {
		t.Errorf("prime type = %T, want datetime for %#v", b, a)
	}

	a = NewFloat(2121181845)
	b = ConvPrimeToDatetime(a)
	if _, ok := b.(Datetime); !ok {
		t.Errorf("prime type = %T, want datetime for %#v", b, a)
	}

	a = NewString("2121181845.980")
	b = ConvPrimeToDatetime(a)
	if _, ok := b.(Datetime); !ok {
		t.Errorf("prime type = %T, want datetime for %#v", b, a)
	}
}

func TestConvPrimeToBool(t *testing.T) {
	var a Prime
	var b Prime

	a = NewBoolean(true)
	b = ConvPrimeToBool(a)
	if _, ok := b.(Boolean); !ok {
		t.Errorf("prime type is %T, want boolean for %#v", b, a)
	}

	a = NewLogicOp(true)
	b = ConvPrimeToBool(a)
	if _, ok := b.(Boolean); !ok {
		t.Errorf("prime type is %T, want boolean for %#v", b, a)
	}

	a = NewString("true")
	b = ConvPrimeToBool(a)
	if _, ok := b.(Boolean); !ok {
		t.Errorf("prime type is %T, want boolean for %#v", b, a)
	}

	a = NewString("error")
	b = ConvPrimeToBool(a)
	if _, ok := b.(Null); !ok {
		t.Errorf("prime type is %T, want boolean for %#v", b, a)
	}
}
