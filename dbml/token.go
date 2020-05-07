package dbml

import "fmt"

type Token struct {
	Pos   Pos
	Type  TokenType
	Value string
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
	return fmt.Sprintf("{Type: %13s\tLine: %3d\tCol: %3d\tValue: %s}", t.Type, t.Pos.Line, t.Pos.Col, t.Value)
}

type Pos struct {
	Line int
	Col  int
}
