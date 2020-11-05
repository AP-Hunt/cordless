package tview

import "github.com/gdamore/tcell/v2"

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o ./fakes/FakeApplication.go . ApplicationInterface
type ApplicationInterface interface {
	SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) *Application
	GetInputCapture() func(event *tcell.EventKey) *tcell.EventKey
	SetScreen(screen tcell.Screen) *Application
	Run() error
	GetComponentAt(x, y int) *Primitive
	Stop()
	Suspend(f func()) bool
	Draw() *Application
	ForceDraw() *Application
	SetBeforeDrawFunc(handler func(screen tcell.Screen) bool) *Application
	GetBeforeDrawFunc() func(screen tcell.Screen) bool
	SetAfterDrawFunc(handler func(screen tcell.Screen)) *Application
	GetAfterDrawFunc() func(screen tcell.Screen)
	SetRoot(root Primitive, fullscreen bool) *Application
	GetRoot() Primitive
	ResizeToFullScreen(p Primitive) *Application
	SetFocus(p Primitive) *Application
	GetFocus() Primitive
	QueueUpdate(f func()) *Application
	QueueUpdateDraw(f func()) *Application
	QueueEvent(event tcell.Event) *Application
}
