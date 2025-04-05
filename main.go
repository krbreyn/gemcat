package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
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
		RunTUI(URL, isURL)
		os.Exit(0)
	}

	switch *interactiveMode {
	case false:
		if isURL {
			_, body, err := Fetch(URL)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(ColorGemtext(body, nil))
			os.Exit(0)
		} else {
			fmt.Println("error: must include URL if not using interactive mode")
			os.Exit(1)
		}

	case true:
		RunCLI(URL, isURL)
		os.Exit(0)
	}
}

func getHostPath(url string) (host, path string) {
	split := strings.SplitN(url, "/", 2)
	if len(split) == 1 {
		host, path = split[0], ""
	} else {
		host, path = split[0], split[1]
	}
	return host, path
}

func isGeminiLink(url string) bool {
	return strings.HasPrefix(url, "gemini://") || !strings.Contains(url, "://")
}
