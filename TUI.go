package TUI

import (
	"log"

	"github.com/marcusolsson/tui-go"
)

// needs to get a client varible that can be used
func StartChatview() {
	chatview := NewChatView()
	chatlogin := NewChatLogin()

	ui, err := tui.New(chatlogin)
	if err != nil {
		// make faltal error log or something like that
	}
	exit := func() {
		// send user logout message to server
		ui.Quit()
	}
	ui.SetKeybinding("Esc", exit)
	ui.SetKeybinding("Ctrl+c", exit)

	chatlogin.Login(func(username string) {
		//username is the new user joining the chat. call
		//the server with the name

		ui.SetWidget(chatview)
	})

	chatview.SendMessage(func(message string) {
		//send user message to server
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
