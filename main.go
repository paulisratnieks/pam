package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	b, err := os.ReadFile("pam.txt")
	if err != nil {
		log.Fatalf("error opening source file: %v", err)
	}

	lexer := NewLexer(string(b))
	tokens, _ := lexer.Scan()
	parser := NewParser(tokens)
	ast := parser.Parse()
	ast.Print()

	input := make([]int, 0)
	b, err = os.ReadFile("input.txt")
	if err == nil {
		for _, s := range strings.Split(string(b), ",") {
			s = strings.TrimSpace(s)
			i, err := strconv.Atoi(s)
			if err != nil {
				log.Fatalf("input contains unparsable int chars: %v", s)
			}
			input = append(input, i)
		}
	}

	interpreter := NewInterpreter(input)
	interpreter.Interpret(ast)

	for _, i := range interpreter.output {
		fmt.Println(i)
	}
}
