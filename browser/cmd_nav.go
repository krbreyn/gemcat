package browser

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type GotoCmd struct{}

func (c GotoCmd) Do(b *Browser, args []string) error {
	if len(args) == 0 {
		return errors.New("must include a link")
	}

	link := args[0]
	if !strings.HasPrefix(link, "gemini://") {
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
