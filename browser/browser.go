package browser

import (
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"

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

func (b *Browser) WasLinkVisited(url string) bool {
	if !strings.HasPrefix(url, "gemini://") {
		url = b.State.CurrURL + "/" + url
	}
	fmt.Println("checking", url)
	return slices.Contains(b.State.Data.History, url)
}

func (b *Browser) IsLinkBookmarked(url string) bool {
	if !strings.HasPrefix(url, "gemini://") {
		url = b.State.CurrURL + "/" + url
	}
	fmt.Println("checking", url)
	return slices.Contains(b.State.Data.Bookmarks, url)
}

func (b *Browser) GotoURL(url *url.URL) error {
	if !slices.Contains(b.State.Data.History, url.String()) {
		b.State.Data.History = append(b.State.Data.History, url.String())
	}

	_, body, err := Fetch(url)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	b.State.CurrURL = url.String()

	_, links := gemtext.DoLinks(body, b.WasLinkVisited, b.IsLinkBookmarked)

	if len(b.State.Stack) != 0 {
		b.State.Pos++
	}
	if b.State.Pos == len(b.State.Stack) {
		b.State.Stack = append(b.State.Stack, gemcat.Page{
			URL:     url.String(),
			Content: body,
			Links:   links,
		})
	} else {
		b.State.Stack = append(b.State.Stack[:b.State.Pos], gemcat.Page{
			URL:     url.String(),
			Content: body,
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
	p := b.GetCurrPage()
	content, _ := gemtext.DoLinks(p.Content, b.WasLinkVisited, b.IsLinkBookmarked)
	p.Content = content
	return RenderPage(p)
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
