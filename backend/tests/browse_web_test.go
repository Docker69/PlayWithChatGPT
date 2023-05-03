package tests_test

import (
	"backend/ai"
	"backend/ai/capabilities"
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

	capability := capabilities.GetCapabilityFactory().Get("browse_web").(*capabilities.BrowseWeb)
	if capability == nil {
		t.Error("Unexpected error: can't get capability")
		t.FailNow()
	}

	err := ai.InitMemory()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Call the scrape function with a fake URL
	//texts, err := capability.Run(&ai.Mem, "https://niyasoft.com/", "What niya soft does?")
	texts, err := capability.Run(&ai.Mem, "https://multi-programming.com/game-development/white-label-solution", "what is white label solution?")
	//texts, err := capability.Run(&ai.Mem, "https://blog.logrocket.com/building-web-scraper-go-colly/", "what is web scraping?")
	//texts, err := capability.Run(&ai.Mem, "https://en.wikipedia.org/wiki/Rome", "What is Rome?")
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
