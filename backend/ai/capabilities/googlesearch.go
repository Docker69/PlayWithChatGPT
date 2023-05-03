package capabilities

import (
	"context"
	"fmt"

	"backend/ai/memory"
	"backend/models"
	"backend/utils"

	googlesearch "github.com/rocketlaunchr/google-search"
)

type GoogleSearch struct{}

func (gs *GoogleSearch) Name() string {
	return "google"
}

func (gs *GoogleSearch) Description() string {
	return "Google Command results: "
}

func (gs *GoogleSearch) Version() string {
	return "1.0"
}

func (gs *GoogleSearch) Run(mem *memory.MemoryCache, args ...interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("GoogleSearch: at least one argument is required")
	}
	command, ok := args[0].(models.ArgsType)
	if !ok {
		return nil, fmt.Errorf("GoogleSearch: input must be a CommandType")
	}
	if command.Input == "" {
		return nil, fmt.Errorf("GoogleSearch: input is empty")
	}

	utils.Logger.Debugf("Performing Google search for '%s'", command.Input)
	results, err := googlesearch.Search(context.Background(), command.Input)
	if err != nil {
		return nil, fmt.Errorf("GoogleSearch: failed to perform Google search: %v", err)
	}

	//convert results to a JSON array
	var jsonResults []map[string]interface{} = make([]map[string]interface{}, 0)
	for _, result := range results {
		// make result a map so it can be converted to JSON
		resultJson := map[string]interface{}{
			"rank":  result.Rank,
			"url":   result.URL,
			"title": result.Title,
			"desc":  result.Description,
		}

		jsonResults = append(jsonResults, resultJson)
	}

	return jsonResults, nil
}

func (gs *GoogleSearch) Stop() error {
	// noop since there are no background processes to stop
	return nil
}
