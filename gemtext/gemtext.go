package gemtext

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/krbreyn/gemcat"
)

func DoLinks(gemtxt string, wasLinkVisited, isLinkBookmarked func(url string) bool) (output string, links []gemcat.Link) {
	scanner := bufio.NewScanner(strings.NewReader(gemtxt))

	var b strings.Builder
	var i int

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "=>") {
			split := strings.Fields(line)
			url, text := split[1], split[2:] // [0] = "=>"
			fmt.Println(split[0], split[1], split[2])

			b.WriteString("=> " + fmt.Sprintf("[%d] ", i))
			if gemcat.IsGeminiLink(url) {
				b.WriteString("(gemini) ")
			} else {
				b.WriteString("(unsupported) ")
			}
			for _, s := range text {
				b.WriteString(s + " ")
			}
			b.WriteString("\n")

			links = append(links, gemcat.Link{No: i, URL: url, Visited: wasLinkVisited(url), Bookmarked: isLinkBookmarked(url)})
			i++
		} else {
			b.WriteString(line + "\n")
		}
	}

	return b.String(), links
}

func ColorGemtext(gemtxt string, links []gemcat.Link) string {
	scanner := bufio.NewScanner(strings.NewReader(gemtxt))

	var b strings.Builder
	// var isInPreformattedBlock bool
	var i int

	for scanner.Scan() {
		line := scanner.Text()

		// if strings.HasPrefix(line, "```") {
		// 	isInPreformattedBlock = !isInPreformattedBlock
		// 	b.WriteString("\033[3m]" + strings.TrimPrefix(line, "```") + "\033[23m\n") // Italics
		// 	continue
		// }

		// if isInPreformattedBlock {
		// 	b.WriteString(line + "\n")
		// 	continue
		// }

		if strings.HasPrefix(line, "#") {
			b.WriteString("\033[34m" + line + "\033[39m\n") // Blue
		} else if strings.HasPrefix(line, "*") {
			b.WriteString("\033[32m" + line + "\033[39m\n") // Green
		} else if strings.HasPrefix(line, ">") {
			b.WriteString("\033[37m" + line + "\033[39m\n") // White
		} else if strings.HasPrefix(line, "=>") {
			if links == nil {
				b.WriteString("\033[36m" + line + "\033[39m\n") // Cyan
				continue
			}

			if links[i].Bookmarked {
				b.WriteString("\033[33m" + line + "\033[39m\n") // Yellow
			} else if links[i].Visited {
				b.WriteString("\033[35m" + line + "\033[39m\n") // Magenta
			} else {
				b.WriteString("\033[36m" + line + "\033[39m\n") // Cyan
			}
			i++
		} else {
			b.WriteString(line + "\n")
		}
	}

	return b.String()
}
