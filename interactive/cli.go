package interactive

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/krbreyn/gemcat/browser"
	"github.com/krbreyn/gemcat/data"
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
		b.GotoURL(URL)
		fmt.Println(b.RenderOutput())
	} else if loadLast {
		fmt.Println(b.RenderOutput())
	}

	scanner := bufio.NewScanner(os.Stdin)

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
