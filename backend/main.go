// package declaration
package main

// import required packages
import (
	"os"

	"backend/ai"
	"backend/router"
	"backend/utils"

	"github.com/joho/godotenv"
)

// init main package
func init() {

	// load the environment variables
	err := godotenv.Load()
	if err != nil {
		utils.Logger.Infof("Error loading .env file, will use default values. Err: %s", err)
	}

	// Get the log level from the environment variables
	logLevelStr, exists := os.LookupEnv("LOG_LEVEL")

	// Set the log level in the utils package if it exists
	if exists {
		// Parse the log level string into a Logrus Level constant
		logLevelBytes := []byte(logLevelStr)
		err = utils.Logger.Level.UnmarshalText(logLevelBytes)
		if err != nil {
			utils.Logger.Fatalf("Failed to parse LOG_LEVEL: %v", err)
		}

		utils.Logger.Infof("Setting log level to: %s", logLevelStr)
	}
	//log main init done
	utils.Logger.Info("main package initialized")
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
		ai.StartConsoleChat()
	}
}
