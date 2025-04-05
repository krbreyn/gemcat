package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
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
			fmt.Println(links)
			b.Stack = []Page{{URL: URL, Content: content, Links: links}}
			fmt.Println(ColorGemtext(content, links))
		} else {
			fmt.Println("welcome to gemcat\ntype help in the future for commands")
		}

		scanner := bufio.NewScanner(os.Stdin)

		url := URL

		gotoURL := func(url string) error {
			host, path := getHostPath(url)
			_, body, err := Fetch(host, path)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return err
			}

			content, links := DoLinks(body)
			if len(b.Stack) != 0 {
				b.Pos++
			}
			if b.Pos == len(b.Stack) {
				b.Stack = append(b.Stack, Page{
					URL:     url,
					Content: content,
					Links:   links,
				})
			} else {
				b.Stack = append(b.Stack[:b.Pos], Page{
					URL:     url,
					Content: content,
					Links:   links,
				})
			}
			return nil
		}

		for {
			fmt.Print("> ")
			scanner.Scan()
			text := scanner.Text()

			if err := scanner.Err(); err != nil {
				panic(err)
			}

			cmd := strings.Fields(text)

			switch cmd[0] {
			case "goto":
				if len(cmd) == 1 {
					fmt.Println("must include a link")
					continue
				}
				url = strings.TrimPrefix(cmd[1], "gemini://")
				err := gotoURL(url)
				if err != nil {
					fmt.Printf("err: %v", err)
					continue
				}
				p := b.Stack[b.Pos]
				fmt.Println(ColorGemtext(p.Content, p.Links))

			case "l", "link":
				if len(cmd) == 1 {
					fmt.Println("must include link number")
					continue
				}

				i, err := strconv.Atoi(cmd[1])
				if err != nil {
					fmt.Println("links are numbers")
				}

				fmt.Println(b.Stack[b.Pos].Links[i])

			case "g":
				if len(cmd) == 1 {
					fmt.Println("must include link number")
					continue
				}

				i, err := strconv.Atoi(cmd[1])
				if err != nil {
					fmt.Println("links are numbers")
				}

				p := b.Stack[b.Pos]
				if i >= len(p.Links) {
					continue
				}

				link := p.Links[i].URL
				p.Links[i].Visited = true

				if strings.HasPrefix(link, "gemini://") {
					link = strings.TrimPrefix(link, "gemini://")
				} else {
					link = strings.TrimPrefix(link, "/")
					url = strings.TrimSuffix(url, "/")
					link = url + "/" + link
				}

				url = link
				err = gotoURL(url)
				if err != nil {
					fmt.Printf("err: %v", err)
					continue
				}
				p = b.Stack[b.Pos]
				fmt.Println(ColorGemtext(p.Content, p.Links))

			case "links":
				for _, s := range b.Stack[b.Pos].Links {
					fmt.Printf("%d %s\n", s.No, s.URL)
				}

			case "stack":
				for i, p := range b.Stack {
					if i == b.Pos {
						fmt.Print("-> ")
					}
					fmt.Printf("%d %s\n", i, p.URL)
				}

			case "history":

			case "b", "back":
				b.Pos--
				url = b.Stack[b.Pos].URL
				p := b.Stack[b.Pos]
				fmt.Println(ColorGemtext(p.Content, p.Links))

			case "f", "forward":
				b.Pos++
				url = b.Stack[b.Pos].URL
				p := b.Stack[b.Pos]
				fmt.Println(ColorGemtext(p.Content, p.Links))

			case "less":

			case "exit", "quit":
				os.Exit(0)
			}
		}
	}
}

type Browser struct {
	Pos     int
	Stack   []Page
	History []string
}

type Page struct {
	URL     string
	Content string
	Links   []Link
}

type Link struct {
	No      int
	URL     string
	Visited bool
}

func DoLinks(gemtxt string) (output string, links []Link) {
	scanner := bufio.NewScanner(strings.NewReader(gemtxt))

	var b strings.Builder
	var i int

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "=>") {
			split := strings.Fields(line)
			b.WriteString("=> " + fmt.Sprintf("[%d] ", i))
			for _, s := range split[2:] {
				b.WriteString(s + " ")
			}
			b.WriteString("\n")
			links = append(links, Link{No: i, URL: split[1]})
			i++
		} else {
			b.WriteString(line + "\n")
		}
	}

	return b.String(), links
}

func ColorGemtext(gemtxt string, links []Link) string {
	scanner := bufio.NewScanner(strings.NewReader(gemtxt))

	var b strings.Builder
	var isInPreformattedBlock bool
	var i int

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "```") {
			isInPreformattedBlock = !isInPreformattedBlock
			// Italics
			b.WriteString("\033[3m]" + strings.TrimPrefix(line, "```") + "\033[23m")
			continue
		}

		if isInPreformattedBlock {
			b.WriteString(line + "\n")
			continue
		}

		if strings.HasPrefix(line, "#") {
			// Blue
			b.WriteString("\033[34m" + line + "\033[39m\n")
		} else if strings.HasPrefix(line, "*") {
			// Green
			b.WriteString("\033[32m" + line + "\033[39m\n")
		} else if strings.HasPrefix(line, ">") {
			// Yellow
			b.WriteString("\033[33m" + line + "\033[39m\n")
		} else if strings.HasPrefix(line, "=>") {
			// Cyan
			if links == nil {
				b.WriteString("\033[36m" + line + "\033[39m\n")
				continue
			}
			split := strings.Fields(line)
			fmt.Println(split)
			if links[i].Visited {
				b.WriteString("\033[35m" + line + "\033[39m\n")
			} else {
				b.WriteString("\033[36m" + line + "\033[39m\n")
			}
			i++
		} else {
			b.WriteString(line + "\n")
		}
	}

	return b.String()
}

func Fetch(host, path string) (status, body string, err error) {
ifRedirect:
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			return handleTOFU(rawCerts, host)
		},
	}

	fmt.Printf("Connecting to gemini://%s/%s\r\n", host, path)
	addr := net.JoinHostPort(host, "1965")
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return "", "", fmt.Errorf("TLS connection failed: %v", err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "gemini://%s/%s\r\n", host, path)

	reader := bufio.NewReader(conn)
	status, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal("Failed to read response:", err)
	}

	status_no, err := strconv.Atoi(strings.Fields(status)[0])
	if err != nil {
		return "", "", fmt.Errorf("weird status err: %v", err)
	}
	if status_no == 30 || status_no == 31 {
		new_url := strings.Fields(status)[1]
		new_url = strings.TrimPrefix(new_url, "gemini://")
		host, path = getHostPath(new_url)
		fmt.Printf("Redirect: gemini://%s/%s\r\n", host, path)
		goto ifRedirect
	}
	if status_no != 20 {
		return "", "", fmt.Errorf("status was not 20 but was %d, status: %s", status_no, status)
	}

	var b strings.Builder
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		b.WriteString(line)
	}

	return status, b.String(), nil
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
