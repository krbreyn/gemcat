package browser

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"slices"
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
		GotoCmd{},
		GotoLinkCmd{},
		BackCmd{},
		ForwardCmd{},
		LinkCmd{},
		LinksCmd{},
		StackCmd{},
		StackGotoCmd{},
		StackCloseCmd{},
		StackEmptyCmd{},
		HistoryCmd{},
		HistoryGotoCmd{},
		BookmarkGotoCmd{},
		BookmarkListCmd{},
		BookmarkAddCurrentCmd{},
		BookmarkAddLinkCmd{},
		BookmarkRemoveCmd{},
		BookmarkRemoveCurrentCmd{},
		LessCmd{},
		ReprintCmd{},
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

// TODO cache list cmd

/*

	Navigation Commands

*/

type GotoCmd struct{}

func (c GotoCmd) Do(b *Browser, args []string) error {
	if len(args) == 0 {
		return errors.New("must include a link")
	}

	link := args[0]
	if strings.HasPrefix(link, "gemini://") {
		link = "gemini://" + link
	}

	u, err := url.Parse(link)
	if err != nil {
		return err
	}

	err = b.GotoURL(u)
	if err != nil {
		return fmt.Errorf("err: %v", err)
	}

	fmt.Println(b.RenderOutput())
	return nil
}

func (c GotoCmd) Help() (words []string, desc string) {
	return []string{"goto", "gt"}, "Open and go to a Gemini link.\nUsage: gt [link]"
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

	if !strings.HasPrefix(link, "gemini://") {
		link = strings.TrimPrefix(link, "/")
		link = strings.TrimSuffix(b.State.CurrURL, "/") + "/" + link
	}

	u, err := url.Parse(link)
	if err != nil {
		return err
	}

	err = b.GotoURL(u)
	if err != nil {
		return fmt.Errorf("err: %v", err)
	}

	fmt.Println(b.RenderOutput())
	return nil
}

func (c GotoLinkCmd) Help() (words []string, desc string) {
	return []string{"gtl", "g"}, "Open and goto the specified link number on the current page.\nUsage: g [no]"
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

type StackGotoCmd struct{}

func (c StackGotoCmd) Do(b *Browser, args []string) error {
	if len(args) == 0 {
		return errors.New("must include stack item number")
	}

	i, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("not a number!")
	}

	if i < 0 || i > len(b.State.Stack)-1 {
		return errors.New("stack item number is out of range")
	}

	b.State.Pos = i
	fmt.Println(b.RenderOutput())
	return nil
}

func (c StackGotoCmd) Help() (words []string, desc string) {
	return []string{"stgt", "stg"}, "Leap to the stack item number."
}

type StackCloseCmd struct{}

func (c StackCloseCmd) Do(b *Browser, args []string) error {
	old := len(b.State.Stack)
	b.State.Stack = b.State.Stack[:b.State.Pos+1]
	fmt.Printf("closed %d pages\n", old-len(b.State.Stack))
	return nil
}

func (c StackCloseCmd) Help() (words []string, desc string) {
	return []string{"stcl"}, "Closes every page beneath the current stack position."
}

type StackEmptyCmd struct{}

func (c StackEmptyCmd) Do(b *Browser, args []string) error {
	l := len(b.State.Stack)
	b.State.Stack = b.State.Stack[:0]
	fmt.Printf("closed %d pages\n", l)
	return nil
}

func (c StackEmptyCmd) Help() (words []string, desc string) {
	return []string{"stem"}, "Empties the stack and closes all pages."
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

type HistoryGotoCmd struct{}

func (c HistoryGotoCmd) Do(b *Browser, args []string) error {
	if len(args) == 0 {
		return errors.New("must include history item number")
	}

	i, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("not a number!")
	}

	if i < 0 || i > len(b.State.Data.History)-1 {
		return errors.New("history item number is out of range")
	}

	u, err := url.Parse(b.State.Data.History[i])
	if err != nil {
		return err
	}
	b.GotoURL(u)
	fmt.Println(b.RenderOutput())
	return nil
}

func (c HistoryGotoCmd) Help() (words []string, desc string) {
	return []string{"hsgt", "hsg"}, "Goto and open an item in your history."
}

/*

	Bookmark Commands

*/

type BookmarkGotoCmd struct{}

func (c BookmarkGotoCmd) Do(b *Browser, args []string) error {
	if len(args) == 0 {
		return errors.New("must include bookmark number")
	}

	i, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("not a number!")
	}

	if i < 0 || i > len(b.State.Data.Bookmarks)-1 {
		return errors.New("bookmark number is out of range")
	}

	u, err := url.Parse(b.State.Data.Bookmarks[i])
	if err != nil {
		return err
	}

	b.GotoURL(u)
	fmt.Println(b.RenderOutput())
	return nil
}

