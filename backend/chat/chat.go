package chat

// import gin framework
import (
	mylogger "backend/utils"
	"context"
	"encoding/json"

	"github.com/sashabaranov/go-openai"
)

type ChatCompletionRequestBody struct {
	Id       string                         `json:"id"`
	Role     string                         `json:"role"`
	Messages []openai.ChatCompletionMessage `json:"messages"`
}

// chat function is the main function of the chat package
func ChatCompletion(apiKey string, reqBody ChatCompletionRequestBody) ([]openai.ChatCompletionMessage, error) {
	mylogger.Logger.WithField("UUID", reqBody.Id).Info("Chat Completion Request")

	// create new client instance with given apiKey
	client := openai.NewClient(apiKey)

	// call OpenAI API to generate response to the user's message
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: reqBody.Messages,
		},
	)

	if err != nil {
		mylogger.Logger.WithField("UUID", reqBody.Id).Errorf("ChatCompletion error: %v\n", err)
		return nil, err
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
