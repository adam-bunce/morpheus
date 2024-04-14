package backend

import (
	"cmp"
	"fmt"
	"github.com/adam-bunce/morpheus/util"
	"os"
	"strings"
)

type Expression interface {
	String() string
	Eval(runtime Runtime) Data
}

type Assign struct {
	Name string
	Expr Expression
}

func (a Assign) String() string {
	return fmt.Sprintf("%s = %s", a.Name, a.Expr)
}

func (a Assign) Eval(r Runtime) Data {
	r.SymbolTable[a.Name] = a.Expr.Eval(r)
	return NoData{}
}

type Block struct {
	Exprs []Expression
}

func (b Block) String() string {
	var sb strings.Builder

	for _, expr := range b.Exprs {
		sb.WriteString(expr.String())
	}

	return sb.String()
}

func (b Block) Eval(r Runtime) Data {
	var last Data

	outsideScope := util.DeepCopyMap(r.SymbolTable)

	for _, expr := range b.Exprs {
		last = expr.Eval(r)
	}

	r.SymbolTable = outsideScope

	return last
}

type Dereference struct {
	Name string
}

func (d Dereference) String() string {
	return d.Name
}

func (d Dereference) Eval(r Runtime) Data {
	val, ok := r.SymbolTable[d.Name]
	if !ok {
		panic(fmt.Sprintf("Attempt to dereference uninitialized variable %s", d.Name))
	}

	return val
}

type ArithOp int

const (
	ADD ArithOp = iota
	SUB
	DIV
	MUL
)

type Arithmetic struct {
	Left  Expression
	Right Expression
	Op    ArithOp
}

var ArithOpToStr = map[ArithOp]string{
	ADD: "+",
	SUB: "-",
	DIV: "/",
	MUL: "*",
}

func (a Arithmetic) String() string {
	return fmt.Sprintf("%s %s %s", a.Left, ArithOpToStr[a.Op], a.Right)
}

func (a Arithmetic) Eval(r Runtime) Data {
	rightValue, ok := a.Right.Eval(r).(IntData)
	if !ok {
		panic("Arithmetic Eval given non IntData for right Expr")
	}
	leftValue, ok := a.Left.Eval(r).(IntData)
	if !ok {
		panic("Arithmetic Eval given non IntData for left Expr")
	}

	switch a.Op {
	case ADD:
		return IntData{Value: leftValue.Value + rightValue.Value, Literal: a.String()}
	case SUB:
		return IntData{Value: leftValue.Value - rightValue.Value, Literal: a.String()}
	case DIV:
		return IntData{Value: leftValue.Value / rightValue.Value, Literal: a.String()}
	case MUL:
		return IntData{Value: leftValue.Value * rightValue.Value, Literal: a.String()}
	default:
		panic(fmt.Sprintf("Unknown operation %d", a.Op))
	}

}

type CmpOp int

const (
	LT CmpOp = iota
	GT
	EQ
	AND
	OR
)

var CmpOpToStr = map[CmpOp]string{
	LT:  "<",
	GT:  ">",
	EQ:  "==",
	AND: "and",
	OR:  "or",
}

type Compare struct {
	Left  Expression
	Right Expression
	Op    CmpOp
}

func (c Compare) String() string {
	return fmt.Sprintf("%s %s %s", c.Left, CmpOpToStr[c.Op], c.Right)
}

// ugly ahh function
func (c Compare) Eval(r Runtime) Data {
	left := c.Left.Eval(r)
	right := c.Right.Eval(r)

	// Both Int
	leftInt, okLeft := left.(IntData)
	rightInt, okRight := right.(IntData)
	if okLeft && okRight {
		return CompareData(leftInt.Value, rightInt.Value, c.Op)
	}

	// Both String
	leftString, okLeft := left.(StringData)
	rightString, okRight := right.(StringData)
	if okLeft && okRight {
		return CompareData(leftString.Value, rightString.Value, c.Op)
	}

	// Both Boolean
	leftBool, okLeft := left.(BooleanData)
	rightBool, okRight := right.(BooleanData)
	if okLeft && okRight {
		switch c.Op {
		case EQ:
			return BooleanData{Value: leftBool == rightBool, Literal: fmt.Sprintf("%t", leftBool == rightBool)}
		case AND:
			result := leftBool.Value && rightBool.Value
			return BooleanData{Value: result, Literal: fmt.Sprintf("%t", result)}
		case OR:
			result := leftBool.Value || rightBool.Value
			return BooleanData{Value: result, Literal: fmt.Sprintf("%t", result)}
		default:
			panic(fmt.Sprintf("Operator %s undefined for BooleanData", CmpOpToStr[c.Op]))
		}

	}

	panic("Comparison not supported for given data types")

	return NoData{}
}

