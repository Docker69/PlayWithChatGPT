package test_utils

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"

	"backend/utils"

	testdataloader "github.com/peteole/testdata-loader"
)

func prepareData() ([]float32, []byte) {
	// Read test data from file
	byteData := testdataloader.GetTestFile("utils/test/test_data.json")

	// Check if byteData is empty
	if len(byteData) == 0 {
		log.Fatal("Test data file is empty")
	}

	// Convert byte array to string
	jsonString := string(byteData)

	// Define a variable to store the parsed JSON data as a map
	var dataMap map[string]interface{}

	// Use the Unmarshal() function from the json package to parse the JSON data into the dataMap variable
	if err := json.Unmarshal([]byte(jsonString), &dataMap); err != nil {
		// Handle any errors that may occur during parsing
		log.Fatal(err)
	}

	// Convert dataMap["input"] to a []float32
	inputIF, ok := dataMap["input"].([]interface{})
	if !ok {
		log.Fatal("Invalid input data type")
	}
	input := make([]float32, len(inputIF))
	for i, v := range inputIF {
		f, ok := v.(float64)
		if !ok {
			log.Fatal("Invalid input element type")
		}
		input[i] = float32(f)
	}

	// Convert dataMap["expectedOutput"] to a []byte
	expectedOutputIF, ok := dataMap["expectedOutput"].([]interface{})
	if !ok {
		log.Fatal("Invalid expected output data type")
	}
	expectedOutput := make([]byte, len(expectedOutputIF))
	for i, v := range expectedOutputIF {
		b, ok := v.(float64)
		if !ok {
			log.Fatal("Invalid input element type")
		}
		expectedOutput[i] = byte(b)
	}

	return input, expectedOutput
}

// Ensure each of the encoding functions produces the same output
func TestEquality(t *testing.T) {

	input, expectedOutput := prepareData()

	outputFloat32ToBytes, err := utils.Float32ToBytes(input)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !bytes.Equal(outputFloat32ToBytes, expectedOutput) {
		t.Errorf("Float32ToBytes(%v) = %v, expected %v", input, outputFloat32ToBytes, expectedOutput)
	}

	outputFloat32ToBytesFastSafe := utils.Float32ToBytesFastSafe(input)

	if !bytes.Equal(outputFloat32ToBytesFastSafe, expectedOutput) {
		t.Errorf("Float32ToBytesFastSafe(%v) = %v, expected %v", input, outputFloat32ToBytesFastSafe, expectedOutput)
	}

	outputFloat32ToBytesAlt, err := utils.Float32ToBytesAlt(input)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !bytes.Equal(outputFloat32ToBytesAlt, expectedOutput) {
		t.Errorf("encodeFloat32ToBytesAlt(%v) = %v, expected %v", input, outputFloat32ToBytesAlt, expectedOutput)
	}

	outputFloat32ToBytesFastUnsafe := utils.Float32ToBytesFastUnsafe(input)

	if !bytes.Equal(outputFloat32ToBytesFastUnsafe, expectedOutput) {
		t.Errorf("Float32ToBytesUnsafe(%v) = %v, expected %v", input, outputFloat32ToBytesFastUnsafe, expectedOutput)
	}
}

func BenchmarkAllFloat32ToBytes(b *testing.B) {
	input, _ := prepareData()

	b.Run("utils.Float32ToBytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			utils.Float32ToBytes(input)
		}
	})

	b.Run("utils.Float32ToBytesFastSafe", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			utils.Float32ToBytesFastSafe(input)
		}
	})

	b.Run("utils.Float32ToBytesAlt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			utils.Float32ToBytesAlt(input)
		}
	})

	b.Run("utils.Float32ToBytesFastUnsafe", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			utils.Float32ToBytesFastUnsafe(input)
		}
	})
}
