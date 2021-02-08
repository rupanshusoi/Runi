package main

import ()

func main() {
  lexer := Lex("test.txt")
  print(lexer.program)
}
