package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/will-wow/larkdown/parser"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		return fmt.Errorf("no input file")
	}

	filename := []byte(args[0])

	source, err := os.ReadFile(string(filename))
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	data, err := parser.MarkdownToTree(source, *spec.NewSpecDocument())
	if err != nil {
		return fmt.Errorf("failed to parse: %w", err)
	}

	fmt.Printf("%+v", data)

	return nil
}
