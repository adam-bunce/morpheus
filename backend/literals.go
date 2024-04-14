package backend

import (
	"fmt"
	"strconv"
	"strings"
)

type IntLiteral struct {
	IntData
}

func (il IntLiteral) String() string {
	return fmt.Sprintf("%d", il.Value)
}

func (il IntLiteral) Eval(r Runtime) Data {
	return il.IntData
}

func NewIntLiteral(value string) IntLiteral {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic("attempt to assign string to integer")
	}

	return IntLiteral{IntData{
		Value:   intValue,
		Literal: value,
	}}
}

type StringLiteral struct {
	StringData
}

func (sl StringLiteral) String() string {
	return sl.Value
}

func (sl StringLiteral) Eval(r Runtime) Data {
	return sl.StringData
}

func NewStringLiteral(value string) StringLiteral {
	return StringLiteral{
		StringData: StringData{
			Value:   strings.Trim(value, "\"'"),
			Literal: value,
		},
	}
}

type BooleanLiteral struct {
	BooleanData
}

func (bl BooleanLiteral) String() string {
	return fmt.Sprintf("%t", bl.Value)
}

func (bl BooleanLiteral) Eval(r Runtime) Data {
	return bl.BooleanData
}

func NewBooleanLiteral(value string) BooleanLiteral {
	return BooleanLiteral{
		BooleanData: BooleanData{
			Value:   value == "true",
			Literal: value,
		},
	}
}
