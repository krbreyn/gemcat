package browser

import (
	"fmt"
	"os"

	"github.com/krbreyn/gemcat/data"
)

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
