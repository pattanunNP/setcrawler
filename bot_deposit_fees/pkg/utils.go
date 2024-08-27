package pkg

import (
	"fmt"
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

func ExtractFee(doc *goquery.Document, classNamce string, col string) string {
	fee := doc.Find(fmt.Sprintf("tr.%s td.%s span", classNamce, col)).Text()
	return CleanString(fee)
}

func ExtractFeePtr(doc *goquery.Document, classNamce string, col string) *string {
	fee := CleanString(doc.Find(fmt.Sprintf("tr.%s td.%s span", classNamce, col)).Text())
	if fee == "" {
		return nil
	}
	return &fee
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
