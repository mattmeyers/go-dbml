package dbml

import (
	"io"
	"io/ioutil"
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

func Tokenize(r io.ReadCloser) []Token {
	defer r.Close()

	f, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	s := []rune(string(f))

	tokens := make([]Token, 0)
	var tok Token
	i := 0

	for i < len(s) {
		switch r := s[i]; {
		case isLetter(r):
			l := i + 1
			for isLetter(s[l]) {
				l++
			}
			v := string(s[i:l])
			tok = Token{Value: v, Type: Lookup(v)}
			i = l
		case isDigit(r):
			l := i + 1
			for isDigit(s[l]) {
				l++
			}
			v := string(s[i:l])
			tok = Token{Value: v, Type: Int}
			i = l
		case r == '"':
			l := i + 1
			for s[l] != '"' {
				l++
			}
			tok = Token{Value: string(s[i+1 : l]), Type: Quote}
			i = l + 1
		case r == '\'':
			l := i + 1
			for s[l] != '\'' {
				l++
			}
			tok = Token{Value: string(s[i+1 : l]), Type: String}
			i = l + 1
		case r == '`':
			l := i + 1
			for s[l] != '`' {
				l++
			}
			tok = Token{Value: string(s[i+1 : l]), Type: FuncExpr}
			i = l + 1
		case r == '(':
			i++
			tok = Token{Value: "(", Type: LParen}
		case r == ')':
			i++
			tok = Token{Value: ")", Type: RParen}
		case r == '[':
			i++
			tok = Token{Value: "[", Type: LBrack}
		case r == ']':
			i++
			tok = Token{Value: "]", Type: RBrack}
		case r == '{':
			i++
			tok = Token{Value: "{", Type: LBrace}
		case r == '}':
			i++
			tok = Token{Value: "}", Type: RBrace}
		case r == ':':
			i++
			tok = Token{Value: ":", Type: Colon}
		case r == ',':
			i++
			tok = Token{Value: ",", Type: Comma}
		default:
			i++
			continue
		}

		tokens = append(tokens, tok)
	}

	return tokens
}
