package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Param struct {
	StageIdx    int
	InputURL    string
	MatchURL    string
	MatchWord   string
	DoMatch     bool
	CountTreads int
	Steps       int
	Idx         int
	Stage       map[string]interface{}
}

type Parser struct {
	par Param
}

func NewParser(par Param) *Parser {
	return &Parser{
		par: par,
	}
}

func (v *Parser) ChangeUrl(url string) error {
	v.par.StageIdx++
	if val, found := v.par.Stage[strconv.Itoa(v.par.StageIdx)]; found {
		part := strings.SplitAfter(v.par.InputURL, "/")
		v.par.InputURL = strings.Join(part[:len(part)-1], "") + val.(string)
		fmt.Println(v.par.InputURL)
	}
	return nil
}

func (v *Parser) ParserUrl(url string) error {
	var client http.Client
	client.Timeout = 10 * time.Second

	response, err := client.Get(url)
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
		if !v.par.DoMatch {
			fmt.Println("DoMatch is false. Stopping execution.")
			return nil
		}
		if strings.Index(line, "/wiki/") >= 6 {
			v.collectorMap(line)
		}
	}

	return nil
}

func (v *Parser) collectorMap(s string) *Param {
	re := regexp.MustCompile(`href="/wiki/([^"]+)"`)
	match := re.FindAllStringSubmatch(s, -1)
	for _, element := range match {
		word := element[1]
		if v.par.MatchWord == word {
			v.par.DoMatch = false
		} else {
			v.par.Idx++
			id := strconv.Itoa(v.par.Idx)
			v.par.Stage[id] = word
		}
	}
	return nil
}

func main() {

	p := NewParser(Param{
		InputURL:    "https://en.wikipedia.org/wiki/World",
		StageIdx:    0,
		DoMatch:     true,
		CountTreads: 4,
		MatchURL:    "https://en.wikipedia.org/wiki/Nuclear_reactor",
		Stage:       make(map[string]interface{}),
	})
	part := strings.Split(p.par.MatchURL, "/")
	p.par.MatchWord = part[len(part)-1]
	for p.par.DoMatch {
		p.ParserUrl(p.par.InputURL)
		// for i := 0; i <= p.par.Idx; i++ {
		// 	if val, found := p.par.Stage[strconv.Itoa(i)]; found {
		// 		fmt.Printf("%d  %s\n", i, val)
		// 	}
		// }
		fmt.Printf("%t\n", p.par.DoMatch)
		p.ChangeUrl(p.par.InputURL)
	}
}
