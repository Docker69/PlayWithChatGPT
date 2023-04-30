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

	"backend/ai/helpers"
	"backend/ai/memory"
	"backend/db/mongodb"
	"backend/models"
	"backend/utils"

	"github.com/sashabaranov/go-openai"
)

// construct the template for main context
func constructTemplate(autoAI *models.AutoAI) (string, error) {

	//get the template from the collection
	template, err := mongodb.TemplatesCollection.GetByName("CONTEXT_DEFAULT")
	if err != nil {
		utils.Logger.Errorf("GetTemplateByName error: %v\n", err)
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
			utils.Logger.Errorf("GetTemplateByName error: %v\n", err)
			return "", err
		}
		//replace the string with the template content
		content = strings.Replace(content, match, template.Content, -1)
	}

	return content, nil
}

// construct the context for ChatGPT from templates collection
func constructContext(autoAI *models.AutoAI, chatContext *string, fullHistory *models.ChatCompletionRequestBody, mem *memory.MemoryCache) ([]openai.ChatCompletionMessage, error) {

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

	// add reminder of the past
	pastThoughts := "This reminds you of these events from your past:\n"

	const REMINDER_LIMIT = 5

	//check depth of history, if empty or only the first context don't bother
	n := len(fullHistory.Messages)

	if n > 1 {
		// get the relevant past from the memory
		// Extract the content from upto last REMINDER_LIMIT messages and add it to a string array
		var contents []string
		var past []string
		// get only the assistant messages from full history
		var assistantMessages []openai.ChatCompletionMessage = []openai.ChatCompletionMessage{}
		for _, msg := range fullHistory.Messages {
			if msg.Role == openai.ChatMessageRoleAssistant {
				assistantMessages = append(assistantMessages, msg)
			}
		}

		m := len(assistantMessages)
		if m > 0 {
			var forEmbeddings []openai.ChatCompletionMessage = []openai.ChatCompletionMessage{}
			if m > REMINDER_LIMIT {
				//last REMINDER_LIMIT messages
				forEmbeddings = assistantMessages[m-REMINDER_LIMIT:]
			} else {
				//get all messages except first one (main context)
				forEmbeddings = assistantMessages
			}

			for _, msg := range forEmbeddings {
				contents = append(contents, msg.Content)
			}

			//join the string array to a single string and search the memory
			past = (*mem).GetRelevantMemories(strings.Join(contents, ", "), REMINDER_LIMIT)

		}

		// if no past, no point in the appending
		if len(past) > 0 {
			pastMessageToAppend := openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: pastThoughts + strings.Join(past, " ") + "\n",
			}
			// get the token count from the chat messages
			curTokens := helpers.NumTokensFromMessages(messages, currentConfig.Model)
			// tokens to be added
			numTokensToAppend := helpers.NumTokensFromMessages([]openai.ChatCompletionMessage{pastMessageToAppend}, currentConfig.Model)

			//if too many tokens start removing past messages until we have 2500 tokens
			for numTokensToAppend+curTokens > 2500 {
				//remove first message in contents (oldest)
				contents = contents[1:]
				//join the string array to a single string and search the memory
				past = (*mem).GetRelevantMemories(strings.Join(contents, " "), REMINDER_LIMIT)

				pastMessageToAppend = openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleSystem,
					Content: pastThoughts + strings.Join(past, " ") + "\n",
				}

				numTokensToAppend = helpers.NumTokensFromMessages([]openai.ChatCompletionMessage{pastMessageToAppend}, currentConfig.Model)

			}
			//finally add the past message to the context
			messages = append(messages, pastMessageToAppend)
		}

	}

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

// init function of the memory storage
func initMemory() (memory.MemoryCache, error) {
	// get memory type from env
	memType, exists := os.LookupEnv("MEMORY_STORAGE")
	if !exists {
		utils.Logger.Info("MEMORY_STORAGE is not defined in env, setting to \"local\"")
		memType = "local"
	}

	var mem memory.MemoryCache
	if memType == "redis" {
		redisMem, err := memory.NewRedisMem()
		if err != nil {
			utils.Logger.Errorf("NewRedisMem error: %v\n", err)
			return nil, err
		}
		mem = redisMem
		utils.Logger.Info("Memory storage is Redis")
	} else {
		//TODO: add support for local memory type
		utils.Logger.Panic("NewLocalStorageMem panic: NOT READY!!!!!\n")
		localMem, err := memory.NewLocalStorageMem(".")
		if err != nil {
			utils.Logger.Panicf("NewLocalStorageMem panic: %v\n", err)
			return nil, err
		}
		mem = localMem
		utils.Logger.Info("Memory storage is Local")
	}

	memory.Init(client)
	return mem, nil
}

