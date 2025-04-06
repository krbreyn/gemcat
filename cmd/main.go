package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/krbreyn/gemcat"
	"github.com/krbreyn/gemcat/browser"
	"github.com/krbreyn/gemcat/gemtext"
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
	var u *url.URL
	if len(args) == 0 {
		isURL = false
	} else {
		URL := args[0]
		if !strings.HasPrefix(URL, "gemini://") {
			URL = "gemini://" + URL
		}
		var err error
		u, err = url.Parse(URL)
		if err != nil {
			fmt.Printf("url parse error: %v", err)
			os.Exit(1)
		}
		isURL = true
	}

	if !*cliMode && !*tuiMode {
		if isURL {
			_, body, err := browser.FetchGemini(u, true)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				os.Exit(1)
			}
			_, links := gemtext.DoLinks(body, func(url string) bool { return false }, func(url string) bool { return false })
			p := gemcat.Page{
				URL:     u.String(),
				Content: body,
				Links:   links,
			}
			width, _, err := term.GetSize(int(os.Stdout.Fd()))
			if err != nil {
				width = 80
			}
			fmt.Println(wordwrap.String(browser.RenderPage(p), width))
			os.Exit(0)
		} else {
			fmt.Println("error: must include URL if not using interactive mode")
			os.Exit(1)
		}
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
