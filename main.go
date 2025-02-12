package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Param struct {
	InputURL    string
	NumberMap   int
	MatchURL    string
	MatchWord   string
	BoolMatch   bool
	CountTreads int
	Storage     chan interface{}
	output      chan string
}

type Parser struct {
	param Param
}

func NewParser(param Param) *Parser {
	return &Parser{
		param: param,
	}
}

func (v *Parser) work(data map[int]string) {

	fmt.Println("New Work")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	buffer := make(map[int]string)
	done := make(chan struct{})
	bridge := make(chan string)
	stringParseChan := make(chan string)
	wg := &sync.WaitGroup{}

	time.Sleep(2 * time.Second)
	for i := 0; i < v.param.CountTreads; i++ {
		fmt.Printf("Start worker ------------>  %d\n\n", i)
		wg.Add(1)
		go v.run(ctx, wg, stringParseChan, bridge)
	}
	go v.accumulator(buffer, wg, done, bridge)
	go func() {
		fmt.Println("Sending string in workerpool...")
		defer fmt.Println("Sended all string in workerpool ---> stop pusher")
		for _, val := range data {
			stringParseChan <- val
		}
		var stop struct{}
		done <- stop
		close(stringParseChan)
	}()

	wg.Wait()
	time.Sleep(1 * time.Second)
	fmt.Println("Accumulator stop work")
	close(done)
	close(bridge)

	select {
	case result := <-v.param.output:
		fmt.Printf("\n\n\nMatched ->>\n\n%s\n", result)
		return
	default:
	}
	newData := v.accretion(v.param.Storage)
	time.Sleep(5 * time.Second)
	fmt.Println("Finish work\n\n\n\n")

	v.work(newData)
}

func (v *Parser) accretion(input chan interface{}) map[int]string {
	fmt.Println("Start accretion")
	data := <-input
	val := data.(map[int]string)
	return val
}

func (v *Parser) accumulator(m map[int]string, wg *sync.WaitGroup, done chan struct{}, input chan string) {
	fmt.Println("Accumulator start work")
	idx := 0
	for {
		select {
		case <-done:
			v.param.Storage <- m
			fmt.Println("SENDED map in Storage")
			break
		case str := <-input:
			m[idx] = str
			idx++
		}
	}
}

func (v *Parser) run(ctx context.Context, wg *sync.WaitGroup, input chan string, bridge chan string) {
	defer fmt.Println("Stop worker")
	defer wg.Done()
	for str := range input {
		select {
		case <-ctx.Done():
			return
		default:
			v.parserUrl(ctx, str, bridge)
		}
	}
}

func (v *Parser) parserUrl(ctx context.Context, s string, bridge chan string) {
	// fmt.Println(s)
	v.param.NumberMap++

	var client http.Client

	response, err := client.Get(s)
	if err != nil {
		fmt.Printf("Request error -> %s\n", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Response processing error -> %s\n", err)
	}

	res := string(body)

	scanner := bufio.NewScanner(strings.NewReader(res))
	for scanner.Scan() {
		line := scanner.Text()
		v.finder(ctx, line, bridge)
	}
	// fmt.Println("End parse inputURL")
}

func (v *Parser) finder(ctx context.Context, s string, buffer chan string) {
	if strings.Index(s, "/wiki/") >= 6 {
		re := regexp.MustCompile(`href="/wiki/([^"]+)"`)
		match := re.FindAllStringSubmatch(s, -1)
		for _, element := range match {
			if element[1] == v.param.MatchWord {
				v.param.BoolMatch = true
				v.param.output <- "https://en.wikipedia.org/wiki/" + element[1]
				close(v.param.Storage)
				close(v.param.output)
			} else {
				// fmt.Println("https://en.wikipedia.org/wiki/" + element[1])
				buffer <- "https://en.wikipedia.org/wiki/" + element[1]
			}
		}
	}
}

func (v *Parser) setupInitialData() map[int]string {
	field := strings.Split(v.param.MatchURL, "/")
	v.param.MatchWord = field[len(field)-1]
	InitMap := make(map[int]string)
	InitMap[0] = v.param.InputURL
	fmt.Printf("Init URL ---> %s\n\n", v.param.InputURL)
	return InitMap
}

func main() {

	p := NewParser(Param{
		NumberMap:   0,
		InputURL:    "https://en.wikipedia.org/wiki/World",
		MatchURL:    "https://en.wikipedia.org/wiki/War",
		CountTreads: 4,
		Storage:     make(chan interface{}, 1),
		output:      make(chan string),
		BoolMatch:   false,
	})
	start := p.setupInitialData()
	p.work(start)
	p.param.Storage <- start

	select {
	case result := <-p.param.output:
		fmt.Printf("\n\n\nMatched ->>\n\n%s\n", result)
	case <-time.After(10 * time.Minute):
		fmt.Println("Timeout reached, stopping...")
	}
}
