package formatter

import (
	"regexp"
	"strings"
)

// CleanString trims and cleans up a string
func CleanString(input string) string {
	trimmed := strings.TrimSpace(strings.ReplaceAll(input, "\n", " "))
	return strings.Join(strings.Fields(trimmed), " ")
}

// SplitAndFormatArray splits a string into an array using "-" as a delimiter
func SplitAndFormatArray(input string) []string {
	parts := strings.Split(input, "-")
	var formatted []string
	for _, part := range parts {
		trimmed := CleanString(part)
		if trimmed != "" {
			formatted = append(formatted, trimmed)
		}
	}
	return formatted
}

func SplitAndFormatArrayBy(input string, delimeter string) []string {
	parts := strings.Split(input, delimeter)
	var formatted []string
	for _, part := range parts {
		trimmed := CleanString(part)
		if trimmed != "" {
			formatted = append(formatted, trimmed)
		}
	}
	return formatted
}

func SplitAndFormatArrayCustom(input string) []string {
	// Check if the string contains numbered lists
	if regexp.MustCompile(`\d\.`).MatchString(input) {
		// Use numbered list as delimiter
		re := regexp.MustCompile(`(\d\.\s?)`)
		parts := re.Split(input, -1)
		var formatted []string
		for _, part := range parts {
			trimmed := CleanString(part)
			if trimmed != "" {
				formatted = append(formatted, trimmed)
			}
		}
		return formatted
	} else {
		// Use "-" as delimiter if no numbered list is found
		return SplitAndFormatArray(input)
	}
}
