package main

import (
	"log"

	"github.com/alexflint/go-arg"
	"github.com/kartmos/wiki-link-explorer.git/internal/parser"
)

func main() {
	var params parser.Param
	arg.MustParse(&params)
	if params.CountTreads <= 0 {
		log.Fatalf("Number of threads must be greater than %d", params.CountTreads)
	}

	p := parser.NewParser(params)

	if err := p.Start(); err != nil {
		log.Fatalln(err)
	}
}
