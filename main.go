package main

import (
	"bufio"
	"encoding/csv"
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
	//logs = pick(logs, 32)
	countTransition(logs)

	fmt.Printf("ユーザー数: %d\n", seq)
	save()
}

func pick(logs []*Log, id int) []*Log {
	var picked []*Log
	for _, log := range logs {
		if log.ID == id {
			picked = append(picked, log)
		}
	}
	return picked
}

func save() {
	records := createRecords()
	f, err := os.Create("data.csv")
	if err != nil {
		panic(err)
	}

	w := csv.NewWriter(f)
	err = w.WriteAll(records)
	if err != nil {
		panic(err)
	}
}

func createRecords() [][]string {
	records := [][]string{
		{"source", "target", "value"},
	}
	for t, count := range transitionCount {
		records = append(records, []string{
			t.From.String(),
			t.To.String(),
			fmt.Sprintf("%d", count),
		})
	}
	return records
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
		if log == nil {
			continue
		}
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

	if !validateRequest(r) {
		return nil
	}

	id := getID(ua)
	req := &Request{
		Method: extractMethod(r),
		URI:    shortenURI(extractURI(r)),
	}

	// IDとしてUser-Agentを使用
	return &Log{id, req}
}
