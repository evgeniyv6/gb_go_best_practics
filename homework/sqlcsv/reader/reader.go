package reader

import (
	"bufio"
	"fmt"
	"gb_go_best_practics/homework/sqlcsv/parser"
	"io"
)

type Reader struct {
	Delimiter    rune
	reader       *bufio.Reader
	line         int
	column       int
	FieldsPerRec int
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		Delimiter:    ',',
		reader:       bufio.NewReader(r),
		line:         1,
		column:       0,
		FieldsPerRec: 0,
	}
}

func (r *Reader) newError(s string) error {
	return fmt.Errorf("line %d, col %d, err: %s", r.line, r.column, s)
}

func (r *Reader) ReadHeader() ([]string, error) {
	record, err := r.parseRecord()
	if err != nil {
		return nil, err
	}

	header := make([]string, len(record))
	for i, v := range record {
		header[i] = v.(parser.String).String()
	}
	return header, nil
}

func (r *Reader) Read() ([]parser.Prime, error) {
	record, err := r.parseRecord()
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (r *Reader) ReadAll() ([][]parser.Prime, error) {
	records := [][]parser.Prime{}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func (r *Reader) parseRecord() ([]parser.Prime, error) {
	var record []parser.Prime
	if r.FieldsPerRec < 1 {
		record = []parser.Prime{}
	} else {
		record = make([]parser.Prime, r.FieldsPerRec)
	}

	fieldIndex := 0
	for {
		if 0 < r.FieldsPerRec && r.FieldsPerRec <= fieldIndex {
			return nil, r.newError("wrong number of fields in line")
		}

		field, eol, err := r.parseField()
		if err == io.EOF {
			if fieldIndex < 1 {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}

		if r.FieldsPerRec < 1 {
			record = append(record, field)
		} else {
			record[fieldIndex] = field
		}

		if eol {
			break
		}

		fieldIndex++
	}

	if r.FieldsPerRec < 1 {
		r.FieldsPerRec = fieldIndex + 1
	}

	if r.FieldsPerRec != fieldIndex+1 {
		return nil, r.newError("wrong number of fields in line")
	}

	return record, nil
}

func (r *Reader) parseField() (parser.Prime, bool, error) {
	var eof error
	eol := false

	quoted := false
	escaped := false

	field := []rune{}

Read:
	for {
		ch, err := r.readRune()
		if err == io.EOF {
			if !escaped && quoted {
				return nil, eol, r.newError("bad symbol \" in field")
			}
			eof = err
			eol = true
			break
		}
		if err != nil {
			return nil, eol, err
		}

		if escaped {
			switch ch {
			case '"':
				escaped = false
				field = append(field, ch)
				continue
			case r.Delimiter:
				break Read
			case '\n':
				eol = true
				break Read
			default:
				r.column--
				return nil, eol, r.newError("bad symbol \" in field")
			}
		}

		if quoted {
			if ch == '"' {
				escaped = true
				continue
			}
			field = append(field, ch)
			continue
		}

		switch ch {
		case '\n':
			eol = true
			break Read
		case r.Delimiter:
			break Read
		case '"':
			if len(field) < 1 {
				quoted = true
			} else {
				field = append(field, ch)
			}
		default:
			field = append(field, ch)
		}
	}

	s := string(field)

	return parser.NewString(s), eol, eof
}

func (r *Reader) readRune() (rune, error) {
	ch, _, err := r.reader.ReadRune()
	r.column++

	if err != nil {
		return ch, err
	}

	if r.isNewLine(ch) {
		r.line++
		r.column = 0
		return '\n', nil
	}
	return ch, nil
}

func (r *Reader) isNewLine(ch rune) bool {
	if ch == '\n' {
		return true
	}
	return false
}
