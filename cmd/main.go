package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/krbreyn/gemcat/browser"
	"github.com/krbreyn/gemcat/gemtext"
	"github.com/krbreyn/gemcat/interactive"
)

func main() {
	cliMode := flag.Bool("i", false, "CLI mode")
	tuiMode := flag.Bool("t", false, "TUI mode")
	loadLast := flag.Bool("ll", false, "Load last session")

	help := flag.Bool("help", false, "Help")
	flag.Parse()

	if *tuiMode && *cliMode {
		fmt.Println("Pick only one!")
		os.Exit(1)
	}

	if *loadLast && (!*cliMode && !*tuiMode) {
		fmt.Println("'-ll' cannot be used outside of interactive mode!")
		os.Exit(1)
	}

	if (len(os.Args) > 2 && os.Args[1] == "help") || *help {
		fmt.Println("todo")
		return
	}

	args := flag.Args()

	var isURL bool
	var URL string
	if len(args) == 0 {
		isURL = false
	} else {
		URL = args[0]
		URL = strings.TrimPrefix(URL, "gemini://")
		isURL = true
	}

	if *tuiMode {
		interactive.RunTUI(URL, isURL, *loadLast)
		os.Exit(0)
	}

	switch *cliMode {
	case false:
		if isURL {
			_, body, err := browser.Fetch(URL)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(gemtext.ColorGemtext(body, nil))
			os.Exit(0)
		} else {
			fmt.Println("error: must include URL if not using interactive mode")
			os.Exit(1)
		}

	case true:
		interactive.RunCLI(URL, isURL, *loadLast)
		os.Exit(0)
	}
}
