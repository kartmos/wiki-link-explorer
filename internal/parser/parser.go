package parser

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

const (
	hrefRegexPattern = `href="/wiki/([^"]+)"`
	wikiURL          = "https://en.wikipedia.org/wiki/"
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
	//запускается в горутине функция work
	//создается buffer для хранения ссылок полученных с воркеров
	// создается канал bridge по которому ассumulator собирает ссылки с воркеров и добавляет в buffer
	// создается канал stringParseChan через который отсылаются ссылки из мапы (buffer) в воркеров для парсинга новых
	buffer := make(map[int]string)
	bridge := make(chan string, v.Param.CountTreads)
	stringParseChan := make(chan string, v.Param.CountTreads)
	wg := &sync.WaitGroup{}

	//запускаем заданное количество воркеров
	wg.Add(v.Param.CountTreads)
	for i := 0; i < v.Param.CountTreads; i++ {
		go v.run(ctx, cancel, wg, stringParseChan, bridge)
	}
	//ждем пока они завершатся и закрываем канал bridge для остановки accumulator
	go func() {
		wg.Wait()
		close(bridge)
	}()
	//в отдельной горутине запускаем accumulator, чтобы сразу аккумулировал ссылки с воркеров
	go v.accumulator(buffer, bridge)
	// пушер ссылок в воркеров
	go func() {
		defer close(stringParseChan)
		for _, val := range data {
			stringParseChan <- val
		}
	}()
	//тут реализуем механизм закрытия горутин воркеров в случае нахождения совпадения в ссылках
	// в случае не нахождения совпадения ждем и достаем по мере появления новой мапы с свежими ссылками из accumulator в канале Storage
	select {
	case <-ctx.Done():
		return
		//заупскаем функцию accertion для реализации type accertion interface{} -> map[int]string
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
		//механизм отмены горутины в случая нахождения совпадения в любой запущенной горутине
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
	//получаем ссылку и делаем запрос
	response, err := client.Get(s)
	if err != nil {
		log.Fatal("Request error", err)
	}

	defer response.Body.Close()
	//читаем тело ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Response processing error", err)
	}
	//превращаем ответ в строку и отправляем в finder
	res := string(body)
	v.finder(cancel, res, bridge)
}

func (v *Parser) finder(cancel context.CancelFunc, s string, bridge chan string) {
	//проверяем есть ли в строке хоть 1 совпадение с интересующим нас хендлером
	if strings.Index(s, wikiURL) >= 6 {
		//выделяем группу захвата в ссылках последние слово в ссылке для его дальнейшего сравнения
		re := regexp.MustCompile(hrefRegexPattern)
		match := re.FindAllStringSubmatch(s, -1)
		//по массиву с отобранными словами из интересующих нас ссылок из тела ответа ищем совпадения
		//совпадений нет кидаем в канал, который читает accumulator и создает мапу для нового запуска функции work
		//совпадения есть выводим результат уровень ссылки и активируем context отмены работы горутин и выход функции work из рекурсии
		for _, element := range match {
			if element[1] == v.Param.MatchWord && !v.Param.BoolMatch {
				v.Param.BoolMatch = true
				result := wikiURL + element[1]
				fmt.Printf("\nMatched on level %d:\n---> %s\n", v.Param.NumberMap+1, result)
				cancel()
				return
			} else {
				bridge <- wikiURL + element[1]
			}
		}
	}
}

// вычленяем слово для поиска создаем мапу для запуска рекусрионной функции work
func (v *Parser) SetupInitialData() map[int]string {
	field := strings.Split(v.Param.MatchURL, "/")
	v.Param.MatchWord = field[len(field)-1]
	InitMap := make(map[int]string)
	InitMap[0] = v.Param.InputURL
	fmt.Printf("Init URL ---> %s\n\n", v.Param.InputURL)
	return InitMap
}
