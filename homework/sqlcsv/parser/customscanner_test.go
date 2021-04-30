package parser

import "testing"

type scanStruct struct {
	token  int
	str    string
	quoted bool
}

var scanTests = []struct {
	name  string
	input string
	want  []scanStruct
	err   string
}{
	{
		name:  "id",
		input: "id",
		want: []scanStruct{
			{
				token: ID,
				str:   "id",
			},
		},
	},
	{
		name:  "backticks id",
		input: "`id`",
		want: []scanStruct{
			{
				token:  ID,
				str:    "id",
				quoted: true,
			},
		},
	},
	{
		name:  "quoted string",
		input: "\"string\\\"\"",
		want: []scanStruct{
			{
				token: STRING,
				str:   "string\"",
			},
		},
	},
	{
		name:  "integer",
		input: "1",
		want: []scanStruct{
			{
				token: INTEGER,
				str:   "1",
			},
		},
	},
	{
		name:  "float",
		input: "1.234",
		want: []scanStruct{
			{
				token: FLOAT,
				str:   "1.234",
			},
		},
	},
	{
		name:  "Datetime",
		input: "\"2021-04-28 00:00:00\"",
		want: []scanStruct{
			{
				token: STRING,
				str:   "2021-04-28 00:00:00",
			},
		},
	},
	{
		name:  "comparison operator",
		input: ">",
		want: []scanStruct{
			{
				token: COMPARISON_OP,
				str:   ">",
			},
		},
	},
	{
		name:  "limit keyword",
		input: "limit",
		want: []scanStruct{
			{
				token: LIMIT,
				str:   "limit",
			},
		},
	},
	{
		name:  "where keyword",
		input: "where",
		want: []scanStruct{
			{
				token: WHERE,
				str:   "where",
			},
		},
	},
}

func TestCustomScanner_Scan(t *testing.T) {
	for _, test := range scanTests {
		s := new(CustomScanner).New(test.input)

		tokenCount := 0
		for {
			token, str, quoted, err := s.Scan()
			tokenCount++

			if err != nil {
				if test.err == "" {
					t.Errorf("test name %s: token %v: unknown error %v", test.name, tokenCount, err.Error())
				} else if test.err != err.Error() {
					t.Errorf("test name %s: token %v: error %v, want %v", test.name, tokenCount, err.Error(), test.err)
				}
				break
			}
			if test.err != "" {
				t.Errorf("test name %s:, token %v: no error, want %v", test.name, tokenCount, test.err)
				break
			}

			if token == EOF {
				break
			}

			if len(test.want) < tokenCount {
				break
			}
			expect := test.want[tokenCount-1]
			if token != expect.token {
				t.Errorf("test name %s, token %v: token = %v, want %v", test.name, tokenCount, GetToken(token), GetToken(expect.token))
			}
			if str != expect.str {
				t.Errorf("test name %s, token %v: str = %v, want %v", test.name, tokenCount, str, expect.str)
			}
			if quoted != expect.quoted {
				t.Errorf("test name %s, token %v: quoted = %v, want %v", test.name, tokenCount, quoted, expect.quoted)
			}
		}

		tokenCount--
		if tokenCount != len(test.want) {
			t.Errorf("test name %s: scan %v tokens in a statement, want %v", test.name, tokenCount, len(test.want))
		}
	}
}
