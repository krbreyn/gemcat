package shell

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/krbreyn/gemcat/browser"
)

type (
	ExitCmd struct{}
	TestCmd struct{}

	GotoCmd    struct{}
	ForwardCmd struct{}
	BackCmd    struct{}

	LinkCmd        struct{}
	LinksCmd       struct{}
	LinkCurrentCmd struct{}
	LinkGotoCmd    struct{}

	StackCmd         struct{}
	StackPosCmd      struct{}
	StackRmCmd       struct{} // TODO
	StackCloseCmd    struct{}
	StackCompressCmd struct{}
	StackEmptyCmd    struct{}
	StackGotoCmd     struct{}

	HistoryCmd         struct{}
	HistoryRmCmd       struct{}
	HistoryClearAllCmd struct{}
	HistoryGotoCmd     struct{}

	// TODO
	BookmarkListCmd       struct{}
	BookmarkRmCmd         struct{}
	BookmarkAddLinkCmd    struct{}
	BookmarkAddCurrentCmd struct{}
	BookmarkClearAllCmd   struct{}
	BookmarkSwapCmd       struct{}
	BookmarkGotoCmd       struct{}

	RefreshCmd      struct{} // TODO
	ReprintCmd      struct{}
	CloseCurrentCmd struct{} // TODO
	CacheListCmd    struct{} // TODO
	JustCatCmd      struct{} // TODO
)

func (_ ExitCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	out.RecvMsg("Goodbye!")
	//data.SaveDataFile
	os.Exit(0)
	return nil
}
func (_ ExitCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"exit", "quit"},
		Desc:  "Exit gemcat.",
	}
}

// Test
func (c TestCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	out.RecvMsg("This is a test command!")
	return nil
}
func (c TestCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"test", "tst"},
		Desc:  "A test command!",
	}
}

// Test End

// Nav
func (_ GotoCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	if len(args) == 0 {
		return errors.New("must include a link)")
	}
	link := args[0]
	if !strings.HasPrefix(link, "gemini://") {
		link = "gemini://" + link
	}

	u, err := url.Parse(link)
	if err != nil {
		return err
	}

	out.RecvMsg(fmt.Sprintln("connecting to", u, "..."))
	err = b.GotoURL(u, true)
	if err != nil {
		return err
	}

	out.RecvPage(b.S.CurrPage())
	return nil
}
func (_ GotoCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"goto", "gt"},
		Desc:  "Open and goto a gemini link.\n\tUsage: gt [link]",
	}
}

func (_ ForwardCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	if b.S.Pos == len(b.S.Stack)-1 {
		return errors.New("you can't go forward")
	}
	b.S.GoForward()
	out.RecvPage(b.S.CurrPage())
	return nil
}
func (_ ForwardCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"forward", "fd", "f"},
		Desc:  "Go forward one page.",
	}
}

func (_ BackCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	if b.S.Pos == 0 {
		return errors.New("you can't go back")
	}
	b.S.GoBack()
	out.RecvPage(b.S.CurrPage())
	return nil
}
func (_ BackCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"back", "b"},
		Desc:  "Go back one page.",
	}
}

// Nav End

// Links
func (_ LinkCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	if len(args) == 0 {
		return errors.New("must include link number")
	}

	i, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("links are numbers")
	}

	out.RecvMsg(b.S.CurrPage().Links[i].URL)
	return nil
}
func (_ LinkCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"link", "l"},
		Desc:  "Print the link belonging to the specified number on the current page.\n\tUsage: l [i]",
	}
}

func (_ LinksCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	if len(b.S.Stack) == 0 {
		return errors.New("you have no current page")
	}

	links := b.S.CurrPage().Links
	if len(links) == 0 {
		return errors.New("there are no links on the current page")
	} else {
		for _, l := range links {
			fmt.Println(l.No, l.URL)
		}
	}
	return nil
}
func (_ LinksCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"links", "ls"},
		Desc:  "Print the links accessible from the current page.",
	}
}

func (_ LinkCurrentCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	if b.S.CurrURL() == "" {
		return errors.New("you have no current page")
	} else {
		out.RecvMsg(b.S.CurrURL())
		return nil
	}
}
func (_ LinkCurrentCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"lc"},
		Desc:  "Print the current link.",
	}
}

func (_ LinkGotoCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	if len(args) == 0 {
		return errors.New("must include link number")
	}

	i, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("links are numbers")
	}

	p := b.S.CurrPage()
	if i >= len(p.Links) || i < 0 {
		return errors.New("invalid link number")
	}

	link := p.Links[i].URL
	if !strings.Contains(link, "://") && !strings.HasPrefix(link, "gemini://") {
		if !strings.HasPrefix(link, "/") {
			link = strings.TrimPrefix(b.S.CurrURL(), "/") + "/" + link
		} else {
			u, err := url.Parse(b.S.CurrURL())
			if err != nil {
				return err
			}
			u.Path = ""
			link = strings.TrimPrefix(u.String(), "/") + link
		}
	}

	u, err := url.Parse(link)
	if err != nil {
		return err
	}

	out.RecvMsg(fmt.Sprintf("connecting to %s...", link))
	err = b.GotoURL(u, true)
	if err != nil {
		return err
	}

	out.RecvPage(b.S.CurrPage())
	return nil
}
func (_ LinkGotoCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"lgoto", "lgt"},
		Desc:  "Goto the specified link number on the current page.\n\tUsage: lg [i]",
	}
}

