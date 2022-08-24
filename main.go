package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

var IDToUA map[int]string
var UAToID map[string]int
var count int

type Request struct {
	Method string
	URI    string
}

type Log struct {
	ID      int
	Request *Request
}

func main() {
	IDToUA = make(map[int]string)
	UAToID = make(map[string]int)

	logs := parse(os.Stdin)
	for _, log := range logs {
		fmt.Printf("%d %s %s\n", log.ID, log.Request.Method, log.Request.URI)
	}
}

func parse(r io.Reader) []*Log {
	scanner := bufio.NewScanner(r)
	var logs []*Log
	for scanner.Scan() {
		log := parseLine(scanner.Text())
		logs = append(logs, log)
	}

	if err := scanner.Err(); err != nil {
		println(err.Error())
		panic(err)
	}
	return logs
}

func parseLine(line string) *Log {
	// ダブルクォートで囲まれた文字列を抽出
	pat := regexp.MustCompile(`"(.*?)"`)
	match := pat.FindAllString(line, -1)

	// 最初がHTTPリクエスト、最後がUser-Agent
	r := trimQuotes(match[0])
	ua := trimQuotes(match[len(match)-1])

	id := getID(ua)
	req := &Request{
		Method: extractMethod(r),
		URI:    shortenURI(extractURI(r)),
	}

	// IDとしてUser-Agentを使用
	return &Log{id, req}
}
