package shell

import (
	"fmt"

	"github.com/krbreyn/gemcat/browser"
)

type ShellOut interface {
	RecvMsg(msg string)
	RecvPage(page browser.Page)
	ShowHelp(help []HelpInfo)
	// GetInput() string
	// GetCert() *x509.Certificate
}

type HelpInfo struct {
	Words []string
	Desc  string
}

type Shell struct {
	Out     ShellOut
	cmd_map map[string]ShellCmd
	help    []HelpInfo
}

type ShellCmd interface {
	Do(b *browser.Browser, out ShellOut, args []string) error
	Help() HelpInfo
}

func NewShell(out ShellOut) Shell {
	cmd_map, help := makeCmdMap()
	return Shell{
		Out:     out,
		cmd_map: cmd_map,
		help:    help,
	}
}

func (sh *Shell) HandleInput(b *browser.Browser, cmd []string) {
	if len(cmd) == 0 {
		return
	}

	opt := cmd[0]
	var args []string
	if len(cmd) > 1 {
		args = cmd[1:]
	}

	if opt == "help" {
		if len(args) != 0 {
			if cmd, ok := sh.cmd_map[args[0]]; ok {
				sh.Out.ShowHelp([]HelpInfo{cmd.Help()})
			} else {
				fmt.Printf("cmd %s does not exist\n", args[0])
			}
			return
		}
		sh.Out.ShowHelp(sh.help)
		return
	}

	if cmd, ok := sh.cmd_map[opt]; ok {
		err := cmd.Do(b, sh.Out, args)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	} else {
		fmt.Printf("error: cmd not recognized: '%s'\n", opt)
	}
}

func makeCmdMap() (map[string]ShellCmd, []HelpInfo) {
	cm := make(map[string]ShellCmd)
	cmds := []ShellCmd{
		ExitCmd{},
		TestCmd{},

		GotoCmd{},
		ForwardCmd{},
		BackCmd{},

		LinkCmd{},
		LinksCmd{},
		LinkCurrentCmd{},
		LinkGotoCmd{},

		StackCmd{},
		StackPosCmd{},
		StackRmCmd{},
		StackCloseCmd{},
		StackCompressCmd{},
		StackEmptyCmd{},
		StackGotoCmd{},

		HistoryCmd{},
		HistoryRmCmd{},
		HistoryClearAllCmd{},
		HistoryGotoCmd{},
	}
	var help []HelpInfo

	for _, cmd := range cmds {
		hi := cmd.Help()
		for _, w := range hi.Words {
			cm[w] = cmd
		}
		help = append(help, hi)
	}

	return cm, help
}
