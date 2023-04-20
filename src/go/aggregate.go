package main

import (
	"encoding/csv"
	"fmt"
	"github.com/harsssh/sankey/model"
	"github.com/harsssh/sankey/parser"
	"os"
)

type Transition struct {
	From model.Request
	To   model.Request
}

const (
	HEADER_SRC = "source"
	HEADER_DST = "target"
	HEADER_CNT = "value"
)

var (
	transitionCount map[Transition]int
)

func init() {
}

func main() {
	transitionCount = make(map[Transition]int)

	logs := parser.Parse(os.Stdin)
	countTransition(logs)

	saveToFile()
}

func saveToFile() {
	records := createRecords()

	w := csv.NewWriter(os.Stdout)
	err := w.WriteAll(records)
	if err != nil {
		panic(err)
	}
}

func createRecords() [][]string {
	records := [][]string{
		{HEADER_SRC, HEADER_DST, HEADER_CNT},
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

func countTransition(logs []model.Log) {
	prevTransition := make(map[string]Transition)
	for _, log := range logs {
		prev, ok := prevTransition[log.UA]
		if ok {
			from := prev.To
			to := log.Request

			t := Transition{
				From: from,
				To:   to,
			}
			transitionCount[t]++
			prevTransition[log.UA] = t
		} else {
			t := Transition{
				To: log.Request,
			}
			prevTransition[log.UA] = t
		}
	}
}
