package models

import "github.com/sashabaranov/go-openai"

// create struct that inherits openai.ClientConfig that has additional fields of ID, Name, Description
type OpenAIConfig struct {
	Id   string `json:"id" bson:"_id,omitempty"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	//Actual configuration of the ChatGPT, see: https://platform.openai.com/docs/api-reference/completions/create
	Model            string         `json:"model"`
	Suffix           string         `json:"suffix,omitempty" bson:"suffix,omitempty"`                       //The suffix that comes after a completion of inserted text.
	MaxTokens        int            `json:"max_tokens,omitempty" bson:"max_tokens,omitempty"`               //The maximum number of tokens to generate in the completion. Old 2048, new 4096
	Temperature      float32        `json:"temperature,omitempty" bson:"temperature,omitempty"`             //What sampling temperature to use, between 0 and 2. default 1.0
	TopP             float32        `json:"top_p,omitempty" bson:"top_p,omitempty"`                         //0..1, what percentage of the most likely tokens to keep when sampling. default 1.0
	N                int            `json:"n,omitempty" bson:"n,omitempty"`                                 //The number of completions to generate. default 1
	Stream           bool           `json:"stream,omitempty" bson:"stream,omitempty"`                       //Whether to use streaming mode. default false
	Stop             []string       `json:"stop,omitempty" bson:"stop,omitempty"`                           //A list of strings to use as the end of the sentence.
	PresencePenalty  float32        `json:"presence_penalty,omitempty" bson:"presence_penalty,omitempty"`   //The penalty to apply to the presence of a word. -2.0..2.0, default 0.0
	FrequencyPenalty float32        `json:"frequency_penalty,omitempty" bson:"frequency_penalty,omitempty"` //The penalty to apply to the frequency of a word. -2.0..2.0, default 0.0
	LogitBias        map[string]int `json:"logit_bias,omitempty" bson:"logit_bias,omitempty"`               //The bias to apply to a particular word.
	User             string         `json:"user,omitempty" bson:"user,omitempty"`                           //The user to assign to the completion.
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
