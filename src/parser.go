package main

import (
	"fmt"
)

type Parser struct {
	lexer *Lexer
	token *Token // current token
}

func (p *Parser) consumeToken() *Token {
	//fmt.Printf("consuming %v\n", *p.token)
	p.token = p.lexer.NextToken()
	return p.token
}

func (p *Parser) peekNextToken() *Token {
	lexer_copy := *p.lexer
	return lexer_copy.NextToken()
}

func (p *Parser) unexpectedToken(type_ string) {
	panic(fmt.Sprintf("unexpected token %v, expected %v", *p.token, type_))
}

func (p *Parser) term(type_ string) {
	if p.token.type_ == type_ {
		p.consumeToken()
	} else {
		p.unexpectedToken(type_)
	}
}

func (p *Parser) ntType() {
	if p.token.type_ == "TYPE" {
		p.consumeToken()
	} else {
		p.unexpectedToken("TYPE")
	}
}

func (p *Parser) ntArgs() {
	p.ntType()
	p.term(IDENT)
	if p.token.type_ == COMMA {
		p.term(COMMA)
		p.ntArgList()
	}
}

func (p *Parser) ntArgList() {
	if p.token.type_ == "TYPE" {
		p.ntArgs()
	}
}

// TODO
func (p *Parser) ntParameterList() {

}

func (p *Parser) ntArrIdent() {
	p.term(IDENT)
	p.term(LBRACKET)
	p.ntExpr()
	p.term(RBRACKET)
}

func (p *Parser) ntFunctionCall() {
	p.term(IDENT)
	p.term(LPAREN)
	p.ntParameterList()
	p.term(RPAREN)
}

// TODO: char
func (p *Parser) ntFactor() {
	if p.token.type_ == LPAREN {
		p.term(LPAREN)
		p.ntExpr()
		p.term(RPAREN)
	} else if p.peekNextToken().type_ == LPAREN {
		p.ntFunctionCall()
	} else if p.token.type_ == IDENT {
		p.term(IDENT)
	} else if p.token.type_ == INTEGER {
		p.term(INTEGER)
	} else if p.peekNextToken().type_ == LBRACKET {
		p.ntArrIdent()
	}
}

func (p *Parser) ntTerm() {
	p.ntFactor()
	if p.token.type_ == STAR {
		p.term(STAR)
		p.ntTerm()
	} else if p.token.type_ == SLASH {
		p.term(SLASH)
		p.ntTerm()
	}
}

func (p *Parser) ntExpr() {
	p.ntTerm()
	if p.token.type_ == PLUS {
		p.term(PLUS)
		p.ntExpr()
	} else if p.token.type_ == MINUS {
		p.term(MINUS)
		p.ntExpr()
	}
}

func (p *Parser) ntReturnStmt() {
	p.term("RETURN_KW")
	if p.token.type_ != SEMICOLON {
		p.ntExpr()
	}
	p.term(SEMICOLON)
}

func (p *Parser) ntCompStmt() {
	p.term(IDENT)
	p.term(COMP_OP)
	p.term(IDENT)
}

func (p *Parser) ntIfElse() {
	p.term("IF_KW")
	p.term(LPAREN)
	p.ntCompStmt()
	p.term(RPAREN)
	p.term(LBRACE)
	p.ntBody()
	p.term(RBRACE)
	p.term("ELSE_KW")
	p.term(LBRACE)
	p.ntBody()
	p.term(RBRACE)
}

func (p *Parser) ntForLoop() {
	p.term("FOR_KW")
	p.term(LPAREN)
	p.ntAssignStmt()
	p.term(SEMICOLON)
	p.ntCompStmt()
	p.term(SEMICOLON)
	p.ntAssignStmt()
	p.term(LBRACE)
	p.ntBody()
	p.term(RBRACE)
}

func (p *Parser) ntAssignStmt() {
	if p.token.type_ == "TYPE" {
		p.ntType()
	}
	p.term(IDENT)
	if p.token.type_ == LBRACKET {
		p.term(LBRACKET)
		p.ntExpr()
		p.term(RBRACKET)
	}
	p.term(ASSIGN)
	p.ntExpr()
	p.term(SEMICOLON)
}

func (p *Parser) ntStmt() {
	if p.token.type_ == "IF_KW" {
		p.ntIfElse()
	} else if p.token.type_ == "FOR_KW" {
		p.ntForLoop()
	} else if p.token.type_ == "RETURN_KW" {
		p.ntReturnStmt()
	} else {
		p.ntAssignStmt()
	}
}

func (p *Parser) ntStmtList() {
	// If the current token is an RBRACE then
	// we've already parsed the last Stmt
	if p.token.type_ != RBRACE {
		p.ntStmt()
		p.ntStmtList()
	}
}

func (p *Parser) ntBody() {
	p.ntStmtList()
}

func (p *Parser) ntFunction() {
	p.ntType()
	p.term(IDENT)
	p.term(LPAREN)
	p.ntArgList()
	p.term(RPAREN)
	p.term(LBRACE)
	p.ntBody()
	p.term(RBRACE)
}

func (p *Parser) ntFunctionList() {
	p.ntFunction()
	if p.token.type_ != EOF {
		p.ntFunctionList()
	}
}

func (p *Parser) ntProgram() {
	p.ntFunctionList()
}

func (p *Parser) Parse() {
	p.ntProgram()
	fmt.Printf("\nParsing done.\n")
}

func Parse(lexer *Lexer) *Parser {
	return &Parser{lexer, lexer.NextToken()}
}
