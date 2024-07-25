package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// PageContent represents the content of a single page
type PageContent struct {
	Page int    `json:"page"`
	Text string `json:"text"`
}



func extractTextByPage(pdfPath string, page int) (string, error) {
	cmd := exec.Command("java", "-jar", "pdfbox-app-3.0.2.jar", "export:text",
		"-encoding", "UTF-8", "-startPage", strconv.Itoa(page), "-endPage", strconv.Itoa(page),
		"-i", pdfPath, "-console")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running PDFBox CLI for page %d: %v", page, err)
	}

	return string(output), nil
}

func getTotalPages(pdfPath string) (int, error) {
	ctx, err := api.ReadContextFile(pdfPath)
	if err != nil {
		return 0, err
	}
	return ctx.PageCount, nil
}

// CleanText makes the provided text more readable
func CleanText(s string) string {
	// Replace escaped newline characters with actual newline characters
	s = strings.ReplaceAll(s, "\\n", "\n")

	// Remove leading and trailing spaces from each line
	s = strings.TrimSpace(s)

	s = strings.ReplaceAll(s, "The encoding parameter is ignored when writing to the console.", "")

	// Replace multiple spaces with a single space
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")

	// Replace misplaced spaces before punctuation
	s = regexp.MustCompile(`\s([,;.!?])`).ReplaceAllString(s, "$1")

	// Remove unnecessary control characters
	unnecessaryChars := regexp.MustCompile(`[\x00-\x1F\x7F-\x9F]`)
	s = unnecessaryChars.ReplaceAllString(s, "")

	// Remove zero-width spaces and joiners, and byte order mark
	s = strings.Map(func(r rune) rune {
		if r == 0x200B || r == 0x200C || r == 0x200D || r == 0xFEFF {
			return -1 // Drop the character
		}
		return r
	}, s)

	// Replace special quotation marks with standard ones
	s = strings.ReplaceAll(s, "“", "\"")
	s = strings.ReplaceAll(s, "”", "\"")
	s = strings.ReplaceAll(s, "‘", "'")
	s = strings.ReplaceAll(s, "’", "'")

	return  s
}
func main() {
	pdfPath := "100000_t_25660331.pdf"

	// Get total number of pages using pdfcpu
	totalPages, err := getTotalPages(pdfPath)
	if err != nil {
		log.Fatalf("Error getting total pages: %v", err)
	}

	// Prepare a slice to store page contents
	var pagesContent []PageContent

	// Iterate through all pages
	for pageIndex := 1; pageIndex <= totalPages; pageIndex++ {
		text, err := extractTextByPage(pdfPath, pageIndex)

		if err != nil {
			log.Printf("Error extracting text from page %d: %v", pageIndex, err)
			continue
		}
		log.Printf("extracting text from page %d\r", pageIndex)

		cleanedText := CleanText(text)
	

		// Add the page content to the slice
		pagesContent = append(pagesContent, PageContent{
			Page: pageIndex,
			Text: cleanedText,
		})
	}

	// Convert the slice to JSON
	jsonData, err := json.MarshalIndent(pagesContent, "", "  ")
	if err != nil {
		log.Fatalf("Error creating JSON: %v", err)
	}

	// Write JSON to file
	err = os.WriteFile("cleaned_pdf_text.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing JSON to file: %v", err)
	}

	fmt.Println("Cleaned text has been saved to cleaned_pdf_text.json")
}
