// local_storage_mem.go
package memory

import (
	"encoding/json"
	"fmt"
	"os"
)

type LocalStorageMem struct {
	data     map[string]string
	filePath string
}

func NewLocalStorageMem(filePath string) (*LocalStorageMem, error) {
	ls := &LocalStorageMem{data: make(map[string]string), filePath: filePath}

	if _, err := os.Stat(ls.filePath); os.IsNotExist(err) {
		return ls, nil
	}

	err := ls.loadFromFile()

	if err != nil {
		return nil, fmt.Errorf("error loading data from file: %v", err)
	}

	return ls, nil
}

func (l *LocalStorageMem) Get(key string) string {
	return l.data[key]
}

func (l *LocalStorageMem) AddMemory(value string) error {
	l.saveToFile()
	return nil
}

func (l *LocalStorageMem) Clear() error {
	l.data = make(map[string]string)
	l.saveToFile()
	return nil
}

func (l *LocalStorageMem) GetRelevantMemories(query string) []string {
	result := make([]string, 0)
	for k, v := range l.data {
		if v == query {
			result = append(result, k)
		}
	}
	return result
}

func (l *LocalStorageMem) GetStats() int {
	return len(l.data)
}

func (l *LocalStorageMem) saveToFile() error {
	jsonData, err := json.Marshal(l.data)
	if err != nil {
		return fmt.Errorf("error marshalling data to json: %v", err)
	}

	err = os.WriteFile(l.filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing data to file: %v", err)
	}

	return nil
}

func (l *LocalStorageMem) loadFromFile() error {
	jsonData, err := os.ReadFile(l.filePath)
	if err != nil {
		return fmt.Errorf("error reading data from file: %v", err)
	}

	err = json.Unmarshal(jsonData, &l.data)
	if err != nil {
		return fmt.Errorf("error unmarshalling data from json: %v", err)
	}

	return nil
}
