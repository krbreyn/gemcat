package browser

import "fmt"

type InputHandler struct {
	cmd_map  map[string]BrowserCmd
	cmds     []BrowserCmd
	help_cmd HelpCmd
}

func NewInputHandler() InputHandler {
	cmd_map, cmds := makeCmdMap()
	return InputHandler{cmd_map, cmds, HelpCmd{}}
}

func (ih *InputHandler) HandleInput(b *Browser, cmd []string) {
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
			if cmd, ok := ih.cmd_map[args[0]]; ok {
				ih.help_cmd.Do([]BrowserCmd{cmd})
			} else {
				fmt.Printf("cmd %s does not exist\n", args[0])
			}
			return
		}
		ih.help_cmd.Do(ih.cmds)
		return
	}

	if cmd, ok := ih.cmd_map[opt]; ok {
		if len(args) != 0 && args[0] == "help" {
			ih.help_cmd.Do([]BrowserCmd{cmd})
			return
		}
		err := cmd.Do(b, args)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	} else {
		fmt.Printf("error: cmd not recognized: %s\n", opt)
	}
}
