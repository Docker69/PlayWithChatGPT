package ai

// import gin framework
import (
	redisclient "backend/db/redis"
	"backend/models"
	"backend/utils"
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/sashabaranov/go-openai"
)

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

	jsonBytes, _ := json.Marshal(currentConfig)
	err := redisclient.SetJson("chat_config", ".", string(jsonBytes))
	if err != nil {
		utils.Logger.Error("Error setting chat config in redis,  error: ", err)
	}
}

// ChatCompletion function is the main function of the chat package
func ChatCompletion(reqBody models.ChatCompletionRequestBody) ([]openai.ChatCompletionMessage, error) {
	utils.Logger.WithField("UUID", reqBody.Id).Info("Chat Completion Request")

	numTokens := utils.NumTokensFromMessages(reqBody.Messages, currentConfig.Model)
	//TODO read token limit from .env file
	allowedTokens := 4000 - numTokens
	utils.Logger.WithField("UUID", reqBody.Id).Debugf("Allowed Tokens for response: %d", allowedTokens)
	// call OpenAI API to generate response to the user's message
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     currentConfig.Model,
			Messages:  reqBody.Messages,
			MaxTokens: allowedTokens,
		},
	)

	if err != nil {
		utils.Logger.WithField("UUID", reqBody.Id).Errorf("ChatCompletion error: %v\n", err)
		return nil, err
	}

	if len(resp.Choices) == 0 {
		utils.Logger.WithField("UUID", reqBody.Id).Error("Empty response from OpenAI CreateChatCompletion API")
		return nil, errors.New("empty response from OpenAI CreateChatCompletion API")
	}

	// get the generated response from OpenAI API
	content := resp.Choices[0].Message.Content

	// add the response to the list of messages
	reqBody.Messages = append(reqBody.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: content,
	})

	utils.Logger.WithField("UUID", reqBody.Id).Debugf("Model: %s", resp.Model)

	jsonStr, _ := json.Marshal(reqBody.Messages)
	utils.Logger.WithField("UUID", reqBody.Id).Debugf("Messages: %s", jsonStr)

	jsonStr, _ = json.Marshal(resp.Usage)
	utils.Logger.WithField("UUID", reqBody.Id).Debugf("Tokens: %s", jsonStr)

	utils.Logger.WithField("UUID", reqBody.Id).Info("Chat Completion Ended!")

	return reqBody.Messages, nil
}
