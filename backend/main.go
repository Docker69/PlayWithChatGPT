// package declaration
package main

// import required packages
import (
	"os"

	"backend/chat"
	mylogger "backend/utils"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// main function of the application
func main() {
	// get command line arguments
	args := os.Args

	// load the environment variables
	err := godotenv.Load()
	if err != nil {
		mylogger.Logger.Panicf("Error loading .env file. Err: %s", err)
	}

	// Get the log level from the environment variables
	logLevelStr, exists := os.LookupEnv("LOG_LEVEL")

	// Set the log level in the mylogger package if it exists
	if exists {
		// Parse the log level string into a Logrus Level constant
		logLevelBytes := []byte(logLevelStr)
		err = mylogger.Logger.Level.UnmarshalText(logLevelBytes)
		if err != nil {
			mylogger.Logger.Fatalf("Failed to parse LOG_LEVEL: %v", err)
		}

		mylogger.Logger.Infof("Setting log level to: %s", logLevelStr)
	}

	// extract and save the OpenAI api key from environment variables
	apiKey, exists := os.LookupEnv("OPENAI_API_KEY")

	if !exists {
		mylogger.Logger.Panic("OpenAI API Key not found, panicking!!!")
	}

	var frontend bool = false
	// check command line arguments and compare them
	for _, arg := range args {
		if arg == "frontend" {
			frontend = true
		}
	}

	if frontend {
		//start the server
		runServer()
	} else {
		//generate a random uuid
		uuid := uuid.New().String()

		// start chat via console
		chat.StartConsoleChat(apiKey, uuid)
	}
}
