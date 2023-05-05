package capabilities

import (
	"backend/ai/memory"
	"backend/models"
	"fmt"
	"os"
)

type WriteToFile struct{}

func (wtf *WriteToFile) Name() string {
	return "write_to_file"
}

func (wtf *WriteToFile) Description() string {
	return "Write to file Command results: "
}

func (wtf *WriteToFile) Version() string {
	return "1.0"
}

func (wtf *WriteToFile) Run(mem *memory.MemoryCache, args ...interface{}) (interface{}, error) {
	// check if at least one argument was passed
	if len(args) == 0 {
		return nil, fmt.Errorf("WriteToFile: at least one argument is required")
	}

	// check if the first argument is of the expected type
	fileArgs, ok := args[0].(models.ArgsType)
	if !ok {
		return nil, fmt.Errorf("WriteToFile: input must be an ArgsType")
	}

	// check if file name is empty
	if fileArgs.File == "" {
		return nil, fmt.Errorf("WriteToFile: file name is empty")
	}

	// check if text is empty
	if fileArgs.Text == "" {
		return nil, fmt.Errorf("WriteToFile: text is empty")
	}

	// check if path is empty
	if fileArgs.Path == "" {
		return nil, fmt.Errorf("WriteToFile: path is empty")
	}

	// create/open the file with write-only mode
	f, err := os.OpenFile(fileArgs.Path+"/"+fileArgs.File, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return nil, fmt.Errorf("WriteToFile: file creation error: %v", err)
	}
	defer f.Close()

	_, err = f.WriteString(fileArgs.Text)
	if err != nil {
		return nil, fmt.Errorf("WriteToFile: file writing error: %v", err)
	}

	// flush the buffer to ensure all the data is written to the file
	err = f.Sync()
	if err != nil {
		return nil, fmt.Errorf("WriteToFile: file flush error: %v", err)
	}

	return "File written successfully", nil
}

func (wtf *WriteToFile) Stop() error {
	// noop since there are no background processes to stop
	return nil
}
