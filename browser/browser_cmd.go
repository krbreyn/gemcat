package browser

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/krbreyn/gemcat"
	"github.com/krbreyn/gemcat/data"
)

type BrowserCmd interface {
	Do(s *Browser, args []string) error
	Help() (words []string, desc string)
}

func makeCmdMap() (map[string]BrowserCmd, []BrowserCmd) {
	cm := make(map[string]BrowserCmd)
	cmds := []BrowserCmd{
		ReprintCmd{},
		GotoCmd{},
		LinkCmd{},
		GotoLinkCmd{},
		LinksCmd{},
		StackCmd{},
		HistoryCmd{},
		BackCmd{},
		ForwardCmd{},
		LessCmd{},
		ExitCmd{},
	}

	for _, cmd := range cmds {
		words, _ := cmd.Help()
		for _, w := range words {
			cm[w] = cmd
		}
	}

	return cm, cmds
}

/*

	Navigation Commands

*/

type GotoCmd struct{}

func (c GotoCmd) Do(b *Browser, args []string) error {
	if len(args) == 0 {
		return errors.New("must include a link")
	}

	err := b.GotoURL(strings.TrimPrefix(args[0], "gemini://"))
	if err != nil {
		return fmt.Errorf("err: %v", err)
	}

	fmt.Println(b.RenderOutput())
	return nil
}

func (c GotoCmd) Help() (words []string, desc string) {
	return []string{"goto", "gt"}, "Open and go to a Gemini link.\nUsage: gt [link]"
}

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

type GotoLinkCmd struct{}

func (c GotoLinkCmd) Do(b *Browser, args []string) error {
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

	if strings.HasPrefix(link, "gemini://") {
		link = strings.TrimPrefix(link, "gemini://")
	} else {
		link = strings.TrimPrefix(link, "/")
		link = strings.TrimSuffix(b.State.CurrURL, "/") + "/" + link
	}

	err = b.GotoURL(link)
	if err != nil {
		return fmt.Errorf("err: %v", err)
	}

	fmt.Println(b.RenderOutput())
	return nil
}

func (c GotoLinkCmd) Help() (words []string, desc string) {
	return []string{"g"}, "Open and goto the specified link number on the current page.\nUsage: g [no]"
}

type BackCmd struct{}

func (c BackCmd) Do(b *Browser, args []string) error {
	b.GoBack()
	fmt.Println(b.RenderOutput())
	return nil
}

func (c BackCmd) Help() (words []string, desc string) {
	return []string{"back", "b"}, "Go back one page."
}

type ForwardCmd struct{}

func (c ForwardCmd) Do(b *Browser, args []string) error {
	b.GoForward()
	fmt.Println(b.RenderOutput())
	return nil
}

func (c ForwardCmd) Help() (words []string, desc string) {
	return []string{"forward", "fd", "f"}, "Go forward one page."
}

/*

	Info Commands

*/

type LinksCmd struct{}

func (c LinksCmd) Do(b *Browser, args []string) error {
	if len(b.State.Stack) == 0 {
		return errors.New("There are no links!")
	}

	links := b.GetCurrPage().Links
	if len(links) == 0 {
		fmt.Println("No links!")
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

type StackCmd struct{}

func (c StackCmd) Do(b *Browser, args []string) error {
	if len(b.State.Stack) == 0 {
		fmt.Println("Your stack is empty!")
		return nil
	}

	for i, p := range b.State.Stack {
		if i == b.State.Pos {
			fmt.Print("-> ")
		}
		fmt.Println(i, p.URL)
	}
	return nil
}

func (c StackCmd) Help() (words []string, desc string) {
	return []string{"stack", "st"}, "Print the link stack and your position in it."
}

type HistoryCmd struct{}

func (c HistoryCmd) Do(b *Browser, args []string) error {
	if len(b.State.Data.History) == 0 {
		fmt.Println("Your history is empty!")
		return nil
	}

	for i, h := range b.State.Data.History {
		fmt.Println(i, h)
	}
	return nil
}

func (c HistoryCmd) Help() (words []string, desc string) {
	return []string{"history", "hs"}, "Print the history of visited pages."
}

/*

	Misc Commands

*/

// TODO
type LessCmd struct{}

func (c LessCmd) Do(b *Browser, args []string) error {
	return nil
}

func (c LessCmd) Help() (words []string, desc string) {
	return []string{"less"}, "Will pipe the current page to less when implemented."
}

type ReprintCmd struct{}

func (c ReprintCmd) Do(b *Browser, args []string) error {
	fmt.Println(b.RenderOutput())
	return nil
}

func (c ReprintCmd) Help() (words []string, desc string) {
	return []string{"reprint", "rp"}, "Reprint the current page's contents."
}

type ExitCmd struct{}

func (c ExitCmd) Do(b *Browser, args []string) error {
	data.SaveDataFile(b.State)
	os.Exit(0)
	return nil // should never happen
}

func (c ExitCmd) Help() (words []string, desc string) {
	return []string{"exit", "quit"}, "Exit the program."
}

type HelpCmd struct{}

func (c HelpCmd) Do(cmds []BrowserCmd) {
	for _, cmd := range cmds {
		words, desc := cmd.Help()

		fmt.Print("Command: ")
		cap := len(words) - 1
		for i, w := range words {
			fmt.Print(w)
			if i != cap {
				fmt.Print(", ")
			}
		}
		fmt.Printf("\nDesc: %s\n\n", desc)
	}
}
