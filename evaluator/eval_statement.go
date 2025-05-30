package evaluator

import (
	"main/ast"
	"main/object"
)

func EvalStatement(s ast.Statement, env Environment) object.Object {
	switch stmt := s.(type) {
	case ast.BlockStatement:
		return evalBlockStatement(stmt, env)
	case ast.ExpressionStatement:
		return EvalExpression(stmt.Expression, env)
	case ast.LetStatement:
		return evalLetStatement(stmt, env)
	case ast.ReturnStatement:
		return evalReturnStatement(stmt, env)
	default:
		panic("unknown statement type")
	}
}

func evalBlockStatement(block ast.BlockStatement, env Environment) object.Object {
	prog := ast.Program{Statements: block.Statements}
	return Eval(prog, env)
}

func evalLetStatement(stmt ast.LetStatement, env Environment) object.Object {
	// Check if the variable already exists in the environment
	ident := stmt.Identifier.TokenLiteral()
	existingVar := env.Get(ident)
	if existingVar != nil {
		panic("variable: " + ident + " already exists")
	}

	env.Create(ident, EvalExpression(stmt.Expression, env))
	return object.NullObj{}
}

func evalReturnStatement(stmt ast.ReturnStatement, env Environment) object.Object {
	if stmt.Expression == nil {
		return object.NullObj{}
	}

	return EvalExpression(stmt.Expression, env)
}
