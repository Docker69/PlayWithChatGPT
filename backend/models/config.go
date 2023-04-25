package models

import "github.com/sashabaranov/go-openai"

// create struct that inherits openai.ClientConfig that has additional fields of ID, Name, Description
type OpenAIConfig struct {
	Id   string `json:"id" bson:"_id,omitempty"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	//Actual configuration of the ChatGPT
	Model            string         `json:"model"`
	MaxTokens        int            `json:"max_tokens,omitempty" bson:"max_tokens,omitempty"`
	Temperature      float32        `json:"temperature,omitempty" bson:"temperature,omitempty"`
	TopP             float32        `json:"topP,omitempty" bson:"top_p,omitempty"`
	N                int            `json:"n,omitempty" bson:"n,omitempty"`
	Stream           bool           `json:"stream,omitempty" bson:"stream,omitempty"`
	Stop             []string       `json:"stop,omitempty" bson:"stop,omitempty"`
	PresencePenalty  float32        `json:"presencePenalty,omitempty" bson:"presence_penalty,omitempty"`
	FrequencyPenalty float32        `json:"frequencyPenalty,omitempty" bson:"frequency_penalty,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty" bson:"logit_bias,omitempty"`
	User             string         `json:"user,omitempty" bson:"user,omitempty"`
}

// create default struct for OpenAIConfig that accepts API key and instantiates OpenAIConfig with openai.DefaultConfig
func NewOpenAIConfig() OpenAIConfig {
	return OpenAIConfig{
		Id:    "",
		Name:  "Default GPT3.5-Turbo",
		Desc:  "The default configuration for GPT3.5-Turbo, only the Model is set to GPT3.5-Turbo",
		Model: openai.GPT3Dot5Turbo,
		//Actual configuration of the ChatGPT
	}
}
