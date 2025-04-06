package browser

import (
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

func removeVisitedFromStack(b *Browser, removedURL string) {
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
				link.Visited = false
			}
		}
	}
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
	err = b.GotoURL(u)
	if err != nil {
		return fmt.Errorf("err: %v", err)
	}
	fmt.Println(b.RenderOutput())
	return nil
}

func (c HistoryGotoCmd) Help() (words []string, desc string) {
	return []string{"hsgt", "hsg"}, "Goto and open an item in your history."
}

type HistoryRemoveCmd struct{}

func (c HistoryRemoveCmd) Do(b *Browser, args []string) error {
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

	removedURL := b.State.Data.History[i]
	fmt.Println("deleting", removedURL, "...")
	b.State.Data.History = slices.Delete(b.State.Data.History, i, i+1)
	removeVisitedFromStack(b, removedURL)
	return nil
}

func (c HistoryRemoveCmd) Help() (words []string, desc string) {
	return []string{"hsrm"}, "Remove an item from your history."
}

type HistoryClearAllCmd struct{}

func (c HistoryClearAllCmd) Do(b *Browser, desc string) error {
	l := len(b.State.Data.Bookmarks)
	b.State.Data.Bookmarks = b.State.Data.Bookmarks[:0]
	fmt.Printf("deleted %d bookmarks\n", l)
	return nil
}

func (c HistoryClearAllCmd) Help() (words []string, desc string) {
	return []string{"hscla"}, "Clear your history."
}
