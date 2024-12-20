package parser

import (
	"fmt"
	"main/ast"
	"main/token"
)

const (
	_           int = iota
	LOWEST          // _ (black identifier)
	EQUALS          // ==
	LESSGREATER     // > or <
	SUM             // +
	PRODUCT         // *
	PREFIX          // -x or !x
	CALL            // myFunc(x)
)

func (p *Parser) badTokenTypeError(expected token.TokenType) error {
	return fmt.Errorf("error - expected: %s - got: %s", expected, p.currToken.Type)
}

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

func (p *Parser) currTokenIsLegalPrefix() bool {
	return ast.IsLegalPrefixOperator(p.currToken.Type)
}

//	func (p *Parser) peekTokenIs(t token.TokenType) bool {
//		return p.peekToken.Type == t
//	}
func (p *Parser) skipToSemicolon() {
	for !p.currTokenIs(token.SEMICOLON) && !p.currTokenIs(token.EOF) {
		p.nextToken()
	}
}