func (c BookmarkGotoCmd) Help() (words []string, desc string) {
	return []string{"bmgt", "bmg"}, "Goto the specified bookmark number.\nUsage bmg [i]"
}

type BookmarkListCmd struct{}

func (c BookmarkListCmd) Do(b *Browser, args []string) error {
	if len(b.State.Data.Bookmarks) == 0 {
		fmt.Println("You have no bookmarks!")
		return nil
	}
	for i, b := range b.State.Data.Bookmarks {
		fmt.Println(i, b)
	}
	return nil
}

func (c BookmarkListCmd) Help() (words []string, desc string) {
	return []string{"bml"}, "List your bookmarks."
}

// TODO
type BookmarkAddCmd struct{}

type BookmarkAddCurrentCmd struct{}

func (c BookmarkAddCurrentCmd) Do(b *Browser, args []string) error {
	url := b.State.CurrURL
	b.State.Data.Bookmarks = append(b.State.Data.Bookmarks, url)
	fmt.Printf("added %s to bookmarks\n", url)
	return nil
}

func (c BookmarkAddCurrentCmd) Help() (words []string, desc string) {
	return []string{"bmac"}, "Add the current page to your bookmarks."
}

type BookmarkAddLinkCmd struct{}

func (c BookmarkAddLinkCmd) Do(b *Browser, args []string) error {
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
	p.Links[i].Bookmarked = true

	if !strings.HasPrefix(link, "gemini://") {
		link = strings.TrimPrefix(link, "/")
		link = strings.TrimSuffix(b.State.CurrURL, "/") + "/" + link
	}

	b.State.Data.Bookmarks = append(b.State.Data.Bookmarks, link)
	fmt.Printf("added %s to bookmarks\n", link)
	return nil
}

func (c BookmarkAddLinkCmd) Help() (words []string, desc string) {
	return []string{"bmal"}, "Add a link from the page to your bookmarks"
}

type BookmarkRemoveCmd struct{}

func (c BookmarkRemoveCmd) Do(b *Browser, args []string) error {
	if len(args) == 0 {
		return errors.New("must include bookmark number")
	}

	i, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("not a number!")
	}

	if i < 0 || i > len(b.State.Data.Bookmarks)-1 {
		return errors.New("bookmark number is out of range")
	}

	removedURL := b.State.Data.Bookmarks[i]
	fmt.Println("deleting", removedURL, "...")
	b.State.Data.Bookmarks = slices.Delete(b.State.Data.Bookmarks, i, i+1)

	removeBookmarkFromStack(b, removedURL)
	return nil
}

func removeBookmarkFromStack(b *Browser, removedURL string) {
	for pageIndex := range b.State.Stack {
		for linkIndex := range b.State.Stack[pageIndex].Links {
			link := &b.State.Stack[pageIndex].Links[linkIndex]
			linkURL := link.URL

			if !strings.HasPrefix(linkURL, "gemini://") && !strings.Contains(linkURL, "://") {
				baseURL := b.State.Stack[pageIndex].URL
				linkURL = strings.TrimPrefix(linkURL, "/")
				linkURL = strings.TrimSuffix(baseURL, "/") + "/" + linkURL
			}

			if linkURL == removedURL {
				link.Bookmarked = false
			}
		}
	}
}

func (c BookmarkRemoveCmd) Help() (words []string, desc string) {
	return []string{"bmrm"}, "Removes a bookmark.\nUsage: bmrm [i]"
}

type BookmarkRemoveCurrentCmd struct{}

func (c BookmarkRemoveCurrentCmd) Do(b *Browser, args []string) error {
	var deletedOnce bool
rerun:
	if slices.Contains(b.State.Data.Bookmarks, b.State.CurrURL) {
		for i, u := range b.State.Data.Bookmarks {
			if u == b.State.CurrURL {
				removedURL := b.State.Data.Bookmarks[i]
				fmt.Println("deleting", removedURL, "...")
				b.State.Data.Bookmarks = slices.Delete(b.State.Data.Bookmarks, i, i+1)
				if !deletedOnce {
					removeBookmarkFromStack(b, removedURL)
					deletedOnce = true
				}
				goto rerun
			}
		}
	}
	return nil
}

func (c BookmarkRemoveCurrentCmd) Help() (words []string, desc string) {
	return []string{"bmrmc"}, "Removes the current page the bookmarks, if it is so.\nUsage: bmrmc [i]"
}

/*

	Misc Commands

*/

// TODO - refresh and download a fresh page
type RefreshCmd struct{}

// TODO
type ListCacheCmd struct{}

// TODO - gemcat within gemcat, just recieve print a link without adding it to the stack
// maybe include an option for the link from a page
type JustCatCmd struct{}

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
