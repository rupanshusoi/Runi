package main

import (
	//"fmt"
	"os"
)

/*
func main() {
	var lexer *Lexer
	if len(os.Args) > 1 {
		lexer = Lex(os.Args[1])
	} else {
		lexer = Lex("test.txt")
	}
	var token *Token
	for true { // lexer.char != 0 {
		token = lexer.NextToken()
        if token.type_ == ILLEGAL && token.literal = "" {
            print("OK")
          } else {
            print("F")
          }
		fmt.Printf("%s, %s, %d\n", token.type_, token.literal, token.line_num)
	}
}
*/

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
