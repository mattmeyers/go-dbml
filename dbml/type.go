package dbml

import "strconv"

// Type represents a token type.
type Type int

// All token types in DBML.
const (
	Illegal Type = iota
	Comment
	EOF

	literalBeg
	Ident     // e.g. table name, column name, etc
	Int       // 123
	Float     // 123.45
	String    // 'abc'
	Quote     // "abc"
	FuncExpr  // `abc`
	Multiline // '''\nabc\n'''
	literalEnd

	relBeg
	MTO // >
	OTO // -
	OTM // <
	relEnd

	delimBeg
	LParen // (
	LBrack // [
	LBrace // {

	RParen // )
	RBrack // ]
	RBrace // }

	Colon     // :
	Semicolon // ;
	Period    // .
	Comma     // ,
	delimEnd

	keywordBeg
	Project
	DatabaseType
	Table
	As
	Indexes
	Ref
	Enum
	Note
	keywordEnd
)

var types = [...]string{
	Illegal: "ILLEGAL",
	Comment: "COMMENT",
	EOF:     "EOF",

	Ident:     "IDENT",
	Int:       "INT",
	Float:     "FLOAT",
	String:    "STRING",
	Quote:     "QUOTE",
	FuncExpr:  "FUNCEXPR",
	Multiline: "MULTILINE",

	MTO: ">",
	OTO: "-",
	OTM: "<",

	LParen: "(",
	LBrack: "[",
	LBrace: "{",

	RParen: ")",
	RBrack: "]",
	RBrace: "}",

	Colon:     ":",
	Semicolon: ";",
	Period:    ".",
	Comma:     ",",

	Project:      "Project",
	DatabaseType: "database_type",
	Table:        "Table",
	As:           "as",
	Indexes:      "Indexes",
	Ref:          "Ref",
	Enum:         "Enum",
	Note:         "Note",
}

func (t Type) String() string {
	s := ""
	if 0 <= t && t < Type(len(types)) {
		s = types[t]
	}
	if s == "" {
		s = "type(" + strconv.Itoa(int(t)) + ")"
	}
	return s
}

var keywords map[string]Type

func init() {
	keywords = make(map[string]Type)
	for i := keywordBeg + 1; i < keywordEnd; i++ {
		keywords[types[i]] = i
	}
}

func Lookup(s string) Type {
	if v, ok := keywords[s]; ok {
		return v
	}
	return Ident
}

func (t Type) IsLiteral() bool { return (literalBeg < t) && (t < literalEnd) }

func (t Type) IsRel() bool { return (relBeg < t) && (t < relEnd) }

func (t Type) IsDelim() bool { return (delimBeg < t) && (t < delimEnd) }

func (t Type) IsKeyword() bool { return (keywordBeg < t) && (t < keywordEnd) }

func IsKeyword(s string) bool { _, ok := keywords[s]; return ok }
