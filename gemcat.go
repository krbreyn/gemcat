package gemcat

import (
	"strings"
)

type BrowserData struct {
	Bookmarks []string `json:"bookmarks"`
	History   []string `json:"history"`
}

type BrowserState struct {
	CurrURL string
	Pos     int
	Stack   []Page
	Data    BrowserData
}
type Link struct {
	No      int
	URL     string
	Visited bool
}

type Page struct {
	URL     string
	Content string
	Links   []Link
}

func GetHostPath(url string) (host, path string) {
	split := strings.SplitN(url, "/", 2)
	if len(split) == 1 {
		host, path = split[0], ""
	} else {
		host, path = split[0], split[1]
	}
	return host, path
}

func IsGeminiLink(url string) bool {
	return strings.HasPrefix(url, "gemini://") || !strings.Contains(url, "://")
}
