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

	"backend/ai/capabilities"
	"backend/ai/helpers"
	"backend/ai/memory"
	"backend/db/mongodb"
	"backend/models"
	"backend/utils"

	"github.com/sashabaranov/go-openai"
)

// How many past contexts to store
const MEMORY_DEPTH = 5

// construct the template for main context
func constructTemplate(autoAI *models.AutoAI) (string, error) {

	//get the template from the collection
	template, err := mongodb.TemplatesCollection.GetByName("CONTEXT_DEFAULT")
	if err != nil {
		utils.Logger.Errorf("GetTemplateByName error: %v", err)
		return "", err
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
			utils.Logger.Errorf("GetTemplateByName error: %v", err)
			return "", err
		}
		//replace the string with the template content
		content = strings.Replace(content, match, template.Content, -1)
	}

	return content, nil
}

// construct the context for ChatGPT from templates collection
func constructContext(autoAI *models.AutoAI, chatContext *string, fullHistory *models.ChatCompletionRequestBody, mem *memory.MemoryCache, memoryToAdded string) ([]openai.ChatCompletionMessage, error) {

	// Initialize as an empty slice
	messages := make([]openai.ChatCompletionMessage, 0)

	// add the message to a list of messages
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: *chatContext,
	})

	// add now the time and date in the following format: 'Wed Apr 26 01:15:31 2023' to the context
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: fmt.Sprintf("The current time and date is %s", time.Now().Format(time.UnixDate)),
	})
	/*
		//get last assitant message from full history, iterate from last to first
		var lastAssistantMessage string = ""

		for i := len(fullHistory.Messages) - 1; i >= 0; i-- {
			if fullHistory.Messages[i].Role == openai.ChatMessageRoleAssistant {
				lastAssistantMessage = fullHistory.Messages[i].Content
				break
			}
		}
	*/
	/*
		//loop over all fullHistory and find the assistant messages
		var searchEmbedding string = ""

			re := regexp.MustCompile(`,\s*|\s+`)
			resultMap := make(map[string]bool)

			for _, msg := range fullHistory.Messages {
				if msg.Role == openai.ChatMessageRoleAssistant {
					responseData, err := unMarshalMsgContent(&msg.Content, fullHistory)
					if err != nil {
						continue
					}

					//searchEmbedding += responseData.Memorize.Subject + " " + responseData.Memorize.Information + " "

					//tokenize key words, insert into map and don't repeat same key word
					for _, s := range re.Split(responseData.Thoughts.Keywords, -1) {
						if !resultMap[s] && s != "" {
							resultMap[s] = true
							searchEmbedding += s + " "
						}
					}
				}
			}
	*/

	//current token count
	tokenCount := helpers.NumTokensFromMessages(messages, currentConfig.Model)
	//current memory token count
	memoriesTokenCount := 0
	//current memory string
	memoriesStr := ""
	memoriesToAdd := openai.ChatCompletionMessage{}

	//join the string array to a single string and search the memory
	past := (*mem).GetRelevantMemories(memoryToAdded, MEMORY_DEPTH)

	//iterate on past memories and add them to the context
	for _, m := range past {

		memoriesStr += m + "\n"
		memoriesToAdd = openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: fmt.Sprintf("This reminds you of these events from your past: %s", memoriesStr),
		}

		//get token count of memory to add
		memoriesTokenCount = helpers.NumTokensFromMessages([]openai.ChatCompletionMessage{memoriesToAdd}, currentConfig.Model)

		//break if going over 3000 tokens over all
		if (tokenCount + memoriesTokenCount) > 2700 {
			break
		}
	}

	if len(past) > 0 && len(memoriesStr) > 0 {
		//add the memories to the context
		messages = append(messages, memoriesToAdd)
	}
	/*
		//iterate over fullHistory and add system messages except the first two
		if len(fullHistory.Messages) > 2 {
			for _, msg := range fullHistory.Messages[2:] {
				if msg.Role == openai.ChatMessageRoleSystem || msg.Role == openai.ChatMessageRoleUser {
					messages = append(messages, msg)
				}
			}
		}
	*/
	return messages, nil
}

