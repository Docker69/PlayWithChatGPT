package tests_test

import (
	"backend/ai"
	"backend/ai/capabilities"
	"backend/models"
	"encoding/json"
	"os"
	"testing"

	"github.com/gocolly/colly/v2"
)

// First, we define a mock Collector to use in our tests
type MockCollector struct {
	OnHTMLFunc func(goquerySelector string, f colly.HTMLCallback)
}

func (mc *MockCollector) OnHTML(selector string, cb colly.HTMLCallback) {
	// Call the OnHTMLFunc callback with the given selector and callback function
	mc.OnHTMLFunc(selector, cb)
}

// Then, we can write our tests
func TestScrape(t *testing.T) {

	capability := capabilities.GetCapabilityFactory().Get("browse_website").(*capabilities.BrowseWeb)
	if capability == nil {
		t.Error("Unexpected error: can't get capability")
		t.FailNow()
	}

	err := ai.InitMemory()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	//texts, err := capability.Run(&ai.Mem, "https://niyasoft.com/", "What niya soft does?")
	//texts, err := capability.Run(&ai.Mem, "https://multi-programming.com/game-development/white-label-solution", "what is white label solution?")
	//texts, err := capability.Run(&ai.Mem, "https://blog.logrocket.com/building-web-scraper-go-colly", "what is web scraping?")
	//texts, err := capability.Run(&ai.Mem, "https://en.wikipedia.org/wiki/Rome", "What is Rome?")
	//texts, err := capability.Run(&ai.Mem, "https://www.washingtonpost.com/books/2022/11/17/best-sci-fi-fantasy", "What are the top 5 books of 2022?")

	//params := []string{"https://niyasoft.com/", "What niya soft does?"}
	//params := []string{"https://multi-programming.com/game-development/white-label-solution", "what is white label solution?"}
	//params := []string{"https://blog.logrocket.com/building-web-scraper-go-colly", "what is web scraping?"}
	//parmas := []string{"https://en.wikipedia.org/wiki/Rome", "What is Rome?"}
	params := []string{"https://www.washingtonpost.com/books/2022/11/17/best-sci-fi-fantasy", "What are the top 5 books of 2022?"}

	args := models.ArgsType{
		URL:      params[0],
		Question: params[1],
		Input:    "",
		Reason:   "",
	}

	// Call the scrape function with a fake URL
	command := models.CommandType{
		Name: "browse_web",
		Args: args,
	}

	texts, err := capability.Run(&ai.Mem, command)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		t.FailNow()
	}
	// Build json response.
	jsonResponse, err := json.Marshal(texts)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	f, err := os.Create("test-results-scrape.txt")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	defer f.Close()

	_, err = f.Write(jsonResponse)
	if err != nil {
		t.Error("Unexpected error: can't write to file")
	}

	// Flush the buffer to ensure all the data is written to the file
	err = f.Sync()
	if err != nil {
		t.Error("Unexpected error: can't flush buffer")
	}

}
