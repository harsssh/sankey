package parser

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strings"

	"github.com/harsssh/sankey/model"
)

func Parse(r io.Reader) []model.Log {
	var logs []model.Log

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		log, err := parseLine(scanner.Text())
		if err != nil {
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

func parseLine(line string) (model.Log, error) {
	// Extract double-quoted string
	pat := regexp.MustCompile(`"(.*?)"`)
	match := pat.FindAllString(line, -1)

	// HTTP Request at the beginning,
	// User-Agent at the end
	r := trimQuotes(match[0])
	ua := trimQuotes(match[len(match)-1])

	if !validateRequest(r) {
		return model.Log{}, errors.New("parse error")
	}

	req := model.Request{
		Method: extractMethod(r),
		URI:    shortenURI(extractURI(r)),
	}

	return model.Log{ua, req}, nil
}

func trimQuotes(s string) string {
	return s[1 : len(s)-1]
}

func validateRequest(r string) bool {
	pat := regexp.MustCompile(`(GET|POST|PUT|DELETE|HEAD|OPTIONS|TRACE|CONNECT)\s+`)
	return pat.MatchString(r)
}

func extractMethod(r string) string {
	pat := regexp.MustCompile(`\s+`)
	match := pat.Split(r, -1)
	return match[0]
}

func extractURI(r string) string {
	pat := regexp.MustCompile(`\s+`)
	match := pat.Split(r, -1)
	return match[1]
}

func shortenURI(uri string) string {
	// remove params
	pat := regexp.MustCompile(`\?.*`)
	uri = pat.ReplaceAllString(uri, "")

	pat = regexp.MustCompile(`/`)
	match := pat.Split(uri, -1)

	// Shorten strings containing numbers or extensions
	var replaced []string
	for _, s := range match {
		pat = regexp.MustCompile(`(\d|\..+$)+`)
		if pat.MatchString(s) {
			replaced = append(replaced, "*")
		} else {
			replaced = append(replaced, s)
		}
	}

	return strings.Join(replaced, "/")
}
