package chat

// import gin framework
import (
	"backend/models"
	mylogger "backend/utils"
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
	mylogger.Logger.Info("Init Chat Package")

	// extract and save the OpenAI api key from environment variables
	exists := false
	apiKey, exists = os.LookupEnv("OPENAI_API_KEY")

	if !exists {
		mylogger.Logger.Panic("OpenAI API Key not found, panicking!!!")
	}

	currentConfig = models.NewOpenAIConfig()

	// create new client instance with given apiKey
	client = openai.NewClientWithConfig(openai.DefaultConfig(apiKey))
	//client := openai.NewClient(apiKey)

	//check that client  is not nil
	if client == nil {
		mylogger.Logger.Panic("OpenAI Client is nil, panicking!!!")
		return
	}

	mylogger.Logger.Info("Chat Package Initialized")
}

// ChatCompletion function is the main function of the chat package
func ChatCompletion(reqBody models.ChatCompletionRequestBody) ([]openai.ChatCompletionMessage, error) {
	mylogger.Logger.WithField("UUID", reqBody.Id).Info("Chat Completion Request")

	// call OpenAI API to generate response to the user's message
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    currentConfig.Model,
			Messages: reqBody.Messages,
		},
	)

	if err != nil {
		mylogger.Logger.WithField("UUID", reqBody.Id).Errorf("ChatCompletion error: %v\n", err)
		return nil, err
	}

	if len(resp.Choices) == 0 {
		mylogger.Logger.WithField("UUID", reqBody.Id).Error("Empty response from OpenAI CreateChatCompletion API")
		return nil, errors.New("empty response from OpenAI CreateChatCompletion API")
	}

	// get the generated response from OpenAI API
	content := resp.Choices[0].Message.Content

	// add the response to the list of messages
	reqBody.Messages = append(reqBody.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: content,
	})

	mylogger.Logger.WithField("UUID", reqBody.Id).Debugf("Model: %s", resp.Model)

	jsonStr, _ := json.Marshal(reqBody.Messages)
	mylogger.Logger.WithField("UUID", reqBody.Id).Debugf("Messages: %s", jsonStr)

	jsonStr, _ = json.Marshal(resp.Usage)
	mylogger.Logger.WithField("UUID", reqBody.Id).Debugf("Tokens: %s", jsonStr)

	mylogger.Logger.WithField("UUID", reqBody.Id).Info("Chat Completion Ended!")

	return reqBody.Messages, nil
}
