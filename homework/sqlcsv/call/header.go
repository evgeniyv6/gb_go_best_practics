package call

import (
	"errors"
	"gb_go_best_practics/homework/sqlcsv/parser"
)

type HeaderField struct {
	Reference  string
	Column     string
	Alias      string
	FromTable  bool
	IsGroupKey bool
}

func (hf HeaderField) Label() string {
	if 0 < len(hf.Alias) {
		return hf.Alias
	}
	return hf.Column
}

type Header []HeaderField

func NewHeader(ref string, words []string) Header {
	h := make([]HeaderField, len(words))

	for i, v := range words {
		h[i].Reference = ref
		h[i].Column = v
		h[i].FromTable = true
	}
	return h
}

func NewEmptyHeader(len int) Header {
	return make([]HeaderField, len)
}

func AddHeaderField(h Header, alias string) Header {
	return append(h, HeaderField{
		Alias: alias,
	})
}

func (h Header) Len() int {
	return len(h)
}

func (h Header) TableColumns() []parser.Expression {
	columns := []parser.Expression{}
	for _, f := range h {
		if !f.FromTable {
			continue
		}

		var lit string
		if 0 < len(f.Reference) {
			lit = f.Reference + "." + f.Column
		} else {
			lit = f.Column
		}
		columns = append(columns, parser.Identifier{Str: lit})
	}
	return columns
}

func (h Header) Contains(ref string, column string) (int, error) {
	idx := -1

	for i, f := range h {
		if 0 < len(ref) {
			if f.Reference != ref || f.Column != column {
				continue
			}
		} else {
			if f.Column != column && f.Alias != column {
				continue
			}
		}

		if -1 < idx {
			return -1, errors.New("field absent")
		}
		idx = i
	}

	if idx < 0 {
		return -1, errors.New("field absent")
	}

	return idx, nil
}
