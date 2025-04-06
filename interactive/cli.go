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
		err = b.GotoURLCache(u)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(b.RenderOutput())
	} else if loadLast || (u != nil && u.String() == b.State.CurrURL) {
		fmt.Println(b.RenderOutput())
	}

	scanner := bufio.NewScanner(os.Stdin)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		err = data.SaveDataFile(b.State)
		if err != nil {
			panic(err)
		}
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
		b.Sh.HandleInput(b, cmd)
	}
}
