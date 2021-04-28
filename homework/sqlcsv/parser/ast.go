package parser

import (
	"fmt"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Statement interface{}

type Expression interface {
	String() string
}

type Prime interface {
	String() string
	Logic() bool
}

func IsPrime(e Expression) bool {
	if e == nil {
		return false
	}
	t := reflect.TypeOf(e)
	v := reflect.TypeOf((*Prime)(nil)).Elem()
	return t.Implements(v)
}

func IsNull(v Prime) bool {
	_, ok := v.(Null)
	return ok
}

// String
type String struct {
	str string
}

func (s String) String() string {
	return s.str
}

func NewString(s string) String {
	return String{
		s,
	}
}

func (s String) Logic() bool {
	if b, err := strconv.ParseBool(s.str); err == nil {
		return b
	}
	return false
}

// Integer
type Integer struct {
	str string
	num int64
}

func NewIntegerFromString(s string) Integer {
	i, _ := strconv.ParseInt(s, 10, 64)
	return Integer{
		str: s,
		num: i,
	}
}

func NewInteger(i int64) Integer {
	return Integer{
		str: strconv.FormatInt(i, 10),
		num: i,
	}
}

func (i Integer) String() string {
	return i.str
}

func (i Integer) Value() int64 {
	return i.num
}

func (i Integer) Logic() bool {
	if i.Value() == 1 {
		return true
	}
	return false
}

// Float
type Float struct {
	str string
	num float64
}

func NewFloatFromString(s string) Float {
	f, _ := strconv.ParseFloat(s, 64)
	return Float{
		str: s,
		num: f,
	}
}

func NewFloat(f float64) Float {
	return Float{
		str: strconv.FormatFloat(f, 'f', -1, 64),
		num: f,
	}
}

func (f Float) String() string {
	return f.str
}

func (f Float) Value() float64 {
	return f.num
}

func (f Float) Logic() bool {
	if f.Value() == 1 {
		return true
	}
	return false
}

// boolean
type Boolean struct {
	str     string
	boolean bool
}

func NewBoolean(b bool) Boolean {
	return Boolean{
		str:     strconv.FormatBool(b),
		boolean: b,
	}
}

func (b Boolean) String() string {
	return b.str
}

func (b Boolean) Logic() bool {
	return b.boolean
}

// logic op
type LogicOp struct {
	str   string
	value bool
}

func NewLogicOp(t bool) LogicOp {
	return LogicOp{
		value: t,
	}
}

func (t LogicOp) String() string {
	return t.str
}

func (t LogicOp) Bool() bool {
	if t.Logic() != true {
		return false
	}
	return true
}

func (t LogicOp) Logic() bool {
	return t.value
}

// Datetime
type Datetime struct {
	str string
	dt  time.Time
}

func NewDatetime(t time.Time) Datetime {
	return Datetime{
		str: t.Format(time.RFC850),
		dt:  t,
	}
}

func (dt Datetime) String() string {
	return dt.str
}

func (dt Datetime) Value() time.Time {
	return dt.dt
}

func (dt Datetime) Logic() bool {
	return !dt.Value().IsZero()
}

// null
type Null struct {
	str string
}

func NewNull() Null {
	return Null{}
}

func (n Null) String() string {
	if len(n.str) < 1 {
		return "NULL"
	}
	return n.str
}

func (n Null) Logic() bool {
	return false
}

// identifier
type Identifier struct {
	Str string
}

func (i Identifier) String() string {
	return i.Str
}

func (i *Identifier) Field() (string, string, error) {
	words := strings.Split(i.Str, ".")
	ref := ""
	field := ""
	if 2 < len(words) {
		// should replace `errors.New(fmt.Sprintf(...))` with `fmt.Errorf(...)` (golint)
		return "", "", fmt.Errorf("incorrect field identifier: %s", i.Str)
	} else if len(words) < 2 {
		field = words[0]
	} else {
		ref = words[0]
		field = words[1]
	}
	return ref, field, nil
}

type Parentheses struct {
	Expr Expression
}

func (p Parentheses) String() string {
	return p.Expr.String()
}

type SelectReq struct {
	SelectExpr Expression
	FromExpr   Expression
	WhereExpr  Expression
	LimitExpr  Expression
}

func (sq SelectReq) String() string {
	s := []string{sq.SelectExpr.String()}
	if sq.FromExpr != nil {
		s = append(s, sq.FromExpr.String())
	}
	if sq.WhereExpr != nil {
		s = append(s, sq.WhereExpr.String())
	}
	if sq.LimitExpr != nil {
		s = append(s, sq.LimitExpr.String())
	}
	return joinSpace(s)
}

type SelectExpr struct {
	Select string
	Fields []Expression
}

func (sc SelectExpr) String() string {
	s := []string{sc.Select}
	s = append(s, joinExpr(sc.Fields))
	return joinSpace(s)
}

type FromExpr struct {
	From   string
	Tables []Expression
}

func (f FromExpr) String() string {
	s := []string{f.From, joinExpr(f.Tables)}
	return joinSpace(s)
}

type WhereExpr struct {
	Where  string
	Filter Expression
}

func (w WhereExpr) String() string {
	s := []string{w.Where, w.Filter.String()}
	return joinSpace(s)
}

type LimitExpr struct {
	Limit  string
	Number int64
}

func (l LimitExpr) String() string {
	s := []string{l.Limit, strconv.FormatInt(l.Number, 10)}
	return joinSpace(s)
}

type Comparison struct {
	LHS      Expression
	Operator Token
	RHS      Expression
}

func (c Comparison) String() string {
	s := []string{c.LHS.String(), c.Operator.Str, c.RHS.String()}
	return joinSpace(s)
}

// Logic
type Logic struct {
	LHS      Expression
	Operator Token
	RHS      Expression
}

func (l Logic) String() string {
	s := []string{}
	if l.LHS != nil {
		s = append(s, l.LHS.String())
	}
	s = append(s, l.Operator.Str, l.RHS.String())
	return joinSpace(s)
}

type Table struct {
	Object Expression
	As     Token
	Alias  Expression
}

func (t Table) String() string {
	s := []string{t.Object.String()}
	if !t.As.IsEmpty() {
		s = append(s, t.As.Str)
	}
	if t.Alias != nil {
		s = append(s, t.Alias.String())
	}
	return joinSpace(s)
}

func (t *Table) Name() string {
	if t.Alias != nil {
		return t.Alias.(Identifier).Str
	}

	if file, ok := t.Object.(Identifier); ok {
		return strings.TrimSuffix(path.Base(file.Str), filepath.Ext(file.Str))
	}

	return t.Object.String()
}

// Field
type Field struct {
	Object Expression
	As     Token
	Alias  Expression
}

func (f Field) String() string {
	s := []string{f.Object.String()}
	if !f.As.IsEmpty() {
		s = append(s, f.As.Str)
	}
	if f.Alias != nil {
		s = append(s, f.Alias.String())
	}
	return joinSpace(s)
}

func (f *Field) Name() string {
	if f.Alias != nil {
		return f.Alias.(Identifier).Str
	}
	return f.Object.String()
}

type Asterisk struct {
}

func (as Asterisk) String() string {
	return "*"
}

type Variable struct {
	Name string
}

func (v Variable) String() string {
	return v.Name
}

func joinSpace(s []string) string {
	return strings.Join(s, " ")
}

func joinExpr(exprs []Expression) string {
	s := make([]string, len(exprs))
	for i, v := range exprs {
		s[i] = v.String()
	}
	return strings.Join(s, ", ")
}
