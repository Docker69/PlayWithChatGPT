// package declaration
package main

// import required packages
import (
	"flag"
	"fmt"
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
		utils.Logger.Infof(".env file not found, using OS ENV variables. Err: %s", err)
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
	var chat bool = false
	var autoai bool = false
	var showHelp bool = false

	flag.BoolVar(&chat, "chat", false, "Start in Chat console mode")
	flag.BoolVar(&autoai, "autoai", false, "Start in Auto AI console mode")
	flag.BoolVar(&showHelp, "help", false, "Show usage information")

	flag.Parse()

	if showHelp {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		flag.PrintDefaults()
	} else if autoai {
		// start chat via console
		ai.StartConsoleAuto()
	} else if chat {
		// start chat via console
		ai.StartConsoleChat()
	} else {
		//start the server
		router.RunServer()
	}
}
