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

func NeedsOneNum(args []string) (int, error) {
	if len(args) == 0 {
		return 0, errors.New("must include one number")
	}

	i, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, errors.New("arg[0] not a number")
	}
	return i, nil
}

func NeedsTwoNums(args []string) (int, int, error) {
	if len(args) != 2 {
		return 0, 0, errors.New("must include two numbers")
	}

	i1, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, 0, errors.New("item 1 is not a number")
	}
	i2, err := strconv.Atoi(args[1])
	if err != nil {
		return 0, 0, errors.New("item 2 is not a number")
	}

	return i1, i2, nil
}

func NormalizeRelativeLink(link string, b *browser.Browser) (string, error) {
	if !strings.Contains(link, "://") && !strings.HasPrefix(link, "gemini://") {
		if !strings.HasPrefix(link, "/") {
			link = strings.TrimPrefix(b.S.CurrURL(), "/") + "/" + link
		} else {
			u, err := url.Parse(b.S.CurrURL())
			if err != nil {
				return "", err
			}
			u.Path = ""
			link = strings.TrimPrefix(u.String(), "/") + link
		}
	}
	return link, nil
}

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

	BookmarksCmd          struct{}
	BookmarkRmCmd         struct{}
	BookmarkAddLinkCmd    struct{}
	BookmarkAddCurrentCmd struct{}
	BookmarkSwapCmd       struct{}
	BookmarkClearAllCmd   struct{}
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
	i, err := NeedsOneNum(args)
	if err != nil {
		return err
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
	i, err := NeedsOneNum(args)
	if err != nil {
		return err
	}

	p := b.S.CurrPage()
	if i >= len(p.Links) || i < 0 {
		return errors.New("invalid link number")
	}

	link := p.Links[i].URL
	link, err = NormalizeRelativeLink(link, b)
	if err != nil {
		return err
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
		Desc:  "Remove the specified item from the stack.\n\tUsage: strm [i]",
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
	i, err := NeedsOneNum(args)
	if err != nil {
		return err
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
		Desc:  "Leap to the specified stack item number.\n\tUsage: stgt [i]",
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
	i, err := NeedsOneNum(args)
	if err != nil {
		return err
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
		Desc:  "Remove an item from your history\n\tUsage: hsrm [i]",
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
	i, err := NeedsOneNum(args)
	if err != nil {
		return err
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
		Desc:  "Open and goto and item in your history.\n\tUsage:hsgt [i]",
	}
}

// History End

// Bookmarks
func (_ BookmarksCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	if len(b.D.Bookmarks) == 0 {
		return errors.New("bookmarks is empty")
	}

	for i, b := range b.D.Bookmarks {
		out.RecvMsg(fmt.Sprintf("%d %s", i, b))
	}
	return nil
}
func (_ BookmarksCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"bml"},
		Desc:  "List your bookmarks.",
	}
}

func (_ BookmarkRmCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	i, err := NeedsOneNum(args)
	if err != nil {
		return err
	}

	if i < 0 || i > len(b.D.Bookmarks)-1 {
		return errors.New("bookmark number is out of range")
	}

	removedURL := b.D.Bookmarks[i]
	out.RecvMsg(fmt.Sprintf("deleting %s ...", removedURL))
	b.D.Bookmarks = slices.Delete(b.D.Bookmarks, i, i+1)
	return nil
}
func (_ BookmarkRmCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"bmrm"},
		Desc:  "Removes a bookmark.\n\tUsage: bmrm [i]",
	}
}

func (_ BookmarkAddLinkCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	i, err := NeedsOneNum(args)
	if err != nil {
		return err
	}

	p := b.S.CurrPage()
	if i >= len(p.Links) || i < 0 {
		return errors.New("invalid link number")
	}

	link := p.Links[i].URL
	link, err = NormalizeRelativeLink(link, b)
	if err != nil {
		return err
	}

	if slices.Contains(b.D.Bookmarks, link) {
		return errors.New("bookmarks already contains this url")
	}

	b.D.Bookmarks = append(b.D.Bookmarks, link)
	out.RecvMsg(fmt.Sprintf("added %s to bookmarks", link))
	return nil
}
func (_ BookmarkAddLinkCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"bmal"},
		Desc:  "Add a link from the current page to your bookmarks.",
	}
}

func (_ BookmarkAddCurrentCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	url := b.S.CurrURL()
	if url == "" {
		return errors.New("current page is empty")
	}
	if slices.Contains(b.D.Bookmarks, url) {
		return errors.New("bookmarks already contains this url")
	}
	b.D.Bookmarks = append(b.D.Bookmarks, url)
	out.RecvMsg(fmt.Sprintf("added %s to bookmarks", url))
	return nil
}
func (_ BookmarkAddCurrentCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"bmac"},
		Desc:  "Add the current page to your bookmarks.",
	}
}

func (_ BookmarkSwapCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	i1, i2, err := NeedsTwoNums(args)
	if err != nil {
		return err
	}
	if i1 < 0 || i1 > len(b.D.Bookmarks)-1 {
		return errors.New("bookmark number 1 is out of range")
	}
	if i2 < 0 || i2 > len(b.D.Bookmarks)-1 {
		return errors.New("bookmark number 2 is out of range")
	}

	temp := b.D.Bookmarks[i1]
	b.D.Bookmarks[i1] = b.D.Bookmarks[i2]
	b.D.Bookmarks[i2] = temp
	out.RecvMsg(fmt.Sprintf("swapped %d and %d", i1, i2))
	return nil
}
func (_ BookmarkSwapCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"bmsw"},
		Desc:  "Swap the places of two bookmark items\n\tUsage: bmsw [i1] [i2]",
	}
}

func (_ BookmarkClearAllCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	l := len(b.D.Bookmarks)
	b.D.Bookmarks = b.D.Bookmarks[:0]
	out.RecvMsg(fmt.Sprintf("deleted %d bookmarks", l))
	return nil
}
func (_ BookmarkClearAllCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"bmcla"},
		Desc:  "Remove all items from your bookmarks.",
	}
}

func (_ BookmarkGotoCmd) Do(b *browser.Browser, out ShellOut, args []string) error {
	i, err := NeedsOneNum(args)
	if err != nil {
		return err
	}

	if i < 0 || i > len(b.D.Bookmarks)-1 {
		return errors.New("bookmark number is out of range")
	}

	u, err := url.Parse(b.D.Bookmarks[i])
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
func (_ BookmarkGotoCmd) Help() HelpInfo {
	return HelpInfo{
		Words: []string{"bmgt"},
		Desc:  "Goto the specified bookmark numbers.\n\tUsage: bmgt [i]",
	}
}

// Bookmarks End

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
