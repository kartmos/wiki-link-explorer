<a id="readme-top"></a>

<h1 align="center">WikiLinkExplorer</h1>

<p align="center">
  A multi-threaded Wikipedia pathfinder in Go!
  <br />
  <br />
  <a href="https://github.com/kartmos/WikiLinkExplorer/issues/new?labels=bug&template=bug-report.md">Report a Bug</a>
  &middot;
  <a href="https://github.com/kartmos/WikiLinkExplorer/issues/new?labels=enhancement&template=feature-request.md">Request a Feature</a>
</p>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#about-the-project">About the Project</a></li>
    <li><a href="#features">Features</a></li>
    <li><a href="#getting-started">Getting Started</a></li>
    <li><a href="#usage">Usage</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About the Project

WikiLinkExplorer is a multithreaded web parser written in Go that searches for a path from one Wikipedia page to another by following hyperlinks. The project demonstrates efficient use of multithreading, channels, and contexts in Go for performing asynchronous tasks. Key features include:

* Parallel crawling with goroutines
* Timeout handling via contexts
* Communication through channels
* HTML parsing with regular expressions

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Features

- Multi-threaded search
- Customizable timeout (default: 5 minutes)
- Progress visualization
- Match tracking through article levels
- Error handling with retries

<!-- GETTING STARTED -->
## Getting Started

### Prerequisites

- Go 1.21+
- Internet connection

### Installation

1. Clone the repository
```sh
   git clone https://github.com/kartmos/WikiLinkExplorer.git
```
<p align="right">(<a href="#readme-top">back to top</a>)</p><!-- USAGE -->

### Usage

1. Navigate to the project directory
```sh
cd WikiLinkExplorer
```

2. Build the binary
```sh
go build -o wiki cmd/wikilinkexplorer/main.go
```
3. Enter the launch parameters
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
Example Usage:

```sh
$ ./wiki World War --count 4 --timeout 5m
```

Example output:

```sh
Matched on level 3:

---> https://en.wikipedia.org/wiki/War
```