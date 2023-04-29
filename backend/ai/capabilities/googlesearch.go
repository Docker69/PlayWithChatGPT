package capabilities

import (
	"context"
	"fmt"

	"backend/utils"

	googlesearch "github.com/rocketlaunchr/google-search"
)

type GoogleSearch struct{}

func (gs *GoogleSearch) Name() string {
	return "Google Search"
}

func (gs *GoogleSearch) Description() string {
	return "Performs a Google search using the google-search library"
}

func (gs *GoogleSearch) Version() string {
	return "1.0"
}

func (gs *GoogleSearch) Get() interface{} {
	// return nil since we don't have any state to get
	return nil
}

func (gs *GoogleSearch) Set(interface{}) error {
	// noop since we don't have any state to set
	return nil
}

func (gs *GoogleSearch) Run() error {
	// Perform a Google search for the term "golang"
	utils.Logger.Debug("Performing Google search for 'golang'")
	results, err := googlesearch.Search(context.Background(), "golang")
	if err != nil {
		return fmt.Errorf("failed to perform Google search: %v", err)
	}
	// Print the top 5 search results
	for i, result := range results[:5] {
		fmt.Printf("%d. %s\n", i+1, result.Title)
		fmt.Println(result.URL)
		fmt.Println(result.Description)
		fmt.Println()
	}
	return nil
}

func (gs *GoogleSearch) Stop() error {
	// noop since there are no background processes to stop
	return nil
}
