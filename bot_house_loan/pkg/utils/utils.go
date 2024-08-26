package utils

import (
	"strconv"
	"strings"
)

func CleanString(input string) string {
	trimmed := strings.TrimSpace(strings.ReplaceAll(input, "\n", " "))
	return strings.Join(strings.Fields(trimmed), " ")
}

func ParseInterestRateConditions(conditions string) []string {
	parts := strings.Split(conditions, "ปีที่")
	var cleanedParts []string
	for _, part := range parts {
		trimmedPart := strings.TrimSpace(part)
		if trimmedPart != "" {
			cleanedParts = append(cleanedParts, "ปีที่ "+trimmedPart)
		}
	}
	return cleanedParts
}

func ExtractNumericValue(input string) (*int, error) {
	cleaned := strings.TrimSpace(input)
	cleaned = strings.ReplaceAll(cleaned, "ปีขึ้นไป", "")
	ageStr := strings.TrimSpace(strings.Join(strings.Fields(cleaned), ""))
	if ageStr == "" {
		return nil, nil
	}
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		return nil, err
	}
	return &age, nil
}

func ParseOptionalCondition(condition string) *string {
	if condition == "" || condition == "-" {
		return nil
	}
	return &condition
}
