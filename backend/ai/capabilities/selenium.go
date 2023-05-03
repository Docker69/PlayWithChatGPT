package capabilities

import (
	"backend/ai/memory"
	"fmt"

	"github.com/tebeka/selenium"
)

// SeleniumWebDriver implements Capable interface and represents a Selenium WebDriver capability.
type SeleniumWebDriver struct {
	// Fields for the individual capability
	name        string
	description string
	version     string

	// Fields specific to Selenium WebDriver
	driver selenium.WebDriver
}

// NewSeleniumWebDriver creates a new instance of SeleniumWebDriver with the specified URL.
func NewSeleniumWebDriver(url string) (*SeleniumWebDriver, error) {
	capability := &SeleniumWebDriver{
		name:        "selenium-webdriver",
		description: "A capability that scrapes and memorize the web site text.",
		version:     "1.0",
	}

	driverCapabilities := selenium.Capabilities{
		"browserName": "chrome",
	}
	// Start a Selenium WebDriver session
	driver, err := selenium.NewRemote(driverCapabilities, url)
	if err != nil {
		return nil, fmt.Errorf("failed to start a new Selenium WebDriver session: %v", err)
	}
	capability.driver = driver

	return capability, nil
}

// Name returns the name of the capability.
func (c *SeleniumWebDriver) Name() string {
	return c.name
}

// Description returns a short description of what the capability does.
func (c *SeleniumWebDriver) Description() string {
	return c.description
}

// Version returns the version number of the capability.
func (c *SeleniumWebDriver) Version() string {
	return c.version
}

// Run runs the capability.
func (c *SeleniumWebDriver) Run(mem *memory.MemoryCache, args ...interface{}) (interface{}, error) {
	var ok bool

	//minimum 2 args
	if len(args) < 2 {
		return nil, fmt.Errorf("SeleniumWebDriver: at least one argument is required")
	}

	url, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("SeleniumWebDriver: url must be a string")
	}
	if url == "" {
		return nil, fmt.Errorf("SeleniumWebDriver: input is empty")
	}

	quetion, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("SeleniumWebDriver: quetion must be a string")
	}
	if quetion == "" {
		return nil, fmt.Errorf("SeleniumWebDriver: quetion is empty")
	}

	//scrape the web site
	text, err := c.scrape(url)
	if err != nil {
		return nil, fmt.Errorf("SeleniumWebDriver, error scraping the web site: %v", err)
	}

	//scrape the web site
	links, err := c.links(url)
	if err != nil {
		return nil, fmt.Errorf("SeleniumWebDriver, error getting links from the web site: %v", err)
	}

	resultJson := map[string]interface{}{
		"summary": text,
		"links":   links,
	}

	return resultJson, nil
}

// Stop stops the capability.
func (c *SeleniumWebDriver) Stop() error {
	// Close the WebDriver session and quit the driver
	err := c.driver.Quit()
	if err != nil {
		return fmt.Errorf("failed to quit the Selenium WebDriver session: %v", err)
	}
	return nil
}

// scrapes the web site and memorizes the text into the memory
func (c *SeleniumWebDriver) scrape(url string) (string, error) {

	return "", nil
}

// scrapes the web site and memorizes the text into the memory
func (c *SeleniumWebDriver) links(url string) (string, error) {

	return "", nil
}
