package chat

// import gin framework
import (
	mylogger "backend/utils"
)

// chat function is the main function of the chat package
func StartChat(apiKey string, uuid string) {
	mylogger.Logger.WithField("UUID", uuid).Info("New chat started!")

	mylogger.Logger.WithField("UUID", uuid).Info("Chat Ended!")
}
