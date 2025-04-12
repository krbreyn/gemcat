package browser

import (
	"bufio"
	"net/url"
	"slices"
	"strings"
)

type Browser struct {
	S State
	D Data
}

// TODO
type Settings struct {
	UseCache  bool
	PageWidth bool
}

func (b *Browser) GotoURL(url *url.URL, doCache bool) error {
	u := url.String()

	if !slices.Contains(b.D.History, u) {
		b.D.History = append(b.D.History, u)
	}

	_, body, err := FetchGemini(url, true)
	if err != nil {
		return err
	}

	links := ParseLinks(body)

	if len(b.S.Stack) != 0 {
		b.S.Pos++
	}

	p := Page{
		URL:     u,
		Content: body,
		Links:   links,
	}

	if b.S.Pos == len(b.S.Stack) {
		b.S.Stack = append(b.S.Stack, p)
	} else {
		b.S.Stack = append(b.S.Stack[:b.S.Pos], p)
	}

	return nil
}

type State struct {
	Pos   int
	Stack []Page
}

func (s *State) CurrPage() Page {
	if len(s.Stack) == 0 || s.Pos > len(s.Stack)-1 {
		return Page{
			URL:     "",
			Content: "",
		}
	}
	return s.Stack[s.Pos]
}

func (s *State) CurrURL() string {
	return s.CurrPage().URL
}

func (s *State) GoForward() {
	if len(s.Stack) == 0 {
		return
	}
	if s.Pos < len(s.Stack)-1 {
		s.Pos++
	}
}

func (s *State) GoBack() {
	if len(s.Stack) == 0 {
		return
	}
	if s.Pos > 0 {
		s.Pos--
	}
}

// TODO
func (s State) ToJson() []byte {
	return nil
}

// TODO
func StateFromJson(b []byte) (State, error) {
	return State{}, nil
}

type Data struct {
	Bookmarks []string
	History   []string
}

// TODO
func (d Data) ToJson() []byte {
	return nil
}

// TODO
func DataFromJson(b []byte) (Data, error) {
	return Data{}, nil
}

type Page struct {
	URL     string
	Content string
	Links   []Link
}

type Link struct {
	No  int
	URL string
}

func ParseLinks(body string) []Link {
	scanner := bufio.NewScanner(strings.NewReader(body))

	var i int
	var links []Link

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "=>") {
			split := strings.Fields(line)
			url := split[1]

			links = append(links, Link{
				No:  i,
				URL: url,
			})
		}
	}

	return links
}
