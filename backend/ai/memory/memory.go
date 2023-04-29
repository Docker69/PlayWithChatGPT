package memory

import (
	"backend/utils"
	"context"

	"github.com/sashabaranov/go-openai"
)

//write interface for storing memory data either in local storage or redis
//the following capabilities are required:
//- get
//- add
//- clear
//- get relevant
//- get stats

type MemoryCache interface {
	GetRelevantMemories(data string, max int) []string
	AddMemory(text string) error
	Clear() error
	GetStats() int
}

var aiClient *openai.Client = nil

func Init(client *openai.Client) {
	// create new client instance with given apiKey
	aiClient = client
	//client := openai.NewClient(apiKey)

	//check that client  is not nil
	if client == nil {
		utils.Logger.Panic("OpenAI Client is nil, panicking!!!")
		return
	}

}

// create function that accepts string and creates embeddings with ada model of openai
func createAdaEmbeddings(text string) []float32 {
	// create embeddings
	response, err := aiClient.CreateEmbeddings(
		context.Background(),
		openai.EmbeddingRequest{
			Model: openai.AdaEmbeddingV2,
			Input: []string{text},
		},
	)
	if err != nil {
		utils.Logger.Error("Error creating embeddings: ", err)
		return nil
	}

	return response.Data[0].Embedding
}
