package updatedialog

import (
	"fmt"
	"sync"

	"github.com/Bios-Marcel/cordless/config"
	"github.com/Bios-Marcel/cordless/tview"
	"github.com/Bios-Marcel/cordless/version"
	"github.com/Bios-Marcel/cordless/windowman"
	tcell "github.com/gdamore/tcell/v2"
	messagebus "github.com/vardius/message-bus"
)

var (
	buttonOk                            = "Thanks for the info"
	buttonDontRemindAgainForThisVersion = fmt.Sprintf("Skip reminders for %s", version.GetLatestRemoteVersion())
	buttonNeverRemindMeAgain            = "Never remind me again"
)

type UpdateDialog struct {
	modal         *tview.Modal
	configuration config.Config
}

func NewUpdateDialog(configuration config.Config) UpdateDialog {
	return UpdateDialog{
		modal:         createModal(),
		configuration: configuration,
	}
}

func (ud *UpdateDialog) HandleKeyEvent(evt *tcell.EventKey) *tcell.EventKey {
	return evt
}

// OnRegister is called when the window is registered in a window manager.
// The messages parameter allows the window to subscribe to messages on the message bus,
// or to publish its own
func (ud *UpdateDialog) OnRegister(messages messagebus.MessageBus) {}

func (ud *UpdateDialog) Open(close windowman.DialogCloser) error {
	var wg sync.WaitGroup
	wg.Add(1)
	ud.modal.SetDoneFunc(func(index int, label string) {
		if label == buttonDontRemindAgainForThisVersion {
			ud.configuration.DontShowUpdateNotificationFor = version.GetLatestRemoteVersion()
			config.PersistConfig()
		} else if label == buttonNeverRemindMeAgain {
			ud.configuration.ShowUpdateNotifications = false
			config.PersistConfig()
		}

		wg.Done()
	})
	wg.Wait()
	close()
	return nil
}

// Show resets the window state and returns the tview.Primitive that the caller should show.
// The setFocus argument is used by the Window to change the focus
func (ud *UpdateDialog) Show(appCtl windowman.ApplicationControl) error {
	appCtl.SetRoot(ud.modal, true)
	return nil
}

func createModal() *tview.Modal {
	dialog := tview.NewModal()
	dialog.SetText(
		fmt.Sprintf(
			"Version %s of cordless is available!\nYou are currently running version %s.\n\nUpdates have to be installed manually or via your package manager.\n\nThe snap package manager isn't supported by cordless anymore!",
			version.GetLatestRemoteVersion(),
			version.Version,
		),
	)
	dialog.AddButtons([]string{
		buttonOk,
		buttonDontRemindAgainForThisVersion,
		buttonNeverRemindMeAgain,
	})

	return dialog
}
