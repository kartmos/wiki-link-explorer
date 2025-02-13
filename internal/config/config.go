package config

import (
	"flag"
	"time"
)

var (
	StartURL  string
	TargetURL string
	Threads   int
	Timeout   time.Duration
)

func init() {
	flag.StringVar(&StartURL, "start", "https://en.wikipedia.org/wiki/World", "Start Wikipedia URL")
	flag.StringVar(&TargetURL, "target", "https://en.wikipedia.org/wiki/War", "Target Wikipedia URL")
	flag.IntVar(&Threads, "threads", 4, "Number of worker threads")
	flag.DurationVar(&Timeout, "timeout", 5*time.Minute, "Maximum search duration")
}