func CompareData[T cmp.Ordered](a, b T, op CmpOp) BooleanData {
	var result bool
	switch op {
	case LT:
		result = a < b
	case GT:
		result = a > b
	case EQ:
		result = a == b
	default:
		panic(fmt.Sprintf("unhandled comparison operator %s", CmpOpToStr[op]))
	}

	return BooleanData{Value: result, Literal: fmt.Sprintf("%t", result)}
}

type Concat struct {
	Left  Expression
	Right Expression
}

func (c Concat) String() string {
	return fmt.Sprintf("%s ++ %s", c.Right, c.Left)
}

func (c Concat) Eval(r Runtime) Data {
	leftStringData, ok := c.Left.Eval(r).(StringData)
	if !ok {
		panic("Concat left given non StringData")
	}
	rightStringData, ok := c.Right.Eval(r).(StringData)
	if !ok {
		panic("Concat right given non StringData")
	}

	value := fmt.Sprintf("%s%s", leftStringData.Value, rightStringData.Value)
	return StringData{
		Value:   value,
		Literal: value,
	}
}

type Loop struct {
	Iterator string
	Start    Expression
	Stop     Expression
	Step     Expression
	Body     Block
}

func (l Loop) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("for %s in ", l.Iterator))
	sb.WriteString(fmt.Sprintf("(%s, %s %s) {\n", l.Start, l.Stop, l.Step))
	for _, expr := range l.Body.Exprs {
		sb.WriteString(expr.String())
		sb.WriteString("\n")
	}
	sb.WriteString("}\n")

	return sb.String()
}

func (l Loop) Eval(r Runtime) Data {
	// will panic if wrong type btw
	startInt := l.Start.Eval(r).(IntData).Value
	stopInt := l.Stop.Eval(r).(IntData).Value
	stepInt := l.Step.Eval(r).(IntData).Value

	if stepInt < 0 {
		for i := startInt; i > stopInt; i += stepInt {
			r.SymbolTable[l.Iterator] = IntData{Value: i, Literal: fmt.Sprintf("%d", i)}
			l.Body.Eval(r)
		}
	} else {
		for i := startInt; i < stopInt; i += stepInt {
			r.SymbolTable[l.Iterator] = IntData{Value: i, Literal: fmt.Sprintf("%d", i)}
			l.Body.Eval(r)
		}
	}

	delete(r.SymbolTable, l.Iterator)

	return NoData{} // loop don't return stuff right?
}

type Print struct {
	ToPrint Expression
}

func (p Print) String() string {
	return fmt.Sprintf("print(%s)", p.ToPrint)
}

func (p Print) Eval(r Runtime) Data {
	fmt.Println(p.ToPrint.Eval(r))
	return NoData{}
}

type Declare struct {
	Name string
	Args []string
	Body Expression
}

func (d Declare) String() string {
	var sb strings.Builder

	sb.WriteString(d.Name)
	sb.WriteString("(")

	for _, arg := range d.Args {
		sb.WriteString(arg)
		sb.WriteString(", ")
	}
	sb.WriteString(") {\n")
	sb.WriteString(d.Body.String())
	sb.WriteString("\n")

	sb.WriteString("\n}")

	return sb.String()
}

func (d Declare) Eval(r Runtime) Data {
	r.SymbolTable[d.Name] = FunctionData{
		Name: d.Name,
		Args: d.Args,
		Body: d.Body,
	}

	return NoData{}
}

type FunctionCall struct {
	Name string
	Args []Expression
}

func (fc FunctionCall) String() string {
	var sb strings.Builder
	sb.WriteString(fc.Name)
	sb.WriteString("(")
	for _, expr := range fc.Args {
		sb.WriteString(expr.String())
		sb.WriteString(", ")
	}
	sb.WriteString(")")

	return sb.String()
}

