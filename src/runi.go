package main

import (
	"os"
)

func main() {
	var lexer *Lexer
	if len(os.Args) > 1 {
		lexer = Lex(os.Args[1])
	} else {
		lexer = Lex("test.txt")
	}
	var parser = Parse(lexer)
	parser.Parse()
}
