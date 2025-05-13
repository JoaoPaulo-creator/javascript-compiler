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

	p := parser.New(lxr)
	program := p.ParseProgram()

	// Check for parsing errors
	if len(p.Errors()) > 0 {
		fmt.Printf("Parser has %d errors:\n", len(p.Errors()))
		for _, msg := range p.Errors() {
			fmt.Printf("\t%s\n", msg)
		}
		return
	}

	fmt.Println(program)
}
