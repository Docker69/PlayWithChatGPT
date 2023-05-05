package capabilities

import (
	"backend/ai/memory"
	"backend/models"
	"fmt"
	"os"
)

type ReadFile struct{}

func (rf *ReadFile) Name() string {
	return "read_file"
}

func (rf *ReadFile) Description() string {
	return "Read from file Command results: "
}

func (rf *ReadFile) Version() string {
	return "1.0"
}

func (rf *ReadFile) Run(mem *memory.MemoryCache, args ...interface{}) (interface{}, error) {
	// check if at least one argument was passed
	if len(args) == 0 {
		return nil, fmt.Errorf("ReadFile: at least one argument is required")
	}

	// check if the first argument is of the expected type
	fileArgs, ok := args[0].(models.ArgsType)
	if !ok {
		return nil, fmt.Errorf("ReadFile: input must be an ArgsType")
	}

	// check if file name is empty
	if fileArgs.File == "" {
		return nil, fmt.Errorf("ReadFile: file name is empty")
	}

	// check if path is empty
	if fileArgs.Path == "" {
		return nil, fmt.Errorf("ReadFile: path is empty")
	}

	// open the file with read-only mode
	f, err := os.OpenFile(fileArgs.Path+"/"+fileArgs.File, os.O_RDONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("ReadFile: file opening error: %v", err)
	}
	defer f.Close()

	// get file size and allocate byte slice accordingly
	fileStat, _ := f.Stat()
	bytes := make([]byte, fileStat.Size())

	// read the entire file contents
	_, err = f.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("ReadFile: file reading error: %v", err)
	}

	// check if file content is empty
	if len(bytes) == 0 {
		return nil, fmt.Errorf("ReadFile: file content is empty")
	}

	//return the text
	return string(bytes), nil
}

func (rf *ReadFile) Stop() error {
	// noop since there are no background processes to stop
	return nil
}
