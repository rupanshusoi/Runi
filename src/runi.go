package main

import (
	"fmt"
	"os"
)

func main() {
	var lexer *Lexer
	if len(os.Args) > 1 {
		lexer = Lex(os.Args[1])
	} else {
		lexer = Lex("test.txt")
	}
	var token *Token
	for lexer.char != 0 {
		token = lexer.NextToken()
		fmt.Printf("%s, %s, %d\n", token.type_, token.literal, token.line_num)
	}
}
