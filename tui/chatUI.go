package tui

import (
	tui "github.com/marcusolsson/tui-go"
)

type SendMessageHandler func(string)

type ChatView struct {
	tui.Box
	view           *tui.Box
	messageHistory *tui.Box
	newMessage     SendMessageHandler
}

func NewChatView() *ChatView {
	chatview := &ChatView{}

	// create the message history field
	historyfield := tui.NewVBox()
	chatview.messageHistory = historyfield
	historyScroll := tui.NewScrollArea(historyfield)
	historyScroll.SetAutoscrollToBottom(true)

	historybox := tui.NewVBox(historyScroll)
	historybox.SetBorder(true)

	// create the chat input field

	chatinput := tui.NewEntry()
	chatinput.SetFocused(true)
	chatinput.SetSizePolicy(tui.Expanding, tui.Maximum)

	chatinput.OnSubmit(func(entry *tui.Entry) {
		if entry.Text() != "" {
			if chatview.newMessage != nil {
				chatview.newMessage(entry.Text())
			}
			entry.SetText("")
		}

	})

	chatbox := tui.NewVBox(chatinput)
	chatbox.SetBorder(true)
	chatbox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chatview.view = tui.NewVBox(historybox, chatbox)
	return chatview
}

func (v *ChatView) SendMessage(handler SendMessageHandler) {
	v.newMessage = handler
}

func (v *ChatView) ReciveMessage(message string) {

	v.messageHistory.Append(tui.NewLabel(message))

}
