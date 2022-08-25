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
var transitionCount map[Transition]int
var seq int

type Request struct {
	Method string
	URI    string
}

type Log struct {
	ID      int
	Request *Request
}

type Transition struct {
	From Request
	To   Request
}

func (r *Request) String() string {
	return r.Method + " " + r.URI
}

func init() {
	IDToUA = make(map[int]string)
	UAToID = make(map[string]int)
	transitionCount = make(map[Transition]int)
}

func main() {
	logs := parse(os.Stdin)
	countTransition(logs)
	for t, count := range transitionCount {
		fmt.Println(t, count)
	}
}

func countTransition(logs []*Log) {
	memo := make(map[int]*Transition)
	for _, log := range logs {
		prev, ok := memo[log.ID]
		if ok {
			from := prev.To
			to := *log.Request

			t := Transition{
				From: from,
				To:   to,
			}
			transitionCount[t]++
			memo[log.ID] = &t
		} else {
			t := &Transition{
				To: *log.Request,
			}
			memo[log.ID] = t
		}
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
