package dbml

import "strconv"

// TokenType represents a token type.
type TokenType int

// All token types in DBML.
const (
	TokenIllegal TokenType = iota
	TokenComment
	TokenEOF

	literalBeg
	TokenIdent     // e.g. table name, column name, etc
	TokenInt       // 123
	TokenFloat     // 123.45
	TokenString    // 'abc'
	TokenQuote     // "abc"
	TokenFuncExpr  // `abc`
	TokenMultiline // '''\nabc\n'''
	literalEnd

	relBeg
	TokenMTO // >
	TokenOTO // -
	TokenOTM // <
	relEnd

	delimBeg
	TokenLParen // (
	TokenLBrack // [
	TokenLBrace // {

	TokenRParen // )
	TokenRBrack // ]
	TokenRBrace // }

	TokenColon     // :
	TokenSemicolon // ;
	TokenPeriod    // .
	TokenComma     // ,
	delimEnd

	keywordBeg
	TokenProject
	TokenDatabaseType
	TokenTable
	TokenAs
	TokenIndexes
	TokenRef
	TokenEnum
	TokenNote
	keywordEnd
)

var types = [...]string{
	TokenIllegal: "ILLEGAL",
	TokenComment: "COMMENT",
	TokenEOF:     "EOF",

	TokenIdent:     "IDENT",
	TokenInt:       "INT",
	TokenFloat:     "FLOAT",
	TokenString:    "STRING",
	TokenQuote:     "QUOTE",
	TokenFuncExpr:  "FUNCEXPR",
	TokenMultiline: "MULTILINE",

	TokenMTO: ">",
	TokenOTO: "-",
	TokenOTM: "<",

	TokenLParen: "(",
	TokenLBrack: "[",
	TokenLBrace: "{",

	TokenRParen: ")",
	TokenRBrack: "]",
	TokenRBrace: "}",

	TokenColon:     ":",
	TokenSemicolon: ";",
	TokenPeriod:    ".",
	TokenComma:     ",",

	TokenProject:      "Project",
	TokenDatabaseType: "database_type",
	TokenTable:        "Table",
	TokenAs:           "as",
	TokenIndexes:      "Indexes",
	TokenRef:          "Ref",
	TokenEnum:         "Enum",
	TokenNote:         "Note",
}

func (t TokenType) String() string {
	s := ""
	if 0 <= t && t < TokenType(len(types)) {
		s = types[t]
	}
	if s == "" {
		s = "type(" + strconv.Itoa(int(t)) + ")"
	}
	return s
}

var keywords map[string]TokenType

func init() {
	keywords = make(map[string]TokenType)
	for i := keywordBeg + 1; i < keywordEnd; i++ {
		keywords[types[i]] = i
	}
}

func Lookup(s string) TokenType {
	if v, ok := keywords[s]; ok {
		return v
	}
	return TokenIdent
}

func (t TokenType) IsLiteral() bool { return (literalBeg < t) && (t < literalEnd) }

func (t TokenType) IsRel() bool { return (relBeg < t) && (t < relEnd) }

func (t TokenType) IsDelim() bool { return (delimBeg < t) && (t < delimEnd) }

func (t TokenType) IsKeyword() bool { return (keywordBeg < t) && (t < keywordEnd) }

func IsKeyword(s string) bool { _, ok := keywords[s]; return ok }
