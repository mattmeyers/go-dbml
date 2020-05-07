package dbml

import (
	"fmt"
)

type Node interface {
	Pos() Pos
}

type AST struct {
	Decls []Node
}

func (a *AST) appendDecl(d Node) {
	a.Decls = append(a.Decls, d)
}

func (a *AST) Print() {
	fmt.Println("[")

	for _, d := range a.Decls {
		switch v := d.(type) {
		case *ProjectDecl:
			fmt.Printf("\tProject: %+v\n", v)
		case *TableDecl:
			fmt.Printf("\tTable: %+v\n", v)
			for _, c := range v.Decls {
				switch u := c.(type) {
				case *ColDecl:
					fmt.Printf("\t\tCol: %+v\n", u)
				}
			}
		case *EnumDecl:
			fmt.Printf("\tEnum: %+v\n", v)
		}
	}

	fmt.Println("]")
}

type ProjectDecl struct {
	pos    Pos
	Name   string
	DBType string
	Note   string
}

type TableDecl struct {
	pos   Pos
	Name  string
	Alias string
	Decls []Node
}
type ColDecl struct {
	pos  Pos
	Name string
	Type string
}
type SettingDecl struct{}
type IndexDecl struct{}

type RefDecl struct{}

type EnumDecl struct {
	pos    Pos
	Name   string
	Values []string
}

func (d *ProjectDecl) Pos() Pos { return d.pos }
func (d *TableDecl) Pos() Pos   { return d.pos }
func (d *ColDecl) Pos() Pos     { return d.pos }
func (d *EnumDecl) Pos() Pos    { return d.pos }

func newAST() AST {
	return AST{Decls: make([]Node, 0)}
}

func Parse(tokens Tokens) AST {
	ast := newAST()
	current := 0

	var d Node

	for current < len(tokens) {
		switch tokens[current].Type {
		case TokenProject:
			d, current = parseProject(tokens, current)
			ast.appendDecl(d)
		case TokenTable:
			d, current = parseTable(tokens, current)
			ast.appendDecl(d)
		case TokenEnum:
			d, current = parseEnum(tokens, current)
			ast.appendDecl(d)
		case TokenEOF:
			break
		default:
			panic("expected one of Project, Table, Ref, or Enum, got " + tokens[current].Value)
		}
		current++
	}

	return ast
}

func parseProject(tokens Tokens, current int) (Node, int) {
	project := &ProjectDecl{}

	current++
	if tokens[current].Type != TokenIdent {
		panic("expected project name, got " + tokens[current].Value)
	}
	project.Name = tokens[current].Value

	current++
	if tokens[current].Type != TokenLBrace {
		panic("expected {, got " + tokens[current].Value)
	}
	project.Name = tokens[current].Value

	for current++; tokens[current].Type != TokenRBrace; current++ {
		switch tokens[current].Type {
		case TokenDatabaseType:
			current++
			if tokens[current].Type != TokenColon {
				panic("expected :, got " + tokens[current].Value)
			}

			current++
			if tokens[current].Type != TokenString {
				panic("expected string, got " + tokens[current].Value)
			}
			project.DBType = tokens[current].Value
		default:
			panic("expected one of database_type or Note, got " + tokens[current].Value)
		}
	}

	return project, current
}

func parseTable(tokens Tokens, current int) (Node, int) {
	table := &TableDecl{Decls: make([]Node, 0)}

	current++
	if tokens[current].Type != TokenIdent {
		panic("expected table name, got " + tokens[current].Value)
	}
	table.Name = tokens[current].Value

	current++
	if tokens[current].Type == TokenAs {
		current++
		if tokens[current].Type != TokenIdent {
			panic("expected table alias, got " + tokens[current].Value)
		}
		table.Alias = tokens[current].Value
	}

	current++
	if tokens[current].Type != TokenLBrace {
		panic("expected {, got " + tokens[current].Value)
	}

	var d Node
	for current++; tokens[current].Type != TokenRBrace; current++ {
		d, current = parseCol(tokens, current)
		table.Decls = append(table.Decls, d)
	}

	return table, current
}

func parseCol(tokens Tokens, current int) (Node, int) {
	col := &ColDecl{}

	if tokens[current].Type != TokenIdent {
		panic("expected col name, got " + tokens[current].Value)
	}
	col.Name = tokens[current].Value

	current++
	if tokens[current].Type != TokenIdent {
		panic("expected col type, got " + tokens[current].Value)
	}
	col.Type = tokens[current].Value

	return col, current
}

func parseEnum(tokens Tokens, current int) (Node, int) {
	enum := &EnumDecl{Values: make([]string, 0)}

	current++
	if tokens[current].Type != TokenIdent {
		panic("expected enum name, got " + tokens[current].Value)
	}
	enum.Name = tokens[current].Value

	current++
	if tokens[current].Type != TokenLBrace {
		panic("expected { name, got " + tokens[current].Value)
	}

	for current++; tokens[current].Type != TokenRBrace; current++ {
		if tokens[current].Type != TokenIdent {
			panic("expected an Ident, got " + tokens[current].Value)
		}
		enum.Values = append(enum.Values, tokens[current].Value)
	}

	return enum, current
}