func (fc FunctionCall) Eval(r Runtime) Data {
	f := r.SymbolTable[fc.Name]
	if f == nil {
		panic(fmt.Sprintf("function %s doesn't exist", fc.Name))
	}
	funcData, ok := f.(FunctionData)
	if !ok {
		panic(fmt.Sprintf("function %s is not type FunctionData is %T", fc.Name, f))
	}
	if len(funcData.Args) != len(fc.Args) {
		panic(fmt.Sprintf("function %s expects %d args got %d", fc.Name, len(funcData.Args), len(fc.Args)))
	}

	var functionArgs = map[string]Data{}
	for i, arg := range fc.Args {
		functionArgs[funcData.Args[i]] = arg.Eval(r) // assignment to en
	}

	if len(functionArgs) > 0 {
		return funcData.Body.Eval(r.SubScope(functionArgs))
	}

	return funcData.Body.Eval(r)
}

// Conditional is util not an expr
type Conditional struct {
	Condition Expression
	Body      Expression
}

func (c Conditional) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("cond(%s)", c.Condition))
	sb.WriteString(c.Body.String())

	return sb.String()
}

type IfElifElse struct {
	If     Conditional
	ElseIf []Conditional
	Else   Expression
}

func (iee IfElifElse) String() string {
	var sb strings.Builder

	sb.WriteString(iee.If.String())
	for _, elseIf := range iee.ElseIf {
		sb.WriteString(elseIf.String() + "\n")
	}

	if iee.Else != nil {
		sb.WriteString(iee.Else.String())
	}

	return sb.String()
}

func (iee IfElifElse) Eval(r Runtime) Data {
	// if
	ifConditionResult, ok := iee.If.Condition.Eval(r).(BooleanData)
	if !ok {
		panic("if condition should return BooleanData")
	}
	if ifConditionResult.Value {
		return iee.If.Body.Eval(r)
	}

	// elif's
	for i, elseIf := range iee.ElseIf {
		elseIfConditionResult, ok := elseIf.Condition.Eval(r).(BooleanData)
		if !ok {
			panic(fmt.Sprintf("%d'th elif condition should return BooleanDat", i))
		}

		if elseIfConditionResult.Value {
			return iee.ElseIf[i].Body.Eval(r)
		}
	}

	// else
	if iee.Else != nil {
		return iee.Else.Eval(r)
	}

	return NoData{}
}

type List struct {
	Values []Expression
}

func (l List) String() string {
	var sb strings.Builder

	sb.WriteString("[ ")
	for _, val := range l.Values {
		sb.WriteString(val.String() + " ")
	}
	sb.WriteString("]")

	return sb.String()
}

func (l List) Eval(r Runtime) Data {
	var list ListData
	for _, val := range l.Values {
		list.Values = append(list.Values, val.Eval(r))
	}

	return list
}

type ListIndex struct {
	List     Expression
	Position Expression
}

func (li ListIndex) String() string {
	return fmt.Sprintf("%s.get(%d)", li.List.String(), li.Position)
}

func (li ListIndex) Eval(r Runtime) Data {
	pos := li.Position.Eval(r)

	posInt, ok := pos.(IntData)
	if !ok {
		panic("list index position must be IntData")
	}

	valList, ok := li.List.Eval(r).(ListData)
	if !ok {
		panic("attempt to index non-ListData type")
	}

	if len(valList.Values) < posInt.Value {
		panic("attempt to index position outside of list")
	}

	return valList.Values[posInt.Value]
}

type ListDelete struct {
	List     Expression
	Position Expression
}

func (ld ListDelete) String() string {
	return fmt.Sprintf("%s.del(%d)", ld.List, ld.Position)
}

