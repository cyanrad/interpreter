package ast

import (
	"main/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}
type Statement interface {
	Node
	statementNode()
}
type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

// returns the token literal of the first statement
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var sb strings.Builder
	for i, s := range p.Statements {
		sb.WriteString(s.String())

		if i < len(p.Statements)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

type BlockStatement struct {
	Token      token.Token // token.LBracket
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var sb strings.Builder

	sb.WriteString("{\n")
	for _, s := range bs.Statements {
		sb.WriteString("\t" + s.String() + "\n")
	}
	sb.WriteString("}")

	return sb.String()
}

type LetStatement struct {
	// Statement
	Token      token.Token // token.LET
	Identifier IdentifierExpression
	Expression Expression
}

func (ls LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls LetStatement) statementNode()       {}
func (ls LetStatement) String() string {
	var sb strings.Builder

	sb.WriteString(ls.TokenLiteral())
	sb.WriteString(" ")
	sb.WriteString(ls.Identifier.TokenLiteral())
	sb.WriteString(" = ")
	sb.WriteString(ls.Expression.TokenLiteral())
	sb.WriteString(";")

	return sb.String()
}

type ReturnStatement struct {
	// Statement
	Token      token.Token // token.RETURN
	Expression Expression
}

func (rs ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs ReturnStatement) statementNode()       {}
func (rs ReturnStatement) String() string {
	var sb strings.Builder

	sb.WriteString(rs.TokenLiteral())
	sb.WriteString(" ")
	sb.WriteString(rs.Expression.String())
	sb.WriteString(";")

	return sb.String()
}

type IfExpression struct {
	// Statement
	Token      token.Token // token.IF
	Conditions []Expression
	Blocks     []BlockStatement // len of Blocks should be Conditions+1 (in case last block is else)
}

func (ie IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie IfExpression) statementNode()       {}
func (ie IfExpression) String() string {
	var sb strings.Builder

	i := 0
	for ; i < len(ie.Conditions); i++ {
		if i != 0 {
			sb.WriteString(" else ")
		}
		sb.WriteString("if ")
		sb.WriteString(ie.Conditions[i].String())
		sb.WriteString(" ")
		sb.WriteString(ie.Blocks[i].String())
	}
	if len(ie.Conditions) < len(ie.Blocks) {
		sb.WriteString(" else ")
		sb.WriteString(ie.Blocks[i].String())
	}
	sb.WriteString(";")

	return sb.String()
}

type ExpressionStatement struct {
	// Statement
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es ExpressionStatement) statementNode()       {}
func (es ExpressionStatement) String() string       { return es.Expression.String() + ";" }

type IdentifierExpression struct {
	// Expression
	Token token.Token // token.Identifier + name
}

func (ie IdentifierExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie IdentifierExpression) expressionNode()      {}
func (ie IdentifierExpression) String() string       { return ie.TokenLiteral() }

type BooleanExpression struct {
	// Expression
	Token token.Token // token.True or token.False
}

func (be BooleanExpression) TokenLiteral() string { return be.Token.Literal }
func (be BooleanExpression) expressionNode()      {}
func (be BooleanExpression) String() string       { return be.TokenLiteral() }

type IntExpression struct {
	// Expression
	Token token.Token // token.INT + value
}

func (ie IntExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie IntExpression) expressionNode()      {}
func (ie IntExpression) String() string       { return ie.TokenLiteral() }

type PrefixExpression struct {
	// Expression
	Token      token.Token // operator toke (e.g. MINUS, EXCLAMATION)
	Expression Expression
}

func (pe PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe PrefixExpression) expressionNode()      {}
func (pe PrefixExpression) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(pe.TokenLiteral())
	sb.WriteString(pe.Expression.String())
	sb.WriteString(")")

	return sb.String()
}

type InfixExpression struct {
	// Expression
	Token token.Token // operator toke (e.g. MINUS, EXCLAMATION)
	Left  Expression
	Right Expression
}

func (ie InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie InfixExpression) expressionNode()      {}
func (ie InfixExpression) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(ie.Left.String())
	sb.WriteString(" ")
	sb.WriteString(ie.TokenLiteral())
	sb.WriteString(" ")
	sb.WriteString(ie.Right.String())
	sb.WriteString(")")

	return sb.String()
}

type CallExpression struct {
	// Expression
	Token      token.Token // the function identifier
	Identifier IdentifierExpression
	Args       []Expression
}

func (ce CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce CallExpression) expressionNode()      {}
func (ce CallExpression) String() string {
	var sb strings.Builder

	sb.WriteString(ce.TokenLiteral() + "(" + ce.Args[0].String())
	if len(ce.Args) != 0 {
		for i := 1; i < len(ce.Args); i++ {
			sb.WriteString(", " + ce.Args[i].String())
		}
	}
	sb.WriteString(")")

	return sb.String()
}
