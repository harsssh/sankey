package main

import (
	"regexp"
	"strings"
)

func trimQuotes(s string) string {
	return s[1 : len(s)-1]
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
	// クエリパラメータを削除
	pat := regexp.MustCompile(`\?.*`)
	uri = pat.ReplaceAllString(uri, "")

	// パスパラメータを短縮
	// スラッシュで分割
	pat = regexp.MustCompile(`/`)
	match := pat.Split(uri, -1)

	// 数字を含む部分を短縮
	var replaced []string
	for _, s := range match {
		pat = regexp.MustCompile(`\d+`)
		if pat.MatchString(s) {
			replaced = append(replaced, "*")
		} else {
			replaced = append(replaced, s)
		}
	}

	// replacedをスラッシュで結合
	return strings.Join(replaced, "/")
}

func getID(ua string) int {
	id, ok := UAToID[ua]
	if !ok {
		id = seq
		UAToID[ua] = id
		IDToUA[id] = ua
		seq++
	}
	return id
}
