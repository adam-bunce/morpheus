package backend

type constraintType int

const (
	Below constraintType = iota
	Above
	Left
	Right
)

type Constraint struct {
	LeftItemName   string
	RightItemName  string
	ConstraintType constraintType
}
