package main

import (
	"io/ioutil"
)

type Lexer struct {
	program  string
	position int
	char     byte
}

type Token struct {
	type_   string
	literal string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (lexer *Lexer) peekChar() byte {
  return lexer.program[position + 1]
}

func (lexer *Lexer) NextToken() *Token {
	var token Token

	switch lexer.char {
	case '(':
		token = Token{"LPAREN", string(lexer.char)}
	case ')':
		token = Token{"RPAREN", string(lexer.char)}
	case '[':
		token = Token{"LBRACKET", string(lexer.char)}
	case ']':
		token = Token{"RBRACKET", string(lexer.char)}
	case '{':
		token = Token{"LBRACE", string(lexer.char)}
	case '}':
		token = Token{"RBRACE", string(lexer.char)}
	case ';':
		token = Token{"SEMICOLON", string(lexer.char)}
	case ',':
		token = Token{"COMMA", string(lexer.char)}
	case '=':
		token = Token{"ASSIGN", string(lexer.char)}
	case '<':
		token = Token{"COMP_OP", string(lexer.char)}
	case '>':
		token = Token{"COMP_OP", string(lexer.char)}
	case '!':
		token = Token{"COMP_OP", string(lexer.char)}
	case '=':
        if lexer.peekChar() == '=' {
            token = Token{"COMP_OP", "=="}
          }
	}
	return &token
}

func Lex(file string) *Lexer {
	program, err := ioutil.ReadFile(file)
	check(err)

	return &Lexer{string(program), 0, program[0]}
}
