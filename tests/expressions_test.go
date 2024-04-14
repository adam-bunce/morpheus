package tests

import (
	"github.com/adam-bunce/morpheus/backend"
	exec "github.com/adam-bunce/morpheus/execute"
	"reflect"
	"testing"
)

func TestAssignExpr(t *testing.T) {
	table := []struct {
		program  string
		name     string
		expected backend.Data
	}{
		{`x = 5;`, "x", backend.IntData{Value: 5, Literal: "5"}},
		{`str = "hello";`, "str", backend.StringData{Value: "hello", Literal: "\"hello\""}},
		{`boolean = true;`, "boolean", backend.BooleanData{Value: true, Literal: "true"}},
	}

	for _, test := range table {
		rt := exec.RunProgram(test.program)

		if !(rt.SymbolTable[test.name] == test.expected) {
			t.Fatalf("expected %s to be %s, got %s", test.name, test.expected, rt.SymbolTable[test.name])
		}
	}
}

func TestDereferenceExpr(t *testing.T) {
	table := []struct {
		program  string
		name     string
		expected backend.Data
	}{
		{`x = 5; y = x;`, "y", backend.IntData{Value: 5, Literal: "5"}},
		{`one = "two"; str = one;`, "str", backend.StringData{Value: "two", Literal: "\"two\""}},
		{`boolean = true; other=boolean`, "other", backend.BooleanData{Value: true, Literal: "true"}},
	}

	for _, test := range table {
		rt := exec.RunProgram(test.program)

		if !(rt.SymbolTable[test.name] == test.expected) {
			t.Fatalf("expected %s to be %s, got %s", test.name, test.expected, rt.SymbolTable[test.name])
		}
	}
}

func TestArithmeticExpr(t *testing.T) {
	table := []struct {
		program  string
		name     string
		expected backend.Data
	}{
		{`x = 5 + 3;`, "x", backend.IntData{Value: 8}},
		{`x = 5 - 3;`, "x", backend.IntData{Value: 2}},
		{`x = 6 * 2;`, "x", backend.IntData{Value: 12}},
		{`x = 6 / 2;`, "x", backend.IntData{Value: 3}},
		{`x = 6 - 10;`, "x", backend.IntData{Value: -4}},
		{`x = 10 * -1;`, "x", backend.IntData{Value: -10}},
		{`x = -10 / -2;`, "x", backend.IntData{Value: 5}},
		{`x = -1 * (1 + 1);`, "x", backend.IntData{Value: -2}},
		{`a = 2; x = a - 2;`, "x", backend.IntData{Value: 0}},
	}

	for i, test := range table {
		rt := exec.RunProgram(test.program)

		if !(rt.SymbolTable[test.name].(backend.IntData).Value == test.expected.(backend.IntData).Value) {
			t.Fatalf("[test %d] expected %s to be %s, got %s", i+1, test.name, test.expected, rt.SymbolTable[test.name])
		}
	}
}

func TestCompareExpr(t *testing.T) {
	table := []struct {
		program  string
		name     string
		expected backend.Data
	}{
		{`x = (5 < 2);`, "x", backend.BooleanData{Value: false}},
		{`x = (5 > 2);`, "x", backend.BooleanData{Value: true}},
		{`x = (true == true);`, "x", backend.BooleanData{Value: true}},
		{`x = (5 == 0);`, "x", backend.BooleanData{Value: false}},
		{`x = ("hi" == "hi");`, "x", backend.BooleanData{Value: true}},
		{`x = ("bye" == "hi");`, "x", backend.BooleanData{Value: false}},
		{`x = ((5 + 8) > 3);`, "x", backend.BooleanData{Value: true}},
		{`x = (false or true)`, "x", backend.BooleanData{Value: true}},
		{`x = (false or false)`, "x", backend.BooleanData{Value: false}},
		{`x = (false and true)`, "x", backend.BooleanData{Value: false}},
		{`x = (true and true)`, "x", backend.BooleanData{Value: true}},
	}

	for i, test := range table {
		rt := exec.RunProgram(test.program)

		if !(rt.SymbolTable[test.name].(backend.BooleanData).Value == test.expected.(backend.BooleanData).Value) {
			t.Fatalf("[test %d] expected %s to be %s, got %s", i+1, test.name, test.expected, rt.SymbolTable[test.name])
		}
	}
}

