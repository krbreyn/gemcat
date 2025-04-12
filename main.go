package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/krbreyn/gemcat/browser"
	"github.com/krbreyn/gemcat/gemtxt"
	"github.com/krbreyn/gemcat/interactive"
	"github.com/muesli/reflow/wordwrap"
	"golang.org/x/term"
)

func main() {
	cliMode := flag.Bool("i", false, "CLI mode")
	tuiMode := flag.Bool("t", false, "TUI mode")
	loadLast := flag.Bool("ll", false, "Load last session")
	help := flag.Bool("help", false, "Help")

	flag.Parse()
	args := flag.Args()
	argc := len(args)

	if *tuiMode && *cliMode {
		die("Pick only CLI mode or TUI mode!")
	}

	if *loadLast && (!*cliMode && !*tuiMode) {
		die("'-ll' cannot be used outside of interactive mode!")
	}

	if argc == 0 && (!*cliMode && !*tuiMode) {
		die("Must include URL if not using interactive mode!")
	}

	if *help || (argc > 2 && args[1] == "help") {
		die("todo")
	}

	var isURL bool
	var u *url.URL
	if argc == 0 {
		isURL = false
	} else {
		URL := args[0]
		if !strings.HasPrefix(URL, "gemini://") {
			URL = "gemini://" + URL
		}
		var err error
		u, err = url.Parse(URL)
		if err != nil {
			die(err.Error())
		}
		isURL = true
	}

	if !*cliMode && !*tuiMode {
		if !isURL {
			fmt.Println("err: must include URL if not using interactive mode")
			os.Exit(1)
		}

		_, body, err := browser.FetchGemini(u, true)
		if err != nil {
			die(err.Error())
		}

		width, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			width = 80
		}

		fmt.Println(wordwrap.String(gemtxt.ColorPlain(body), width))
		os.Exit(0)
	}

	if *cliMode {
		interactive.RunCLI(u, isURL, *loadLast)
		os.Exit(0)
	}

	if *tuiMode {
		interactive.RunTUI(u, isURL, *loadLast)
		os.Exit(0)
	}
}

func die(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
