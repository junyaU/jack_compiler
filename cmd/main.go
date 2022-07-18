package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/jack_compiler"
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

	jack_compiler.NewTokenizer(analyzer.Files()[0])

	fmt.Println("compile success")
}
