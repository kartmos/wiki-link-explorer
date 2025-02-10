package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type Param struct {
	NumberMap   int
	InputURL    string
	MatchURL    string
	MatchWord   string
	DoMatch     bool
	CountTreads int
	Storage     chan interface{}
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

func (v *Parser) pusher(input chan interface{}) {
	fmt.Println("Start pusher")
	for msg := range input {
		v.accretionMap(msg)
	}

}

func (v *Parser) accretionMap(data interface{}) {
	if m, ok := data.(map[int]string); ok {
		v.mapRunner(m)
	}
}

func (v *Parser) mapRunner(data map[int]string) {
	PusherIdx := 0
	if x, found := data[PusherIdx]; found {
		fmt.Printf("URL in pusher:\n%s\n\n", x)
		v.param.input <- x
		PusherIdx++
	}

}
func (v *Parser) run(input chan string) {
	fmt.Println("Start worker")
	for msg := range input {
		v.parserUrl(msg)
	}
}

func (v *Parser) parserUrl(s string) {
	fmt.Println("Create Goroutine")
	v.param.NumberMap++
	linksBuffer := make(map[int]string)
	mapIndex := 0

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
		v.finder(line, linksBuffer, &mapIndex)
	}
	v.param.Storage <- linksBuffer
	fmt.Printf("\n\n-> %d Map in Chan\n\n", v.param.NumberMap)
}

func (v *Parser) finder(s string, buffer map[int]string, idx *int) error {
	if !v.param.DoMatch {
		return nil
	}
	if strings.Index(s, "/wiki/") >= 6 {
		re := regexp.MustCompile(`href="/wiki/([^"]+)"`)
		match := re.FindAllStringSubmatch(s, -1)
		for _, element := range match {
			if !v.match(element[1]) {
				buffer[*idx] = "https://en.wikipedia.org/wiki/" + string(element[1])
				*idx++
			} else {
				close(v.param.Storage) // тут надо все остановить
			}
		}
	}
	return nil
}

func (v *Parser) match(s string) bool {
	if s == v.param.MatchWord {
		v.param.DoMatch = true
		fmt.Printf("\n\n\n\n FOUND --------> %t\n\n\n\n", v.param.DoMatch)
	}
	return v.param.DoMatch
}

func (v *Parser) setupInitialData() {
	field := strings.Split(v.param.MatchURL, "/")
	v.param.MatchWord = field[len(field)-1]
	InitMap := make(map[int]string)
	InitMap[0] = v.param.InputURL
	v.param.Storage <- InitMap
	fmt.Printf("Init URL ---> %s\n\n", v.param.InputURL)
}

func main() {

	p := NewParser(Param{
		NumberMap:   0,
		InputURL:    "https://en.wikipedia.org/wiki/World",
		MatchURL:    "https://en.wikipedia.org/wiki/Académie_Française",
		CountTreads: 4,
		Storage:     make(chan interface{}),
		input:       make(chan string),
		output:      make(chan string),
		DoMatch:     false,
	})

	p.setupInitialData()

	for i := 0; i < p.param.CountTreads; i++ {
		go p.run(p.param.input)
	}

	go p.pusher(p.param.Storage)

	close(p.param.output)

	for result := range p.param.output {
		fmt.Printf("\n\n\nMatched ->>\n\n%s", result)
	}
	fmt.Scan()
}