func (ld ListDelete) Eval(r Runtime) Data {
	pos := ld.Position.Eval(r)
	posInt, ok := pos.(IntData)
	if !ok {
		panic("list delete position must be IntData")
	}

	val := ld.List.Eval(r)
	valList, ok := val.(ListData)
	if !ok {
		panic("attempt to delete list index non-ListData type")
	}

	if len(valList.Values) < posInt.Value {
		panic("attempt to index position outside of list")
	}

	var newList ListData
	for i, item := range valList.Values {
		if i != posInt.Value {
			newList.Values = append(newList.Values, item)
		}
	}

	switch list := ld.List.(type) {
	case Dereference:
		// update runtime
		r.SymbolTable[list.Name] = newList
	case List:
		// don't update runtime
	default:
		panic(fmt.Sprintf("ListDelete List is not Defrefrence or List type got:%T", list))
	}

	return newList
}

type ListAdd struct {
	List  Expression
	Value Expression
}

func (la ListAdd) String() string {
	return fmt.Sprintf("%s.add(%s)", la.List, la.Value)
}

func (la ListAdd) Eval(r Runtime) Data {
	valueData := la.Value.Eval(r)

	runtimeData := la.List.Eval(r)
	runtimeListData, ok := runtimeData.(ListData)
	if !ok {
		panic("attempt to .add to non-ListData type")
	}

	var newList ListData
	for _, value := range append(runtimeListData.Values, valueData) {
		newList.Values = append(newList.Values, value)
	}

	switch list := la.List.(type) {
	case Dereference:
		// update runtime
		r.SymbolTable[list.Name] = newList
	case List:
		// don't update runtime
	default:
		panic(fmt.Sprintf("ListAdd List is not Defrefrence or List type got:%T", list))
	}

	return newList
}

type ListLength struct {
	List Expression
}

func (ll ListLength) String() string {
	return fmt.Sprintf("%s.len", ll.List)
}

func (ll ListLength) Eval(r Runtime) Data {
	runtimeData := ll.List.Eval(r)

	runtimeListData, ok := runtimeData.(ListData)
	if !ok {
		panic("attempt get length of non-ListData type")
	}

	return IntData{
		Value: len(runtimeListData.Values),
	}
}

type BoxExpr struct {
	Id string
}

func (b BoxExpr) String() string {
	return fmt.Sprintf("box: %s", b.Id)
}

func (b BoxExpr) Eval(r Runtime) Data {
	return NewBox(r.Solver, b.Id)
}

type GroupExpr struct {
	Items       Expression
	Constraints []Constraint
}

func (g GroupExpr) String() string {
	return g.String()
}

func (g GroupExpr) Eval(r Runtime) Data {
	for _, c := range g.Constraints {
		right := Dereference{strings.Trim(c.RightItemName, "*")}.Eval(r)
		left := Dereference{strings.Trim(c.LeftItemName, "*")}.Eval(r)

		// might get expression that doesnt instantly give us a a layout item
		// variable -> function -> returns a layout item
		switch right.(type) {
		case FunctionData:
			right = right.(FunctionData).Body.Eval(r)
		}
		switch left.(type) {
		case FunctionData:
			left = left.(FunctionData).Body.Eval(r)
		}

		switch c.ConstraintType {
		case Below:
			left.(LayoutItem).IsBelow(right.(LayoutItem))
		case Above:
			left.(LayoutItem).IsAbove(right.(LayoutItem))
		case Left:
			left.(LayoutItem).IsLeftOf(right.(LayoutItem))
		case Right:
			left.(LayoutItem).IsRightOf(right.(LayoutItem))
		}
	}

	items := g.Items.Eval(r).(ListData)
	var layoutItems []LayoutItem
	for _, item := range items.Values {
		layoutItems = append(layoutItems, item.(LayoutItem))
	}

	return NewGroup(r.Solver, layoutItems...)
}

type Htmlify struct {
	Layout Expression
	File   string
}

func (h Htmlify) String() string {
	return fmt.Sprintf("Htmlify( %s )", h.Layout)
}

func (h Htmlify) Eval(r Runtime) Data {
	top := `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Layout</title>

</head>
<body>
`
	bottom := `
</body>
</html>
`

	layout := h.Layout.Eval(r)
	html := layout.(LayoutItem).AsHtml()

	file, err := os.Create(strings.Trim(h.File, "\"") + ".html")
	if err != nil {
		panic("error htmlifying")
	}
	defer file.Close()

	_, err = file.WriteString(top + html + bottom)
	if err != nil {
		panic("error writing to file")
	}

	return NoData{}
}
