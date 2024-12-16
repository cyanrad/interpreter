package parser

import (
	"fmt"
	"main/ast"
	"main/lexer"
	"main/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func CreateParser(l *lexer.Lexer) Parser {
	p := Parser{l: l}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.GetNextToken()
}

func (p *Parser) ParseProgram() (ast.Program, []error) {
	prog := ast.Program{Statements: []ast.Statement{}}
	errors := []error{}

	for !p.currTokenIs(token.EOF) {
		var s ast.Statement
		var err error
		if p.currTokenIs(token.LET) {
			s, err = p.parseLetStatement()
		} else if p.currTokenIs(token.RETURN) {
			s, err = p.parseReturnStatement()
		}

		if err != nil {
			errors = append(errors, err)
			p.skipToSemicolon()
		} else {
			prog.Statements = append(prog.Statements, s)
		}

		p.nextToken()
	}

	return prog, errors
}

// you can assume that the parse functions have the currToken as the first token in it
// you can assume that by the end that the currtoken should equal to ; if viable
func (p *Parser) parseLetStatement() (ast.LetStatement, error) {
	letToken := p.currToken

	p.nextToken()
	if !p.currTokenIs(token.IDENTIFIER) {
		return ast.LetStatement{}, p.badTokenTypeError(token.IDENTIFIER)
	}
	identExp := ast.IdentifierExpression{Token: p.currToken}

	p.nextToken()
	if !p.currTokenIs(token.EQUAL) {
		return ast.LetStatement{}, p.badTokenTypeError(token.EQUAL)
	}

	p.nextToken()
	var valueExp ast.IntExpression
	if p.currTokenIs(token.INT) {
		valueExp, _ = p.parseMathExpression()
	} else {
		return ast.LetStatement{}, fmt.Errorf("error - expected: expression - got: %s", p.currToken.Type)
	}

	p.nextToken()
	if !p.currTokenIs(token.SEMICOLON) {
		return ast.LetStatement{}, p.badTokenTypeError(token.SEMICOLON)
	}

	return ast.LetStatement{
			Token:      letToken,
			Identifier: identExp,
			Expression: valueExp,
		},
		nil
}

func (p *Parser) parseReturnStatement() (ast.ReturnStatement, error) {
	returnToken := p.currToken

	p.nextToken()
	var exp ast.Expression
	if p.currTokenIs(token.IDENTIFIER) {
		exp = ast.IdentifierExpression{Token: p.currToken}
	} else if p.currTokenIs(token.INT) {
		exp, _ = p.parseMathExpression()
	} else {
		return ast.ReturnStatement{}, fmt.Errorf("error - expected: expression - got: %s", p.currToken.Type)
	}

	p.nextToken()
	if !p.currTokenIs(token.SEMICOLON) {
		return ast.ReturnStatement{}, p.badTokenTypeError(token.SEMICOLON)
	}

	return ast.ReturnStatement{
			Token:      returnToken,
			Expression: exp,
		},
		nil
}

func (p *Parser) parseMathExpression() (ast.IntExpression, error) {
	intToken := p.currToken
	return ast.IntExpression{Token: intToken}, nil
}
