package chat

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	mylogger "backend/utils"

	"github.com/sashabaranov/go-openai"
)

// set quit command string
const quitStr = "!quit"

// StartChat starts an infinite loop that will keep asking for user input until !quit command is entered
func StartChat(apiKey string) {

	//declare messages
	messages := make([]openai.ChatCompletionMessage, 0)

	// create a buffered reader to read input from the console
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("ChatGPT Role -> ")
	// read input from console
	text, _ := reader.ReadString('\n')

	// replace CRLF with LF in the text
	text = strings.Replace(text, "\n", "", -1)

	// add the message to a list of messages
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: text,
	})

	fmt.Println("Conversation")
	fmt.Println("---------------------")

	// start an infinite loop that will keep asking for user input until !quit command is entered
	for {
		fmt.Print("-> ")
		// read input from console
		text, _ := reader.ReadString('\n')

		// check if quit command entered, if so exit the loop
		if strings.Contains(text, quitStr) {
			fmt.Println("Goodbye !!")
			break
		}

		// replace CRLF with LF in the text
		text = strings.Replace(text, "\n", "", -1)

		// add the message to a list of messages
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: text,
		})

		// create new client instance with given apiKey
		client := openai.NewClient(apiKey)

		// call OpenAI API to generate response to the user's message
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    openai.GPT3Dot5Turbo,
				Messages: messages,
			},
		)

		if err != nil {
			mylogger.Logger.Errorf("ChatCompletion error: %v\n", err)
			continue
		}

		// get the generated response from OpenAI API
		content := resp.Choices[0].Message.Content

		// add the response to the list of messages
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: content,
		})

		// print the generated response to console
		fmt.Println(content)

		mylogger.Logger.Debugf("Model: %s", resp.Model)

		jsonStr, _ := json.Marshal(messages)
		mylogger.Logger.Debugf("Messages: %s", jsonStr)

		jsonStr, _ = json.Marshal(resp.Usage)
		mylogger.Logger.Debugf("Tokens: %s", jsonStr)
	}
	reader.Reset(os.Stdin)
}
