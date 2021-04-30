package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const (
	EOF = iota
	UNDEFINED
)

var keywords = []int{
	SELECT,
	FROM,
	WHERE,
	LIMIT,
	AND,
	OR,
}

var tokenIds = map[int]string{
	ID:            "ID",
	STRING:        "STRING",
	INTEGER:       "INTEGER",
	FLOAT:         "FLOAT",
	BOOLEAN:       "BOOLEAN",
	DATETIME:      "DATETIME",
	VARIABLE:      "VARIABLE",
	SELECT:        "SELECT",
	FROM:          "FROM",
	WHERE:         "WHERE",
	LIMIT:         "LIMIT",
	AND:           "AND",
	OR:            "OR",
	COMPARISON_OP: "COMPARISON_OP",
}

var comparisonOp = []string{
	"=",
	">",
	"<",
	">=",
	"<=",
	"<>",
	"!=",
}

type CustomScanner struct {
	src    []rune
	srcPos int
	offset int
	err    error
}

func (s *CustomScanner) New(src string) *CustomScanner {
	s.src = []rune(src)
	s.srcPos = 0
	s.offset = 0
	s.err = nil
	return s
}

func GetToken(token int) string {
	if val, ok := tokenIds[token]; ok {
		return val
	}
	return fmt.Sprint(token) // string(token)
}

func (s *CustomScanner) pos() rune {
	if len(s.src) <= s.srcPos {
		return EOF
	}
	return s.src[s.srcPos]
}

func (s *CustomScanner) inc() rune {
	ch := s.pos()
	if ch != EOF {
		s.srcPos++
		s.offset++
	}
	return ch
}

func (s *CustomScanner) runes() []rune {
	return s.src[(s.srcPos - s.offset):s.srcPos]
}

func (s *CustomScanner) toStr() string {
	return string(s.runes())
}

func (s *CustomScanner) escapedTokenStr() string {
	runes := s.runes()
	quote := runes[0]
	switch quote {
	case '"', '\'', '`':
		if runes[len(runes)-1] == quote {
			runes = runes[1:(len(runes) - 1)]
		} else {
			runes = runes[1:]
		}

		escaped := []rune{}
		for i := 0; i < len(runes); i++ {
			if runes[i] == '\\' && (i+1) < len(runes) && runes[i+1] == quote {
				i++
			}
			escaped = append(escaped, runes[i])
		}
		runes = escaped
	}
	return string(runes)
}

func (s *CustomScanner) Scan() (int, string, bool, error) {
	ch := s.pos()

	for s.isIndent(ch) {
		s.inc()
		ch = s.pos()
	}

	s.offset = 0
	ch = s.inc()
	token := ch
	literal := string(ch)
	quoted := false

	switch {
	case s.isIdentRune(ch):
		s.scanID()

		literal = s.toStr()
		if _, e := strconv.ParseInt(literal, 10, 64); e == nil {
			token = INTEGER
		} else if _, e := strconv.ParseFloat(literal, 64); e == nil {
			token = FLOAT
		} else if t, e := s.isKeyword(literal); e == nil {
			token = rune(t)
		} else {
			token = ID
		}
	case s.isOperatorRune(ch):
		s.scanOp()

		literal = s.toStr()
		if e := s.searchComparisonOp(literal); e == nil {
			token = COMPARISON_OP
		} else if 1 < len(literal) {
			token = UNDEFINED
		}
	default:
		switch ch {
		case EOF:
			break
		case '"', '\'':
			s.scanString(ch)
			literal = s.escapedTokenStr()
			token = STRING
		case '`':
			s.scanString(ch)
			literal = s.escapedTokenStr()
			token = ID
			quoted = true
		}
	}

	return int(token), literal, quoted, s.err
}

func (s *CustomScanner) isIndent(ch rune) bool {
	switch ch {
	case '\t', '\n', '\r', ' ':
		return true
	}
	return false
}

func (s *CustomScanner) scanString(quote rune) {
	for {
		ch := s.inc()
		if ch == EOF {
			s.err = errors.New("string not term")
			break
		} else if ch == quote {
			break
		} else if ch == '\\' {
			s.inc()
		}
	}

	return
}

func (s *CustomScanner) scanID() {
	for s.isIdentRune(s.pos()) {
		s.inc()
	}
	return
}

func (s *CustomScanner) isIdentRune(ch rune) bool {
	return ch == '_' || ch == '$' || ch == '.' || unicode.IsLetter(ch) || unicode.IsDigit(ch)
}

func (s *CustomScanner) scanOp() {
	for s.isOperatorRune(s.pos()) {
		s.inc()
	}
	return
}

func (s *CustomScanner) isOperatorRune(ch rune) bool {
	switch ch {
	case '=', '>', '<', '!', '|', ':':
		return true
	}
	return false
}

func (s *CustomScanner) isKeyword(str string) (int, error) {
	for _, v := range keywords {
		if strings.EqualFold(GetToken(v), str) {
			return v, nil
		}
	}
	// should replace `errors.New(fmt.Sprintf(...))` with `fmt.Errorf(...)` (golint)
	return ID, fmt.Errorf("%q not keyword", str)
}

func (s *CustomScanner) searchComparisonOp(str string) error {
	for _, v := range comparisonOp {
		if v == str {
			return nil
		}
	}
	// should replace `errors.New(fmt.Sprintf(...))` with `fmt.Errorf(...)` (golint)
	return fmt.Errorf("%q not comparison operator", str)
}
