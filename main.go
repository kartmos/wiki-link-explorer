package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Param struct {
	PusherIdx   int
	MemoryIdx   int
	InputURL    string
	MatchURL    string
	MatchWord   string
	DoMatch     bool
	CountTreads int
	Idx         int
	Memory      map[string]interface{}
	input       chan string
	output      chan string
	sleep       time.Duration
	signal      *sync.WaitGroup
}

type Parser struct {
	param Param
}

func NewParser(param Param) *Parser {
	return &Parser{
		param: param,
	}
}

func (v *Parser) pusher() chan string {
	defer close(v.param.input)
	defer close(v.param.output)
	v.param.PusherIdx = 0
	for v.param.DoMatch {
		if x, found := v.param.Memory[strconv.Itoa(v.param.PusherIdx)]; found {
			fmt.Printf("URL in pusher:\n%s\n\n", x.(string))
			v.param.input <- x.(string)
			v.param.PusherIdx++
		}
	}
	return nil
}

func (v *Parser) ChangerUrlGet(url string) error {
	v.param.MemoryIdx++
	if val, found := v.param.Memory[strconv.Itoa(v.param.MemoryIdx)]; found {
		part := strings.SplitAfter(v.param.InputURL, "/")
		v.param.InputURL = strings.Join(part[:len(part)-1], "") + val.(string)
		fmt.Println(v.param.InputURL)
	}
	return nil
}

func (v *Parser) ParserUrl(input chan string) chan string {
	v.param.output = make(chan string)
	go func() {
		defer close(v.param.output)

		var client http.Client
		response, err := client.Get(<-input)
		fmt.Println("Get response")
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		res := string(body)

		scanner := bufio.NewScanner(strings.NewReader(res))
		fmt.Println("Create scanner")
		for scanner.Scan() {
			line := scanner.Text()
			v.finder(line)
		}
	}()

	return nil
}
func (v *Parser) finder(s string) *Param {
	if !v.param.DoMatch {
		fmt.Println("DoMatch is false. Stopping execution.")
		return nil
	}
	if strings.Index(s, "/wiki/") >= 6 {
		re := regexp.MustCompile(`href="/wiki/([^"]+)"`)
		match := re.FindAllStringSubmatch(s, -1)
		for _, element := range match {
			if v.Matcher(element[1]) {
				v.mapAccumulator(element[1])
			} else {
				v.param.output <- "https://en.wikipedia.org/wiki/" + element[1]
				break
			}
		}
	}
	return nil
}

func (v *Parser) Matcher(s string) bool {
	if s == v.param.MatchWord {
		v.param.DoMatch = false
	}
	return v.param.DoMatch
}

func (v *Parser) mapAccumulator(s string) *Param {
	v.param.Idx++
	id := strconv.Itoa(v.param.Idx)
	v.param.Memory[id] = "https://en.wikipedia.org/wiki/" + s
	fmt.Printf("Matched\n\n %s", v.param.Memory[id].(string))
	return nil
}

func main() {

	p := NewParser(Param{
		InputURL:    "https://en.wikipedia.org/wiki/World",
		MemoryIdx:   0,
		DoMatch:     true,
		CountTreads: 4,
		MatchURL:    "https://en.wikipedia.org/wiki/Nuclear_reactor",
		Memory:      make(map[string]interface{}),
		input:       make(chan string),
		sleep:       1 * time.Second,
	})
	part := strings.Split(p.param.MatchURL, "/")
	p.param.MatchWord = part[len(part)-1]
	p.param.Memory[strconv.Itoa(p.param.MemoryIdx)] = p.param.InputURL
	fmt.Printf("Start URL:\n%s\n\n", p.param.InputURL)
	p.param.input = p.pusher()
	fmt.Println("Start pusher")

	for result := range p.ParserUrl(p.param.input) {
		fmt.Printf("Matched\n\n%s", result)
	}
}
