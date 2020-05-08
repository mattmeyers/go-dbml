package dbml

type Project struct {
	Name   string
	DBType string
	Note   string
	Tables []Table
}

type Table struct {
	Name    string
	Alias   string
	Note    string
	Columns []Col
}

type Col struct {
	Name string
	Type string
}

type Enum struct {
	Name   string
	Values []string
}
