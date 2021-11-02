package TUI

import (
	tui "github.com/marcusolsson/tui-go"
)

type LoginHandler func(string)

type ChatLogin struct {
	tui.Box
	view     *tui.Box
	username LoginHandler
}

func NewChatLogin() *ChatLogin {

	lable := tui.NewLabel("Please write your user name")
	newUser := tui.NewEntry()
	newUser.SetFocused(true)
	newUser.SetSizePolicy(tui.Expanding, tui.Maximum)

	loginbox := tui.NewVBox(lable, newUser)
	loginbox.SetBorder(true)

	chatview := &ChatLogin{}

	chatview.view = loginbox

	newUser.OnSubmit(func(entry *tui.Entry) {
		if entry.Text() != "" {
			if chatview.username != nil {
				chatview.username(entry.Text())
			}
		}

	})
	return chatview

}
func (v *ChatLogin) Login(handler LoginHandler) {
	v.username = handler
}
