package main

import (
	"fmt"
	codegen "js-compiler/condegen"
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

	// code generationr
	ir, err := codegen.CompileToLLVM(program)
	if err != nil {
		fmt.Errorf("code generation error: %v", err)
		os.Exit(1)
	}

	irFile := "something.ll"
	err = os.WriteFile(irFile, []byte(ir), 0664)
	if err != nil {
		fmt.Printf("Error writing IR file: %v\n", err)
		os.Exit(1)
	}

}
