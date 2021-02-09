package main

var keywords = map[string]string{
	"int":    "TYPE",
	"char":   "TYPE",
	"for":    "KEYWORD",
	"if":     "KEYWORD",
	"else":   "KEYWORD",
	"return": "KEYWORD",
}

const (
	LPAREN   = "LPAREN"
	RPAREN   = "RPAREN"
	LBRACKET = "LBRACKET"
	RBRACKET = "RBRACKET"
	LBRACE   = "LBRACE"
	RBRACE   = "RBRACE"

	SEMICOLON = "SEMICOLON"
	COMMA     = "COMMA"

	PLUS    = "PLUS"
	MINUS   = "MINUS"
	STAR    = "STAR"
	SLASH   = "SLASH"
	COMP_OP = "COMP_OP"

	ASSIGN = "ASSIGN"
	IDENT  = "IDENT"

	INTEGER = "INTEGER"

	ILLEGAL = "ILLEGAL"
)
