package capabilities

import (
	"backend/ai/memory"
	"backend/models"
	"backend/utils"
	"encoding/json"
	"fmt"
)

// Capable defines the interface that each capability must implement.
type Capable interface {
	// Name returns the name of the capability.
	Name() string
	// Description returns a short description of what the capability does.
	Description() string
	// Version returns the version number of the capability.
	Version() string
	// Run runs the capability.
	Run(mem *memory.MemoryCache, args ...interface{}) (interface{}, error)
	// Stop stops the capability.
	Stop() error
}

type CommandStruct struct {
	Name string
	Args interface{}
}

// CapabilityFactory stores and manages all the registered capabilities.
type CapabilityFactory struct {
	capabilitiesMap map[string]Capable // A map of capabilities keyed by their names
}

// Add adds a new capability to the factory.
func (f *CapabilityFactory) Add(capability Capable) error {
	f.capabilitiesMap[capability.Name()] = capability
	return nil
}

// List returns a slice of all the registered capabilities.
func (f *CapabilityFactory) List() []Capable {
	capabilitiesSlice := make([]Capable, 0)
	for _, value := range f.capabilitiesMap {
		capabilitiesSlice = append(capabilitiesSlice, value)
	}
	return capabilitiesSlice
}

// Get returns a capability with the given name.
func (f *CapabilityFactory) Get(name string) Capable {
	//check if capability exists
	return f.capabilitiesMap[name]
}

// Remove removes a capability with the given name from the factory.
func (f *CapabilityFactory) Remove(name string) error {
	delete(f.capabilitiesMap, name)
	return nil
}

// Execute executes a command from the factory, and returns the results as JSON
func (f *CapabilityFactory) Execute(Command models.CommandType, mem *memory.MemoryCache) (string, error) {

	//nothing to do
	nothingToDo := map[string]bool{
		"do_nothing": true,
		"user_input": true,
	}
	//check if Command is in nothing array
	if nothingToDo[Command.Name] {
		return "", nil
	}

	// Get the capability.
	capability := f.Get(Command.Name)

	// Check if capability is nil.
	if capability == nil {
		return "", fmt.Errorf("capability not found: %s", Command.Name)
	}

	// Run the capability and get results.
	results, err := capability.Run(mem, Command.Args)
	if err != nil {
		return "", err
	}

	// Build json response.
	jsonResponse, err := json.Marshal(results)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("{\"%s\": %s}", capability.Name(), string(jsonResponse))
	return result, nil
}

// run capability from CapabilityFactory based on command input and arguments (interface) passed in

// GetCapabilityFactory returns a pointer to a new instance of the CapabilityFactory type.
func GetCapabilityFactory() *CapabilityFactory {
	return capabilityFactory
}

var capabilityFactory *CapabilityFactory

func init() {
	capabilityFactory = &CapabilityFactory{
		capabilitiesMap: make(map[string]Capable),
	}

	//Add the google search capability
	capabilityFactory.Add(&GoogleSearch{})
	//Add the selenium web driver capability
	//capabilityFactory.Add(&SeleniumWebDriver{})
	//Add the web browse
	browseWeb, err := NewBrowseWeb()
	if err != nil {
		// handle error
		utils.Logger.Errorf("CapabilityFactory: error creating browes web capability %v", err)
	}
	capabilityFactory.Add(browseWeb)
}
