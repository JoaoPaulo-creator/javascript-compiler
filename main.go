package main

import (
	"fmt"
	"js-compiler/lexer"
	"js-compiler/parser"
	"os"
)

func main() {
	args := os.Args[1]
	file, err := os.ReadFile(args)
	if err != nil {
		panic(err.Error())
	}
	lxr := lexer.New(string(file))

	// tokens := lxr.Tokenize()
	p := parser.New(lxr)
	program := p.ParseProgram()

	fmt.Println(program)

	// for _, tok := range tokens {
	// 	fmt.Printf("Token %s\n", tok.Type)
	// }
}
