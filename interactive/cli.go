package interactive

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/krbreyn/gemcat"
	"github.com/krbreyn/gemcat/browser"
	"github.com/krbreyn/gemcat/data"
	"github.com/krbreyn/gemcat/gemtext"
)

func RunCLI(URL string, isURL bool, loadLast bool) {
	b := browser.NewBrowser()

	state, err := data.LoadDataFile()
	if err != nil {
		panic(err)
	}

	if loadLast {
		b.State = state
	} else {
		b.State.Data = state.Data
	}

	fmt.Println("welcome to gemcat\ntype help to see the available commands")

	if isURL {
		_, body, err := browser.Fetch(URL)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}

		content, links := gemtext.DoLinks(body)
		b.State.Stack = []gemcat.Page{{URL: URL, Content: content, Links: links}}

		fmt.Println(b.RenderOutput())
	} else if loadLast {
		fmt.Println(b.RenderOutput())
	}

	scanner := bufio.NewScanner(os.Stdin)

	b.State.CurrURL = URL

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "Error reading input:", err)
			}
			break
		}

		text := scanner.Text()
		cmd := strings.Fields(text)

		b.IH.HandleInput(b, cmd)
	}
}
