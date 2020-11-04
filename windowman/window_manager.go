package windowman

import (
	tcell "github.com/gdamore/tcell/v2"

	"fmt"

	"github.com/Bios-Marcel/cordless/config"
	"github.com/Bios-Marcel/cordless/shortcuts"
	"github.com/Bios-Marcel/cordless/tview"
	messagebus "github.com/vardius/message-bus"
)

var (
	windowManagerSingleton WindowManagerInterface = nil
)

type EventHandler func(*tcell.EventKey) *tcell.EventKey

type WindowManagerInterface interface {
	GetMessageBus() messagebus.MessageBus
	GetVisibleWindow() Window
	Dialog(dialog Dialog) error
	ShowWindow(identifier string) error
	RegisterWindow(identifier string, window Window) error
	Run() error

	// FIXME Temporary solution.
	GetUnderlyingApp() *tview.Application
}

type ApplicationControl interface {
	SetFocus(p tview.Primitive) *tview.Application
	SetRoot(root tview.Primitive, fullscreen bool) *tview.Application
	Draw() *tview.Application
	QueueUpdate(f func()) *tview.Application
	QueueUpdateDraw(f func()) *tview.Application
}

type WindowManager struct {
	tviewApp *tview.Application

	windowRegistry map[string]Window
	visibleWindow  Window
	messages       messagebus.MessageBus
}

func GetWindowManager() WindowManagerInterface {
	if windowManagerSingleton == nil {
		windowManagerSingleton = NewWindowManager()
	}

	return windowManagerSingleton
}

func NewWindowManager() WindowManagerInterface {
	return NewWindowManagerWithTViewApp(tview.NewApplication())
}

func NewWindowManagerWithTViewApp(app *tview.Application) WindowManagerInterface {
	wm := &WindowManager{
		tviewApp:       app,
		windowRegistry: make(map[string]Window),
		visibleWindow:  nil,
		messages:       messagebus.New(100),
	}

	wm.tviewApp.MouseEnabled = config.Current.MouseEnabled

	// WindowManagerInterface sets the root input handler.
	// It captures exit application shortcuts, and exits the application,
	// or otherwise allows the event to bubble down.
	wm.tviewApp.SetInputCapture(wm.exitApplicationEventHandler)

	return wm
}

func (wm *WindowManager) GetMessageBus() messagebus.MessageBus {
	return wm.messages
}

func (wm *WindowManager) GetVisibleWindow() Window {
	return wm.visibleWindow
}

func (wm *WindowManager) GetUnderlyingApp() *tview.Application {
	return wm.tviewApp
}

func (wm *WindowManager) ShowWindow(identifier string) error {
	if w, exists := wm.windowRegistry[identifier]; exists {
		return wm.makeWindowVisible(w)
	} else {
		return fmt.Errorf("'%s' is not a registered window", identifier)
	}
}

func (wm *WindowManager) Dialog(dialog Dialog) error {
	panic("not implemented")
}

func (wm *WindowManager) RegisterWindow(identifier string, window Window) error {
	if _, exists := wm.windowRegistry[identifier]; exists {
		return fmt.Errorf("another window is already registered under the name '%s'", identifier)
	} else {
		wm.windowRegistry[identifier] = window
		window.OnRegister(wm.messages)
		return nil
	}
}

func (wm *WindowManager) Run() error {
	if wm.visibleWindow == nil {
		return fmt.Errorf("no window has been made visible before running the application")
	}
	return wm.tviewApp.Run()
}

func createSetFocusCallback(app *tview.Application) func(tview.Primitive) error {
	return func(primitive tview.Primitive) error {
		app.SetFocus(primitive)
		return nil
	}
}

func (wm *WindowManager) exitApplicationEventHandler(event *tcell.EventKey) *tcell.EventKey {
	if shortcuts.ExitApplication.Equals(event) {
		wm.tviewApp.Stop()
		return nil
	}
	return event
}

func (wm *WindowManager) makeWindowVisible(window Window) error {
	err := window.Show(wm.tviewApp)

	if err != nil {
		return err
	}

	passThroughHandler := stackEventHandler(
		wm.exitApplicationEventHandler,
		func(evt *tcell.EventKey) *tcell.EventKey {
			return window.HandleKeyEvent(evt)
		},
	)

	wm.tviewApp.SetInputCapture(passThroughHandler)
	wm.visibleWindow = window
	return nil
}

func stackEventHandler(root EventHandler, new EventHandler) EventHandler {
	return func(event *tcell.EventKey) *tcell.EventKey {
		rootEvt := root(event)

		if rootEvt == nil {
			return nil
		}

		return new(rootEvt)
	}
}
