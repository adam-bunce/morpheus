package tests

import (
	"github.com/adam-bunce/morpheus/backend"
	exec "github.com/adam-bunce/morpheus/execute"
	"github.com/lithdew/casso"
	"testing"
)

// NOTE: i have no idea how to test layout stuff, so this is a mess o7

func TestBox_IsLeftOf(t *testing.T) {
	s := casso.NewSolver()

	b1 := backend.NewBox(s, "b1")
	b2 := backend.NewBox(s, "b2")

	b1.IsLeftOf(b2)

	b1x := s.Val(b1.CX)
	b1w := s.Val(b1.CW)
	b2x := s.Val(b2.CX)

	if !(b1x+b1w <= b2x) {
		t.Fatalf("%f+%f (%f) should be <= %f", b1x, b1w, b1x+b1w, b2x)
	}
}

func TestBox_IsRightOf(t *testing.T) {
	s := casso.NewSolver()

	b1 := backend.NewBox(s, "b3")
	b2 := backend.NewBox(s, "b4")

	b1.IsRightOf(b2)

	b1x := s.Val(b1.CX)
	b2x := s.Val(b2.CX)
	b2w := s.Val(b2.CW)

	if !(b1x >= b2x+b2w) {
		t.Fatalf("%f should be >= %f + %f (%f)", b1x, b2x, b2w, b2x+b2w)
	}
}

func TestBox_IsAbove(t *testing.T) {
	s := casso.NewSolver()

	b1 := backend.NewBox(s, "b1")
	b2 := backend.NewBox(s, "b2")

	b1.IsAbove(b2)

	if !(s.Val(b1.CY)+s.Val(b1.CH) >= s.Val(b2.CY)) {
		t.Fatalf("b1 should be above b2")
	}
}

func TestBox_IsBelow(t *testing.T) {

	s := casso.NewSolver()

	b1 := backend.NewBox(s, "b1")
	b2 := backend.NewBox(s, "b2")

	b1.IsBelow(b2)

	if !(s.Val(b1.CY) <= s.Val(b2.CY)+s.Val(b2.CH)) {
		t.Fatalf("b1 should be below b2")
	}

}

func TestBoxCreation(t *testing.T) {

	program := `
g = Box("box1");
`

	exec.RunProgram(program)
}

func TestGroupCreation(t *testing.T) {

	program := `
a = Box("box1");
b = Box("box2");
c = Box("box3");
d = Box("box4");

g = Group([a, b, c, d] : [
	*a is below *b,
	*c is below *a,
	*d is right of *a,
	*d is below *b
]);

`

	exec.RunProgram(program)

}

func TestGroupedGroupCreation(t *testing.T) {

	program := `
function create() {
	a = Box("box1");
	b = Box("box2");

	Group([a, b] : [*a is below *b])
}

function create_other() {
	c = Box("C");
	d = Box("D");
	e = Box("E");

	Group([c,d, e] : [
		*c is below *e,
		*e is left of *c,
		*d is left of *c
	])
}

function group_again() {
	group_a = create();
	group_b = create_other();
	group_c = create_other();

	Group([group_a, group_b, group_c] : [
		*group_a is right of *group_b,
		*group_a is below *group_b,
		*group_c is below *group_a
	])
}

super_group_a = group_again();
super_group_b = group_again();

super_group = Group([super_group_a, super_group_b] : [
	*super_group_a is below *super_group_b,
	*super_group_a is right of *super_group_b
])

super_group.htmlify("test")
`

	exec.RunProgram(program)
}