func initChats(autoAI *models.AutoAI, human *models.Human, chatFullHistory *models.ChatCompletionRequestBody, chatContext *string, mem *memory.MemoryCache) error {
	//if autoAI has non empty ChatID, load the chat
	if autoAI.ChatId != "" {
		var err error = nil
		*chatFullHistory, err = mongodb.ChatsCollection.GetById(autoAI.ChatId)
		if err != nil {
			utils.Logger.Errorf("GetChatById error: %v\n", err)
			return err
		}

		if len(chatFullHistory.Messages) == 0 {
			chatFullHistory.Messages, err = constructContext(autoAI, chatContext, chatFullHistory, mem)
			if err != nil {
				utils.Logger.Errorf("full chat history init error: %v\n", err)
				return err
			}
		}
	} else {
		chatFullHistory.Role = autoAI.Role
		chatFullHistory.HumanId = human.Id
		//pass the chat as pointer to the function
		_id, err := mongodb.ChatsCollection.Insert(chatFullHistory)
		if err != nil {
			utils.Logger.Errorf("InitNewChatDocument error: %v\n", err)
			return err
		}

		chatFullHistory.Id = _id
		fmt.Println("Chat ID: ", chatFullHistory.Id)

		//Create ChatRecord with the chat id and role
		chatRecord := models.ChatRecord{Id: _id, Role: chatFullHistory.Role}
		human.ChatIds = append(human.ChatIds, chatRecord)
		err = mongodb.HumansCollection.UpdateChats(human)
		if err != nil {
			utils.Logger.Errorf("UpdateHumanChats error: %v\n", err)
			return err
		}

		//update the autoAI chat id
		autoAI.ChatId = _id
		err = mongodb.AutoAIsCollection.Update(autoAI)

		if err != nil {
			utils.Logger.Errorf("UpdateAutoAIChatId error: %v\n", err)
			return err
		}

		// Init full chat history
		chatFullHistory.Messages, err = constructContext(autoAI, chatContext, chatFullHistory, mem)
		if err != nil {
			utils.Logger.Errorf("full chat history init error: %v\n", err)
			return err
		}

	}

	return nil
}

