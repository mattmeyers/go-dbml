package dbml

import "github.com/mattmeyers/go-dbml/dbml/token"

type AST struct {
	Project ProjectDecl
	Tables  []TableDecl
	Refs    []RefDecl
	Enums   []EnumDecl
}

type ProjectDecl struct {
	Name   string
	DBType string
	Note   string
}

type TableDecl struct {
	Name    string
	Alias   string
	Cols    []ColDecl
	Indexes []IndexDecl
	Note    string
}
type ColDecl struct {
	Name string
	Type string
}
type SettingDecl struct{}
type IndexDecl struct{}

type RefDecl struct{}

type EnumDecl struct {
	Name   string
	Values []string
}

func newAST() AST {
	return AST{
		Tables: make([]TableDecl, 0),
		Refs:   make([]RefDecl, 0),
		Enums:  make([]EnumDecl, 0),
	}
}

func Parse(tokens token.Tokens) AST {
	ast := newAST()
	current := 0

	var table TableDecl
	var enum EnumDecl

	for current < len(tokens) {
		switch tokens[current].Type {
		case token.Project:
			ast.Project, current = parseProject(tokens, current)
		case token.Table:
			table, current = parseTable(tokens, current)
			ast.Tables = append(ast.Tables, table)
		case token.Enum:
			enum, current = parseEnum(tokens, current)
			ast.Enums = append(ast.Enums, enum)
		case token.EOF:
			break
		default:
			panic("expected one of Project, Table, Ref, or Enum, got " + tokens[current].Value)
		}
		current++
	}

	return ast
}

func parseProject(tokens token.Tokens, current int) (ProjectDecl, int) {
	project := ProjectDecl{}

	current++
	if tokens[current].Type != token.Ident {
		panic("expected project name, got " + tokens[current].Value)
	}
	project.Name = tokens[current].Value

	current++
	if tokens[current].Type != token.LBrace {
		panic("expected {, got " + tokens[current].Value)
	}
	project.Name = tokens[current].Value

	for current++; tokens[current].Type != token.RBrace; current++ {
		switch tokens[current].Type {
		case token.DatabaseType:
			current++
			if tokens[current].Type != token.Colon {
				panic("expected :, got " + tokens[current].Value)
			}

			current++
			if tokens[current].Type != token.String {
				panic("expected string, got " + tokens[current].Value)
			}
			project.DBType = tokens[current].Value
		default:
			panic("expected one of database_type or Note, got " + tokens[current].Value)
		}
	}

	return project, current
}

func parseTable(tokens token.Tokens, current int) (TableDecl, int) {
	table := TableDecl{Cols: make([]ColDecl, 0)}

	current++
	if tokens[current].Type != token.Ident {
		panic("expected table name, got " + tokens[current].Value)
	}
	table.Name = tokens[current].Value

	current++
	if tokens[current].Type == token.As {
		current++
		if tokens[current].Type != token.Ident {
			panic("expected table alias, got " + tokens[current].Value)
		}
		table.Alias = tokens[current].Value
	}

	current++
	if tokens[current].Type != token.LBrace {
		panic("expected {, got " + tokens[current].Value)
	}

	var col ColDecl
	for current++; tokens[current].Type != token.RBrace; current++ {
		col, current = parseCol(tokens, current)
		table.Cols = append(table.Cols, col)
	}

	return table, current
}

func parseCol(tokens token.Tokens, current int) (ColDecl, int) {
	col := ColDecl{}

	if tokens[current].Type != token.Ident {
		panic("expected col name, got " + tokens[current].Value)
	}
	col.Name = tokens[current].Value

	current++
	if tokens[current].Type != token.Ident {
		panic("expected col type, got " + tokens[current].Value)
	}
	col.Type = tokens[current].Value

	return col, current
}

func parseEnum(tokens token.Tokens, current int) (EnumDecl, int) {
	enum := EnumDecl{Values: make([]string, 0)}

	current++
	if tokens[current].Type != token.Ident {
		panic("expected enum name, got " + tokens[current].Value)
	}
	enum.Name = tokens[current].Value

	current++
	if tokens[current].Type != token.LBrace {
		panic("expected { name, got " + tokens[current].Value)
	}

	for current++; tokens[current].Type != token.RBrace; current++ {
		if tokens[current].Type != token.Ident {
			panic("expected an Ident, got " + tokens[current].Value)
		}
		enum.Values = append(enum.Values, tokens[current].Value)
	}

	return enum, current
}
