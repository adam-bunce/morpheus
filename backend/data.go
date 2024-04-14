package backend

import (
	"fmt"
	"strings"
)

type Data interface {
	String() string
}

type IntData struct {
	Value   int
	Literal string
}

func (id IntData) String() string { return fmt.Sprintf("IntData:%v", id.Value) }

type BooleanData struct {
	Value   bool
	Literal string
}

func (bd BooleanData) String() string { return fmt.Sprintf("BooleanData:%v", bd.Value) }

type StringData struct {
	Value   string
	Literal string
}

func (sd StringData) String() string { return fmt.Sprintf("StringData:'%v'", sd.Value) }

type FunctionData struct {
	Name string
	Args []string
	Body Expression
}

func (f FunctionData) String() string {
	var sb strings.Builder

	sb.WriteString(f.Name)
	sb.WriteString("(")

	for _, arg := range f.Args {
		sb.WriteString(arg)
		sb.WriteString(", ")
	}
	sb.WriteString(") {\n")
	sb.WriteString(f.Body.String())
	sb.WriteString("\n")

	sb.WriteString("\n}")

	return sb.String()
}

type NoData struct{}

func (nd NoData) String() string { return "NoData" }

type ListData struct {
	Values []Data
}

func (ld ListData) String() string {
	var sb strings.Builder

	sb.WriteString("ListData:[ ")
	for _, val := range ld.Values {
		sb.WriteString(val.String() + " ")
	}
	sb.WriteString("]")

	return sb.String()
}

func (ld ListData) Index(r Runtime, position Expression) Data {
	return ld.Values[position.Eval(r).(IntData).Value]
}
