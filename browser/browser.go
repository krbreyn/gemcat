package browser

import (
	"fmt"
	"os"
	"slices"

	"github.com/krbreyn/gemcat"
	"github.com/krbreyn/gemcat/gemtext"
	"github.com/muesli/reflow/wordwrap"
	"golang.org/x/term"
)

func RenderPage(p gemcat.Page) string {
	return gemtext.ColorGemtext(p.Content, p.Links)
}

type Browser struct {
	State gemcat.BrowserState
	IH    InputHandler
}

func NewBrowser() *Browser {
	return &Browser{IH: NewInputHandler()}
}

func (b *Browser) GotoURL(url string) error {
	if !slices.Contains(b.State.Data.History, url) {
		b.State.Data.History = append(b.State.Data.History, url)
	}

	_, body, err := Fetch(url)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	b.State.CurrURL = url

	content, links := gemtext.DoLinks(body)

	if len(b.State.Stack) != 0 {
		b.State.Pos++
	}
	if b.State.Pos == len(b.State.Stack) {
		b.State.Stack = append(b.State.Stack, gemcat.Page{
			URL:     url,
			Content: content,
			Links:   links,
		})
	} else {
		b.State.Stack = append(b.State.Stack[:b.State.Pos], gemcat.Page{
			URL:     url,
			Content: content,
			Links:   links,
		})
	}

	return nil
}

func (b *Browser) GetCurrPage() gemcat.Page {
	return b.State.Stack[b.State.Pos]
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
	if len(b.State.Stack) == 0 {
		return "You have no current page!"
	}
	return RenderPage(b.GetCurrPage())
}

func (b *Browser) GoForward() {
	if len(b.State.Stack) == 0 {
		return
	}
	if b.State.Pos < len(b.State.Stack)-1 {
		b.State.Pos++
	}
	b.State.CurrURL = b.GetCurrPage().URL
}

func (b *Browser) GoBack() {
	if len(b.State.Stack) == 0 {
		return
	}
	if b.State.Pos > 0 {
		b.State.Pos--
	}
	b.State.CurrURL = b.GetCurrPage().URL
}
