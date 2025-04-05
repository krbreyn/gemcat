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
	interactiveMode := flag.Bool("i", false, "Interactive mode")
	tuiMode := flag.Bool("t", false, "TUI mode")

	help := flag.Bool("help", false, "Help")
	flag.Parse()

	if *tuiMode && *interactiveMode {
		fmt.Println("Pick only one!")
		os.Exit(1)
	}

	if (len(os.Args) > 2 && os.Args[1] == "help") || *help {
		fmt.Println("todo")
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
		interactive.RunTUI(URL, isURL)
		os.Exit(0)
	}

	switch *interactiveMode {
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
		interactive.RunCLI(URL, isURL)
		os.Exit(0)
	}
}
