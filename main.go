package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Param struct {
	InputURL    string
	MatchURL    string
	MatchWord   string
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

func (v *Parser) ParserUrl(url string) error {
	var client http.Client

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
			break
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
		CountTreads: 4,
		MatchURL:    "https://en.wikipedia.org/wiki/Russia",
		Stage:       make(map[string]interface{}),
	})

	p.par.MatchWord = strings.Join(strings.SplitAfterN(p.par.MatchURL, "/", 4), "")
	p.ParserUrl(p.par.InputURL)

	for i := 0; i <= p.par.Idx; i++ {
		if val, found := p.par.Stage[strconv.Itoa(i)]; found {
			fmt.Printf("%d  %s\n", i, val)
		}
	}
}
