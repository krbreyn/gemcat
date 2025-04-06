package browser

import (
	"errors"
	"fmt"
	"strconv"
)

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
	b.State.CurrURL = b.State.Stack[b.State.Pos].URL
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

type StackCompressCmd struct{}

func (c StackCompressCmd) Do(b *Browser, args []string) error {
	old := len(b.State.Stack)
	b.State.Stack = b.State.Stack[b.State.Pos:]
	b.State.Pos = 0
	fmt.Printf("closed %d pages\n", old-len(b.State.Stack))
	return nil
}

func (c StackCompressCmd) Help() (words []string, desc string) {
	return []string{"stcmp"}, "Closes every page above the current stack position."
}

type StackEmptyCmd struct{}

func (c StackEmptyCmd) Do(b *Browser, args []string) error {
	l := len(b.State.Stack)
	b.State.Stack = b.State.Stack[:0]
	b.State.CurrURL = ""
	b.State.Pos = 0
	fmt.Printf("closed %d pages\n", l)
	return nil
}

func (c StackEmptyCmd) Help() (words []string, desc string) {
	return []string{"stem"}, "Empties the stack and closes all pages."
}
