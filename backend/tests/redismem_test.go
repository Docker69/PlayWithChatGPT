package tests_test

import (
	"backend/ai/memory"
	"backend/utils"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

func TestRelevantResults(t *testing.T) {

	// load the environment variables
	err := godotenv.Load()
	if err != nil {
		utils.Logger.Infof(".env file not found, using OS ENV variables. Err: %s", err)
	}

	mem, err := memory.NewRedisMem()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// extract and save the OpenAI api key from environment variables
	apiKey, exists := os.LookupEnv("OPENAI_API_KEY")

	if !exists {
		utils.Logger.Panic("OpenAI API Key not found, panicking!!!")
	}

	// create new client instance with given apiKey
	client := openai.NewClientWithConfig(openai.DefaultConfig(apiKey))

	memory.Init(client)

	//results := mem.GetRelevantMemories("{    \"thoughts\": {        \"text\": \"I need to ask the user what type of SciFi books they enjoy reading. This will help me narrow down my search results and provide a better recommendation.\",        \"plan\": [            \"Ask the user what type of SciFi books they enjoy reading using the 'user_input' command.\"        ],        \"reasoning\": \"By asking the user about their preferences, I can provide a more personalized recommendation that they are more likely to enjoy.\",        \"criticism\": \"I need to ensure that I don't rely solely on the user's input and use other resources such as current best seller lists and book reviews to come up with a comprehensive list of recommendations.\"    },    \"command\": {        \"name\": \"user_input\",        \"args\": {            \"input\": \"Please tell me what type of SciFi books you enjoy reading. Are you more interested in hard science fiction, space operas or fantasy-based science fiction?\"        }    }}", 5)
	results := mem.GetRelevantMemories("thoughts command browse_website", 5)

	//print results
	f, err := os.Create("redismem_test_results.txt")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	defer f.Close()

	for _, result := range results {
		//convert string to bytes
		_, err = f.WriteString("===================================================================================\n" +
			result +
			"\n===================================================================================\n")
		if err != nil {
			t.Error("Unexpected error: can't write to file")
		}
	}

	// Flush the buffer to ensure all the data is written to the file
	err = f.Sync()
	if err != nil {
		t.Error("Unexpected error: can't flush buffer")
	}

}
