package main

import (
	"fmt"
	"io/ioutil"

	ggv "github.com/awalterschulze/gographviz"
)

type Parser struct {
	lexer *Lexer
	token *Token // current token
	tree  *ggv.Graph
}

func (p *Parser) consumeToken() *Token {
	//fmt.Printf("consuming %v\n", *p.token)
	token := p.token
	p.token = p.lexer.NextToken()
	return token
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

func (p *Parser) editNodeName(id, name string) {
	attrs, _ := ggv.NewAttrs(make(map[string]string))
	attrs.Add("label", name)
	p.tree.Nodes.Lookup[id].Attrs.Extend(attrs)
}

func (p *Parser) ntType(parent string) {
	p.addNode(fmt.Sprintf("\"Type (%v)\"", p.token.literal), parent)
	if p.token.type_ == "TYPE" {
		p.consumeToken()
	} else {
		p.unexpectedToken("TYPE")
	}
}

func (p *Parser) ntArgs(parent string) {
	id := p.addNode("Args", parent)
	p.ntType(id)
	p.editNodeName(id, "\"Args ("+p.term(IDENT).literal+")\"")
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

func (p *Parser) ntParameters(parent string) {
	id := p.addNode("Param", parent)
	p.editNodeName(id, "\"Args ("+p.term(IDENT).literal+")\"")
	if p.token.type_ == COMMA {
		p.term(COMMA)
		p.ntParameterList(id)
	}
}

func (p *Parser) ntParameterList(parent string) {
	id := p.addNode("ParameterList", parent)
	if p.token.type_ == IDENT {
		p.ntParameters(id)
	}
}

func (p *Parser) ntArrIdent(parent string) {
	id := p.addNode("\"ArrIdent ("+p.term(IDENT).literal+")\"", parent)
	p.term(LBRACKET)
	p.ntExpr(id)
	p.term(RBRACKET)
}

func (p *Parser) ntFunctionCall(parent string) {
	id := p.addNode("\"FunctionCall ("+p.term(IDENT).literal+")\"", parent)
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
		p.editNodeName(id, "\"Factor ("+p.term(IDENT).literal+")\"")
	} else if p.token.type_ == INTEGER {
		p.editNodeName(id, "\"Factor ("+p.term(INTEGER).literal+")\"")
	} else if p.peekNextToken().type_ == LBRACKET {
		p.ntArrIdent(id)
	} else if p.token.type_ == STRING {
		literal := p.term(STRING).literal
		literal = literal[1 : len(literal)-1]
		p.editNodeName(id, "\"Factor String ("+literal+")\"")
	}
}

func (p *Parser) ntTerm(parent string) {
	id := p.addNode("Term", parent)
	p.ntFactor(id)
	if p.token.type_ == STAR {
		p.term(STAR)
		p.editNodeName(id, "\"Term (*)\"")
		p.ntTerm(id)
	} else if p.token.type_ == SLASH {
		p.term(SLASH)
		p.editNodeName(id, "\"Term (/)\"")
		p.ntTerm(id)
	}
}

func (p *Parser) ntCompExpr(parent string) {
	id := p.addNode("CompExpr", parent)
	p.ntExpr(id)
	if p.token.type_ == COMP_OP {
		symbol := p.token.literal
		p.term(COMP_OP)
		p.editNodeName(id, "\"CompExpr ("+fmt.Sprintf("%s", symbol)+")\"")
		p.ntExpr(id)
	}
}
func (p *Parser) ntExpr(parent string) {
	id := p.addNode("Expr", parent)
	p.ntTerm(id)
	if p.token.type_ == PLUS {
		p.term(PLUS)
		p.editNodeName(id, "\"Expr (+)\"")
		p.ntExpr(id)
	} else if p.token.type_ == MINUS {
		p.term(MINUS)
		p.editNodeName(id, "\"Expr (-)\"")
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

func (p *Parser) ntIfElse(parent string) {
	id := p.addNode("IfElse", parent)
	p.term("IF_KW")
	p.term(LPAREN)
	p.ntCompExpr(id)
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
	p.ntCompExpr(id)
	p.term(SEMICOLON)
	p.ntAssignStmt(id)
	p.term(RPAREN)
	p.term(LBRACE)
	p.ntBody(id)
	p.term(RBRACE)
}

func (p *Parser) ntAssignStmt(parent string) {
	id := p.addNode("AssignStmt", parent)
	if p.token.type_ == "TYPE" {
		p.ntType(id)
	}
	p.editNodeName(id, "\"AssignStmt ("+p.term(IDENT).literal+")\"")
	if p.token.type_ == LBRACKET {
		p.term(LBRACKET)
		p.ntExpr(id)
		p.term(RBRACKET)
	}
	p.term(EQUALS)
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
	p.ntType(id)
	p.editNodeName(id, "\"Function ("+p.term(IDENT).literal+")\"")
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
	p.tree.AddNode("G", "Program", map[string]string{"peripheries": "2"})
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
