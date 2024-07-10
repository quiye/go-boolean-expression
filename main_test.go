package main

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestExpr_eval(t *testing.T) {
	type fields struct {
		op        Operator
		childs    []*Expr
		leafValue int
	}
	type args struct {
		input []int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "And",
			fields: fields{
				op: And,
				childs: []*Expr{
					{
						op:        Leaf,
						childs:    nil,
						leafValue: 1,
					},
					{
						op:        Leaf,
						childs:    nil,
						leafValue: 2,
					},
				},
				leafValue: 0,
			},
			args: args{
				input: []int{1, 2},
			},
			want: true,
		},
		{
			name: "Or",
			fields: fields{
				op: Or,
				childs: []*Expr{
					{
						op:        Leaf,
						childs:    nil,
						leafValue: 1,
					},
					{
						op:        Leaf,
						childs:    nil,
						leafValue: 2,
					},
				},
				leafValue: 0,
			},
			args: args{
				input: []int{1},
			},
			want: true,
		},
		{
			name: "Not",
			fields: fields{
				op: Not,
				childs: []*Expr{
					{
						op:        Leaf,
						childs:    nil,
						leafValue: 1,
					},
				},
				leafValue: 0,
			},
			args: args{
				input: []int{1},
			},
			want: false,
		},
		{
			name: "Not (success)",
			fields: fields{
				op: Not,
				childs: []*Expr{
					{
						op:        Leaf,
						childs:    nil,
						leafValue: 1,
					},
				},
				leafValue: 0,
			},
			args: args{
				input: []int{2},
			},
			want: true,
		},
		{
			name: "Leaf",
			fields: fields{
				op:        Leaf,
				childs:    nil,
				leafValue: 1,
			},
			args: args{
				input: []int{1},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Expr{
				op:        tt.fields.op,
				childs:    tt.fields.childs,
				leafValue: tt.fields.leafValue,
			}
			if got := e.Eval(tt.args.input); got != tt.want {
				t.Errorf("Expr.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}

// function for random generation of test data
// use random for generating random data and random operators
func generateRandomExpr(depth int) *Expr {
	if depth == 0 {
		return NewLeafExpr(rand.Intn(10))
	}
	switch rand.Intn(4) {
	case 0:
		return NewAndExpr([]*Expr{generateRandomExpr(depth - 1), generateRandomExpr(depth - 1)})
	case 1:
		return NewOrExpr([]*Expr{generateRandomExpr(depth - 1), generateRandomExpr(depth - 1)})
	case 2:
		return NewNotExpr(generateRandomExpr(depth - 1))
	default:
		return NewLeafExpr(rand.Intn(10))
	}
}

func TestExpr_evalRandom(t *testing.T) {
	for range 10 {
		e := generateRandomExpr(3)
		input := []int{rand.Intn(10), rand.Intn(10), rand.Intn(10)}
		result := e.Eval(input)
		// show the generated expression and input and the result
		t.Logf("\ninput: %v\nresult: %v\nexpr: ", input, result)
		printExpr(e)
		println()
		println()
	}
}

func printExpr(e *Expr) {
	switch e.op {
	case And:
		print("(")
		for i, child := range e.childs {
			printExpr(child)
			if i < len(e.childs)-1 {
				print(" AND ")
			}
		}
		print(")")
	case Or:
		print("(")
		for i, child := range e.childs {
			printExpr(child)
			if i < len(e.childs)-1 {
				print(" OR ")
			}
		}
		print(")")
	case Not:
		print("NOT ")
		printExpr(e.childs[0])
	case Leaf:
		if e.isNegatedLeaf {
			print("Neg_")
		}
		print(e.leafValue)
	}
}

func TestApplyDeMorgansLaw(t *testing.T) {
	type args struct {
		e    *Expr
		want *Expr
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "And",
			args: args{
				e: NewNotExpr(NewAndExpr([]*Expr{
					NewLeafExpr(1),
					NewLeafExpr(2),
					NewLeafExpr(3),
				})),
				want: NewOrExpr([]*Expr{
					NewNegatedLeafExpr(1),
					NewNegatedLeafExpr(2),
					NewNegatedLeafExpr(3),
				}),
			},
		},
		{
			name: "Or",
			args: args{
				e: NewNotExpr(NewOrExpr([]*Expr{
					NewLeafExpr(1),
					NewLeafExpr(2),
					NewLeafExpr(3),
				})),
				want: NewAndExpr([]*Expr{
					NewNegatedLeafExpr(1),
					NewNegatedLeafExpr(2),
					NewNegatedLeafExpr(3),
				}),
			},
		},
		{
			name: "Complex",
			args: args{
				e: NewNotExpr(NewAndExpr([]*Expr{
					NewLeafExpr(1),
					NewOrExpr([]*Expr{
						NewLeafExpr(2),
						NewLeafExpr(3),
					}),
				})),
				want: NewOrExpr([]*Expr{
					NewNegatedLeafExpr(1),
					NewAndExpr([]*Expr{
						NewNegatedLeafExpr(2),
						NewNegatedLeafExpr(3),
					}),
				}),
			},
		},
		{
			name: "Not_Not_Not",
			args: args{
				e:    NewNotExpr(NewNotExpr(NewNotExpr(NewLeafExpr(1)))),
				want: NewNegatedLeafExpr(1),
			},
		},
		{
			name: "Not_Not",
			args: args{
				e:    NewNotExpr(NewNotExpr(NewLeafExpr(1))),
				want: NewLeafExpr(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ApplyDeMorgansLaw(tt.args.e)
			if !reflect.DeepEqual(tt.args.e, tt.args.want) {
				t.Errorf("ApplyDeMorgansLaw() = %+v, want %+v", tt.args.e, tt.args.want)
			}
		})
	}
}

func TestExpr_evalRandom2(t *testing.T) {
	for range 10 {
		e := generateRandomExpr(5)
		input := []int{rand.Intn(10), rand.Intn(10), rand.Intn(10)}
		result := e.Eval(input)
		// show the generated expression and input and the result
		t.Logf("\ninput: %v\nresult: %v\nexpr: ", input, result)
		printExpr(e)
		println()
		ApplyDeMorgansLaw(e)
		printExpr(e)
		println()
	}
}
