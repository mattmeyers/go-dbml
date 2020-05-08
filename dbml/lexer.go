package dbml

import (
	"io"
	"io/ioutil"
)

type Lexer struct {
	src    string
	pos    int
	start  int
	width  int
	tokens chan Token
	line   int
}

func Tokenize(r io.Reader) Tokens {
	f, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	s := []rune(string(f))

	tokens := make(Tokens, 0)
	var tok Token
	i := 0
	line, col := 1, 1

	for i < len(s) {
		switch r := s[i]; {
		case isLetter(r):
			l := i + 1
			for isLetter(s[l]) {
				l++
			}
			v := string(s[i:l])
			tok = Token{Value: v, Type: Lookup(v), Pos: Pos{Line: line, Col: col}}
			col += l - i
			i = l
		case isDigit(r):
			l := i + 1
			for isDigit(s[l]) {
				l++
			}
			v := string(s[i:l])
			tok = Token{Value: v, Type: TokenInt, Pos: Pos{Line: line, Col: col}}
			col += l - i
			i = l
		case r == '"':
			l := i + 1
			for s[l] != '"' {
				l++
			}
			tok = Token{Value: string(s[i+1 : l]), Type: TokenQuote, Pos: Pos{Col: col, Line: line}}
			i = l + 1
		case r == '\'':
			l := i + 1
			for s[l] != '\'' {
				l++
			}
			tok = Token{Value: string(s[i+1 : l]), Type: TokenString, Pos: Pos{Col: col, Line: line}}
			i = l + 1
		case r == '`':
			l := i + 1
			for s[l] != '`' {
				l++
			}
			tok = Token{Value: string(s[i+1 : l]), Type: TokenFuncExpr, Pos: Pos{Col: col, Line: line}}
			i = l + 1
		case r == '(':
			i++
			tok = Token{Value: "(", Type: TokenLParen, Pos: Pos{Line: line, Col: col}}
			col++
		case r == ')':
			i++
			tok = Token{Value: ")", Type: TokenRParen, Pos: Pos{Line: line, Col: col}}
			col++
		case r == '[':
			i++
			tok = Token{Value: "[", Type: TokenLBrack, Pos: Pos{Line: line, Col: col}}
			col++
		case r == ']':
			i++
			tok = Token{Value: "]", Type: TokenRBrack, Pos: Pos{Line: line, Col: col}}
			col++
		case r == '{':
			i++
			tok = Token{Value: "{", Type: TokenLBrace, Pos: Pos{Line: line, Col: col}}
			col++
		case r == '}':
			i++
			tok = Token{Value: "}", Type: TokenRBrace, Pos: Pos{Line: line, Col: col}}
			col++
		case r == ':':
			i++
			tok = Token{Value: ":", Type: TokenColon, Pos: Pos{Line: line, Col: col}}
			col++
		case r == ',':
			i++
			tok = Token{Value: ",", Type: TokenComma, Pos: Pos{Line: line, Col: col}}
			col++
		case isNewline(r):
			line++
			col = 1
			i++
			continue
		default:
			i++
			col++
			continue
		}

		tokens = append(tokens, tok)
	}

	return append(tokens, Token{Type: TokenEOF})
}

func toLower(r rune) rune {
	return ('a' - 'A') | r
}

func isLetter(r rune) bool {
	return ('a' <= toLower(r) && toLower(r) <= 'z') || r == '_'
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isAlphanumeric(r rune) bool {
	return isLetter(r) || isDigit(r)
}

func isNewline(r rune) bool {
	return r == '\n'
}
