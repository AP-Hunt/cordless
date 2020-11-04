package windowman_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/Bios-Marcel/cordless/windowman"
	"github.com/Bios-Marcel/cordless/windowman/fakes"
)

var _ = Describe("WindowManager", func() {
	Describe("Dialog", func() {
		It("calls the Open method of the dialog", func() {
			wm := windowman.NewWindowManager()
			window := fakes.FakeWindow{}
			dialog := fakes.FakeDialog{}

			// Ensure the dialog closes by automatically closing it
			dialog.OpenCalls(func(closer windowman.DialogCloser) error {
				closer()
				return nil
			})

			wm.RegisterWindow("root", &window)
			wm.ShowWindow("root")

			wm.Dialog(&dialog)
			Expect(dialog.OpenCallCount()).To(Equal(1))
		})

		It("shows the previous window when the dialog has been closed", func() {
			wm := windowman.NewWindowManager()
			window := fakes.FakeWindow{}
			dialog := fakes.FakeDialog{}
			var dialogCloser windowman.DialogCloser

			// Ensure the dialog closes by automatically closing it
			dialog.OpenCalls(func(closer windowman.DialogCloser) error {
				dialogCloser = closer
				return nil
			})

			wm.RegisterWindow("root", &window)
			wm.ShowWindow("root")

			Expect(wm.GetVisibleWindow()).To(BeIdenticalTo(&window), "the root window was not visible before the dialog was shown")
			wm.Dialog(&dialog)
			Expect(wm.GetVisibleWindow()).To(BeIdenticalTo(&dialog), "the dialog was not the visible window whilst it was open")
			dialogCloser()
			Expect(wm.GetVisibleWindow()).To(BeIdenticalTo(&window), "the root window was not visible after the dialog was closed")
		})
	})

	Describe("RegisterWindow", func() {
		It("returns an error if an identifier has already been used", func() {
			wm := windowman.NewWindowManager()
			window1 := fakes.FakeWindow{}
			window2 := fakes.FakeWindow{}

			err := wm.RegisterWindow("id", &window1)
			Expect(err).ToNot(HaveOccurred())

			err = wm.RegisterWindow("id", &window2)
			Expect(err).To(HaveOccurred())
		})

		It("calls the 'OnRegister' method of the window", func() {
			wm := windowman.NewWindowManager()
			window := fakes.FakeWindow{}

			err := wm.RegisterWindow("id", &window)
			Expect(err).ToNot(HaveOccurred())

			Expect(window.OnRegisterCallCount()).To(Equal(1))
		})
	})

	Describe("Run", func() {
		It("returns an error if no window has been shown", func() {
			wm := windowman.NewWindowManager()

			err := wm.Run()
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("ShowWindow", func() {
		It("cannot show an unregistered window", func() {
			wm := windowman.NewWindowManager()

			err := wm.ShowWindow("unregistered-window")
			Expect(err).To(HaveOccurred())
		})

		It("shows a registered window", func() {
			wm := windowman.NewWindowManager()
			window := fakes.FakeWindow{}

			err := wm.RegisterWindow("a-window", &window)
			Expect(err).ToNot(HaveOccurred())
			err = wm.ShowWindow("a-window")
			Expect(err).ToNot(HaveOccurred())

			shownWindow := wm.GetVisibleWindow()
			Expect(shownWindow).To(BeIdenticalTo(&window))
		})

		It("calls the 'Show' method of the window", func() {
			wm := windowman.NewWindowManager()
			window := fakes.FakeWindow{}

			err := wm.RegisterWindow("a-window", &window)
			Expect(err).ToNot(HaveOccurred())
			err = wm.ShowWindow("a-window")
			Expect(err).ToNot(HaveOccurred())

			Expect(window.ShowCallCount()).To(Equal(1))
		})
	})
})
