package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/kartmos/wiki-link-explorer.git/internal/parser"
)

func main() {
	var params parser.Param
	arg.MustParse(&params)

	params.NumberMap = 1
	params.BoolMatch = false

	if params.CountTreads <= 0 {
		log.Fatalf("Number of threads must be greater than %d", params.CountTreads)
	}
	p := parser.NewParser(params)

	p.Param.Storage = make(chan interface{})

	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(p.Param.Timeout, func() {
		fmt.Printf("\nTimeout reached after %v, stopping...\n", p.Param.Timeout)
		cancel()
	})
	defer close(p.Param.Storage)
	start := p.SetupInitialData()
	p.Work(ctx, cancel, start)

	<-ctx.Done()
	timer.Stop()
}
