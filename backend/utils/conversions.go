package utils

import (
	"bytes"
	"encoding/binary"
	"math"
	"unsafe"
)

func Float32ToBytes(input []float32) ([]byte, error) {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, input)
	if err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

func Float32ToBytesFastSafe(fs []float32) []byte {
	buf := make([]byte, len(fs)*4)
	for i, f := range fs {
		u := math.Float32bits(f)
		binary.LittleEndian.PutUint32(buf[i*4:], u)
	}
	return buf
}

func Float32ToBytesAlt(fs []float32) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, fs)

	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

// fast but unsafe
func Float32ToBytesFastUnsafe(fs []float32) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(&fs[0])), len(fs)*4)
}

/*
func decodeUnsafe(bs []byte) []float32 {
	return unsafe.Slice((*float32)(unsafe.Pointer(&bs[0])), len(bs)/4)
}
*/
