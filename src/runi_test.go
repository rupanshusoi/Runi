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
		{EQUALS, "=", 2},
		{INTEGER, "0", 2},
		{SEMICOLON, ";", 2},
		{RBRACE, "}", 3},
	}

	checkOutput(t, "tests/simple.txt", output[:])
}

func TestBig(t *testing.T) {
	output := [...]Token{
		{keywords["int"], "int", 1},
		{IDENT, "main", 1},
		{LPAREN, "(", 1},
		{RPAREN, ")", 1},
		{LBRACE, "{", 1},
		{keywords["int"], "int", 2},
		{IDENT, "x", 2},
		{EQUALS, "=", 2},
		{INTEGER, "5", 2},
		{SEMICOLON, ";", 2},
		{keywords["int"], "int", 3},
		{IDENT, "j", 3},
		{EQUALS, "=", 3},
		{INTEGER, "0", 3},
		{SEMICOLON, ";", 3},
		{keywords["for"], "for", 4},
		{LPAREN, "(", 4},
		{keywords["int"], "int", 4},
		{IDENT, "i", 4},
		{EQUALS, "=", 4},
		{INTEGER, "0", 4},
		{SEMICOLON, ";", 4},
		{IDENT, "i", 4},
		{COMP_OP, "<", 4},
		{INTEGER, "10", 4},
		{SEMICOLON, ";", 4},
		{IDENT, "i", 4},
		{EQUALS, "=", 4},
		{IDENT, "i", 4},
		{PLUS, "+", 4},
		{INTEGER, "1", 4},
		{RPAREN, ")", 4},
		{LBRACE, "{", 4},
		{keywords["if"], "if", 5},
		{LPAREN, "(", 5},
		{IDENT, "x", 5},
		{COMP_OP, "<", 5},
		{IDENT, "i", 5},
		{RPAREN, ")", 5},
		{LBRACE, "{", 5},
		{IDENT, "j", 6},
		{EQUALS, "=", 6},
		{IDENT, "j", 6},
		{PLUS, "+", 6},
		{INTEGER, "1", 6},
		{SEMICOLON, ";", 6},
		{RBRACE, "}", 7},
		{RBRACE, "}", 8},
		{RBRACE, "}", 9},
	}

	checkOutput(t, "tests/big.txt", output[:])
}

func TestCheckLine(t *testing.T) {
	output := [...]Token{
		{keywords["int"], "int", 5},
		{IDENT, "main", 5},
		{LPAREN, "(", 5},
		{RPAREN, ")", 5},
		{LBRACE, "{", 5},
		{keywords["int"], "int", 8},
		{IDENT, "x", 8},
		{EQUALS, "=", 8},
		{INTEGER, "5", 9},
		{SEMICOLON, ";", 9},
		{keywords["int"], "int", 10},
		{IDENT, "j", 10},
		{EQUALS, "=", 10},
		{INTEGER, "0", 10},
		{SEMICOLON, ";", 10},
		{keywords["for"], "for", 11},
		{LPAREN, "(", 11},
		{keywords["int"], "int", 11},
		{IDENT, "i", 11},
		{EQUALS, "=", 11},
		{INTEGER, "0", 11},
		{SEMICOLON, ";", 11},
		{IDENT, "i", 12},
		{COMP_OP, "<", 12},
		{INTEGER, "10", 12},
		{SEMICOLON, ";", 12},
		{IDENT, "i", 12},
		{EQUALS, "=", 12},
		{IDENT, "i", 12},
		{PLUS, "+", 12},
		{INTEGER, "1", 12},
		{RPAREN, ")", 12},
		{LBRACE, "{", 12},
		{keywords["if"], "if", 13},
		{LPAREN, "(", 13},
		{IDENT, "x", 13},
		{COMP_OP, "<", 13},
		{IDENT, "i", 13},
		{RPAREN, ")", 13},
		{LBRACE, "{", 13},
		{IDENT, "j", 14},
		{EQUALS, "=", 14},
		{IDENT, "j", 14},
		{PLUS, "+", 14},
		{INTEGER, "1", 14},
		{SEMICOLON, ";", 14},
		{RBRACE, "}", 15},
		{RBRACE, "}", 16},
		{RBRACE, "}", 19},
	}

	checkOutput(t, "tests/checkLine.txt", output[:])
}

func TestRandom(t *testing.T) {
	output := [...]Token{
		{STRING, "\"hello i am a string\"", 1},
		{STRING, "\"another one\"", 2},
		{INTEGER, "1234", 3},
		{INTEGER, "99192", 4},
		{PLUS, "+", 5},
		{EQUALS, "=", 5},
		{MINUS, "-", 5},
		{STAR, "*", 5},
		{SLASH, "/", 5},
		{COMP_OP, "!=", 5},
	}

	checkOutput(t, "tests/random.txt", output[:])
}

func TestIllegal(t *testing.T) {
	output := [...]Token{
		{ILLEGAL, "&", 1},
		{IDENT, "aa", 1},
		{ILLEGAL, "!", 2},
		{INTEGER, "123", 3},
	}

	checkOutput(t, "tests/illegal.txt", output[:])
}

func TestComments(t *testing.T) {
	output := [...]Token{
		{COMMENT, COMMENT, 1},
		{keywords["int"], "int", 3},
		{IDENT, "x", 3},
		{EQUALS, "=", 3},
		{INTEGER, "1", 3},
		{SEMICOLON, ";", 3},
	}

	checkOutput(t, "tests/comments.txt", output[:])
}
