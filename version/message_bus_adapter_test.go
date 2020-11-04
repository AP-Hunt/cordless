package version_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	messagebus "github.com/vardius/message-bus"

	"github.com/Bios-Marcel/cordless/version"
)

var _ = Describe("MessageBusAdapter", func() {
	Describe("EmitVersionCheckingToMessageBus", func() {
		It("emits a message when the channel receives a true value", func() {
			bus := messagebus.New(10)
			channel := make(chan bool, 2)
			eventHasBeenPublished := false

			bus.Subscribe(version.MessageUpdateAvailable, func(_ bool) {
				eventHasBeenPublished = true
			})

			version.EmitVersionCheckingToMessageBus(bus, channel)
			channel <- true

			Eventually(func() bool {
				return eventHasBeenPublished
			}, 3*time.Second).Should(BeTrue())
		})
	})
})
