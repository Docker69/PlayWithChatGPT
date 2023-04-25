// package declaration
package main

// import required packages
import (
	"os"

	"backend/chat"
	router "backend/router"
	mylogger "backend/utils"

	"github.com/joho/godotenv"
)

// init main package
func init() {

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
	//log main init done
	mylogger.Logger.Info("main package initialized")
}

// main function of the application
func main() {
	// get command line arguments
	args := os.Args

	var frontend bool = false
	// check command line arguments and compare them
	for _, arg := range args {
		if arg == "frontend" {
			frontend = true
		}
	}

	if frontend {
		//start the server
		router.RunServer()
	} else {
		// start chat via console
		chat.StartConsoleChat()
	}
}
