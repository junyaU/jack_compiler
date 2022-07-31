package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/jack_compiler"
	"log"
	"os"
	"strings"
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

	for _, file := range analyzer.Files() {
		t := jack_compiler.NewTokenizer(file)

		e, err := jack_compiler.NewComplicationEngine(file.Name()[:strings.Index(file.Name(), ".jack")])
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		e.Compile(t)

		if err := e.Write(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	fmt.Println("compile success")
}
