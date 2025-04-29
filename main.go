package main

import (
	"log"
	"os"
)

func main() {
	b, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	lexer := NewLexer(string(b))
	tokens, _ := lexer.Scan()
	parser := NewParser(tokens)
	log.Println(parser.Parse().print())
}