func TestConcatExpr(t *testing.T) {
	table := []struct {
		program  string
		name     string
		expected backend.Data
	}{
		{`x = "hello" ++ "world";`, "x", backend.StringData{Value: "helloworld"}},
		{`x = "hello" ++ " " ++ "world";`, "x", backend.StringData{Value: "hello world"}},
		{`x = "idk" ++ " " ++ "how else " ++ "to test this";`, "x", backend.StringData{Value: "idk how else to test this"}},
	}

	for i, test := range table {
		rt := exec.RunProgram(test.program)

		if !(rt.SymbolTable[test.name].(backend.StringData).Value == test.expected.(backend.StringData).Value) {
			t.Fatalf("[test %d] expected %s to be %s, got %s", i+1, test.name, test.expected, rt.SymbolTable[test.name])
		}
	}
}

func TestLoopExpr(t *testing.T) {
	table := []struct {
		program  string
		name     string
		expected backend.Data
	}{
		{`
			acc = 0;
			for i in (0, 5, 1) {
				acc = acc + i
			}
			`,
			"acc",
			backend.IntData{Value: 10}},
		{
			`
			acc = 0;
			for i in (0, -10, -1) {
				acc = acc + i
			}
			`,
			"acc",
			backend.IntData{Value: -45},
		},
	}

	for i, test := range table {
		rt := exec.RunProgram(test.program)

		if !(rt.SymbolTable[test.name].(backend.IntData).Value == test.expected.(backend.IntData).Value) {
			t.Fatalf("[test %d] expected %s to be %s, got %s", i+1, test.name, test.expected, rt.SymbolTable[test.name])
		}
	}
}

func TestDeclareExpr(t *testing.T) {
	table := []struct {
		program  string
		name     string
		expected backend.Data
	}{
		{
			program: `function x(one, two, three) {
						  	y = 5;	
						  }`,
			name: "x",
			expected: backend.FunctionData{
				Name: "x",
				Args: []string{"one", "two", "three"},
				Body: backend.Block{
					Exprs: []backend.Expression{
						backend.Assign{
							Name: "y",
							Expr: backend.NewIntLiteral("5"),
						},
					},
				},
			},
		},
	}

	for _, test := range table {
		rt := exec.RunProgram(test.program)

		if !reflect.DeepEqual(rt.SymbolTable[test.name], test.expected) {
			t.Fatalf("actual %v didn't match expected %v", rt.SymbolTable[test.name], test.expected)
		}
	}
}

func TestFunctionCallExpr(t *testing.T) {
	table := []struct {
		program  string
		name     string
		expected backend.Data
	}{
		{
			program: `function add(a, b) {
							a + b;
						  }
					 x = add(5, 6);`,
			name: "x", expected: backend.IntData{Value: 11, Literal: "a + b"}},
		{
			program: `function concat(a, b) {
							x = 0;
							for i in (0, 5, 1) {
								x = x + 1;
							}

							a ++ b;
						  }
					 x = concat("hello ", "world");`,
			name: "x", expected: backend.StringData{Value: "hello world", Literal: "hello world"}},
		{
			program:  `function hi() { "hello" } x = hi();`,
			name:     "x",
			expected: backend.StringData{Value: "hello", Literal: "\"hello\""},
		},
	}

	for _, test := range table {
		rt := exec.RunProgram(test.program)

		if !(rt.SymbolTable[test.name] == test.expected) {
			t.Fatalf("actual %v didn't match expected %v", rt.SymbolTable[test.name], test.expected)
		}
	}
}

