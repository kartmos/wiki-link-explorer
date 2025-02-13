package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/kartmos/wiki-link-explorer/internal/config"
	"github.com/kartmos/wiki-link-explorer/internal/parser"
)

func main() {
	flag.Parse()

	if config.Threads <= 0 {
		log.Println("Number of threads must be greater than 0")
	}

	ctx, cancel := context.WithCancel(context.Background())
	timeout := config.Timeout
	timer := time.AfterFunc(timeout, func() {
		fmt.Printf("\nTimeout reached after %v, stopping...\n", timeout)
		cancel()
	})
	p := parser.NewParser(parser.Param{
		NumberMap:   1,
		InputURL:    config.StartURL,
		MatchURL:    config.TargetURL,
		CountTreads: config.Threads,
		Storage:     make(chan interface{}),
		BoolMatch:   false,
	})
	defer close(p.Param.Storage)
	start := p.SetupInitialData()
	p.Work(ctx, cancel, start)

	<-ctx.Done()
	timer.Stop()
}
