package tui

import (
	"log"

	"github.com/marcusolsson/tui-go"
)

var chatview *ChatView
var Ui tui.UI

// needs to get a client varible that can be used
func StartChatview() {
	chatview := NewChatView()
	chatlogin := NewChatLogin()

	ui, err := tui.New(chatlogin.view)
	if err != nil {
		log.Fatal(err)
	}
	exit := func() {
		// send user logout message to server
		ui.Quit()
	}
	ui.SetKeybinding("Esc", exit)
	ui.SetKeybinding("Ctrl+c", exit)
	Ui = ui
	chatlogin.Login(func(username string) {
		//username is the new user joining the chat. call
		//the server with the name

		ui.SetWidget(chatview.view)
	})

	chatview.SendMessage(func(message string) {
		chatview.ReciveMessage(message)
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
func ReciveMessage(msg string) {
	Ui.Update(func() { chatview.ReciveMessage(msg) })
}
