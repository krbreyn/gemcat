package main

import (
	"fmt"
	"slices"
)

type Page struct {
	URL     string
	Content string
	Links   []Link
}

func (p Page) Render() string {
	return ColorGemtext(p.Content, p.Links)
}

type Link struct {
	No      int
	URL     string
	Visited bool
}

type Browser struct {
	CurrURL string
	Pos     int
	Stack   []Page
	History []string
}

func (b *Browser) GotoURL(url string) error {
	host, path := getHostPath(url)
	_, body, err := Fetch(host, path)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	b.CurrURL = url

	content, links := DoLinks(body)
	if len(b.Stack) != 0 {
		b.Pos++
	}

	if b.Pos == len(b.Stack) {
		b.Stack = append(b.Stack, Page{
			URL:     url,
			Content: content,
			Links:   links,
		})
	} else {
		b.Stack = append(b.Stack[:b.Pos], Page{
			URL:     url,
			Content: content,
			Links:   links,
		})
	}

	if !slices.Contains(b.History, url) {
		b.History = append(b.History, url)
	}

	return nil
}

func (b *Browser) GetCurrPage() Page {
	return b.Stack[b.Pos]
}

func (b *Browser) RenderCurrPage() string {
	if len(b.Stack) == 0 {
		return "You have no current page!"
	}
	return b.GetCurrPage().Render()
}

func (b *Browser) GoForward() {
	if len(b.Stack) == 0 {
		return
	}

	if b.Pos < len(b.Stack)-1 {
		b.Pos++
	}
	b.CurrURL = b.GetCurrPage().URL
}

func (b *Browser) GoBack() {
	if len(b.Stack) == 0 {
		return
	}

	if b.Pos > 0 {
		b.Pos--
	}
	b.CurrURL = b.GetCurrPage().URL
}
