package chat

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	mongodb "backend/db"
	models "backend/models"
	mylogger "backend/utils"

	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

// set quit command string
const quitStr = "!quit"

// StartConsoleChat starts an infinite loop that will keep asking for user input until !quit command is entered
func StartConsoleChat() {

	mylogger.Logger.Info("New console chat started!")

	//declare pointer to chat struct and initialize it
	var chat *models.ChatCompletionRequestBody = new(models.ChatCompletionRequestBody)

	// create a buffered reader to read input from the console
	reader := bufio.NewReader(os.Stdin)

	//Read Human Nickname from console
	fmt.Print("Enter your nickname -> ")
	// read input from console
	text, _ := reader.ReadString('\n')
	// replace CRLF with LF in the text
	text = strings.Replace(text, "\n", "", -1)
	humnan, err := mongodb.HumansCollection.GetByNickname(text)
	if err != nil {
		mylogger.Logger.Errorf("GetHumanByNickname error: %v\n", err)
	}

	//check if human exists
	if humnan.Id == "" {
		mylogger.Logger.Infof("Human with nickname %s not found!\n", text)
		//ask for human name
		fmt.Print("Enter your name -> ")
		// read input from console
		text, _ = reader.ReadString('\n')
		// replace CRLF with LF in the text
		text = strings.Replace(text, "\n", "", -1)
		humnan.Name = text
		//ask for human nickname
		fmt.Print("Enter your nickname -> ")
		// read input from console
		text, _ = reader.ReadString('\n')
		// replace CRLF with LF in the text
		text = strings.Replace(text, "\n", "", -1)
		humnan.NickName = text
		//insert human to db
		_id, err := mongodb.HumansCollection.Insert(&humnan)
		if err != nil {
			mylogger.Logger.Errorf("InsertHuman error: %v\n", err)
			return
		}
		humnan.Id = _id
	}

	fmt.Print("ChatGPT Role -> ")
	// read input from console
	text, _ = reader.ReadString('\n')

	// replace CRLF with LF in the text
	text = strings.Replace(text, "\n", "", -1)

	// add the message to a list of messages
	chat.Messages = append(chat.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: text,
	})
	chat.Role = text

	//pass the chat as pointer to the function
	_id, err := mongodb.ChatsCollection.Insert(chat)
	if err != nil {
		mylogger.Logger.Errorf("InitNewChatDocument error: %v\n", err)
	}

	chat.Id = _id
	fmt.Println("Chat ID: ", chat.Id)

	//Create ChatRecord with the chat id and role
	chatRecord := models.ChatRecord{Id: chat.Id, Role: chat.Role}
	humnan.ChatIds = append(humnan.ChatIds, chatRecord)
	err = mongodb.HumansCollection.UpdateChats(&humnan)
	if err != nil {
		mylogger.Logger.Errorf("UpdateHumanChats error: %v\n", err)
	}

	mylogger.Logger.WithFields(
		logrus.Fields{
			"role": text,
			"UUID": _id,
		}).Info("Setting ChatGPT role")

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
		chat.Messages = append(chat.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: text,
		})

		//Update the chat document in the database
		err := mongodb.ChatsCollection.Update(chat)

		if err != nil {
			mylogger.Logger.WithField("UUID", chat.Id).Errorf("UpdateChat error: %v\n", err)
			continue
		}

		// create new client instance with given apiKey
		client := openai.NewClient(apiKey)

		// call OpenAI API to generate response to the user's message
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    openai.GPT3Dot5Turbo,
				Messages: chat.Messages,
			},
		)

		if err != nil {
			mylogger.Logger.WithField("UUID", _id).Errorf("ChatCompletion error: %v\n", err)
			continue
		}

		// get the generated response from OpenAI API
		content := resp.Choices[0].Message.Content

		// add the response to the list of messages
		chat.Messages = append(chat.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: content,
		})

		//Update the chat document in the database
		err = mongodb.ChatsCollection.Update(chat)
		if err != nil {
			mylogger.Logger.WithField("UUID", chat.Id).Errorf("UpdateChat error: %v\n", err)
			continue
		}

		// print the generated response to console
		fmt.Println(content)

		mylogger.Logger.WithField("UUID", _id).Debugf("Model: %s", resp.Model)

		jsonStr, _ := json.Marshal(chat.Messages)
		mylogger.Logger.WithField("UUID", _id).Debugf("Messages: %s", jsonStr)

		jsonStr, _ = json.Marshal(resp.Usage)
		mylogger.Logger.WithField("UUID", _id).Debugf("Tokens: %s", jsonStr)
	}
	reader.Reset(os.Stdin)

	// The context is used to inform the server it has 20 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := mongodb.Shutdown(ctx); err != nil {
		if ctx.Err() != nil {
			// context already cancelled
			mylogger.Logger.Info("MongoDB shutdown cancelled")
			return
		}
		mylogger.Logger.Fatalf("MongoDB forced to shutdown: %s", err.Error())
	}

	mylogger.Logger.WithField("UUID", _id).Info("Console Chat Ended!")
}