// StartConsoleChat starts an infinite loop that will keep asking for user input until !quit command is entered
func StartConsoleAuto() {

	// call initMemory function to get memory storage
	mem, err := initMemory()
	if err != nil {
		utils.Logger.Errorf("initMemory error: %v\n", err)
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
		utils.Logger.Panicf("GetHumanByNickname panic: %v\n", err)
	}

	var autoAI models.AutoAI
	//init AutoAI
	autoAI, err = initAutoAI(reader, &human, &mem)
	if err != nil {
		utils.Logger.Errorf("InitAutoAI error: %v\n", err)
		return
	}

	//construct the template for the system context
	chatContext, err := constructTemplate(&autoAI)
	if err != nil {
		utils.Logger.Errorf("ConstructTemplate error: %v\n", err)
		return
	}

	// get user directive from templates  collection
	userDirective, err := mongodb.TemplatesCollection.GetByName("USER_DIRECTIVE")
	if err != nil {
		utils.Logger.Errorf("GetTemplateByName error: %v\n", err)
	}

	//init Chats
	err = initChats(&autoAI, &human, chatFullHistory, &chatContext, &mem)
	if err != nil {
		utils.Logger.Errorf("InitChats error: %v\n", err)
		return
	}

	fmt.Println("Conversation")
	fmt.Println("---------------------")

	// start an infinite loop that will keep asking for user input until !quit command is entered
	for {

		// Create a channel to signal when the spinner should stop
		done := make(chan bool)

		// Start the spinner
		go utils.Spinner(done)

		// construct the context for ChatGPT
		messagesToSend, err := constructContext(&autoAI, &chatContext, chatFullHistory, &mem)
		if err != nil {
			utils.Logger.Errorf("ConstructContext error: %v\n", err)
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

		/*
			// fill additional mesages from full history starting from the last message till more than 2500 tokens
			n := len(chatFullHistory.Messages)
			if n > 2 {
				messageTokens := helpers.NumTokensFromMessages(chatFullHistory.Messages[1:], currentConfig.Model)
				//iterate in reverse order over chatFullHistory, don't include the 3 first messages that are already there
				var idx int = n - 2
				for idx = 2; idx < n && numTokens+messageTokens > 2500; idx++ {
					//get token count of the message
					messageTokens = helpers.NumTokensFromMessages(chatFullHistory.Messages[idx:], currentConfig.Model)
				}
				//add the slice of messages to the messagesToSend
				messagesToSend = append(messagesToSend, chatFullHistory.Messages[idx:]...)
			}
		*/
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
			utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("ChatCompletion error: %v\n", err)
			continue
		}

		//get usage of the tokens
		jsonStr, _ := json.Marshal(assistant_response.Usage)
		utils.Logger.WithField("UUID", chatFullHistory.Id).Debugf("Tokens: %s", jsonStr)

		// get the generated response from OpenAI API
		content := assistant_response.Choices[0].Message.Content

		var jsonValid bool = true
		if !json.Valid([]byte(content)) {

			// Handle invalid JSON error
			utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("Invalid JSON response: %v\n", err)

			//jsonStr, _ := json.Marshal(chat.Messages)
			utils.Logger.WithField("UUID", chatFullHistory.Id).Debugf("content dump: %s", content)

			content, err = utils.FixJSON(content)
			if err != nil {
				utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("FixJSON error: %v\n", err)
				jsonValid = false
			}
			//check if valid after the fix
			if !json.Valid([]byte(content)) {
				utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("Invalid JSON response after fix: %v\n", err)
				jsonValid = false
			}
		}

		if jsonValid {
			// add the response to the list of messages
			chatFullHistory.Messages = append(chatFullHistory.Messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: content,
			})

			//add the response to the memory
			err = mem.AddMemory(content)
			if err != nil {
				utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("AddMemory error: %v\n", err)
			}

			//convert string to byte array
			contentByte := []byte(content)
			//unmarshal the byte array to struct
			var responseData models.Response
			err = json.Unmarshal(contentByte, &responseData)
			if err != nil {
				utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("==========================================================================================\n"+
					"Error unmarsheling response: %v\n,"+
					"Content Dump: %s\n"+
					"==========================================================================================\n", err, content)
				//try to unmarshel then to a map
				var dataMap map[string]interface{}
				if err = json.Unmarshal(contentByte, &dataMap); err != nil {
					// Handle any errors that may occur during parsing
					utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("Error unmarsheling to a map: %v\n", err)
					return
				}

				//output the map
				utils.Logger.WithField("UUID", chatFullHistory.Id).Info("==========================================================================================")
				utils.Logger.WithField("UUID", chatFullHistory.Id).Info("Map:")
				utils.Logger.WithField("UUID", chatFullHistory.Id).Info(dataMap)
				utils.Logger.WithField("UUID", chatFullHistory.Id).Info("==========================================================================================")

				continue
			}
			//output the response in human readible format
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
			fmt.Printf("- Keywords: %v\n", responseData.Thoughts.Keywords)

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
			// Output what to memorize
			fmt.Println("Memorize:")
			fmt.Printf("- Subject: %v\n", responseData.Memorize.Subject)
			fmt.Printf("- Information: %v\n", responseData.Memorize.Information)
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

				//Update the chat document in the database befoew exiting
				err = mongodb.ChatsCollection.Update(chatFullHistory)
				if err != nil {
					utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("on exit UpdateChat error: %v\n", err)
				}
				break
			} else if isY != "y" {
				//add user input
				chatFullHistory.Messages = append(chatFullHistory.Messages, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Human feedback: " + input,
				})
				if err != nil {
					utils.Logger.Errorf("Adding user input error: %v\n", err)
				}
			}

			//Update the chat document in the database
			err = mongodb.ChatsCollection.Update(chatFullHistory)
			if err != nil {
				utils.Logger.WithField("UUID", chatFullHistory.Id).Errorf("UpdateChat error: %v\n", err)
				continue
			}
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
