package version

import messagebus "github.com/vardius/message-bus"

const MessageUpdateAvailable = "update.available"

func EmitVersionCheckingToMessageBus(bus messagebus.MessageBus, updatesAvailableChannel chan bool) {
	go func() {
		availability := <-updatesAvailableChannel
		bus.Publish(MessageUpdateAvailable, availability)
	}()
}
