package browser

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/krbreyn/gemcat"
)

type LinkCmd struct{}

func (c LinkCmd) Do(b *Browser, args []string) error {
	if len(args) == 0 {
		return errors.New("must include link number")
	}

	i, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("links are numbers")
	}

	fmt.Println(b.GetCurrPage().Links[i].URL)
	return nil
}

func (c LinkCmd) Help() (words []string, desc string) {
	return []string{"link", "l"}, "Print the link of the corresponding link number from the current page.\nUsage: l [no]"
}

type LinksCmd struct{}

func (c LinksCmd) Do(b *Browser, args []string) error {
	if len(b.State.Stack) == 0 {
		return errors.New("you have no page")
	}

	links := b.GetCurrPage().Links
	if len(links) == 0 {
		return errors.New("there are no links")
	} else {
		for _, s := range b.GetCurrPage().Links {
			fmt.Println(s.No, s.URL)
		}
	}
	return nil
}

func (c LinksCmd) Help() (words []string, desc string) {
	return []string{"links", "ls"}, "List the links accessable from the current page."
}

type LinkCurrentCmd struct{}

func (c LinkCurrentCmd) Do(b *Browser, args []string) error {
	if b.State.CurrURL == "" {
		return errors.New("you have no current page")
	} else {
		fmt.Println(b.State.CurrURL)
	}
	return nil
}

func (c LinkCurrentCmd) Help() (words []string, desc string) {
	return []string{"lc"}, "Print the current link."
}

type LinkGotoCmd struct{}

func (c LinkGotoCmd) Do(b *Browser, args []string) error {
	if len(args) == 0 {
		return errors.New("must include link number")
	}

	i, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("links are numbers")
	}

	p := b.GetCurrPage()
	if i >= len(p.Links) || i < 0 {
		return errors.New("invalid link number")
	}

	link := p.Links[i].URL
	p.Links[i].Visited = true

	if !gemcat.IsGeminiLink(link) {
		return fmt.Errorf("cannot open link type of %s", link)
	}

	if !strings.HasPrefix(link, "gemini://") {
		link = strings.TrimPrefix(link, "/")
		link = strings.TrimSuffix(b.State.CurrURL, "/") + "/" + link
	}

	u, err := url.Parse(link)
	if err != nil {
		return err
	}

	err = b.GotoURLCache(u)
	if err != nil {
		return fmt.Errorf("err: %v", err)
	}

	fmt.Println(b.RenderOutput())
	return nil
}

func (c LinkGotoCmd) Help() (words []string, desc string) {
	return []string{"gtl", "g"}, "Open and goto the specified link number on the current page.\nUsage: g [no]"
}
