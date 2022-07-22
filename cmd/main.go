package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/jack_compiler"
	"log"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", errors.New("no argument specified"))
		os.Exit(1)
	}

	analyzer, err := jack_compiler.NewAnalyzer(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
		os.Exit(1)
	}

	defer analyzer.Close()

	e, err := jack_compiler.NewComplicationEngine(args[0])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err := e.CompileClass(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	t := jack_compiler.NewTokenizer(analyzer.Files()[0])

	for t.HasMoreTokens() {

		tokenType, err := t.TokenType()
		if err != nil {
			log.Println(err)
		}

		switch tokenType {
		case jack_compiler.KEYWORD:
			keyword, err := t.Keyword()
			if err != nil {
				log.Println(err)
			}
			log.Println(keyword)

		case jack_compiler.SYMBOL:
		case jack_compiler.STRING_CONST:
		case jack_compiler.INT_CONST:
		case jack_compiler.IDENTIFIER:
		default:

		}

		t.Advance()
	}

	fmt.Println("compile success")
}
