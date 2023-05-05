package tests_test

import (
	"backend/ai/capabilities"
	"backend/models"
	"testing"
)

// Then, we can write our tests
func TestWriteToFile(t *testing.T) {

	write_capability := capabilities.GetCapabilityFactory().Get("write_to_file").(*capabilities.WriteToFile)
	if write_capability == nil {
		t.Error("Unexpected error: can't get capability")
		t.FailNow()
	}

	args := models.ArgsType{
		Path: "/home/bennyk/projects/PlayWithChatGPT/backend/ai_workspace",
		File: "SciFi_recommendations.txt",
		Text: "Hello World",
	}

	_, err := write_capability.Run(nil, args)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		t.FailNow()
	}

	args = models.ArgsType{
		Path: "/home/bennyk/projects/PlayWithChatGPT/backend/ai_workspace",
		File: "SciFi_recommendations.txt",
	}

	read_capability := capabilities.GetCapabilityFactory().Get("read_file").(*capabilities.ReadFile)

	if read_capability == nil {
		t.Error("Unexpected error: can't get capability")
		t.FailNow()
	}

	text, err := read_capability.Run(nil, args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		t.FailNow()
	}

	//convert texts to string array and test if the first element is "Hello World"
	if text.(string) != "Hello World" {
		t.Errorf("Unexpected error: %v", err)
		t.FailNow()
	}
}
