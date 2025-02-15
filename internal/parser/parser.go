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
	InputWord   string           `arg:"positional" default:"World" help:"Start Word, where program start parse links"`
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

	buffer := make(map[int]string)
	bridge := make(chan string, v.Param.CountTreads)
	stringParseChan := make(chan string, v.Param.CountTreads)
	wg := &sync.WaitGroup{}

	wg.Add(v.Param.CountTreads)
	for i := 0; i < v.Param.CountTreads; i++ {
		go v.run(ctx, cancel, wg, stringParseChan, bridge)
	}

	go func() {
		wg.Wait()
		close(bridge)
	}()

	go v.accumulator(buffer, bridge)

	go func() {
		defer close(stringParseChan)
		for _, val := range data {
			stringParseChan <- val
		}
	}()

	select {
	case <-ctx.Done():
		return

	case r := <-v.Param.Storage:
		newData := v.accretion(r)
		v.Work(ctx, cancel, newData)
	}
}

func (v *Parser) accretion(input interface{}) map[int]string {
	val := input.(map[int]string)
	return val
}

func (v *Parser) accumulator(m map[int]string, bridge chan string) {
	idx := 0

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

func (v *Parser) SetupInitialData() map[int]string {
	InitMap := make(map[int]string)
	InitMap[0] = v.Param.InputWord
	fmt.Printf("Init URL ---> %s\n\n", InitMap[0])
	return InitMap
}
