package main

import (
	"io/ioutil"
	"strings"
)

type Lexer struct {
	program  string
	next_idx int  // index of next char
	char     byte // current char
	line_num int
}

type Token struct {
	type_    string
	literal  string
	line_num int // need it for parsing
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Peek next char without advancing next_idx
func (lexer *Lexer) peekChar() byte {
	if lexer.next_idx >= len(lexer.program) {
		return 0
	} else {
		return lexer.program[lexer.next_idx]
	}
}

func (lexer *Lexer) readChar() {
	if lexer.next_idx >= len(lexer.program) {
		lexer.char = 0
	} else {
		lexer.char = lexer.program[lexer.next_idx]
		lexer.next_idx += 1
	}
}

func (lexer *Lexer) readCharSkipWhitespace() {
	lexer.readChar()
	for lexer.char == ' ' || lexer.char == '\n' {
		if lexer.char == 10 {
			lexer.line_num += 1
		}
		lexer.readChar()
	}
}

// Read n chars including the current one, and return a string of length n
func (lexer *Lexer) readChars(n int) string {
	start := lexer.next_idx - 1
	for n != 1 {
		lexer.readChar()
		n -= 1
	}
	return lexer.program[start:lexer.next_idx]
}

func (lexer *Lexer) emitIllegalToken() *Token {
	return &Token{ILLEGAL, string(lexer.char), lexer.line_num}
}

func isAlpha(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z'
}

func (lexer *Lexer) emitIdentifierToken() *Token {
	start := lexer.next_idx - 1
	for isAlpha(lexer.peekChar()) {
		lexer.readChar()
	}
	return &Token{IDENT, lexer.program[start:lexer.next_idx], lexer.line_num}
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func (lexer *Lexer) emitIntegerToken() *Token {
	start := lexer.next_idx - 1
	for isDigit(lexer.peekChar()) {
		lexer.readChar()
	}
	return &Token{INTEGER, lexer.program[start:lexer.next_idx], lexer.line_num}
}

func (lexer *Lexer) emitStringToken() *Token {
	// Advance over the opening quote
	lexer.readChar()

	start := lexer.next_idx - 1
	for lexer.peekChar() != '"' {
		lexer.readChar()
		if lexer.char == 0 {
			panic("unexpected EOF while looking for closing quote")
		} else if lexer.char == '\n' {
			lexer.line_num += 1
		}
	}
	token := Token{STRING, lexer.program[start:lexer.next_idx], lexer.line_num}

	// Advance over the closing quote
	lexer.readChar()

	return &token
}

func (lexer *Lexer) NextToken() *Token {
	var token Token

	switch lexer.char {
	case '(':
		token = Token{LPAREN, string(lexer.char), lexer.line_num}
	case ')':
		token = Token{RPAREN, string(lexer.char), lexer.line_num}
	case '[':
		token = Token{LBRACKET, string(lexer.char), lexer.line_num}
	case ']':
		token = Token{RBRACKET, string(lexer.char), lexer.line_num}
	case '{':
		token = Token{LBRACE, string(lexer.char), lexer.line_num}
	case '}':
		token = Token{RBRACE, string(lexer.char), lexer.line_num}
	case ';':
		token = Token{SEMICOLON, string(lexer.char), lexer.line_num}
	case ',':
		token = Token{COMMA, string(lexer.char), lexer.line_num}
	case '+':
		token = Token{PLUS, string(lexer.char), lexer.line_num}
	case '-':
		token = Token{MINUS, string(lexer.char), lexer.line_num}
	case '*':
		token = Token{STAR, string(lexer.char), lexer.line_num}
	case '/':
		if lexer.peekChar() == '/' {
			for lexer.peekChar() != '\n' {
				lexer.readChar()
			}
			// This call will increment the line number without us
			// having to do it directly
			lexer.readCharSkipWhitespace()
			// The caller expects us to return a token every time,
			// so we call NextToken() on our own here
			return lexer.NextToken()
		} else {
			token = Token{SLASH, string(lexer.char), lexer.line_num}
		}
	case '<':
		token = Token{COMP_OP, string(lexer.char), lexer.line_num}
	case '>':
		token = Token{COMP_OP, string(lexer.char), lexer.line_num}
	case '!':
		if lexer.peekChar() == '=' {
			token = Token{COMP_OP, lexer.readChars(2), lexer.line_num}
		} else {
			token = *lexer.emitIllegalToken()
		}
	case '=':
		if lexer.peekChar() == '=' {
			token = Token{COMP_OP, lexer.readChars(2), lexer.line_num}
		} else {
			token = Token{ASSIGN, string(lexer.char), lexer.line_num}
		}
	case '"':
		token = *lexer.emitStringToken()
	default:
		if isAlpha(lexer.char) {
			token = *lexer.emitIdentifierToken()

			// Emit a different token if the identifier is a keyword
			type_, isKeyword := keywords[token.literal]
			if isKeyword {
				token = Token{type_, token.literal, lexer.line_num}
			}
		} else if isDigit(lexer.char) {
			token = *lexer.emitIntegerToken()
		} else {
			token = *lexer.emitIllegalToken()
		}
	}

	lexer.readCharSkipWhitespace()
	return &token
}

func Lex(file string) *Lexer {
	p, err := ioutil.ReadFile(file)
	check(err)

	program := string(p)

	var leadingNewlines int
	for _, value := range program {
		if value == '\n' {
			leadingNewlines += 1
		} else if value != ' ' {
			break
		}
	}

	program = strings.TrimSpace(program)
	return &Lexer{program, 1, program[0], leadingNewlines + 1}
}
