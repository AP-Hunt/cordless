package messagebus_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/Bios-Marcel/cordless/messagebus"
)

var _ = Describe("MessageBus", func() {
	Describe("SubscribeOnce", func() {
		It("invokes a handler once and then unsubscribes it", func() {
			bus := NewMessageBus()
			callCounter := 0
			callback := func() {
				callCounter++
			}

			bus.SubscribeOnce("topic", callback)

			bus.Publish("topic")
			Eventually(func() int {
				return callCounter
			}).Should(Equal(1), "callback counter should have been 1 after the first message")

			bus.Publish("topic")
			Consistently(func() int {
				return callCounter
			}).ShouldNot(BeNumerically(">", 1), "callback should not have been invoked a second time")
		})

		It("will handle functions with arguments and return types", func() {
			bus := NewMessageBus()
			invoked := false

			callback := func(_ string, _ int, _ bool) (int, error) {
				invoked = true

				return 1, nil
			}

			bus.SubscribeOnce("topic", callback)

			bus.Publish("topic", "string", 100, false)

			Eventually(func() bool {
				return invoked
			}).Should(BeTrue(), "callback should eventually back been invoked")
		})
	})
})
