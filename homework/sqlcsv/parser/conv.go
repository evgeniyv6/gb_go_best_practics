package parser

import (
	"strconv"
	"strings"
	"time"
)

func ConvPrimeToFloat(p Prime) Prime {
	switch p.(type) {
	case Integer:
		return NewFloat(float64(p.(Integer).Value()))
	case Float:
		return p
	case String:
		if f, e := strconv.ParseFloat(p.(String).str, 64); e == nil {
			return NewFloat(f)
		}
	}

	return NewNull()
}

func ConvPrimeToDatetime(p Prime) Prime {
	switch p.(type) {
	case Integer:
		dt := time.Unix(p.(Integer).Value(), 0)
		return NewDatetime(dt)
	case Float:
		dt := float64ToTime(p.(Float).Value())
		return NewDatetime(dt)
	case Datetime:
		return p
	case String:
		if f, e := strconv.ParseFloat(p.(String).str, 64); e == nil {
			dt := float64ToTime(f)
			return NewDatetime(dt)
		}
	}

	return NewNull()
}

func ConvPrimeToBool(p Prime) Prime {
	switch p.(type) {
	case Boolean:
		return p
	case LogicOp:
		return NewBoolean(p.(LogicOp).Bool())
	case String:
		if b, e := strconv.ParseBool(p.(String).str); e == nil {
			return NewBoolean(b)
		}
	}
	return NewNull()
}

func float64ToTime(f float64) time.Time {
	s := strconv.FormatFloat(f, 'f', -1, 64)
	ns := strings.Split(s, ".")
	sec, _ := strconv.ParseInt(ns[0], 10, 64)
	var nsec int64
	if 1 < len(ns) {
		nsec, _ = strconv.ParseInt(ns[1]+strings.Repeat("0", 9-len(ns[1])), 10, 64)
	}
	return time.Unix(sec, nsec)
}
