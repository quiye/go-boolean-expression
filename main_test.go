package main

import "testing"

func TestExpr_eval(t *testing.T) {
	type fields struct {
		op        Operator
		childs    []Expr
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
				childs: []Expr{
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
				childs: []Expr{
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
				childs: []Expr{
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
				childs: []Expr{
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
