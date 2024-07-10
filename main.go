package main

import "slices"

func main() {
	// hello world
	println("Hello, World!")
}

type Operator int

const (
	And Operator = iota
	Or
	Not
	Leaf
)

type Expr struct {
	op            Operator
	childs        []*Expr
	leafValue     int
	isNegatedLeaf bool
}

func NewAndExpr(childs []*Expr) *Expr {
	return &Expr{op: And, childs: childs}
}

func NewOrExpr(childs []*Expr) *Expr {
	return &Expr{op: Or, childs: childs}
}

func NewNotExpr(child *Expr) *Expr {
	return &Expr{op: Not, childs: []*Expr{child}}
}

func NewLeafExpr(value int) *Expr {
	return &Expr{op: Leaf, leafValue: value}
}

func NewNegatedLeafExpr(value int) *Expr {
	return &Expr{op: Leaf, leafValue: value, isNegatedLeaf: true}
}

func (e *Expr) Eval(input []int) bool {
	switch e.op {
	case And:
		for _, child := range e.childs {
			if !child.Eval(input) {
				return false
			}
		}
		return true
	case Or:
		for _, child := range e.childs {
			if child.Eval(input) {
				return true
			}
		}
		return false
	case Not:
		return !e.childs[0].Eval(input)
	case Leaf:
		return slices.Index(input, e.leafValue) != -1
	}

	// fixme: error handling
	return false
}

func ApplyDeMorgansLaw(e *Expr) {
	switch e.op {
	case Not:
		child := e.childs[0]
		switch child.op {
		case Or:
			e.op = And
			e.childs = make([]*Expr, len(child.childs))
			for i, c := range child.childs {
				e.childs[i] = NewNotExpr(c)
				ApplyDeMorgansLaw(e.childs[i])
			}
		case And:
			e.op = Or
			e.childs = make([]*Expr, len(child.childs))
			for i, c := range child.childs {
				e.childs[i] = NewNotExpr(c)
				ApplyDeMorgansLaw(e.childs[i])
			}
		case Not:
			*e = *child.childs[0]
			ApplyDeMorgansLaw(e)
		default:
			*e = *child
			e.isNegatedLeaf = true
		}
	case Or, And:
		for _, child := range e.childs {
			ApplyDeMorgansLaw(child)
		}
	default:
		return
	}
}
