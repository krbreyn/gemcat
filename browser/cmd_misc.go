package browser

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/krbreyn/gemcat/data"
)

type CloseCurrentCmd struct{}

type RefreshCmd struct{}

func (c RefreshCmd) Do(b *Browser, args []string) error {
	u, err := url.Parse(b.State.CurrURL)
	if err != nil {
		return err
	}
	if len(b.State.Stack) > 1 {
		b.GoBack()
	} else {
		stem := StackEmptyCmd{}
		if err := stem.Do(b, nil); err != nil {
			return err
		}
	}
	err = b.GotoURLNoCache(u)
	if err != nil {
		return nil
	}
	fmt.Println(b.RenderOutput())
	return nil
}

func (c RefreshCmd) Help() (words []string, desc string) {
	return []string{"rfsh"}, "Refresh the current page."
}

// TODO
type ListCacheCmd struct{}

// TODO - gemcat within gemcat, just recieve print a link without adding it to the stack
// maybe include an option for the link from a page
type JustCatCmd struct{}

type JustCatLess struct{}

type LessCmd struct{}

func (c LessCmd) Do(b *Browser, args []string) error {
	if len(b.State.Stack) == 0 {
		return errors.New("you have no current page")
	}

	cmd := exec.Command("less", "-R")
	cmd.Stdin = strings.NewReader(b.RenderOutput())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (c LessCmd) Help() (words []string, desc string) {
	return []string{"less"}, "Opens the current page using les."
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
	err := data.SaveDataFile(b.State)
	if err != nil {
		panic(err)
	}
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
