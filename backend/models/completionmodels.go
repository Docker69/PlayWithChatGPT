package completionmodels

import "github.com/sashabaranov/go-openai"

type ChatCompletionRequestBody struct {
	Id       string                         `json:"id" bson:"_id,omitempty"`
	Role     string                         `json:"role"`
	Messages []openai.ChatCompletionMessage `json:"messages"`
}
