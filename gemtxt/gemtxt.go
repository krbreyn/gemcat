package gemtxt

import (
	"bufio"
	"fmt"
	"strings"
)

func ColorPlain(body string) string {
	scanner := bufio.NewScanner(strings.NewReader(body))

	var b strings.Builder
	var isInPreformattedBlock bool

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "```") {
			isInPreformattedBlock = !isInPreformattedBlock
			b.WriteString("\033[3m]" + strings.TrimPrefix(line, "```") + "\033[23m\n") // Italics
			continue
		}

		if isInPreformattedBlock {
			b.WriteString(line + "\n")
			continue
		}

		if strings.HasPrefix(line, "#") {
			b.WriteString("\033[34m" + line + "\033[39m\n") // Blue
		} else if strings.HasPrefix(line, "*") {
			b.WriteString("\033[32m" + line + "\033[39m\n") // Green
		} else if strings.HasPrefix(line, ">") {
			b.WriteString("\033[37m" + line + "\033[39m\n") // White
		} else if strings.HasPrefix(line, "=>") {
			b.WriteString("\033[36m" + line + "\033[39m\n") // Cyan
		} else {
			b.WriteString(line + "\n")
		}
	}

	return b.String()
}

func ColorWithLinkNosAndNoURLs(body string) string {
	var li int
	lf := func(line string) string {
		var b strings.Builder
		split := strings.Fields(line)
		b.WriteString("\033[36m" + split[0] + " " + fmt.Sprintf("[%d] ", li))
		for _, str := range split[2:] {
			b.WriteString(str + " ")
		}
		b.WriteString("\033[39m\n")
		li++
		return b.String()
	}
	return ColorLinkFunc(body, lf)
}

func ColorLinkFunc(body string, link_func func(line string) string) string {
	scanner := bufio.NewScanner(strings.NewReader(body))

	var b strings.Builder
	var isInPreformattedBlock bool

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "```") {
			isInPreformattedBlock = !isInPreformattedBlock
			b.WriteString("\033[3m]" + strings.TrimPrefix(line, "```") + "\033[23m\n") // Italics
			continue
		}

		if isInPreformattedBlock {
			b.WriteString(line + "\n")
			continue
		}

		if strings.HasPrefix(line, "#") {
			b.WriteString("\033[34m" + line + "\033[39m\n") // Blue
		} else if strings.HasPrefix(line, "*") {
			b.WriteString("\033[32m" + line + "\033[39m\n") // Green
		} else if strings.HasPrefix(line, ">") {
			b.WriteString("\033[37m" + line + "\033[39m\n") // White
		} else if strings.HasPrefix(line, "=>") {
			b.WriteString(fmt.Sprintf("%s\n", link_func(line)))
		} else {
			b.WriteString(line + "\n")
		}
	}

	return b.String()
}
