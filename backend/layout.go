package backend

import (
	"fmt"
	"github.com/lithdew/casso"
	"strings"
)

type LayoutItem interface {
	RightEdge() float64
	LeftEdge() float64
	Top() float64
	Bottom() float64
	GetX() casso.Symbol
	GetY() casso.Symbol
	GetW() casso.Symbol
	GetH() casso.Symbol

	IsBelow(item LayoutItem)
	IsAbove(item LayoutItem)

	IsRightOf(item LayoutItem)
	IsLeftOf(item LayoutItem)

	String() string
	AsHtml() string
}

type Box struct {
	solver *casso.Solver
	Id     string

	CX casso.Symbol
	CY casso.Symbol
	CW casso.Symbol
	CH casso.Symbol
}

func NewBox(s *casso.Solver, id string) Box {
	// all boxes default 50x50 might add size params later
	// width constraint
	bw := casso.New()
	casso.NewConstraint(casso.EQ, -50, bw.T(1))
	s.AddConstraint(casso.NewConstraint(casso.EQ, -50, bw.T(1)))

	// height constraint
	bh := casso.New()
	casso.NewConstraint(casso.EQ, -50, bh.T(1))
	s.AddConstraint(casso.NewConstraint(casso.EQ, -50, bh.T(1)))

	bx, by := casso.New(), casso.New()
	return Box{
		solver: s,
		Id:     id,
		CX:     bx,
		CY:     by,
		CW:     bw,
		CH:     bh,
	}
}

func (b Box) GetX() casso.Symbol { return b.CX }
func (b Box) GetY() casso.Symbol { return b.CY }
func (b Box) GetW() casso.Symbol { return b.CW }
func (b Box) GetH() casso.Symbol { return b.CH }

func (b Box) AsHtml() string {
	var sb strings.Builder

	sb.WriteString("<div ")
	sb.WriteString("style=\"")
	sb.WriteString("border: solid grey 1px;")
	sb.WriteString("position: absolute;")
	sb.WriteString(fmt.Sprintf("top: %.fpx;", b.Top()))
	sb.WriteString(fmt.Sprintf("left: %.fpx;", b.LeftEdge()))
	sb.WriteString(fmt.Sprintf("width: %.fpx;", b.RightEdge()-b.LeftEdge()))
	sb.WriteString(fmt.Sprintf("height: %.fpx;", b.Bottom()-b.Top()))
	sb.WriteString("\">")
	sb.WriteString(fmt.Sprintf("BOX %s", b.Id))
	sb.WriteString("</div>")

	return sb.String()
}

func (b Box) RightEdge() float64 { return b.solver.Val(b.CX) + b.solver.Val(b.CW) }
func (b Box) LeftEdge() float64  { return b.solver.Val(b.CX) }
func (b Box) Top() float64       { return b.solver.Val(b.CY) }
func (b Box) Bottom() float64    { return b.solver.Val(b.CY) + b.solver.Val(b.CH) }

func (b Box) String() string {
	return fmt.Sprintf("BOX{X:%.2f Y:%.2f W:%.2f H:%.2f, ID:%s}",
		b.solver.Val(b.CX),
		b.solver.Val(b.CY),
		b.solver.Val(b.CW),
		b.solver.Val(b.CH),
		b.Id)
}

func (b Box) IsLeftOf(item LayoutItem) {
	// b.X + b.W <= item.x
	_, err := b.solver.AddConstraint(casso.NewConstraint(casso.LTE, 0, b.CX.T(1), b.CW.T(1), item.GetX().T(-1)))
	if err != nil {
		panic(fmt.Sprintf("failed to add IsLeftOf constraint. err=%v", err))
	}
}

func (b Box) IsRightOf(item LayoutItem) {
	// b.X >= item.x + item.W
	_, err := b.solver.AddConstraint(casso.NewConstraint(casso.GTE, 0, b.CX.T(1), item.GetX().T(-1), item.GetW().T(-1)))
	if err != nil {
		panic(fmt.Sprintf("failed to add IsRightOf constraint. err=%v", err))
	}
}

func (b Box) IsAbove(item LayoutItem) {
	// b.Y + b.H <= item.Y
	_, err := b.solver.AddConstraint(casso.NewConstraint(casso.LTE, 0, b.CY.T(1), b.CH.T(1), item.GetY().T(-1)))
	if err != nil {
		panic(fmt.Sprintf("failed to add IsAbove constraint. err=%v", err))
	}
}

func (b Box) IsBelow(item LayoutItem) {
	// b.Y >= item.Y + item.H
	_, err := b.solver.AddConstraint(casso.NewConstraint(casso.GTE, 0, b.CY.T(1), item.GetY().T(-1), item.GetH().T(-1)))
	if err != nil {
		panic(fmt.Sprintf("failed to add IsBelow constraint. err=%v", err))
	}
}

type Group struct {
	Items  []LayoutItem
	solver *casso.Solver

	X casso.Symbol
	Y casso.Symbol
	H casso.Symbol
	W casso.Symbol
}

