package dbml

import (
	"errors"
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

func (a *AST) Project() Project {
	for _, d := range a.Decls {
		if v, ok := d.(*ProjectDecl); ok {
			return v.Parse()
		}
	}
	return Project{}
}

func (a *AST) Tables() []*TableDecl {
	tables := []*TableDecl{}
	for _, d := range a.Decls {
		if v, ok := d.(*TableDecl); ok {
			tables = append(tables, v)
		}
	}
	return tables
}

func (a *AST) Enums() []*EnumDecl {
	enums := []*EnumDecl{}
	for _, d := range a.Decls {
		if v, ok := d.(*EnumDecl); ok {
			enums = append(enums, v)
		}
	}
	return enums
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

func (d *ProjectDecl) Parse() Project {
	p := Project{Name: d.Name, DBType: d.DBType}
	return p
}

func (d *TableDecl) Pos() Pos { return d.pos }

func (d *TableDecl) Parse() Table {
	return Table{}
}

func (d *ColDecl) Pos() Pos { return d.pos }

func (d *ColDecl) Parse() Col {
	return Col{}
}

func (d *EnumDecl) Pos() Pos { return d.pos }

func (d *EnumDecl) Parse() Enum {
	return Enum{}
}

func newAST() AST {
	return AST{Decls: make([]Node, 0)}
}

func Parse(tokens Tokens) (AST, error) {
	ast := newAST()
	current := 0

	var d Node
	var err error

	for current < len(tokens) {
		switch tokens[current].Type {
		case TokenProject:
			d, current, err = parseProject(tokens, current)
			if err != nil {
				return AST{}, err
			}
			ast.appendDecl(d)
		case TokenTable:
			d, current, err = parseTable(tokens, current)
			if err != nil {
				return AST{}, err
			}
			ast.appendDecl(d)
		case TokenEnum:
			d, current, err = parseEnum(tokens, current)
			if err != nil {
				return AST{}, err
			}
			ast.appendDecl(d)
		case TokenEOF:
			break
		default:
			return AST{}, errors.New(tokens[current].Pos.String() + " expected one of Project, Table, Ref, or Enum, got " + tokens[current].Value)
		}
		current++
	}

	return ast, nil
}

func parseProject(tokens Tokens, current int) (Node, int, error) {
	project := &ProjectDecl{}
	project.pos = tokens[current].Pos

	current++
	if tokens[current].Type != TokenIdent {
		return nil, 0, errors.New(tokens[current].Pos.String() + " expected project name, got " + tokens[current].Value)
	}
	project.Name = tokens[current].Value

	current++
	if tokens[current].Type != TokenLBrace {
		return nil, 0, errors.New(tokens[current].Pos.String() + " expected {, got " + tokens[current].Value)
	}
	project.Name = tokens[current].Value

	for current++; tokens[current].Type != TokenRBrace; current++ {
		switch tokens[current].Type {
		case TokenDatabaseType:
			current++
			if tokens[current].Type != TokenColon {
				return nil, 0, errors.New(tokens[current].Pos.String() + " expected :, got " + tokens[current].Value)
			}

			current++
			if tokens[current].Type != TokenString {
				return nil, 0, errors.New(tokens[current].Pos.String() + " expected string, got " + tokens[current].Value)
			}
			project.DBType = tokens[current].Value
		default:
			return nil, 0, errors.New(tokens[current].Pos.String() + " expected one of database_type or Note, got " + tokens[current].Value)
		}
	}

	return project, current, nil
}

func parseTable(tokens Tokens, current int) (Node, int, error) {
	table := &TableDecl{Decls: make([]Node, 0)}
	table.pos = tokens[current].Pos

	current++
	if tokens[current].Type != TokenIdent {
		return nil, 0, errors.New(tokens[current].Pos.String() + " expected table name, got " + tokens[current].Value)
	}
	table.Name = tokens[current].Value

	current++
	if tokens[current].Type == TokenAs {
		current++
		if tokens[current].Type != TokenIdent {
			return nil, 0, errors.New(tokens[current].Pos.String() + " expected table alias, got " + tokens[current].Value)
		}
		table.Alias = tokens[current].Value
	}

	current++
	if tokens[current].Type != TokenLBrace {
		return nil, 0, errors.New(tokens[current].Pos.String() + " expected {, got " + tokens[current].Value)
	}

	var d Node
	var err error
	for current++; tokens[current].Type != TokenRBrace; current++ {
		d, current, err = parseCol(tokens, current)
		if err != nil {
			return nil, 0, err
		}
		table.Decls = append(table.Decls, d)
	}

	return table, current, nil
}

func parseCol(tokens Tokens, current int) (Node, int, error) {
	col := &ColDecl{}
	col.pos = tokens[current].Pos

	if tokens[current].Type != TokenIdent {
		return nil, 0, errors.New(tokens[current].Pos.String() + " expected col name, got " + tokens[current].Value)
	}
	col.Name = tokens[current].Value

	current++
	if tokens[current].Type != TokenIdent {
		return nil, 0, errors.New(tokens[current].Pos.String() + " expected col type, got " + tokens[current].Value)
	}
	col.Type = tokens[current].Value

	return col, current, nil
}

func parseEnum(tokens Tokens, current int) (Node, int, error) {
	enum := &EnumDecl{Values: make([]string, 0)}
	enum.pos = tokens[current].Pos

	current++
	if tokens[current].Type != TokenIdent {
		return nil, 0, errors.New(tokens[current].Pos.String() + " expected enum name, got " + tokens[current].Value)
	}
	enum.Name = tokens[current].Value

	current++
	if tokens[current].Type != TokenLBrace {
		return nil, 0, errors.New(tokens[current].Pos.String() + " expected { name, got " + tokens[current].Value)
	}

	for current++; tokens[current].Type != TokenRBrace; current++ {
		if tokens[current].Type == TokenIdent {
			enum.Values = append(enum.Values, tokens[current].Value)
		} else if tokens[current].Type == TokenQuote {
			enum.Values = append(enum.Values, `"`+tokens[current].Value+`"`)
		} else {
			return nil, 0, errors.New(tokens[current].Pos.String() + " expected an Ident, got " + tokens[current].Value)
		}
	}

	return enum, current, nil
}
