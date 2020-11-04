package windowman_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	
	"github.com/Bios-Marcel/cordless/windowman"
	"github.com/Bios-Marcel/cordless/windowman/fakes"
)

var _ = Describe("WindowManager", func() {
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

		It("shows a registered window", func(){
			wm := windowman.NewWindowManager()
			window := fakes.FakeWindow{}

			err := wm.RegisterWindow("a-window", &window)
			Expect(err).ToNot(HaveOccurred())
			err = wm.ShowWindow("a-window")
			Expect(err).ToNot(HaveOccurred())

			shownWindow := wm.GetVisibleWindow()
			Expect(shownWindow).To(BeIdenticalTo(&window))
		})

		It("calls the 'Show' method of the window", func(){
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
