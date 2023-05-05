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
	"github.com/gocolly/colly/v2/extensions"
	"github.com/joho/godotenv"
	"github.com/neurosnap/sentences/english"
	"github.com/sashabaranov/go-openai"
)

// BrowseWeb implements Capable interface and represents a colly web scrapper capability.
type BrowseWeb struct {
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
		MAX_TOKENS: max_tokens,
		client:     client,
	}

	return capability, nil
}

// Name returns the name of the capability.
func (c *BrowseWeb) Name() string {
	return "browse_website"
}

// Description returns a short description of what the capability does.
func (c *BrowseWeb) Description() string {
	return "A capability that scrapes and memorize the web site text."
}

// Version returns the version number of the capability.
func (c *BrowseWeb) Version() string {
	return "1.0"
}

// Stop stops the capability.
func (c *BrowseWeb) Stop() error {
	return nil
}

// Run runs the capability.

func (c *BrowseWeb) Run(mem *memory.MemoryCache, args ...interface{}) (interface{}, error) {

	var ok bool

	//minimum 2 args
	if len(args) < 1 {
		return nil, fmt.Errorf("BrowseWeb: at least one argument is required")
	}
	command, ok := args[0].(models.ArgsType)
	if !ok {
		return nil, fmt.Errorf("BrowseWeb: input must be a CommandType")
	}

	if command.URL == "" {
		return nil, fmt.Errorf("BrowseWeb: url is empty")
	}
	if command.Question == "" {
		return nil, fmt.Errorf("BrowseWeb: quetion is empty")
	}

	//scrape the web site
	texts, links, err := c.scrape(command.URL)
	if err != nil {
		return nil, fmt.Errorf("BrowseWeb, error scraping the web site: %v", err)
	}

	//summarize the web site content
	text, err := c.summarize(texts, command.URL, command.Question, mem)
	if err != nil {
		return nil, fmt.Errorf("BrowseWeb, error summarizing the web site: %v", err)
	}

	//no more than 5 links
	if len(links) > 5 {
		links = links[:5]
	}

	resultJson := map[string]interface{}{
		"Answer gathered from website: ": text,
		"Links":                          links,
	}

	return resultJson, nil
}

// scrapes the web site and memorizes the text into the memory
func (c *BrowseWeb) scrape(url string) ([]string, []string, error) {

	//scrape the web site
	text := []string{}
	links := []string{}
	coll := colly.NewCollector()
	//use random user agent
	extensions.RandomUserAgent(coll)

	if coll == nil {
		return []string{}, []string{}, fmt.Errorf("BrowseWeb: collector is nil")
	}
	coll.SetRequestTimeout(120 * time.Second)

	coll.OnHTML("body", func(e *colly.HTMLElement) {
		//text = e.Text
		//text = e.ChildTexts("p, li, h1, h2, h3, h4, h5, h6")
		text = e.ChildTexts("p, li, span")
		hrefs := e.ChildAttrs("a[href]", "href")
		linkTexts := e.ChildTexts("a[href]")
		if len(linkTexts) != len(hrefs) {
			utils.Logger.Error("BrowseWeb: links hrefs not equal to link texts")
			return
		}

		//find if relative link or absolute
		for i, link := range hrefs {
			//if length is greater than 1 and starts with / then it is a relative link
			if strings.HasPrefix(link, "/") && len(link) > 1 {
				links = append(links, "\"text\":\""+linkTexts[i]+"\", \"link\":\""+url+link+"\"")
			}
		}
	})

	err := coll.Visit(url)

	if err != nil {
		return []string{}, []string{}, fmt.Errorf("BrowseWeb: error visiting the web site: %v", err)
	}

	return text, links, nil
}

// scrapes the web site for links
func (c *BrowseWeb) summarize(texts []string, url string, question string, mem *memory.MemoryCache) (string, error) {

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

	//init variables
	summaries := []string{}
	messagesToSend := []openai.ChatCompletionMessage{}
	numTokens := 0
	flattenedText := ""

	utils.Logger.Info("BrowseWeb: Start summarizing the text")
	for _, s := range sentences {

		flattenedText += s.Text

		messagesToSend = []openai.ChatCompletionMessage{summaryMessage(flattenedText+"\n", question)}

		// get the final token count from the chat messages
		numTokens = helpers.NumTokensFromMessages(messagesToSend, models.NewOpenAIConfig().Model)

		if numTokens > 3000 {
			idx := len(summaries)
			header := fmt.Sprintf("Source: %s\nContent summary part#%d: ", url, idx)

			summary, err := c.sendSummaryMessage(numTokens, header, messagesToSend, mem)
			if err != nil {
				return "", fmt.Errorf("BrowseWeb: error sending summary message: %v", err)
			}
			summaries = append(summaries, summary)
			flattenedText = ""
			utils.Logger.Infof("BrowseWeb: Summarizied part #%d", idx)
		}
	}

	//there are always at least some sentences left
	//handle edge case where the last sentence actually brought numTokens > 3000, if this is the case then flattenedText == ""
	if flattenedText != "" {
		idx := len(summaries)
		header := fmt.Sprintf("Source: %s\nContent summary part#%d: ", url, idx)
		summary, err := c.sendSummaryMessage(numTokens, header, messagesToSend, mem)
		if err != nil {
			return "", fmt.Errorf("BrowseWeb: error sending summary message: %v", err)
		}
		summaries = append(summaries, summary)
		utils.Logger.Infof("BrowseWeb: Summarizied part #%d", idx)
	}

	utils.Logger.Info("BrowseWeb: Summarized all individual parts")

	if len(summaries) == 1 {
		return summaries[0], nil
	}

	// join the summaries to one string and request overall summary
	messagesToSend = []openai.ChatCompletionMessage{summaryMessage(strings.Join(summaries, "\n")+"\n", question)}
	numTokens = helpers.NumTokensFromMessages(messagesToSend, models.NewOpenAIConfig().Model)
	summary, err := c.sendSummaryMessage(numTokens, "", messagesToSend, mem)
	if err != nil {
		return "", fmt.Errorf("BrowseWeb: error sending summary message: %v", err)
	}

	utils.Logger.Info("BrowseWeb: Summarized over all summaries")

	return summary, nil
}

func (c *BrowseWeb) sendSummaryMessage(numTokens int, header string, messagesToSend []openai.ChatCompletionMessage, mem *memory.MemoryCache) (string, error) {

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

	//"Source: %s\nContent summary part#%d: %s"
	// get the generated response from OpenAI API
	summaryContent := header + response.Choices[0].Message.Content
	//add the response to the memory
	err = (*mem).AddMemory(summaryContent)
	if err != nil {
		return "", fmt.Errorf("BrowseWeb: error adding memory: %v", err)
	}

	return summaryContent, nil
}

func summaryMessage(content string, question string) openai.ChatCompletionMessage {

	return openai.ChatCompletionMessage{
		Role: "user",
		Content: `"""` + content + `""" \n` + `Using the above text, answer the following` +
			` question: "` + question + `" -- if the question cannot be answered using the text,` +
			" summarize the text.",
	}
}