func (g Group) RightEdge() float64 { return g.solver.Val(g.X) + g.solver.Val(g.W) }
func (g Group) LeftEdge() float64  { return g.solver.Val(g.X) }
func (g Group) Top() float64       { return g.solver.Val(g.Y) }
func (g Group) Bottom() float64    { return g.solver.Val(g.Y) + g.solver.Val(g.H) }
func (g Group) GetX() casso.Symbol { return g.X }
func (g Group) GetY() casso.Symbol { return g.Y }
func (g Group) GetW() casso.Symbol { return g.W }
func (g Group) GetH() casso.Symbol { return g.H }

func (g Group) AsHtml() string {
	var sb strings.Builder

	sb.WriteString("<div>")

	for _, item := range g.Items {
		sb.WriteString(item.AsHtml())
	}

	sb.WriteString("</div>")

	return sb.String()
}

func NewGroup(solver *casso.Solver, items ...LayoutItem) Group {
	var maxX, maxY, minX, minY float64

	for _, item := range items {
		if item.Top() < minY {
			minY = item.Top()
		}

		if item.Bottom() > maxY {
			maxY = item.Bottom()
		}

		if item.RightEdge() > maxX {
			maxX = item.RightEdge()
		}

		if item.LeftEdge() < minX {
			minX = item.LeftEdge()
		}
	}

	// size constraints based on children
	groupX := casso.New()
	groupY := casso.New()
	groupW := casso.New()
	groupH := casso.New()

	// groupX = min_x
	solver.AddConstraint(casso.NewConstraint(casso.GTE, -1*minX, groupX.T(1)))
	// groupY = min_y
	solver.AddConstraint(casso.NewConstraint(casso.GTE, -1*minY, groupY.T(1)))
	// groupW = min_x - max_x max-min 100 -50 = w
	solver.AddConstraint(casso.NewConstraint(casso.GTE, -1*(maxX-minX), groupW.T(1)))
	// groupH = min_y - max_y
	solver.AddConstraint(casso.NewConstraint(casso.GTE, -1*(maxY-minY), groupH.T(1)))

	// keep the group in bounds of the screen
	solver.AddConstraint(casso.NewConstraint(casso.GTE, 0, groupX.T(1)))
	solver.AddConstraint(casso.NewConstraint(casso.GTE, 0, groupY.T(1)))

	// when we do a group, each child gets a new constraint applied to its X/Y
	// they both must be greater than or equal to the parent's x ans y
	for _, item := range items {
		// items X and Y must be >= the min of the groups X/Y
		solver.AddConstraint(casso.NewConstraint(casso.GTE, 0, item.GetX().T(1), groupX.T(-1)))
		solver.AddConstraint(casso.NewConstraint(casso.GTE, 0, item.GetY().T(1), groupY.T(-1)))
	}

	hold := Group{
		Items:  items,
		solver: solver,
		X:      groupX,
		Y:      groupY,
		H:      groupH,
		W:      groupW,
	}

	return hold
}

func (g Group) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("GROUP{X:%.2f Y:%.2f W:%.2f H:%.2f",
		g.solver.Val(g.X),
		g.solver.Val(g.Y),
		g.solver.Val(g.W),
		g.solver.Val(g.H),
	))

	sb.WriteString("\nCHILDREN:{\n")

	for _, child := range g.Items {
		sb.WriteString("\t" + child.String())
		sb.WriteString("\n")
	}
	sb.WriteString("}}")

	return sb.String()
}

func (g Group) IsLeftOf(item LayoutItem) {
	// b.X + b.W <= item.x
	_, err := g.solver.AddConstraint(casso.NewConstraint(casso.LTE, 0, g.X.T(1), g.W.T(1), item.GetX().T(-1)))

	if err != nil {
		panic(fmt.Sprintf("failed to add Group IsLeftOf constraint. err=%v", err))
	}
}

func (g Group) IsRightOf(item LayoutItem) {
	// b.X >= item.x + item.W
	_, err := g.solver.AddConstraint(casso.NewConstraint(casso.GTE, 0, g.X.T(1), item.GetX().T(-1), item.GetW().T(-1)))

	if err != nil {
		panic(fmt.Sprintf("failed to add Group IsRightOf constraint. err=%v", err))
	}
}

func (g Group) IsBelow(item LayoutItem) {
	// b.Y >= item.Y + item.H
	_, err := g.solver.AddConstraint(casso.NewConstraint(casso.GTE, 0, g.Y.T(1), item.GetY().T(-1), item.GetH().T(-1)))
	if err != nil {
		panic(fmt.Sprintf("failed to add Group IsBelow constraint. err=%v", err))
	}
}

func (g Group) IsAbove(item LayoutItem) {
	// b.Y + b.H <= item.Y
	_, err := g.solver.AddConstraint(casso.NewConstraint(casso.LTE, 0, g.GetY().T(1), item.GetH().T(1), item.GetY().T(-1)))
	if err != nil {
		panic(fmt.Sprintf("failed to add Group IsAbove constraint. err=%v", err))
	}
}
