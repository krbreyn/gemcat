package interactive

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/krbreyn/gemcat/browser"
	"github.com/krbreyn/gemcat/data"
)

func RunCLI(u *url.URL, isURL bool, loadLast bool) {
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

	if isURL && u.String() != b.State.CurrURL {
		b.GotoURL(u)
		fmt.Println(b.RenderOutput())
	} else if loadLast || u.String() == b.State.CurrURL {
		fmt.Println(b.RenderOutput())
	}

	scanner := bufio.NewScanner(os.Stdin)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		data.SaveDataFile(b.State)
		os.Exit(0)
	}()

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "Error reading input:", err)
			}
			break
		}

		cmd := strings.Fields(scanner.Text())
		b.IH.HandleInput(b, cmd)
	}
}
