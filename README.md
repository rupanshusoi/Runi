# Runi
Hand-written lexer and parser for C in Go.

## Motivation
This implementation is a part of our laboratory project to implement lexical and syntax analyzers at BITS Pilani. We also wanted to go out of our comfort zone and choose a modern language that we had not tried before.

## Lexer

### Usage
For now, the lexer is hard-coded to receive a file `test.txt` as input. Please do `cd src` and run `go run .` to tokenize the file, which is assumed to contain C source code.

## Authors
Rupanshu Soi & Nipun Wahi, Department of Computer Science, BITS Pilani at Hyderabad, India.

## Why "Runi" ?
`"Rupanshu Soi"[:2] + "Nipun Wahi"[:2]`

## References
- [Writing An Interpreter in Go by Thorsten Ball](https://interpreterbook.com/)
