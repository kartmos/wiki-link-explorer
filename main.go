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
}

type Parser struct {
	param Param
}

func NewParser(param Param) *Parser {
	return &Parser{
		param: param,
	}
}

func pusher[T any, O any](input chan T, wg *sync.WaitGroup, execute func(msg T) O) {
	defer wg.Done()
	for msg := range input {
		execute(msg)
	}

}

func (v *Parser) mapRunner(m map[string]interface{}) {
	PusherIdx := 0
	if x, found := m[strconv.Itoa(PusherIdx)]; found {
		fmt.Printf("URL in pusher:\n%s\n\n", x.(string))
		v.param.output <- x.(string)
		PusherIdx++
	}
}

func (v *Parser) ParserUrl(input chan string) {
	fmt.Println("Create Goroutine")

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
	for scanner.Scan() {
		line := scanner.Text()
		v.finder(line)
	}
	fmt.Println("DONE")
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
				if element[1] == "Special:Random" {
					break
				}
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
		fmt.Printf("\n\n\n\n НАЙДЕН --------> %t\n\n\n\n", v.param.DoMatch)
	}
	return v.param.DoMatch
}

func (v *Parser) mapAccumulator(s string) *Param {
	v.param.Idx++
	id := strconv.Itoa(v.param.Idx)
	v.param.Memory[id] = "https://en.wikipedia.org/wiki/" + s
	// fmt.Printf("Collect --------> %s\n", v.param.Memory[id].(string))
	return nil
}

func main() {

	p := NewParser(Param{
		InputURL:    "https://en.wikipedia.org/wiki/World",
		MemoryIdx:   0,
		DoMatch:     true,
		CountTreads: 4,
		MatchURL:    "https://en.wikipedia.org/wiki/Académie_Française",
		Memory:      make(map[string]interface{}),
		input:       make(chan string),
		output:      make(chan string),
	})
	part := strings.Split(p.param.MatchURL, "/")
	p.param.MatchWord = part[len(part)-1]
	p.param.Memory[strconv.Itoa(p.param.MemoryIdx)] = p.param.InputURL
	fmt.Printf("Start URL:\n%s\n\n", p.param.InputURL)

	for i := 0; i < p.param.CountTreads; i++ {
		go p.ParserUrl(p.param.input)
	}

	p.param.input = p.pusher()
	fmt.Println("Start pusher")
	close(p.param.output)

	for result := range p.param.output {
		fmt.Printf("\n\n\nMatched ->>\n\n%s", result)
	}
	fmt.Scan()
}
