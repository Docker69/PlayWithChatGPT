package capabilities

import (
	"fmt"

	"backend/utils"

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
		description: "A capability that provides access to Selenium WebDriver.",
		version:     "1.0",
	}

	// Start a Selenium WebDriver session
	driver, err := selenium.NewRemote(selenium.Capabilities{}, url)
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

// Get returns the state of the capability.
func (c *SeleniumWebDriver) Get() interface{} {
	return c.driver
}

// Set sets the state of the capability.
func (c *SeleniumWebDriver) Set(value interface{}) error {
	var ok bool
	c.driver, ok = value.(selenium.WebDriver)
	if !ok {
		return fmt.Errorf("invalid type for Selenium WebDriver capability: expected selenium.WebDriver, but got %T", value)
	}
	return nil
}

// Run runs the capability.
func (c *SeleniumWebDriver) Run() error {
	// Navigate to Google using the WebDriver instance
	err := c.driver.Get("https://www.google.com")
	if err != nil {
		return fmt.Errorf("failed to navigate to google.com using the Selenium WebDriver: %v", err)
	}

	// Print the page title
	title, err := c.driver.Title()
	if err != nil {
		return fmt.Errorf("failed to get the page title using the Selenium WebDriver: %v", err)
	}

	utils.Logger.Debugf("Page title: %s\n", title)

	return nil
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
