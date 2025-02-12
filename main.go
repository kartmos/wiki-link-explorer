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

func (v *Parser) work(ctx context.Context, cancel context.CancelFunc, data map[int]string) {

	buffer := make(map[int]string)
	bridge := make(chan string, v.param.CountTreads)
	stringParseChan := make(chan string, v.param.CountTreads)
	wg := &sync.WaitGroup{}

	fmt.Printf("\n\nLevel %d\n\n", v.param.NumberMap)
	wg.Add(v.param.CountTreads)
	for i := 1; i <= v.param.CountTreads; i++ {
		fmt.Printf("Start worker ------------>  %d\n\n", i)
		go v.run(ctx, cancel, wg, stringParseChan, bridge)
	}
	go func() {
		wg.Wait()
		close(bridge)
	}()

	go v.accumulator(ctx, buffer, bridge)

	go func() {
		defer close(stringParseChan)
		fmt.Println("Sending string in workerpool...")
		defer fmt.Println("Sent all string in workerpool ---> Stop pushing")
		for _, val := range data {
			stringParseChan <- val
		}
	}()

	select {
	case <-ctx.Done():
		return
	case r := <-v.param.Storage:
		newData := v.accretion(r)

		if !v.param.BoolMatch {
			fmt.Println("Not matched")
			select {
			case <-ctx.Done():
				return
			default:
				v.work(ctx, cancel, newData)
			}
		}
	}
}

func (v *Parser) accretion(input interface{}) map[int]string {
	val := input.(map[int]string)
	return val
}

func (v *Parser) accumulator(ctx context.Context, m map[int]string, bridge chan string) {
	fmt.Println("Accumulator start working")
	idx := 0

	for url := range bridge {
		if !v.param.BoolMatch && url != "" {
			m[idx] = url
			idx++
		}
	}
	v.param.Storage <- m
	v.param.NumberMap++
}

func (v *Parser) run(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, input chan string, bridge chan string) {
	defer wg.Done()
	for str := range input {
		select {
		case <-ctx.Done():
			fmt.Println("Stop worker")
			return
		default:
			v.parserUrl(cancel, str, bridge)
		}
	}
}

func (v *Parser) parserUrl(cancel context.CancelFunc, s string, bridge chan string) {
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
		v.finder(cancel, line, bridge)
	}
}

func (v *Parser) finder(cancel context.CancelFunc, s string, buffer chan string) {
	if strings.Index(s, "/wiki/") >= 6 {
		re := regexp.MustCompile(`href="/wiki/([^"]+)"`)
		match := re.FindAllStringSubmatch(s, -1)
		for _, element := range match {
			if element[1] == v.param.MatchWord && !v.param.BoolMatch {
				v.param.BoolMatch = true
				result := "https://en.wikipedia.org/wiki/" + element[1]
				fmt.Printf("\n\n\nMatched on level %d:\n\n\n---> %s\n\n\n\n\n", v.param.NumberMap+1, result)
				cancel()
				return
			} else {
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
	ctx, cancel := context.WithCancel(context.Background())
	timeout := 1 * time.Minute
	timer := time.AfterFunc(timeout, func() {
		fmt.Printf("\nTimeout reached after %v, stopping...\n", timeout)
		cancel()
	})
	p := NewParser(Param{
		NumberMap: 1,
		InputURL:  "https://en.wikipedia.org/wiki/World",
		MatchURL:  "https://en.wikipedia.org/wiki/War",
		// MatchURL:    "https://en.wikipedia.org/wiki/Civilian_casualty", // 4 уровень и долгий поиск
		CountTreads: 4,
		Storage:     make(chan interface{}),
		BoolMatch:   false,
	})
	defer close(p.param.Storage)
	start := p.setupInitialData()
	p.work(ctx, cancel, start)

	<-ctx.Done()
	timer.Stop()
}
