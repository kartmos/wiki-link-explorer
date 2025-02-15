<a id="readme-top"></a>

  <h1 align="center">WikiLinkExplorer</h1>

  <p align="center">
    Многопоточный поисковик путей в Википедии на Go!
    <br />
    <br />
    <a href="https://github.com/kartmos/wiki-link-explorer/issues/new?labels=bug&template=bug-report.md">Сообщить об ошибке</a>
    &middot;
    <a href="https://github.com/kartmos/wiki-link-explorer/issues/new?labels=enhancement&template=feature-request.md">Предложить улучшение</a>
  </p>
</div>

<!-- ОГЛАВЛЕНИЕ -->
<details>
  <summary>Оглавление</summary>
  <ol>
    <li><a href="#о-проекте">О проекте</a></li>
    <li><a href="#возможности">Возможности</a></li>
    <li><a href="#начало-работы">Начало работы</a></li>
    <li><a href="#использование">Использование</a></li>
  </ol>
</details>

<!-- О ПРОЕКТЕ -->
## О проекте

WikiLinkExplorer — это многопоточный веб-парсер, написанный на Go, который ищет путь от одной страницы Wikipedia до другой, проходя по гиперссылкам. Проект демонстрирует эффективное использование многопоточности, каналов и контекстов в Go для выполнения асинхронных задач. Основные особенности:

* Параллельный краулинг с горутинами
* Обработка таймаутов через контексты
* Коммуникация через каналы
* Парсинг HTML с регулярными выражениями

<p align="right">(<a href="#readme-top">наверх</a>)</p>

### Возможности

- Многопоточный поиск
- Настраиваемый таймаут (по умолчанию: 5 минут)
- Визуализация прогресса
- Отслеживание совпадений через уровни статей
- Обработка ошибок с повторами

<!-- НАЧАЛО РАБОТЫ -->
## Начало работы

### Требования

- Go 1.21+
- Интернет-соединение

### Установка

1. Клонировать репозиторий
```sh
   git clone https://github.com/kartmos/wiki-link-explorer.git
```
<p align="right">(<a href="#readme-top">наверх</a>)</p><!-- ИСПОЛЬЗОВАНИЕ -->

### Использование

1. Переходим в директорию проекта
```sh
cd WikiLinkExplorer
```

2. Собираем бинарный файл
```sh
go build -o wiki-explorer cmd/wikilinkexplorer/main.go
```
3. Вводим параметры запуска
```sh
Usage: wiki [--count COUNT] [--timeout TIMEOUT] [INPUTWORD [MATCHWORD]]

Positional arguments:
  INPUTWORD              Start Word, where program start parse links
  MATCHWORD              Word, which you want to find in wiki

Options:
  --count COUNT          Count threads will create in worker pool [default: 4]
  --timeout TIMEOUT      Timeout duration (e.g., 5m, 10s) [default: 5m]
  --help, -h             display this help and exit
```
Пример ввода:

```sh
$ ./wiki World War --count 4 --timeout 5m
```

Пример вывода:

```sh
Matched on level 3:


---> https://en.wikipedia.org/wiki/War
```


