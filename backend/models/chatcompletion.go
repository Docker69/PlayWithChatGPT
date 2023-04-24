package models

import "github.com/sashabaranov/go-openai"

type ChatCompletionRequestBody struct {
	Id       string                         `json:"id" bson:"_id,omitempty"`
	Role     string                         `json:"role"`
	HumanId  string                         `json:"humanId"`
	Messages []openai.ChatCompletionMessage `json:"messages"`
}
