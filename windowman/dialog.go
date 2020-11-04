package windowman

type DialogCloser func() error

//counterfeiter:generate -o ./fakes/FakeDialog.go . Dialog
type Dialog interface {
	Window
	Open(close DialogCloser) error
}
