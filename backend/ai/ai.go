package ai

// import gin framework
import (
	"backend/models"
	"backend/utils"
	"os"
	"strconv"

	"github.com/sashabaranov/go-openai"
)

var MAX_TOKENS int = 4000

var currentConfig models.OpenAIConfig = models.OpenAIConfig{}
var apiKey string = ""
var client *openai.Client = nil

// init the chat package
func init() {
	utils.Logger.Info("Init Chat Package")

	// extract and save the OpenAI api key from environment variables
	exists := false
	apiKey, exists = os.LookupEnv("OPENAI_API_KEY")

	if !exists {
		utils.Logger.Panic("OpenAI API Key not found, panicking!!!")
	}

	// read the MAX_TOKENS from environment variables as integer
	exists = false
	max_tokens, exists := os.LookupEnv("MAX_TOKENS")
	if !exists {
		utils.Logger.Error("MAX_TOKENS not found, setting to default value 4000")
		MAX_TOKENS = 4000
	}
	// convert the max_tokens string to MAX_TOKENS int
	var err error = nil
	MAX_TOKENS, err = strconv.Atoi(max_tokens)
	if err != nil {
		utils.Logger.Error("Error converting MAX_TOKENS to int, setting to default value 4000")
		MAX_TOKENS = 4000
	}

	utils.Logger.Info("MAX_TOKENS: ", MAX_TOKENS)

	// init openAi config with default values
	currentConfig = models.NewOpenAIConfig()

	// create new client instance with given apiKey
	client = openai.NewClientWithConfig(openai.DefaultConfig(apiKey))
	//client := openai.NewClient(apiKey)

	//check that client  is not nil
	if client == nil {
		utils.Logger.Panic("OpenAI Client is nil, panicking!!!")
		return
	}

	utils.Logger.Info("Chat Package Initialized")
}