func initAutoAI(reader *bufio.Reader, human *models.Human, mem *memory.MemoryCache) (models.AutoAI, error) {
	autoAI := models.AutoAI{}
	// Retrieve existing AutoAIs for the given Human ID
	autoAIs, err := mongodb.AutoAIsCollection.GetAllByHumanID(human.Id)
	if err != nil {
		return autoAI, fmt.Errorf("error retrieving AutoAIs: %v", err)
	}

	// Print out available AutoAIs and prompt user to choose one
	if len(autoAIs) > 0 {
		stats := (*mem).GetStats()
		if stats != nil {
			fmt.Println("Memory stats:")
			fmt.Println(stats)
			fmt.Println()
		}

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

	//clear memory from all memories
	//(*mem).Clear()

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
	autoAI.Id, err = mongodb.AutoAIsCollection.Insert(autoAI)
	if err != nil {
		return autoAI, fmt.Errorf("error inserting AutoAI: %v", err)
	}

	return autoAI, nil
}

func initChats(autoAI *models.AutoAI, human *models.Human, chatFullHistory *models.ChatCompletionRequestBody, chatContext *string, mem *memory.MemoryCache) error {
	//if autoAI has non empty ChatID, load the chat
	if autoAI.ChatId != "" {
		var err error = nil
		*chatFullHistory, err = mongodb.ChatsCollection.GetById(autoAI.ChatId)
		if err != nil {
			utils.Logger.Errorf("GetChatById error: %v", err)
			return err
		}

		if len(chatFullHistory.Messages) == 0 {
			chatFullHistory.Messages, err = constructContext(autoAI, chatContext, chatFullHistory, mem, "")
			if err != nil {
				utils.Logger.Errorf("full chat history init error: %v", err)
				return err
			}
		}
	} else {
		chatFullHistory.Role = autoAI.Role
		chatFullHistory.HumanId = human.Id
		//pass the chat as pointer to the function
		_id, err := mongodb.ChatsCollection.Insert(chatFullHistory)
		if err != nil {
			utils.Logger.Errorf("InitNewChatDocument error: %v", err)
			return err
		}

		chatFullHistory.Id = _id
		fmt.Println("Chat ID: ", chatFullHistory.Id)

		//Create ChatRecord with the chat id and role
		chatRecord := models.ChatRecord{Id: _id, Role: chatFullHistory.Role}
		human.ChatIds = append(human.ChatIds, chatRecord)
		err = mongodb.HumansCollection.UpdateChats(human)
		if err != nil {
			utils.Logger.Errorf("UpdateHumanChats error: %v", err)
			return err
		}

		//update the autoAI chat id
		autoAI.ChatId = _id
		err = mongodb.AutoAIsCollection.Update(autoAI)

		if err != nil {
			utils.Logger.Errorf("UpdateAutoAIChatId error: %v", err)
			return err
		}

		// Init full chat history
		chatFullHistory.Messages, err = constructContext(autoAI, chatContext, chatFullHistory, mem, "")
		if err != nil {
			utils.Logger.Errorf("full chat history init error: %v", err)
			return err
		}

	}

	return nil
}

func unMarshalMsgContent(msgContent *string, chatFullHistory *models.ChatCompletionRequestBody) (models.Response, error) {

	var responseData models.Response = models.Response{}
	var jsonValid bool = true
	var err error = nil

	jsonStr := *msgContent

	if !json.Valid([]byte(jsonStr)) {

		// Handle invalid JSON error
		utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("Invalid JSON response")

		//jsonStr, _ := json.Marshal(chat.Messages)
		utils.Logger.WithField("UUID", chatFullHistory.Id).Debugf("content dump: %s", *msgContent)

		//find first "{" and truncate everything up to it
		index := strings.Index(jsonStr, "{") // find index of first "{"
		if index != -1 {                     // check if "{" was found
			jsonStr = (jsonStr)[index:] // truncate everything before "{"
		}
		//check if valid after the fix
		if !json.Valid([]byte(jsonStr)) {
			utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("Invalid JSON response after fix, trying second fix")
			//find last "}" and truncate everything after it
			index = strings.LastIndex(jsonStr, "}") // find index of last "}"
			if index != -1 {                        // check if "}" was found
				jsonStr = (jsonStr)[:index+1] // truncate everything after "}"
			}
			//check if valid after the fix
			if !json.Valid([]byte(jsonStr)) {
				utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("invalid JSON response after fix, trying third fix: %v", err)

				jsonStr, err = utils.FixJSON(jsonStr)
				if err != nil {
					utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("FixJSON error: %v", err)
					return responseData, fmt.Errorf("FixJSON error: %v", err)
				}
				//check if valid after the fix
				if !json.Valid([]byte(jsonStr)) {
					utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("invalid JSON response after fix", err)
					jsonValid = false
					return responseData, fmt.Errorf("invalid JSON response after fix")
				}
			}
		}
	}

	if jsonValid {
		//convert string to byte array
		contentByte := []byte(jsonStr)
		//unmarshal the byte array to struct
		err = json.Unmarshal(contentByte, &responseData)
		if err != nil {
			utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("==========================================================================================\n"+
				"Error unmarsheling response: %v,"+
				"Content Dump: %s\n"+
				"==========================================================================================\n", err, msgContent)
			//try to unmarshel then to a map
			var dataMap map[string]interface{}
			if err = json.Unmarshal(contentByte, &dataMap); err != nil {
				// Handle any errors that may occur during parsing
				utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("Error unmarsheling to a map: %v", err)
				return models.Response{}, err
			}

			//output the map
			utils.Logger.WithField("UUID", chatFullHistory.Id).Info("==========================================================================================")
			utils.Logger.WithField("UUID", chatFullHistory.Id).Info("Map:")
			utils.Logger.WithField("UUID", chatFullHistory.Id).Info(dataMap)
			utils.Logger.WithField("UUID", chatFullHistory.Id).Info("==========================================================================================")

			return models.Response{}, err
		}
		//just in case the conent was fixed
		*msgContent = jsonStr
	}
	return responseData, nil
}

func outputResponse(responseData models.Response) {
	// Output the thoughts
	fmt.Println("==========================================================================================")
	fmt.Println("Thoughts:")
	fmt.Printf("- Text: %v\n", responseData.Thoughts.Text)
	fmt.Printf("- Reasoning: %v\n", responseData.Thoughts.Reasoning)
	fmt.Printf("- Plan: %v\n", responseData.Thoughts.Plan)
	fmt.Printf("- Criticism: %v\n", responseData.Thoughts.Criticism)
	if responseData.Thoughts.Speak != "" {
		fmt.Printf("- Speak: %v\n", responseData.Thoughts.Speak)
	}
	if responseData.Thoughts.Keywords != "" {
		fmt.Printf("- Keywords: %v\n", responseData.Thoughts.Keywords)
	}

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
	if responseData.Command.Args.File != "" {
		fmt.Printf("  - File: %v\n", responseData.Command.Args.File)
	}
	if responseData.Command.Args.Text != "" {
		fmt.Printf("  - Text: %v\n", responseData.Command.Args.Text)
	}
	// Output what to memorize
	//fmt.Println("Memorize:")
	//fmt.Printf("- Subject: %v", responseData.Memorize.Subject)
	//fmt.Printf("- Information: %v", responseData.Memorize.Information)
	fmt.Println("==========================================================================================")
}

// StartConsoleChat starts an infinite loop that will keep asking for user input until !quit command is entered
func StartConsoleAuto() {

	// call initMemory function to get memory storage
	err := InitMemory()
	if err != nil {
		utils.Logger.Errorf("initMemory error: %v", err)
		return
	}

	utils.Logger.Info("New console Auto started!")

	//declare pointer to chat struct and initialize it
	var chatFullHistory *models.ChatCompletionRequestBody = new(models.ChatCompletionRequestBody)

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
		utils.Logger.Panicf("GetHumanByNickname panic: %v", err)
	}

	var autoAI models.AutoAI
	//init AutoAI
	autoAI, err = initAutoAI(reader, &human, &Mem)
	if err != nil {
		utils.Logger.Errorf("InitAutoAI error: %v", err)
		return
	}

	//construct the template for the system context
	chatContext, err := constructTemplate(&autoAI)
	if err != nil {
		utils.Logger.Errorf("ConstructTemplate error: %v", err)
		return
	}

	// get user directive from templates  collection
	userDirective, err := mongodb.TemplatesCollection.GetByName("USER_DIRECTIVE")
	if err != nil {
		utils.Logger.Errorf("GetTemplateByName error: %v", err)
	}

	//init Chats
	err = initChats(&autoAI, &human, chatFullHistory, &chatContext, &Mem)
	if err != nil {
		utils.Logger.Errorf("InitChats error: %v", err)
		return
	}

	memoryToAdd := ""

	fmt.Println("Conversation")
	fmt.Println("---------------------")

	// start an infinite loop that will keep asking for user input until !quit command is entered
	for {

		// Create a channel to signal when the spinner should stop
		done := make(chan bool)

		// Start the spinner
		go utils.Spinner(done)

		// construct the context for ChatGPT
		messagesToSend, err := constructContext(&autoAI, &chatContext, chatFullHistory, &Mem, memoryToAdd)
		if err != nil {
			utils.Logger.Errorf("ConstructContext error: %v", err)
			return
		}

		// get the token count from the chat messages
		numTokens := helpers.NumTokensFromMessages(messagesToSend, currentConfig.Model)
		//User Directive Message
		userDirectiveMessage := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: userDirective.Content,
		}
		userDirectiveTokens := helpers.NumTokensFromMessages([]openai.ChatCompletionMessage{userDirectiveMessage}, currentConfig.Model)

		//current potential count of Tokens
		numTokens += userDirectiveTokens

		// add user directive to the context
		messagesToSend = append(messagesToSend, userDirectiveMessage)

		// get the final token count from the chat messages
		numTokens = helpers.NumTokensFromMessages(messagesToSend, currentConfig.Model)

		//TODO read token limit from .env file
		allowedTokens := MAX_TOKENS - numTokens
		utils.Logger.WithField("UUID", chatFullHistory.Id).Debugf("Allowed Tokens for response: %d", allowedTokens)

		// Call OpenAI API to generate response to the user's message
		assistant_response, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:     currentConfig.Model,
				Messages:  messagesToSend,
				MaxTokens: allowedTokens,
			},
		)
		done <- true // Signal the spinner to stop spinning since the operation is done

		if err != nil {
			utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("ChatCompletion error: %v", err)
			continue
		}

		//get usage of the tokens
		jsonStr, _ := json.Marshal(assistant_response.Usage)
		utils.Logger.WithField("UUID", chatFullHistory.Id).Debugf("Tokens: %s", jsonStr)

		// get the generated response from OpenAI API
		content := assistant_response.Choices[0].Message.Content

		//parse the response to get the command
		responseData, err := unMarshalMsgContent(&content, chatFullHistory)
		if err != nil {
			//if parse unsuccessful try again without taking into account current response
			content = ""
			continue
		}

		//output the response in human readible format
		outputResponse(responseData)

		if responseData.Command.Name == "task_complete" {
			// check if quit command entered, if so exit the loop
			fmt.Println("Goodbye !!")

			//Update the chat document in the database befoew exiting
			err = mongodb.ChatsCollection.Update(chatFullHistory)
			if err != nil {
				utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("on exit UpdateChat error: %v", err)
			}
			break
		}

		// read input from console
		fmt.Println("Enter 'y' to authorise command, '!quit' to exit program, or enter feedback for ...")
		fmt.Print("Input -> ")
		// read input from console
		input, _ := reader.ReadString('\n')
		// replace CRLF with LF in the text
		input = strings.Replace(input, "\n", "", -1)

		//check if user authorised the command, equals exactly to 'y'
		isY := strings.ToLower(strings.TrimSpace(input))

		userInput := "GENERATE NEXT COMMAND JSON"
		if strings.Contains(input, quitStr) {
			// check if quit command entered, if so exit the loop
			fmt.Println("Goodbye !!")

			//Update the chat document in the database befoew exiting
			err = mongodb.ChatsCollection.Update(chatFullHistory)
			if err != nil {
				utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("on exit UpdateChat error: %v", err)
			}
			break
		} else if isY != "y" {
			userInput = input
			//add user input
			chatFullHistory.Messages = append(chatFullHistory.Messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Human feedback: " + input,
			})
			if err != nil {
				utils.Logger.Errorf("Adding user input error: %v", err)
			}
		}

		//TODO: get it from environment variables
		responseData.Command.Args.Path = "/home/bennyk/projects/PlayWithChatGPT/backend/ai_workspace"
		// execute any command requested by the AI
		resultStr, err := capabilities.GetCapabilityFactory().Execute(responseData.Command, &Mem)
		if err != nil {
			utils.Logger.Errorf("Execute command error: %v", err)
		} else if resultStr != "" {
			//add user input
			chatFullHistory.Messages = append(chatFullHistory.Messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: resultStr,
			})
			if err != nil {
				utils.Logger.Errorf("Adding command result error: %v", err)
			}
		}

		//only if content isn't empty
		if content != "" {
			memoryToAdd = fmt.Sprintf("%s\nResult: %s\nHuman Feedback: %s", content, resultStr, userInput)
			//add the response to the memory
			err = Mem.AddMemory(memoryToAdd)
			if err != nil {
				utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("AddMemory error: %v", err)
			}

			// add the response to the list of messages
			chatFullHistory.Messages = append(chatFullHistory.Messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: content,
			})

			//try to send only few words to get the relevant memory for next response
			memoryToAdd = "thoughts command " + responseData.Command.Name + " " + responseData.Thoughts.Keywords
		}

		//Update the chat document in the database
		err = mongodb.ChatsCollection.Update(chatFullHistory)
		if err != nil {
			utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("UpdateChat error: %v", err)
			continue
		}
	}

	//close the console
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

	utils.Logger.WithField("UUID", chatFullHistory.Id).Info("Console Chat Ended!")
}

func GetCapabilityFactory() {
	panic("unimplemented")
}
