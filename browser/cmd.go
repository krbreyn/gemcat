package browser

type BrowserCmd interface {
	Do(b *Browser, args []string) error
	Help() (words []string, desc string)
}

func makeCmdMap() (map[string]BrowserCmd, []BrowserCmd) {
	cm := make(map[string]BrowserCmd)
	cmds := []BrowserCmd{
		GotoCmd{},
		BackCmd{},
		ForwardCmd{},
		LinkCmd{},
		LinksCmd{},
		LinkGotoCmd{},
		StackCmd{},
		StackGotoCmd{},
		StackCloseCmd{},
		StackEmptyCmd{},
		HistoryCmd{},
		HistoryGotoCmd{},
		BookmarkGotoCmd{},
		BookmarkListCmd{},
		BookmarkAddCurrentCmd{},
		BookmarkAddLinkCmd{},
		BookmarkRemoveCmd{},
		BookmarkRemoveCurrentCmd{},
		LessCmd{},
		ReprintCmd{},
		ExitCmd{},
	}

	for _, cmd := range cmds {
		words, _ := cmd.Help()
		for _, w := range words {
			cm[w] = cmd
		}
	}

	return cm, cmds
}
