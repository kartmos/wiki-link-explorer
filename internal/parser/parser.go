package parser

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"
)

const (
	hrefRegexPattern = `href="/wiki/([^"]+)"`
	wikiURL          = "https://en.wikipedia.org/wiki/"
)

var re *regexp.Regexp = regexp.MustCompile(hrefRegexPattern)
var client http.Client = *http.DefaultClient

type Param struct {
	StartWord   string        `arg:"positional" default:"World" help:"Start Word, where program start parse links"`
	MatchWord   string        `arg:"positional" help:"Word, which you want to find in wiki"`
	CountTreads int           `arg:"--count" default:"4" help:"Count threads will create in worker pool"`
	Timeout     time.Duration `arg:"--timeout" default:"5m" help:"Timeout duration (e.g., 5m, 10s)"`
}

type Parser struct {
	Param        Param
	BoolMatch    bool
	BackTracking map[string]string
}

type Result struct {
	ParentName string
	PageName   string
}

func NewParser(param Param) *Parser {
	return &Parser{
		Param:        param,
		BoolMatch:    false,
		BackTracking: map[string]string{},
	}
}

func (v *Parser) Work(parent context.Context) error {
	//make chan (outputChan) where workers send parse link and accumulator get them and save in buffer
	//make inputChan where func push links in workerpool
	outputChan := make(chan []Result, v.Param.CountTreads)
	inputChan := make(chan string, v.Param.CountTreads)
	schedulerChan := make(chan []string)
	defer close(schedulerChan)

	// initialize
	inputChan <- v.Param.StartWord

	//Implemented workerpool for parsing new links and research match word in links
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	wg := &sync.WaitGroup{}
	log.Printf("Starting %d working goroutines", v.Param.CountTreads)
	wg.Add(v.Param.CountTreads)
	for i := 0; i < v.Param.CountTreads; i++ {
		go v.run(ctx, wg, inputChan, outputChan)
	}

	//wait all workers and close chan (outputChan)
	go func() {
		wg.Wait()
		close(outputChan)
	}()

	go v.scheduler(ctx, schedulerChan, inputChan)

	links := []string{}
	seen := map[string]bool{}
	//collect new links in buffer

	var result Result
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case schedulerChan <- links:
			links = []string{}
		case results := <-outputChan:
			for _, res := range results {
				word := res.PageName
				if word == "" {
					continue
				}
				if seen[word] {
					continue
				}
				seen[word] = true
				if _, ok := v.BackTracking[word]; !ok {
					v.BackTracking[word] = res.ParentName
				}
				links = append(links, word)

				if word == v.Param.MatchWord {
					result = res
					goto path_restore
				}
			}
		}
	}

path_restore:
	fmt.Printf("\nMatched on level ??:\n---> %s\n", wikiURL+result.PageName)

	cur := result.ParentName
	path := []string{result.PageName, result.ParentName}

	for {
		parent, ok := v.BackTracking[cur]
		if !ok || parent == "" || cur == parent {
			break
		} else {
			cur = parent
			path = append(path, parent)
		}
	}

	fmt.Printf("Path: %+v\n", path)
	return nil
}

func (v *Parser) scheduler(ctx context.Context, schedulerChan <-chan []string, inputChan chan<- string) {
	defer close(inputChan)
	for {
		//before implement new cycle for func Work, wait two signal
		//if program found match, get single and break cycle
		//if program didn't find match word, get buffer with new links and implement new cycle func Work
		var links []string
		select {
		case <-ctx.Done():
			return
		case links = <-schedulerChan:
		}

		for _, val := range links {
			select {
			case <-ctx.Done():
				return
			case inputChan <- val:
			}
		}
	}
}

func (v *Parser) run(ctx context.Context, wg *sync.WaitGroup, input <-chan string, output chan<- []Result) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case pageName := <-input:
			v.parserUrl(ctx, pageName, output)
		}
	}
}

func (v *Parser) parserUrl(ctx context.Context, pageName string, bridge chan<- []Result) {
	log.Printf("Parse URL --->%s", wikiURL+pageName)

	req, err := http.NewRequestWithContext(ctx, "GET", wikiURL+pageName, nil)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		log.Println("[WARN] NewRequestWithContext", err)
		return
	}

	response, err := client.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		log.Println("[WARN] Request error", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		log.Println("[WARN] Response processing error", err)
		return
	}

	res := string(body)
	//if find match word, cancel all workers and func
	//if didn't match, send word in accumulator
	result := []Result{}
	for _, element := range re.FindAllStringSubmatch(res, -1) {
		capture := element[1]
		if capture == "" {
			continue
		}
		result = append(result, Result{PageName: capture, ParentName: pageName})
	}

	select {
	case <-ctx.Done():
	case bridge <- result:
	}
}

func (v *Parser) Start() error {
	ctx, cancel := context.WithDeadlineCause(context.Background(), time.Now().Add(v.Param.Timeout),
		fmt.Errorf("Timeout reached after %v, stopping", v.Param.Timeout),
	)
	defer cancel()

	fmt.Printf("Init URL ---> %s\n\n", wikiURL+v.Param.StartWord)
	return v.Work(ctx)
}