// Links End

// Stack
func (_ StackCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	if len(b.S.Stack) == 0 {
		return errors.New("stack is empty")
	}

	for i, p := range b.S.Stack {
		if i == b.S.Pos {
			out.RecvMsg(fmt.Sprintf("-> %d %s", i, p.URL))
		} else {
			out.RecvMsg(fmt.Sprintf("%d %s", i, p.URL))
		}
	}
	return nil
}
func (_ StackCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"stack", "st"},
		Desc:  "Print the stack and your position in it.",
	}
}

func (_ StackPosCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	out.RecvMsg(strconv.Itoa(b.S.Pos))
	return nil
}
func (_ StackPosCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"stpos"},
		Desc:  "Print the current stack position number.",
	}
}

// TODO
func (_ StackRmCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	out.RecvMsg("Not implemented!")
	return nil
}
func (_ StackRmCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"strm"},
		Desc:  "Remove the specified item from the stack.",
	}
}

func (_ StackCloseCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	old := len(b.S.Stack)
	b.S.Stack = b.S.Stack[:b.S.Pos+1]
	out.RecvMsg(fmt.Sprintf("closed %d pages", old-len(b.S.Stack)))
	return nil
}
func (_ StackCloseCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"stcl"},
		Desc:  "Close every page beneath the current stack position",
	}
}

func (_ StackCompressCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	old := len(b.S.Stack)
	b.S.Stack = b.S.Stack[b.S.Pos:]
	b.S.Pos = 0
	out.RecvMsg(fmt.Sprintf("closed %d pages", old-len(b.S.Stack)))
	return nil
}
func (_ StackCompressCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"stcmp"},
		Desc:  "Closes every page above the current stack position.",
	}
}

func (_ StackEmptyCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	l := len(b.S.Stack)
	b.S.Stack = b.S.Stack[:0]
	b.S.Pos = 0
	out.RecvMsg(fmt.Sprintf("closed %d pages", l))
	return nil
}
func (_ StackEmptyCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"stem"},
		Desc:  "Empties the stack, closing all pages.",
	}
}

func (_ StackGotoCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	if len(args) == 0 {
		return errors.New("must include stack item number")
	}

	i, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("not a number")
	}

	if i < 0 || i > len(b.S.Stack)-1 {
		return errors.New("stack item number is out of range")
	}

	b.S.Pos = i
	out.RecvPage(b.S.CurrPage())
	return nil
}
func (_ StackGotoCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"stgt"},
		Desc:  "Leap to the specified stack item number.",
	}
}

// Stack End

// History
func (_ HistoryCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	if len(b.D.History) == 0 {
		return errors.New("history is empty")
	}

	for i, h := range b.D.History {
		out.RecvMsg(fmt.Sprintf("%d %s", i, h))
	}
	return nil
}
func (_ HistoryCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"hs"},
		Desc:  "Print the history of visited pages.",
	}
}

func (_ HistoryRmCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	if len(args) == 0 {
		return errors.New("must include history item number")
	}

	i, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("not a number")
	}

	if i < 0 || i > len(b.D.History)-1 {
		return errors.New("history item number is out of range")
	}

	removedURL := b.D.History[i]
	out.RecvMsg(fmt.Sprintf("deleting %s...", removedURL))
	b.D.History = slices.Delete(b.D.History, i, i+1)
	return nil
}
func (_ HistoryRmCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"hsrm"},
		Desc:  "Remove an item from your history",
	}
}

func (_ HistoryClearAllCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	l := len(b.D.History)
	b.D.History = b.D.History[:0]
	out.RecvMsg(fmt.Sprintf("deleted %d bookmarks", l))
	return nil
}
func (_ HistoryClearAllCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"hscla"},
		Desc:  "Clear your history.",
	}
}

func (_ HistoryGotoCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	if len(args) == 0 {
		return errors.New("must include history item number")
	}

	i, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("not a number!")
	}

	if i < 0 || i > len(b.D.History)-1 {
		return errors.New("history item number is out of range")
	}

	u, err := url.Parse(b.D.History[i])
	if err != nil {
		return err
	}
	err = b.GotoURL(u, true)
	if err != nil {
		return err
	}

	out.RecvPage(b.S.CurrPage())
	return nil
}
func (_ HistoryGotoCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"hsgt"},
		Desc:  "Open and goto and item in your history.",
	}
}

// History End

// Misc
func (_ ReprintCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	out.RecvPage(b.S.CurrPage())
	return nil
}
func (_ ReprintCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"reprint", "rp"},
		Desc:  "Reprint the current page's contents.",
	}
}

// Misc End
