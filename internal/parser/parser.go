package parser

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type Param struct {
	InputURL    string
	NumberMap   int
	MatchURL    string
	MatchWord   string
	BoolMatch   bool
	CountTreads int
	Storage     chan interface{}
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

	fmt.Printf("\n\nLevel %d\n\n", v.Param.NumberMap)
	wg.Add(v.Param.CountTreads)
	for i := 1; i <= v.Param.CountTreads; i++ {
		fmt.Printf("Start worker ------------>  %d\n\n", i)
		go v.run(ctx, cancel, wg, stringParseChan, bridge)
	}
	go func() {
		wg.Wait()
		close(bridge)
	}()

	go v.accumulator(buffer, bridge)

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
	case r := <-v.Param.Storage:
		newData := v.accretion(r)

		if !v.Param.BoolMatch {
			fmt.Println("Not matched")
			select {
			case <-ctx.Done():
				return
			default:
				v.Work(ctx, cancel, newData)
			}
		}
	}
}

func (v *Parser) accretion(input interface{}) map[int]string {
	val := input.(map[int]string)
	return val
}

func (v *Parser) accumulator(m map[int]string, bridge chan string) {
	fmt.Println("Accumulator start working")
	idx := 0

	for url := range bridge {
		if !v.Param.BoolMatch && url != "" {
			m[idx] = url
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
			if element[1] == v.Param.MatchWord && !v.Param.BoolMatch {
				v.Param.BoolMatch = true
				result := "https://en.wikipedia.org/wiki/" + element[1]
				fmt.Printf("\n\n\nMatched on level %d:\n\n\n---> %s\n\n\n\n\n", v.Param.NumberMap+1, result)
				cancel()
				return
			} else {
				buffer <- "https://en.wikipedia.org/wiki/" + element[1]
			}
		}
	}
}

func (v *Parser) SetupInitialData() map[int]string {
	field := strings.Split(v.Param.MatchURL, "/")
	v.Param.MatchWord = field[len(field)-1]
	InitMap := make(map[int]string)
	InitMap[0] = v.Param.InputURL
	fmt.Printf("Init URL ---> %s\n\n", v.Param.InputURL)
	return InitMap
}
