package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func RunCLI(URL string, isURL bool) {
	var b Browser

	if isURL {
		_, body, err := Fetch(URL)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}

		content, links := DoLinks(body)
		b.Stack = []Page{{URL: URL, Content: content, Links: links}}

		fmt.Println(b.RenderOutput())
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
