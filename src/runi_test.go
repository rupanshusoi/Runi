package main

import (
	"testing"
)

func checkTokens(t *testing.T, t1, t2 *Token) {
	if t1.type_ != t2.type_ {
		t.Errorf("Incorrect type: expected %s, got %s\n", t1.type_, t2.type_)
	}
	if t1.literal != t2.literal {
		t.Errorf("Incorrect literal: expected %s, got %s\n", t1.literal, t2.literal)
	}
	if t1.line_num != t2.line_num {
		t.Errorf("Incorrect line_num: expected %d, got %d\n", t1.line_num, t2.line_num)
	}
}

func checkOutput(t *testing.T, file string, output []Token) {
	lexer := Lex(file)
	var token *Token
	for i := 0; lexer.char != 0; i++ {
		token = lexer.NextToken()
		checkTokens(t, &output[i], token)
	}
}

func TestSimple(t *testing.T) {
	output := [...]Token{
		{keywords["int"], "int", 1},
		{IDENT, "main", 1},
		{LPAREN, "(", 1},
		{RPAREN, ")", 1},
		{LBRACE, "{", 1},
		{keywords["int"], "int", 2},
		{IDENT, "x", 2},
		{ASSIGN, "=", 2},
		{INTEGER, "0", 2},
		{SEMICOLON, ";", 2},
		{RBRACE, "}", 3},
	}

	checkOutput(t, "tests/simple.txt", output[:])
}
