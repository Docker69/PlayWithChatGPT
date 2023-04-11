// package declaration
package main

// import required packages
import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

// main function of the application
func main() {
    // set quit command string
    quitStr := "!quit"

	log.SetOutput(os.Stderr)
    // load the environment variables
    err := godotenv.Load()
    if err != nil {
        log.Printf("Error loading .env file")
    }

    // extract and save the OpenAI api key from environment variables
    apiKey := os.Getenv("OPENAI_API_KEY")

    // create new client instance with given apiKey
    client := openai.NewClient(apiKey)

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

        // call OpenAI API to generate response to the user's message
        resp, err := client.CreateChatCompletion(
            context.Background(),
            openai.ChatCompletionRequest{
                Model:    openai.GPT3Dot5Turbo,
                Messages: messages,
            },
        )

        if err != nil {
            log.Printf("ChatCompletion error: %v\n", err)
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

		jsonStr, _ := json.Marshal(messages)
		log.Printf("Messages: %s", jsonStr)
		}
}
