package utils

import (
	"bot_personal_insurance/pkg/models"
	"encoding/json"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func CleanString(input string) string {
	trimmed := strings.TrimSpace(strings.ReplaceAll(input, "\n", " "))
	return strings.Join(strings.Fields(trimmed), " ")
}

func FilterEmptyStrings(input []string) []string {
	var result []string
	for _, str := range input {
		trimmed := strings.TrimSpace(str)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func ParseAgeRange(ageRange string) (*int, *int) {
	ageRange = strings.TrimSpace(ageRange)
	if ageRange == "ไม่มีกำหนด" || ageRange == "" {
		return nil, nil
	}

	agePattern := regexp.MustCompile(`(\d+)-(\d+)`)
	matches := agePattern.FindStringSubmatch(ageRange)

	if len(matches) == 3 {
		minAge, _ := strconv.Atoi(matches[1])
		maxAge, _ := strconv.Atoi(matches[2])
		return &minAge, &maxAge
	}
	return nil, nil
}

func ParseTextIntoArray(text, delimeter string) []string {
	parts := strings.Split(text, delimeter)
	parts = FilterEmptyStrings(parts)
	if len(parts) == 0 {
		return nil
	}
	return parts
}

func WriteJSON(data []models.PersonalLoan, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, jsonData, 0644)
}

func NullIfEmpty(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
