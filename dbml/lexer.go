package dbml

import (
	"io"
	"io/ioutil"

	"github.com/mattmeyers/go-dbml/dbml/token"
)

func toLower(r rune) rune {
	return ('a' - 'A') | r
}

func isLetter(r rune) bool {
	return ('a' <= toLower(r) && toLower(r) <= 'z') || r == '_'
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func Tokenize(r io.Reader) token.Tokens {
	f, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	s := []rune(string(f))

	tokens := make(token.Tokens, 0)
	var tok token.Token
	i := 0

	for i < len(s) {
		switch r := s[i]; {
		case isLetter(r):
			l := i + 1
			for isLetter(s[l]) {
				l++
			}
			v := string(s[i:l])
			tok = token.Token{Value: v, Type: token.Lookup(v)}
			i = l
		case isDigit(r):
			l := i + 1
			for isDigit(s[l]) {
				l++
			}
			v := string(s[i:l])
			tok = token.Token{Value: v, Type: token.Int}
			i = l
		case r == '"':
			l := i + 1
			for s[l] != '"' {
				l++
			}
			tok = token.Token{Value: string(s[i+1 : l]), Type: token.Quote}
			i = l + 1
		case r == '\'':
			l := i + 1
			for s[l] != '\'' {
				l++
			}
			tok = token.Token{Value: string(s[i+1 : l]), Type: token.String}
			i = l + 1
		case r == '`':
			l := i + 1
			for s[l] != '`' {
				l++
			}
			tok = token.Token{Value: string(s[i+1 : l]), Type: token.FuncExpr}
			i = l + 1
		case r == '(':
			i++
			tok = token.Token{Value: "(", Type: token.LParen}
		case r == ')':
			i++
			tok = token.Token{Value: ")", Type: token.RParen}
		case r == '[':
			i++
			tok = token.Token{Value: "[", Type: token.LBrack}
		case r == ']':
			i++
			tok = token.Token{Value: "]", Type: token.RBrack}
		case r == '{':
			i++
			tok = token.Token{Value: "{", Type: token.LBrace}
		case r == '}':
			i++
			tok = token.Token{Value: "}", Type: token.RBrace}
		case r == ':':
			i++
			tok = token.Token{Value: ":", Type: token.Colon}
		case r == ',':
			i++
			tok = token.Token{Value: ",", Type: token.Comma}
		default:
			i++
			continue
		}

		tokens = append(tokens, tok)
	}

	return append(tokens, token.Token{Type: token.EOF})
}
