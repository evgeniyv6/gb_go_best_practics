%{
package parser
%}

%union{
    stmnt     []Statement
    statement   Statement
    expression  Expression
    expressions []Expression
    prime       Prime
    identifier  Identifier
    text        String
    integer     Integer
    float       Float
    variable    Variable
    token       Token
}

%type<stmnt>     stmnt
%type<statement>   statement
%type<expression>  expression
%type<expression>  select_req
%type<expression>  select_expr
%type<expression>  from_expr
%type<expression>  where_expr
%type<expression>  limit_expr
%type<prime>       prime
%type<expression>  value
%type<expression>  comparison
%type<expression>  logic
%type<expression>  table
%type<expression>  field_object
%type<expression>  field
%type<expressions> tables
%type<expressions> fields
%type<identifier>  identifier
%type<text>        text
%type<integer>     integer
%type<float>       float
%type<variable>    variable
%type<token>       statement_terminal
%token<token> ID STRING INTEGER FLOAT BOOLEAN DATETIME VARIABLE
%token<token> SELECT FROM WHERE
%token<token> LIMIT
%token<token> AND OR
%token<token> COMPARISON_OP

%left OR
%left AND
%left COMPARISON_OP
%left '+' '-'
%left '*' '/' '%'

%%

stmnt
    :
    {
        $$ = nil
        yylex.(*Lexer).stmnt = $$
    }
    | statement stmnt
    {
        $$ = append([]Statement{$1}, $2...)
        yylex.(*Lexer).stmnt = $$
    }

statement
    : expression statement_terminal
    {
        $$ = $1
    }

expression
    : select_req
    {
        $$ = $1
    }

select_req
    : select_expr from_expr where_expr limit_expr
    {
        $$ = SelectReq{
            SelectExpr:  $1,
            FromExpr:    $2,
            WhereExpr:   $3,
            LimitExpr:   $4,
        }
    }

select_expr
    : SELECT fields
    {
        $$ = SelectExpr{Select: $1.Str, Fields: $2}
    }

from_expr
    :
    {
        $$ = nil
    }
    | FROM tables
    {
        $$ = FromExpr{From: $1.Str, Tables: $2}
    }

where_expr
    :
    {
        $$ = nil
    }
    | WHERE value
    {
        $$ = WhereExpr{Where: $1.Str, Filter: $2}
    }

limit_expr
    :
    {
        $$ = nil
    }
    | LIMIT integer
    {
        $$ = LimitExpr{Limit: $1.Str, Number: $2.Value()}
    }

prime
    : text
    {
        $$ = $1
    }
    | integer
    {
        $$ = $1
    }
    | float
    {
        $$ = $1
    }

value
    : identifier
    {
        $$ = $1
    }
    | prime
    {
        $$ = $1
    }
    | comparison
    {
        $$ = $1
    }
    | logic
    {
        $$ = $1
    }
    | variable
    {
        $$ = $1
    }
    | '(' value ')'
    {
        $$ = Parentheses{Expr: $2}
    }

comparison
    : value COMPARISON_OP value
    {
        $$ = Comparison{LHS: $1, Operator: $2, RHS: $3}
    }

logic
    : value OR value
    {
        $$ = Logic{LHS: $1, Operator: $2, RHS: $3}
    }
    | value AND value
    {
        $$ = Logic{LHS: $1, Operator: $2, RHS: $3}
    }

table
    : identifier
    {
        $$ = Table{Object: $1}
    }
    | identifier identifier
    {
        $$ = Table{Object: $1, Alias: $2}
    }

field_object
    : value
    {
        $$ = $1
    }
    | '*'
    {
        $$ = Asterisk{}
    }

field
    : field_object
    {
        $$ = Field{Object: $1}
    }

tables
    : table
    {
        $$ = []Expression{$1}
    }
    | table ',' tables
    {
        $$ = append([]Expression{$1}, $3...)
    }

fields
    : field
    {
        $$ = []Expression{$1}
    }
    | field ',' fields
    {
        $$ = append([]Expression{$1}, $3...)
    }

identifier
    : ID
    {
        $$ = Identifier{Str: $1.Str}
    }

text
    : STRING
    {
        $$ = NewString($1.Str)
    }

integer
    : INTEGER
    {
        $$ = NewIntegerFromString($1.Str)
    }
    | '-' integer
    {
        i := $2.Value() * -1
        $$ = NewInteger(i)
    }

float
    : FLOAT
    {
        $$ = NewFloatFromString($1.Str)
    }
    | '-' float
    {
        f := $2.Value() * -1
        $$ = NewFloat(f)
    }

variable
    : VARIABLE
    {
        $$ = Variable{Name:$1.Str}
    }

statement_terminal
    :
    {
        $$ = Token{}
    }
    | ';'
    {
        $$ = Token{Token: ';', Str: string(';')}
    }

%%

func SetDebugLevel(level int, verbose bool) {
	yyDebug        = level
	yyErrorVerbose = verbose
}

func Parse(s string) ([]Statement, error) {
    l := new(Lexer)
    l.New(s)
    yyParse(l)
    return l.stmnt, l.err
}