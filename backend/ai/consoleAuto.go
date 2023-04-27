package ai

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"backend/db/mongodb"
	"backend/models"
	"backend/utils"

	"github.com/sashabaranov/go-openai"
)

// construct the context for ChatGPT from templates collection
func constructContext(autoAI *models.AutoAI, chat *models.ChatCompletionRequestBody) error {
	//get the template from the collection
	template, err := mongodb.TemplatesCollection.GetByName("CONTEXT_DEFAULT")
	if err != nil {
		utils.Logger.Errorf("GetTemplateByName error: %v\n", err)
		return err
	}

	//construct the context from the template
	content := template.Content

	//replace now {$NAME} with the name
	content = strings.Replace(content, "{$NAME}", autoAI.Name, -1)
	//replace now {$ROLE} with the role
	content = strings.Replace(content, "{$ROLE}", autoAI.Role, -1)
	//replace now {$GOALS} with the numbered list of goals
	var numberedGoals string
	for i, goal := range autoAI.Goals {
		numberedGoals += fmt.Sprintf("%d. %s\n", i+1, goal)
	}
	content = strings.Replace(content, "{$GOALS}", numberedGoals, -1)

	//find in content all strings within curly brackets that start with $
	re := regexp.MustCompile(`\{\$[a-zA-Z_]+\}`)
	matches := re.FindAllString(content, -1)
	for _, match := range matches {
		//remove the $ sign and curly brackets
		template_name := match[2 : len(match)-1]
		//get the template corresponding to the string
		template, err := mongodb.TemplatesCollection.GetByName(template_name)
		if err != nil {
			utils.Logger.Errorf("GetTemplateByName error: %v\n", err)
			return err
		}
		//replace the string with the template content
		content = strings.Replace(content, match, template.Content, -1)
	}

	// add the message to a list of messages
	chat.Messages = append(chat.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: content,
	})
	if err != nil {
		utils.Logger.Errorf("AppendContext error: %v\n", err)
		return err
	}

	// add now the time and date in the following format: 'Wed Apr 26 01:15:31 2023' to the context
	chat.Messages = append(chat.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: fmt.Sprintf("The current time and date is %s", time.Now().Format(time.UnixDate)),
	})
	if err != nil {
		utils.Logger.Errorf("AppendContext error: %v\n", err)
		return err
	}

	//TODO:add reminder of the past

	// add user directive from template to the context
	template, err = mongodb.TemplatesCollection.GetByName("USER_DIRECTIVE")
	if err != nil {
		utils.Logger.Errorf("GetTemplateByName error: %v\n", err)
		return err
	}

	// add now the time and date in the following format: 'Wed Apr 26 01:15:31 2023' to the context
	chat.Messages = append(chat.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: template.Content,
	})
	if err != nil {
		utils.Logger.Errorf("AppendContext error: %v\n", err)
		return err
	}

	return nil
}

func initAutoAI(reader *bufio.Reader, human *models.Human) (models.AutoAI, error) {
	autoAI := models.AutoAI{}
	// Retrieve existing AutoAIs for the given Human ID
	autoAIs, err := mongodb.AutoAIsCollection.GetAllByHumanID(human.Id)
	if err != nil {
		return autoAI, fmt.Errorf("error retrieving AutoAIs: %v", err)
	}

	// Print out available AutoAIs and prompt user to choose one
	if len(autoAIs) > 0 {
		fmt.Println("Existing AutoAIs:")
		for i, autoAI := range autoAIs {
			fmt.Printf("%d: %s\n", i+1, autoAI.Name)
		}
		fmt.Print("Choose an existing AutoAI (0 to create a new one): ")
		text, err := reader.ReadString('\n')
		if err != nil {
			return autoAI, fmt.Errorf("error reading input: %v", err)
		}
		text = strings.TrimSpace(text)
		index, err := strconv.Atoi(text)
		if err != nil {
			return autoAI, fmt.Errorf("invalid input: %s", text)
		}
		if index > 0 && index <= len(autoAIs) {
			return autoAIs[index-1], nil
		}
	}

	// If no existing AutoAI was chosen, create a new one
	fmt.Print("Enter AutoAI name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return autoAI, fmt.Errorf("error reading input: %v", err)
	}
	autoAI.Name = strings.TrimSpace(name)

	fmt.Print("Enter AutoAI role: ")
	role, err := reader.ReadString('\n')
	if err != nil {
		return autoAI, fmt.Errorf("error reading input: %v", err)
	}
	autoAI.Role = strings.TrimSpace(role)

	fmt.Println("Enter AutoAI goals (hit enter with empty line when done):")
	for {
		goal, err := reader.ReadString('\n')
		if err != nil {
			return autoAI, fmt.Errorf("error reading input: %v", err)
		}
		goal = strings.TrimSpace(goal)
		if goal == "" {
			break
		}
		autoAI.Goals = append(autoAI.Goals, goal)
	}

	autoAI.HumanId = human.Id

	// Insert the new AutoAI into the database
	err = mongodb.AutoAIsCollection.Insert(autoAI)
	if err != nil {
		return autoAI, fmt.Errorf("error inserting AutoAI: %v", err)
	}

	return autoAI, nil
}

