package login

import (
	"github.com/Bios-Marcel/cordless/authentication"
	"github.com/Bios-Marcel/cordless/windowman"
	"github.com/Bios-Marcel/discordgo"
	tcell "github.com/gdamore/tcell/v2"
	messagebus "github.com/vardius/message-bus"
)

type LoginWindow struct {
	windowman.Window

	LoginWindowComponent *Login
}

func NewLoginWindow(configDirectory string, authenticator authentication.Authenticator) LoginWindow {
	return LoginWindow{
		LoginWindowComponent: NewLogin(configDirectory, authenticator),
	}
}

// Show resets the window state and returns the tview.Primitive that the caller should show.
// The setFocus argument is used by the Window to change the focus
func (lw *LoginWindow) Show(appCtl windowman.ApplicationControl) error {
	lw.LoginWindowComponent.SetAppControl(appCtl)
	appCtl.SetRoot(lw.LoginWindowComponent, true)
	return nil
}

func (lw *LoginWindow) HandleKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	return event
}

func (lw *LoginWindow) OnRegister(messages messagebus.MessageBus) {

	// Publish a login success message when the underlying component calls back to
	// say that authentication was successful
	lw.LoginWindowComponent.OnLoginSuccess(func(session discordgo.Session, ready discordgo.Ready) {
		messages.Publish(messagebus.TopicDiscordLoginSuccess, session, ready)
	})
}
