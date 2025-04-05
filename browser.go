package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/muesli/reflow/wordwrap"
	"golang.org/x/term"
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
	Data    Data
	History []string
}

func NewBrowser() *Browser {
	return &Browser{}
}

func (b *Browser) GotoURL(url string) error {
	if !slices.Contains(b.History, url) {
		b.History = append(b.History, url)
	}

	_, body, err := Fetch(url)
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

	return nil
}

func (b *Browser) GetCurrPage() Page {
	return b.Stack[b.Pos]
}

// TODO - don't wrap preformatted blocks?
func (b *Browser) RenderOutput() string {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80
	}
	return wordwrap.String(b.RenderCurrPage(), width)
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
