package main

import (
	"fmt"
)

func main() {
	lexer := Lex("test.txt")
	var token *Token
	for lexer.char != 0 {
		token = lexer.NextToken()
		fmt.Printf("%s, %s, %d\n", token.type_, token.literal, token.line_num)
	}
}
