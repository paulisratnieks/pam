package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	log.SetFlags(0)
	input := flag.String("i", "", "input file")
	output := flag.String("o", "", "output file")
	debug := flag.Bool("d", false, "debug mode")
	flag.Parse()
	source := flag.Arg(0)

	if source == "" {
		log.Fatal("usage: pam [-i input file] [-o output file] [-d debug] file")
	}

	b, err := os.ReadFile(source)
	if err != nil {
		log.Fatalf("error opening source file: %v", err)
	}
	var inp []int
	if *input != "" {
		b, err := os.ReadFile(*input)
		if err == nil {
			for _, s := range strings.Split(string(b), ",") {
				s = strings.TrimSpace(s)
				i, err := strconv.Atoi(s)
				if err != nil {
					log.Fatalf("input contains unparsable int chars: %v", s)
				}
				inp = append(inp, i)
			}
		}
	}

	lexer := NewLexer(string(b))
	tokens, errs := lexer.Scan()
	if errs != nil {
		log.Fatalf("errors when tokenizing input: %v", errs)
	}

	parser := NewParser(tokens)
	ast := parser.Parse()
	if *debug {
		log.Println(ast.Print())
	}

	interpreter := NewInterpreter(inp)
	interpreter.Interpret(ast)
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			log.Fatalf("error creating output file: %v", err)
		}
		defer f.Close()

		for _, i := range interpreter.output {
			_, err = fmt.Fprintln(f, strconv.Itoa(i))
			if err != nil {
				log.Fatalf("error writing to output file: %v", err)
			}
		}
	} else {
		for _, i := range interpreter.output {
			log.Println(i)
		}
	}
}
