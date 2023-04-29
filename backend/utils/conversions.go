package utils

import (
	"bytes"
	"encoding/binary"
)

func Float32ToBytes(input []float32) ([]byte, error) {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, input)
	if err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}
