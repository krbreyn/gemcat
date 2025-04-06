package browser

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

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
