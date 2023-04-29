package capabilities

// Capable defines the interface that each capability must implement.
type Capable interface {
	// Name returns the name of the capability.
	Name() string
	// Description returns a short description of what the capability does.
	Description() string
	// Version returns the version number of the capability.
	Version() string
	// Get returns the state of the capability.
	Get() interface{}
	// Set sets the state of the capability.
	Set(interface{}) error
	// Run runs the capability.
	Run() error
	// Stop stops the capability.
	Stop() error
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
	return f.capabilitiesMap[name]
}

// Remove removes a capability with the given name from the factory.
func (f *CapabilityFactory) Remove(name string) error {
	delete(f.capabilitiesMap, name)
	return nil
}

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
	capabilityFactory.Add(&SeleniumWebDriver{})
}
