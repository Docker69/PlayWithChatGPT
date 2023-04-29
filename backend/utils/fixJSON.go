package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func FixJSON(jsonStr string) (string, error) {
	// If the JSON string is already valid, return it unmodified
	if json.Valid([]byte(jsonStr)) {
		return jsonStr, nil
	}

	// Remove any comments from the JSON string
	jsonStr = stripComments(jsonStr)

	// Try to fix any missing commas between object properties or array elements
	jsonStr = fixMissingCommas(jsonStr)

	// Try to fix any unmatched quotes in the JSON string
	jsonStr = fixUnmatchedQuotes(jsonStr)

	// Trim any trailing commas at the end of arrays or objects
	jsonStr = trimTrailingCommas(jsonStr)

	// Check if the resulting JSON string is now valid
	if !json.Valid([]byte(jsonStr)) {
		//create error
		err := fmt.Errorf("failed to fix JSON syntax errors")
		return "", err
	}

	return jsonStr, nil
}

// Helper function to remove single-line comments from a JSON string
func stripComments(jsonStr string) string {
	lines := strings.Split(jsonStr, "\n")
	for i, line := range lines {
		if strings.Contains(line, "//") {
			lines[i] = line[:strings.Index(line, "//")]
		}
	}
	return strings.Join(lines, "\n")
}

// Helper function to add missing commas between object properties or array elements
func fixMissingCommas(jsonStr string) string {
	return regexp.MustCompile(`(?m)^\s*([}\]])\s*(,?)$`).ReplaceAllString(jsonStr, "${1},")
}

// Helper function to escape unmatched quotes in a JSON string
func fixUnmatchedQuotes(jsonStr string) string {
	return regexp.MustCompile(`([^\\]|^)(")`).ReplaceAllStringFunc(jsonStr, func(m string) string {
		if m == `"` {
			return `\"`
		}
		return m
	})
}

// Helper function to trim trailing commas at the end of arrays or objects
func trimTrailingCommas(jsonStr string) string {
	return regexp.MustCompile(`,\s*([\]}])`).ReplaceAllString(jsonStr, "${1}")
}
