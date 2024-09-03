package pkg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CleanString(text string) string {
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")
	text = strings.TrimSpace(text)

	words := strings.Fields(text)
	return strings.Join(words, " ")
}

func ExtractFeeArray(doc *goquery.Document, className string, col string) []string {
	fee := CleanString(doc.Find(fmt.Sprintf("tr.%s td.%s span", className, col)).Text())
	if fee == "" {
		return nil
	}
	return splitAndStore(fee)
}

func ExtractFeeArrayPtr(doc *goquery.Document, className string, col string) []string {
	fee := CleanString(doc.Find(fmt.Sprintf("tr.%s td.%s span", className, col)).Text())
	if fee == "" {
		return nil
	}
	return splitAndStore(fee)
}

func ExtractNumericValue(text string) *float64 {
	re := regexp.MustCompile(`\d+(\.\d+)?`)
	match := re.FindString(text)
	if match == "" {
		return nil
	}
	value, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return nil
	}
	return &value
}

func ExtractURL(doc *goquery.Document, className string, col string) *string {
	href, exists := doc.Find(fmt.Sprintf("tr.%s td.%s a", className, col)).Attr("href")
	if !exists || href == "" {
		return nil
	}
	return &href
}

func DetermineTotalPage(doc *goquery.Document) int {
	totalPages := 1

	doc.Find("ul.pagination li a").Each(func(i int, s *goquery.Selection) {
		pageNum, exists := s.Attr("data-page")
		if exists {
			page, err := strconv.Atoi(pageNum)
			if err == nil && page > totalPages {
				totalPages = page
			}
		}
	})
	return totalPages
}

func splitAndStore(text string) []string {
	if text == "" {
		return nil
	}
	parts := strings.Split(text, "-")
	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
