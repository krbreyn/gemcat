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

func IsGeminiLink(url string) bool {
	return strings.HasPrefix(url, "gemini://") || !strings.Contains(url, "://")
}
