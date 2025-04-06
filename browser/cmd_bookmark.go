package browser

import (
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

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

	err = b.GotoURL(u)
	if err != nil {
		return fmt.Errorf("err: %v", err)
	}
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

type BookmarkClearAllCmd struct{}

func (c BookmarkClearAllCmd) Do(b *Browser, args []string) error {
	l := len(b.State.Data.Bookmarks)
	b.State.Data.Bookmarks = b.State.Data.Bookmarks[:0]
	fmt.Printf("deleted %d bookmarks\n", l)
	return nil
}

func (c BookmarkClearAllCmd) Help() (words []string, desc string) {
	return []string{"bmcla"}, "Remove all items from your bookmarks."
}

type BookmarkSwapCmd struct{}

func (c BookmarkSwapCmd) Do(b *Browser, args []string) error {
	if len(args) != 2 {
		return errors.New("must include bookmark numbers to swap")
	}

	i1, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("item 1 is not a number!")
	}
	i2, err := strconv.Atoi(args[1])
	if err != nil {
		return errors.New("item 2 is not a number!")
	}

	if i1 < 0 || i1 > len(b.State.Data.Bookmarks)-1 {
		return errors.New("bookmark number 1 is out of range")
	}
	if i2 < 0 || i2 > len(b.State.Data.Bookmarks)-1 {
		return errors.New("bookmark number 2 is out of range")
	}

	temp := b.State.Data.Bookmarks[i1]
	b.State.Data.Bookmarks[i1] = b.State.Data.Bookmarks[i2]
	b.State.Data.Bookmarks[i2] = temp
	fmt.Printf("swapped %d and %d\n", i1, i2)
	return nil
}

func (c BookmarkSwapCmd) Help() (words []string, desc string) {
	return []string{"bmsw"}, "Swap the places of two bookmark items."
}