func TestIfElifElseExpr(t *testing.T) {
	table := []struct {
		program  string
		name     string
		expected backend.Data
	}{
		{
			// test elif
			program: `
x = 0;
x_str = "empty";
if (x < 0) { x_str =  "x < 0"; }
elif (x > 0) { x_str =  "x < 0"; }
elif (x == 0) {x_str = "x == 0"; }
else { x_str = "error ig"; }
`,
			name:     "x_str",
			expected: backend.StringData{Value: "x == 0"},
		},

		// test first if
		{
			program: `
x = -10;
x_str = "empty";
if (x < 0) { x_str =  "x < 0"; }
elif (x > 0) { x_str =  "x < 0"; }
elif (x == 0) {x_str = "x == 0"; }
else { x_str = "error ig"; }
`,
			name:     "x_str",
			expected: backend.StringData{Value: "x < 0"},
		},

		// test else
		{
			program: `
x = 0;
x_str = "empty";
if (x < 0) { x_str =  "x < 0"; }
elif (x > 0) { x_str =  "x < 0"; }
else { x_str = "x == 0"; }
`,
			name:     "x_str",
			expected: backend.StringData{Value: "x == 0"},
		},
		{
			program: `
x_str = "empty";
if (10 > 0) { x_str = "two"; }
`,
			name:     "x_str",
			expected: backend.StringData{Value: "two"},
		},
		{
			program: `
x_str = "empty";
if (10 > 100) { x_str = "two"; }
elif (100 > 10) { x_str = "three"; }
`,
			name:     "x_str",
			expected: backend.StringData{Value: "three"},
		},
		{
			program: `
x_str = "empty";
if (10 > 100) { x_str = "two"; }
else { x_str = "three"; }
`,
			name:     "x_str",
			expected: backend.StringData{Value: "three"},
		},
	}

	for _, test := range table {
		rt := exec.RunProgram(test.program)

		if !(rt.SymbolTable[test.name].(backend.StringData).Value == test.expected.(backend.StringData).Value) {
			t.Fatalf("actual %v didn't match expected %v", rt.SymbolTable[test.name], test.expected)
		}
	}
}

func TestListExpr(t *testing.T) {
	tests := []struct {
		program  string
		name     string
		expected backend.Data
	}{
		{
			program: `list = [1,2,3];`,
			name:    "list",
			expected: backend.ListData{Values: []backend.Data{
				backend.IntData{Value: 1},
				backend.IntData{Value: 2},
				backend.IntData{Value: 3},
			},
			},
		},
		{
			program:  `value = [1,2,3]; value = value.get(2);`,
			name:     "value",
			expected: backend.IntData{Value: 3},
		},
		{
			program: `value = [1,2,3]; value = value.del(2);`,
			name:    "value",
			expected: backend.ListData{
				Values: []backend.Data{
					backend.IntData{Value: 1},
					backend.IntData{Value: 2},
				},
			},
		},
		{
			program: `value = [1,2,3]; value = value.add(4);`,
			name:    "value",
			expected: backend.ListData{
				Values: []backend.Data{
					backend.IntData{Value: 1},
					backend.IntData{Value: 2},
					backend.IntData{Value: 3},
					backend.IntData{Value: 4},
				},
			},
		},
		{
			program:  `value = [1,2,3]; length = value.len;`,
			name:     "length",
			expected: backend.IntData{Value: 3},
		},
		{
			program:  `value = []; length = value.len;`,
			name:     "length",
			expected: backend.IntData{Value: 0},
		},
		{
			program: `
value = [];
value = value.add(1);
value = value.add(2);
value = value.add(3);
length = value.len
		`,
			name:     "length",
			expected: backend.IntData{Value: 3},
		},
		{
			program: `
value = [];
value.add(1);
value.add(2);
value.add(3);
length = value.len
		`,
			name:     "length",
			expected: backend.IntData{Value: 3},
		},
		{
			program: `
value = [];
value.add(1);
value.add(2);
value.add(3);
length = value.len
		`,
			name:     "length",
			expected: backend.IntData{Value: 3},
		},
		{
			program:  `length = [].len;`,
			name:     "length",
			expected: backend.IntData{Value: 0},
		},
		{
			program:  ` length = [1,2,3].len; `,
			name:     "length",
			expected: backend.IntData{Value: 3},
		},
		{
			program:  `item = [1,2,3].get(2);`,
			name:     "item",
			expected: backend.IntData{Value: 3},
		},
		{
			program: `item = [1,2,3].del(0);`,
			name:    "item",
			expected: backend.ListData{
				Values: []backend.Data{
					backend.IntData{Value: 2},
					backend.IntData{Value: 3},
				},
			},
		},
	}

	for _, test := range tests {
		rt := exec.RunProgram(test.program)

		if !(rt.SymbolTable[test.name].String() == test.expected.String()) {
			t.Fatalf("actual %v didn't match expected %v", rt.SymbolTable[test.name], test.expected)
		}
	}
}