// StartConsoleChat starts an infinite loop that will keep asking for user input until !quit command is entered
func StartConsoleAuto() {

	utils.Logger.Info("New console Auto started!")

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
	human, err := mongodb.HumansCollection.GetByNickname(text)
	if err != nil {
		utils.Logger.Panicf("GetHumanByNickname panic: %v\n", err)
	}

	var autoAI models.AutoAI
	//init AutoAI
	autoAI, err = initAutoAI(reader, &human)
	if err != nil {
		utils.Logger.Errorf("InitAutoAI error: %v\n", err)
		return
	}

	//if autoAI has non empty ChatID, load the chat
	if autoAI.ChatId != "" {
		*chat, err = mongodb.ChatsCollection.GetById(autoAI.ChatId)
		if err != nil {
			utils.Logger.Errorf("GetChatById error: %v\n", err)
			return

		}
	} else {
		chat.Role = autoAI.Role
		chat.HumanId = human.Id
		//pass the chat as pointer to the function
		_id, err := mongodb.ChatsCollection.Insert(chat)
		if err != nil {
			utils.Logger.Errorf("InitNewChatDocument error: %v\n", err)
			chat.Id = _id
			fmt.Println("Chat ID: ", chat.Id)

			//Create ChatRecord with the chat id and role
			chatRecord := models.ChatRecord{Id: _id, Role: chat.Role}
			human.ChatIds = append(human.ChatIds, chatRecord)
			err = mongodb.HumansCollection.UpdateChats(&human)
			if err != nil {
				utils.Logger.Errorf("UpdateHumanChats error: %v\n", err)
			}

			//update the autoAI chat id
			autoAI.ChatId = _id
			err = mongodb.AutoAIsCollection.Update(autoAI)

			if err != nil {
				utils.Logger.Errorf("UpdateAutoAIChatId error: %v\n", err)
				return
			}

		}

	}
	// construct the context for ChatGPT
	err = constructContext(&autoAI, chat)
	if err != nil {
		utils.Logger.Errorf("ConstructContext error: %v\n", err)
		return
	}

	fmt.Println("Conversation")
	fmt.Println("---------------------")

	// get user directive from templates  collection
	userDirective, err := mongodb.TemplatesCollection.GetByName("USER_DIRECTIVE")
	if err != nil {
		utils.Logger.Errorf("GetTemplateByName error: %v\n", err)
	}

	// start an infinite loop that will keep asking for user input until !quit command is entered
	for {

		//Update the chat document in the database
		err := mongodb.ChatsCollection.Update(chat)

		if err != nil {
			utils.Logger.WithField("UUID", chat.Id).Errorf("UpdateChat error: %v\n", err)
			continue
		}

		done := make(chan bool) // Create a channel to signal when the spinner should stop

		go utils.Spinner(done) // Start the spinner

		// Call OpenAI API to generate response to the user's message
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    openai.GPT3Dot5Turbo,
				Messages: chat.Messages,
			},
		)
		done <- true // Signal the spinner to stop spinning since the operation is done

		if err != nil {
			utils.Logger.WithField("UUID", chat.Id).Errorf("ChatCompletion error: %v\n", err)
			continue
		}

		//get usage of the tokens
		jsonStr, _ := json.Marshal(resp.Usage)
		utils.Logger.WithField("UUID", chat.Id).Debugf("Tokens: %s", jsonStr)

		// get the generated response from OpenAI API
		content := resp.Choices[0].Message.Content

		var jsonValid bool = true
		if !json.Valid([]byte(content)) {
			// Handle invalid JSON error
			utils.Logger.WithField("UUID", chat.Id).Errorf("Invalid JSON response: %v\n", err)
			content, err = utils.FixJSON(content)
			if err != nil {
				utils.Logger.WithField("UUID", chat.Id).Errorf("FixJSON error: %v\n", err)
				jsonValid = false
			}
			//check if valid after the fix
			if !json.Valid([]byte(content)) {
				utils.Logger.WithField("UUID", chat.Id).Errorf("Invalid JSON response after fix: %v\n", err)
				jsonValid = false
			}
		}

		if jsonValid {

			// add the response to the list of messages
			chat.Messages = append(chat.Messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: content,
			})

			//Update the chat document in the database
			err = mongodb.ChatsCollection.Update(chat)
			if err != nil {
				utils.Logger.WithField("UUID", chat.Id).Errorf("UpdateChat error: %v\n", err)
				continue
			}

			//convert string to byte array
			contentByte := []byte(content)
			//unmarshal the byte array to struct
			var responseData models.Response
			err = json.Unmarshal(contentByte, &responseData)
			if err != nil {
				utils.Logger.WithField("UUID", chat.Id).Errorf("Error unmarsheling response: %v\n", err)
				return
			}
			//output the response in human readible format
			// Output the thoughts
			fmt.Println("==========================================================================================")
			fmt.Println("Thoughts:")
			fmt.Printf("- Text: %v\n", responseData.Thoughts.Text)
			fmt.Printf("- Reasoning: %v\n", responseData.Thoughts.Reasoning)
			fmt.Printf("- Plan: %v\n", responseData.Thoughts.Plan)
			fmt.Printf("- Criticism: %v\n", responseData.Thoughts.Criticism)
			fmt.Printf("- Speak: %v\n", responseData.Thoughts.Speak)

			// Output the command
			fmt.Println("Command:")
			fmt.Printf("- Name: %v\n", responseData.Command.Name)
			fmt.Printf("- Args:\n")
			if responseData.Command.Args.URL != "" {
				fmt.Printf("  - URL: %v\n", responseData.Command.Args.URL)
			}
			if responseData.Command.Args.Question != "" {
				fmt.Printf("  - Question: %v\n", responseData.Command.Args.Question)
			}
			if responseData.Command.Args.Input != "" {
				fmt.Printf("  - Input: %v\n", responseData.Command.Args.Input)
			}
			if responseData.Command.Args.Reason != "" {
				fmt.Printf("  - Reason: %v\n", responseData.Command.Args.Reason)
			}
			fmt.Println("==========================================================================================")

			//utils.Logger.WithField("UUID", chat.Id).Debugf("Model: %s", resp.Model)

			//jsonStr, _ := json.Marshal(chat.Messages)
			//utils.Logger.WithField("UUID", chat.Id).Debugf("Messages: %s", jsonStr)

			// read input from console
			fmt.Println("Enter 'y' to authorise command, '!quit' to exit program, or enter feedback for ...")
			fmt.Print("Input -> ")
			// read input from console
			input, _ := reader.ReadString('\n')
			// replace CRLF with LF in the text
			input = strings.Replace(input, "\n", "", -1)

			//check if user authorised the command, equals exactly to 'y'
			isY := strings.ToLower(strings.TrimSpace(input))

			if strings.Contains(input, quitStr) {
				// check if quit command entered, if so exit the loop
				fmt.Println("Goodbye !!")
				break
			} else if isY != "y" {
				//add user input
				chat.Messages = append(chat.Messages, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Human feedback: " + input,
				})
				if err != nil {
					utils.Logger.Errorf("Adding user input error: %v\n", err)
				}

			}

			//add user directive from template to the context
			chat.Messages = append(chat.Messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: userDirective.Content,
			})
			if err != nil {
				utils.Logger.Errorf("Adding user directive error: %v\n", err)
			}
		}

	}
	reader.Reset(os.Stdin)

	// The context is used to inform the server it has 20 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := mongodb.Shutdown(ctx); err != nil {
		if ctx.Err() != nil {
			// context already cancelled
			utils.Logger.Info("MongoDB shutdown cancelled")
			return
		}
		utils.Logger.Fatalf("MongoDB forced to shutdown: %s", err.Error())
	}

	utils.Logger.WithField("UUID", chat.Id).Info("Console Chat Ended!")
}
