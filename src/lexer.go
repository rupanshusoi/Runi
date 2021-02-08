package main

import (
	"io/ioutil"
	"strings"
)

type Lexer struct {
	program  string
	position int  // index of next char
	char     byte // current char
	line_num int
}

type Token struct {
	type_   string
	literal string
}

var keywords = map[string]string{
	"int":    "TYPE",
	"char":   "TYPE",
	"for":    "KEYWORD",
	"if":     "KEYWORD",
	"else":   "KEYWORD",
	"return": "KEYWORD",
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Peek next char without advancing position
func (lexer *Lexer) peekChar() byte {
	return lexer.program[lexer.position]
}

func (lexer *Lexer) readChar() {
	// The last char should be EOF anyway
	if lexer.position >= len(lexer.program)-1 {
		lexer.char = 0
	} else {
		lexer.char = lexer.program[lexer.position]
		lexer.position += 1
	}
}

func (lexer *Lexer) readCharSkipWhitespace() {
	lexer.readChar()
	// ASCII for \n is 10
	for lexer.char == ' ' || lexer.char == 10 {
		if lexer.char == 10 {
			lexer.line_num += 1
		}
		lexer.readChar()
	}
}

// Read n chars including the current one, and return a string of length n
func (lexer *Lexer) readChars(n int) string {
	start := lexer.position - 1
	for n != 1 {
		lexer.readChar()
		n -= 1
	}
	return lexer.program[start:lexer.position]
}

func (lexer *Lexer) emitIllegalToken() Token {
	return Token{"ILLEGAL", string(lexer.char)}
}

func isAlpha(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z'
}

func (lexer *Lexer) emitIdentifierToken() Token {
	start := lexer.position - 1
	for isAlpha(lexer.peekChar()) {
		lexer.readChar()
	}
	return Token{"IDENT", lexer.program[start:lexer.position]}
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func (lexer *Lexer) emitIntegerToken() Token {
	start := lexer.position - 1
	for isDigit(lexer.peekChar()) {
		lexer.readChar()
	}
	return Token{"INTEGER", lexer.program[start:lexer.position]}
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
	case '+':
		token = Token{"PLUS", string(lexer.char)}
	case '-':
		token = Token{"MINUS", string(lexer.char)}
	case '*':
		token = Token{"STAR", string(lexer.char)}
	case '/':
		token = Token{"SLASH", string(lexer.char)}
	case '<':
		token = Token{"COMP_OP", string(lexer.char)}
	case '>':
		token = Token{"COMP_OP", string(lexer.char)}
	case '!':
		if lexer.peekChar() == '=' {
			token = Token{"COMP_OP", lexer.readChars(2)}
		} else {
			token = lexer.emitIllegalToken()
		}
	case '=':
		if lexer.peekChar() == '=' {
			token = Token{"COMP_OP", lexer.readChars(2)}
		} else {
			token = Token{"ASSIGN", string(lexer.char)}
		}
	default:
		if isAlpha(lexer.char) {
			token = lexer.emitIdentifierToken()

			// Emit a different token if the identifier is a keyword
			type_, isKeyword := keywords[token.literal]
			if isKeyword {
				token = Token{type_, token.literal}
			}
		} else if isDigit(lexer.char) {
			token = lexer.emitIntegerToken()
		} else {
			token = lexer.emitIllegalToken()
		}
	}

	lexer.readCharSkipWhitespace()
	return &token
}

func Lex(file string) *Lexer {
	program, err := ioutil.ReadFile(file)
	check(err)

	program_ := strings.TrimSpace(string(program))
	return &Lexer{program_, 1, program_[0], 1}
}
