package evaluator

import (
	"testing"
)

func TestBuiltinFunctions(t *testing.T) {
	InitBuiltins()
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{"let test = {}; push(test, 1, 2); len(test)", 1},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input, t)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		}
	}
}
