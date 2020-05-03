package token

import "fmt"

type Token struct {
	Value string
	Type
}

type Tokens []Token

func (t Tokens) Print() {
	fmt.Println("[")
	for _, tok := range t {
		fmt.Printf("\t%s\n", tok.String())
	}
	fmt.Println("]")
}

func (t Token) String() string {
	return fmt.Sprintf("{Type: %10s\tValue: %s}", t.Type, t.Value)
}
