package main

import (
	"fmt"
    "io/ioutil"
    ggv "github.com/awalterschulze/gographviz"
)

type Parser struct {
	lexer *Lexer
	token *Token // current token
    tree *ggv.Graph
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

func (p *Parser) term(type_ string) *Token {
	if p.token.type_ == type_ {
		return p.consumeToken()
	} else {
		p.unexpectedToken(type_)
	}
    panic("error")
}

func (p *Parser) addNode(name, parent string) string {
    id := RandomString()
    p.tree.AddNode("G", id, map[string]string{"label": name})
    p.tree.AddEdge(parent, id, true, nil)
    return id
}

func (p *Parser) ntType() {
	if p.token.type_ == "TYPE" {
		p.consumeToken()
	} else {
		p.unexpectedToken("TYPE")
	}
}

func (p *Parser) ntArgs(parent string) {
    id := p.addNode("Args", parent)
	p.ntType()
	p.term(IDENT)
	if p.token.type_ == COMMA {
		p.term(COMMA)
		p.ntArgList(id)
	}
}

func (p *Parser) ntArgList(parent string) {
    id := p.addNode("ArgList", parent)
	if p.token.type_ == "TYPE" {
		p.ntArgs(id)
	}
}

// TODO
func (p *Parser) ntParameterList(parent string) {

}

func (p *Parser) ntArrIdent(parent string) {
    id := p.addNode("ArrIdent", parent)
	p.term(IDENT)
	p.term(LBRACKET)
	p.ntExpr(id)
	p.term(RBRACKET)
}

func (p *Parser) ntFunctionCall(parent string) {
    id := p.addNode("FunctionCall", parent)
	p.term(IDENT)
	p.term(LPAREN)
	p.ntParameterList(id)
	p.term(RPAREN)
}

// TODO: char
func (p *Parser) ntFactor(parent string) {
    id := p.addNode("Factor", parent)
	if p.token.type_ == LPAREN {
		p.term(LPAREN)
		p.ntExpr(id)
		p.term(RPAREN)
	} else if p.peekNextToken().type_ == LPAREN {
		p.ntFunctionCall(id)
	} else if p.token.type_ == IDENT {
		p.term(IDENT)
	} else if p.token.type_ == INTEGER {
		p.term(INTEGER)
	} else if p.peekNextToken().type_ == LBRACKET {
		p.ntArrIdent(id)
	}
}

func (p *Parser) ntTerm(parent string) {
    id := p.addNode("Term", parent)
	p.ntFactor(id)
	if p.token.type_ == STAR {
		p.term(STAR)
		p.ntTerm(id)
	} else if p.token.type_ == SLASH {
		p.term(SLASH)
		p.ntTerm(id)
	}
}

func (p *Parser) ntExpr(parent string) {
    id := p.addNode("Expr", parent)
	p.ntTerm(id)
	if p.token.type_ == PLUS {
		p.term(PLUS)
		p.ntExpr(id)
	} else if p.token.type_ == MINUS {
		p.term(MINUS)
		p.ntExpr(id)
	}
}

func (p *Parser) ntReturnStmt(parent string) {
    id := p.addNode("ReturnStmt", parent)
	p.term("RETURN_KW")
	if p.token.type_ != SEMICOLON {
		p.ntExpr(id)
	}
	p.term(SEMICOLON)
}

func (p *Parser) ntCompStmt(parent string) {
    p.addNode("CompStmt", parent)
	p.term(IDENT)
	p.term(COMP_OP)
	p.term(IDENT)
}

func (p *Parser) ntIfElse(parent string) {
    id := p.addNode("IfElse", parent)
	p.term("IF_KW")
	p.term(LPAREN)
	p.ntCompStmt(id)
	p.term(RPAREN)
	p.term(LBRACE)
	p.ntBody(id)
	p.term(RBRACE)
	p.term("ELSE_KW")
	p.term(LBRACE)
	p.ntBody(id)
	p.term(RBRACE)
}

func (p *Parser) ntForLoop(parent string) {
    id := p.addNode("ForLoop", parent)
	p.term("FOR_KW")
	p.term(LPAREN)
	p.ntAssignStmt(id)
	p.term(SEMICOLON)
	p.ntCompStmt(id)
	p.term(SEMICOLON)
	p.ntAssignStmt(id)
	p.term(LBRACE)
	p.ntBody(id)
	p.term(RBRACE)
}

func (p *Parser) ntAssignStmt(parent string) {
    id := p.addNode("AssignStmt", parent)
	if p.token.type_ == "TYPE" {
		p.ntType()
	}
	p.term(IDENT)
	if p.token.type_ == LBRACKET {
		p.term(LBRACKET)
		p.ntExpr(id)
		p.term(RBRACKET)
	}
	p.term(ASSIGN)
	p.ntExpr(id)
	p.term(SEMICOLON)
}

func (p *Parser) ntStmt(parent string) {
    id := p.addNode("Stmt", parent)
	if p.token.type_ == "IF_KW" {
		p.ntIfElse(id)
	} else if p.token.type_ == "FOR_KW" {
		p.ntForLoop(id)
	} else if p.token.type_ == "RETURN_KW" {
		p.ntReturnStmt(id)
	} else {
		p.ntAssignStmt(id)
	}
}

func (p *Parser) ntStmtList(parent string) {
	// If the current token is an RBRACE then
	// we've already parsed the last Stmt
	if p.token.type_ != RBRACE {
        id := p.addNode("StmtList", parent)
		p.ntStmt(id)
		p.ntStmtList(id)
	}
}

func (p *Parser) ntBody(parent string) {
    id := p.addNode("Body", parent)
	p.ntStmtList(id)
}

func (p *Parser) ntFunction(parent string) {
    id := p.addNode("Function", parent)
	p.ntType()
    p.term(IDENT)
	p.term(LPAREN)
	p.ntArgList(id)
	p.term(RPAREN)
	p.term(LBRACE)
	p.ntBody(id)
	p.term(RBRACE)
}

func (p *Parser) ntFunctionList(parent string) {
    id := p.addNode("FunctionList", parent)
	p.ntFunction(id)
	if p.token.type_ != EOF {
		p.ntFunctionList(id)
	}
}

func (p *Parser) ntProgram() {
    p.tree.AddNode("G", "Program", nil)
	p.ntFunctionList("Program")
}

func (p *Parser) Parse() {
	p.ntProgram()
    ioutil.WriteFile("out.dot", []byte(p.tree.String()), 0644)
}

func Parse(lexer *Lexer) *Parser {
    graphAst, _ := ggv.ParseString(`digraph G {}`)
    graph := ggv.NewGraph()
    ggv.Analyse(graphAst, graph)
	return &Parser{lexer, lexer.NextToken(), graph}
}
