package parser

import (
	"context"
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

type Param struct {
	StartWord   string           `arg:"positional" default:"World" help:"Start Word, where program start parse links"`
	MatchWord   string           `arg:"positional" help:"Word, which you want to find in wiki"`
	NumberMap   int              `arg:"-"`
	BoolMatch   bool             `arg:"-"`
	CountTreads int              `arg:"--count" default:"4" help:"Count threads will create in worker pool"`
	Timeout     time.Duration    `arg:"--timeout" default:"5m" help:"Timeout duration (e.g., 5m, 10s)"`
	Storage     chan interface{} `arg:"-"`
}

type Parser struct {
	Param Param
}

func NewParser(param Param) *Parser {
	return &Parser{
		Param: param,
	}
}
func (v *Parser) Work(ctx context.Context, cancel context.CancelFunc, data map[int]string) {
	//make storage (buffer) for link's last word, that will send in func Work newly
	//make chan (bridge) where workers send parse link and accumulator get them and save in buffer
	//make stringParseChan where func push links in workerpool
	buffer := make(map[int]string)
	bridge := make(chan string, v.Param.CountTreads)
	stringParseChan := make(chan string, v.Param.CountTreads)
	wg := &sync.WaitGroup{}
	//Implemented workerpool for parsing new links and research match word in links
	wg.Add(v.Param.CountTreads)
	for i := 0; i < v.Param.CountTreads; i++ {
		go v.run(ctx, cancel, wg, stringParseChan, bridge)
	}
	//wait all workers and close chan (bridge)
	go func() {
		wg.Wait()
		close(bridge)
	}()
	//Implemented in other goroutine for collect new links in buffer
	go v.accumulator(buffer, bridge)
	// func push links in workerpool
	go func() {
		defer close(stringParseChan)
		for _, val := range data {
			stringParseChan <- val
		}
	}()
	//before implement new cycle for func Work, wait two signal
	//if program found match, get single and break cycle
	//if program didn't find match word, get buffer with new links and implement new cycle func Work
	select {
	case <-ctx.Done():
		return

	case r := <-v.Param.Storage:
		newData := v.assertion(r)
		v.Work(ctx, cancel, newData)
	}
}

func (v *Parser) assertion(input interface{}) map[int]string {
	//assert interface{} in map[int]string
	val := input.(map[int]string)
	return val
}

func (v *Parser) accumulator(m map[int]string, bridge chan string) {
	idx := 0
	//collect new links in buffer
	for word := range bridge {
		if !v.Param.BoolMatch && word != "" {
			m[idx] = word
			idx++
		}
	}
	v.Param.Storage <- m
	v.Param.NumberMap++
}

func (v *Parser) run(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, input chan string, bridge chan string) {
	defer wg.Done()
	for str := range input {
		select {
		case <-ctx.Done():
			return
		default:
			v.parserUrl(cancel, str, bridge)
		}
	}
}

func (v *Parser) parserUrl(cancel context.CancelFunc, s string, bridge chan string) {
	var client http.Client

	response, err := client.Get(wikiURL + s)
	if err != nil {
		log.Fatal("Request error", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Response processing error", err)
	}

	res := string(body)
	v.finder(cancel, res, bridge)
}

func (v *Parser) finder(cancel context.CancelFunc, s string, bridge chan string) {

	re := regexp.MustCompile(hrefRegexPattern)
	match := re.FindAllStringSubmatch(s, -1)
	//if find match word, cancel all workers and func
	//if didn't match, send word in accumulator
	for _, element := range match {
		if element[1] == v.Param.MatchWord && !v.Param.BoolMatch {
			v.Param.BoolMatch = true
			result := wikiURL + element[1]
			fmt.Printf("\nMatched on level %d:\n---> %s\n", v.Param.NumberMap+1, result)
			cancel()
			return
		} else {
			bridge <- element[1]
		}
	}
}

// determinate data for fist call func Work with StartWord
func (v *Parser) SetupInitialData() map[int]string {
	InitMap := make(map[int]string)
	InitMap[0] = v.Param.StartWord
	fmt.Printf("Init URL ---> %s\n\n", wikiURL+InitMap[0])
	return InitMap
}
