// package declaration
package main

// import required packages
import (
	"log"
	"os"

	"PlayWithChatGPT/chat"

	"github.com/joho/godotenv"
)

// main function of the application
func main() {
	// load the environment variables
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
	}

	// extract and save the OpenAI api key from environment variables
	apiKey := os.Getenv("OPENAI_API_KEY")

	// start the chat
	chat.StartChat(apiKey)
}
