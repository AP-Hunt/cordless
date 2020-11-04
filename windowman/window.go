package windowman

import (
	tcell "github.com/gdamore/tcell/v2"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o ./fakes/FakeWindow.go . Window
type Window interface {
	// Show resets the window state and returns the tview.Primitive that the caller should show.
	// The setFocus argument is used by the Window to change the focus
	Show(ApplicationControl) error
	HandleKeyEvent(*tcell.EventKey) *tcell.EventKey

	// OnRegister is called when the window is registered in a window manager
	OnRegister()
}
