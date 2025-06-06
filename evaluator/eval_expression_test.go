package evaluator

import (
	"main/object"
	"testing"

	"main/lexer"
	"main/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input, t)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input, t)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalPrefixExpression(t *testing.T) {
	intTests := []struct {
		input    string
		expected int64
	}{
		{"-5", -5},
		{"--5", 4},
		{"++5", 6},
	}
	for _, tt := range intTests {
		evaluated := testEval(tt.input, t)
		testIntegerObject(t, evaluated, tt.expected)
	}

	boolTests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
	}
	for _, tt := range boolTests {
		evaluated := testEval(tt.input, t)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalInfixExpression(t *testing.T) {
	intTests := []struct {
		input    string
		expected int64
	}{
		// Arithmetic Operations
		{"5 + 5", 10},
		{"5 - 5", 0},
		{"5 * 5", 25},
		{"12 / 6", 2},
		{"5 % 2", 1},
		{"5 + 5 - 5", 5},
		{"5 * 2 + 10", 20},
		{"10 - 5 * 2", 0},

		// Additional
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},

		// Bitwise Operations
		{"5 & 3", 1},
		{"5 | 3", 7},
	}

	for _, tt := range intTests {
		evaluated := testEval(tt.input, t)
		testIntegerObject(t, evaluated, tt.expected)
	}

	boolTests := []struct {
		input    string
		expected bool
	}{
		// Boolean Operations
		{"true == true", true},
		{"true == false", false},
		{"true != false", true},
		{"false == false", true},

		// Comparison Operations
		{"5 < 10", true},
		{"10 > 5", true},
		{"5 <= 5", true},
		{"5 >= 5", true},
		{"5 < 5", false},
		{"5 > 5", false},
		{"5 <= 4", false},
		{"5 >= 6", false},
		{"5 == 5", true},
		{"5 != 5", false},
		{"5 == 6", false},
		{"5 != 6", true},

		// Logical Operations
		{"true && true", true},
		{"true && false", false},
		{"false || true", true},
		{"false || false", false},

		// Combined Operations
		{"5 + 5 == 10", true},
		{"5 + 5 != 10", false},
		{"5 * 2 < 10", false},

		// With Parentheses
		{"(5 + 5) == 10", true},
		{"(5 + 5) != 10", false},
		{"(5 * 2) < (10 - 1)", false},
		{"(5 * 2 > 10) || (5 < 8)", true},

		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range boolTests {
		evaluated := testEval(tt.input, t)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalIfElseExpression(t *testing.T) {
	intTests := []struct {
		input    string
		expected int64
	}{
		{"if (true) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},

		// nested if
		{"if (2 > 1) { if (true) { 10 } else { 20 } } else { 30 }", 10},

		// else if
		{"if (1 > 2) { 10 } else if (true) { 20 } else { 30 }", 20},
	}

	for _, tt := range intTests {
		evaluated := testEval(tt.input, t)
		testIntegerObject(t, evaluated, tt.expected)
	}

	nullTests := []struct {
		input string
	}{
		{"if (false) { 10 }"},
		{"if (0) { 10 }"},
	}

	for _, tt := range nullTests {
		evaluated := testEval(tt.input, t)
		testNullObject(t, evaluated)
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		// {"fn(x) { x; }(5)", 5},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input, t)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input, t)

	result, ok := evaluated.(*object.ArrayObj)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}
	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
	}
	for _, tt := range tests {
		// t.Log(i, tt.input)
		evaluated := testEval(tt.input, t)

		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
	"one": 10 - 9,
	two: 1 + 1,
	"thr" + "ee": 6 / 2,
	4: 4,
	true: 5,
	false: 6
	}`
	evaluated := testEval(input, t)

	result, ok := evaluated.(*object.HashObj)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}
	expected := map[object.HashKey]int64{
		(&object.StringObj{Value: "one"}).HashKey():   1,
		(&object.StringObj{Value: "two"}).HashKey():   2,
		(&object.StringObj{Value: "three"}).HashKey(): 3,
		(&object.IntegerObj{Value: 4}).HashKey():      4,
		(&object.BooleanObj{Value: true}).HashKey():   5,
		(&object.BooleanObj{Value: false}).HashKey():  6,
	}
	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}
	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input, t)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.BooleanObj)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.IntegerObj)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	_, ok := obj.(*object.NullObj)
	if !ok {
		t.Errorf("object is not Null. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testEval(input string, t *testing.T) object.Object {
	l := lexer.CreateLexer(input)
	p := parser.CreateParser(l)
	program, errs := p.ParseProgram()

	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}

	env := NewEnvironment()
	val, err := Eval(program, env)
	if !err.Ok() {
		panic("Eval returned an error: " + err.Inspect())
	}

	return val
}
