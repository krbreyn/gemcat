package interactive

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/krbreyn/gemcat/browser"
	"github.com/krbreyn/gemcat/gemtxt"
	"github.com/krbreyn/gemcat/shell"
	"github.com/muesli/reflow/wordwrap"
	"golang.org/x/term"
)

func RunCLI(u *url.URL, isURL bool, loadLast bool) {
	b := &browser.Browser{}
	sh := shell.NewShell(CLIOutput{})
	scanner := bufio.NewScanner(os.Stdin)

	if isURL && u.String() != b.S.CurrURL() {
		err := shell.GotoCmd{}.Do(b, sh.Out, []string{u.String()})
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	fmt.Println("welcome to gemcat\ntype 'help' to see the available commands!")

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "Error reading input:", err)
			}
			break
		}

		cmd := strings.Fields(scanner.Text())
		sh.HandleInput(b, cmd)
	}
	os.Exit(0)
}

type CLIOutput struct{}

func (o CLIOutput) RecvError(err error) {
	fmt.Fprintln(os.Stderr, err)
}

func (o CLIOutput) RecvMsg(msg string) {
	fmt.Println(msg)
}

func (o CLIOutput) RecvPage(page browser.Page) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80
	}
	fmt.Println(wordwrap.String(gemtxt.ColorWithLinkNosAndNoURLs(page.Content), width))
}

func (o CLIOutput) ShowHelp(help []shell.HelpInfo) {
	for _, cmd := range help {
		cap := len(cmd.Words) - 1
		for i, w := range cmd.Words {
			fmt.Print(w)
			if i != cap {
				fmt.Print(", ")
			}
		}
		fmt.Printf("\n\t%s\n", cmd.Desc)
	}
}

// func (o CLIOutput) GetInput() string {

// }

// func (o CLIOutput) GetCert() *x509.Certificate {

// }
