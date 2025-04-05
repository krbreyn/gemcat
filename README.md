# Gemcat
Gemcat is a terminal-based, CLI-focused browser for the Gemini protocol, with Gopher support and RSS feed inclusion aimed for the future, with the goal being to act as a complete mult-tool for exploring and interacting with the non-HTTP web.

Firstly, you can use it to fetch and print Gemini content, such as:

`gemcat geminiprotocol.net/`

Secondly, gemcat has a interactive shell-styled style client. Launch it by using:

`gemcat -i`

with, optionally, a URL following as the initial request. From there you can type `help` to see the available commands that let you navigate and interact with the browser.

TODO: interactive mode demonstration gif (after i implement bookmarks) (also todo: cache listing command)

Thirdly, in the future, gemcat will also have a TUI powered by tcell that will have keyboard-navigation and full mouse support.
