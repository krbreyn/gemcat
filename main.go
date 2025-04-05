package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	interactiveMode := flag.Bool("i", false, "Interactive mode")
	help := flag.Bool("help", false, "Help")
	flag.Parse()

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

	switch *interactiveMode {
	case false:
		if isURL {
			host, path := getHostPath(URL)
			_, body, err := Fetch(host, path)
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
		var b Browser

		if isURL {
			host, path := getHostPath(URL)

			_, body, err := Fetch(host, path)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				os.Exit(1)
			}

			content, links := DoLinks(body)
			b.Stack = []Page{{URL: URL, Content: content, Links: links}}

			fmt.Println(b.RenderCurrPage())
		}
		fmt.Println("welcome to gemcat\ntype help to see the available commands")

		scanner := bufio.NewScanner(os.Stdin)

		b.CurrURL = URL

		ih := NewInputHandler()

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

			ih.HandleInput(&b, cmd)
		}
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
