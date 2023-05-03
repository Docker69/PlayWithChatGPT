package capabilities

import (
	"backend/ai/helpers"
	"backend/ai/memory"
	"backend/models"
	"backend/utils"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/joho/godotenv"
	"github.com/neurosnap/sentences/english"
	"github.com/sashabaranov/go-openai"
)

// BrowseWeb implements Capable interface and represents a colly web scrapper capability.
type BrowseWeb struct {
	// Fields for the individual capability
	name        string
	description string
	version     string

	//openai client
	MAX_TOKENS int
	client     *openai.Client
}

// NewBrowseWeb creates a new instance of BrowseWeb with the specified URL.
func NewBrowseWeb() (*BrowseWeb, error) {

	utils.Logger.Info("Init BrowseWeb capability")

	// load the environment variables
	err := godotenv.Load()
	if err != nil {
		utils.Logger.Infof(".env file not found, using OS ENV variables. Err: %s", err)
	}

	// extract and save the OpenAI api key from environment variables
	exists := false
	apiKey, exists := os.LookupEnv("OPENAI_API_KEY")

	if !exists {
		utils.Logger.Panic("OpenAI API Key not found, panicking!!!")
	}

	// read the MAX_TOKENS from environment variables as integer
	exists = false
	max_tokens_str, exists := os.LookupEnv("MAX_TOKENS")
	if !exists {
		utils.Logger.Error("MAX_TOKENS not found, setting to default value 4000")
		max_tokens_str = "4000"
	}
	// convert the max_tokens string to MAX_TOKENS int
	max_tokens, err := strconv.Atoi(max_tokens_str)
	if err != nil {
		utils.Logger.Error("Error converting MAX_TOKENS to int, setting to default value 4000")
		max_tokens = 4000
	}

	// create new client instance with given apiKey
	client := openai.NewClientWithConfig(openai.DefaultConfig(apiKey))
	//client := openai.NewClient(apiKey)

	//check that client  is not nil
	if client == nil {
		utils.Logger.Panic("OpenAI Client is nil, panicking!!!")
	}

	capability := &BrowseWeb{
		name:        "browse_web",
		description: "A capability that scrapes and memorize the web site text.",
		version:     "1.0",
		MAX_TOKENS:  max_tokens,
		client:      client,
	}

	return capability, nil
}

// Name returns the name of the capability.
func (c *BrowseWeb) Name() string {
	return c.name
}

// Description returns a short description of what the capability does.
func (c *BrowseWeb) Description() string {
	return c.description
}

// Version returns the version number of the capability.
func (c *BrowseWeb) Version() string {
	return c.version
}

// Run runs the capability.
func (c *BrowseWeb) Run(mem *memory.MemoryCache, args ...interface{}) (interface{}, error) {
	var ok bool

	//minimum 2 args
	if len(args) < 2 {
		return nil, fmt.Errorf("BrowseWeb: at least one argument is required")
	}

	url, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("BrowseWeb: url must be a string")
	}
	if url == "" {
		return nil, fmt.Errorf("BrowseWeb: input is empty")
	}

	quetion, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("BrowseWeb: quetion must be a string")
	}
	if quetion == "" {
		return nil, fmt.Errorf("BrowseWeb: quetion is empty")
	}

	//scrape the web site
	texts, err := c.scrape(url)
	if err != nil {
		return nil, fmt.Errorf("BrowseWeb, error scraping the web site: %v", err)
	}

	//summarize the web site content
	text, err := c.summarize(texts, quetion, mem)
	if err != nil {
		return nil, fmt.Errorf("BrowseWeb, error summarizing the web site: %v", err)
	}

	//scrape the web site
	links, err := c.links(url)
	if err != nil {
		return nil, fmt.Errorf("BrowseWeb, error getting links from the web site: %v", err)
	}

	resultJson := map[string]interface{}{
		"summary": text,
		"links":   links,
	}

	return resultJson, nil
}

// Stop stops the capability.
func (c *BrowseWeb) Stop() error {
	return nil
}

// scrapes the web site and memorizes the text into the memory
func (c *BrowseWeb) scrape(url string) ([]string, error) {

	//scrape the web site
	text := []string{}
	coll := colly.NewCollector()
	if coll == nil {
		return []string{}, fmt.Errorf("BrowseWeb: collector is nil")
	}
	coll.SetRequestTimeout(120 * time.Second)

	coll.OnHTML("body", func(e *colly.HTMLElement) {
		//text = e.Text
		text = e.ChildTexts("p, li, h1, h2, h3, h4, h5, h6")
	})

	coll.OnHTML("a[href]", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
	})

	err := coll.Visit(url)

	if err != nil {
		return []string{}, fmt.Errorf("BrowseWeb: error visiting the web site: %v", err)
	}

	return text, nil
}

// scrapes the web site for links
func (c *BrowseWeb) summarize(texts []string, quetion string, mem *memory.MemoryCache) (string, error) {

	text := ""

	for _, line := range texts {
		line = strings.ReplaceAll(line, "\n", " ")
		line = strings.ReplaceAll(line, "\t", " ")
		if strings.Count(line, " ") >= 3 {
			text += line + " "
		}
	}
	//text = strings.Join(filteredText, " ")
	//text = strings.ReplaceAll(text, "\n", " ")
	//text = strings.ReplaceAll(text, "\t", " ")

	tokenizer, err := english.NewSentenceTokenizer(nil)
	if err != nil {
		panic(err)
	}

	sentences := tokenizer.Tokenize(text)
	if sentences == nil {
		return "", fmt.Errorf("BrowseWeb: error tokenizing the text")
	}

	flattenedText := ""
	for _, s := range sentences {
		flattenedText += s.Text
	}

	flattenedText += "\n"

	messagesToSend := []openai.ChatCompletionMessage{}
	messagesToSend = append(messagesToSend, summaryMessage(flattenedText, quetion))

	// get the final token count from the chat messages
	numTokens := helpers.NumTokensFromMessages(messagesToSend, models.NewOpenAIConfig().Model)

	//TODO read token limit from .env file
	allowedTokens := c.MAX_TOKENS - numTokens

	if allowedTokens < 100 {
		return "", fmt.Errorf("BrowseWeb: allowed tokens is less than 100")
	}

	// Call OpenAI API to generate response to the user's message
	response, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     models.NewOpenAIConfig().Model,
			Messages:  messagesToSend,
			MaxTokens: allowedTokens,
		},
	)

	if err != nil {
		return "", fmt.Errorf("BrowseWeb: error calling OpenAI API: %v", err)
	}

	// get the generated response from OpenAI API
	content := response.Choices[0].Message.Content

	//add the response to the memory
	err = (*mem).AddMemory(content)
	if err != nil {
		return "", fmt.Errorf("BrowseWeb: error adding memory: %v", err)
	}
	return "", nil
}

// scrapes the web site for links
func (c *BrowseWeb) links(url string) (string, error) {

	return "", nil
}

func summaryMessage(content string, question string) openai.ChatCompletionMessage {

	return openai.ChatCompletionMessage{
		Role: "user",
		Content: `"""` + content + `""" \n` + `Using the above text, answer the following` +
			` question: "` + question + `" -- if the question cannot be answered using the text,` +
			" summarize the text.",
	}
}
