package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/ledongthuc/pdf"
	"github.com/nguyenthenguyen/docx"
	"github.com/xuri/excelize/v2"
)

type Data struct {
	CompanyName string     `json:"company_name"`
	YearPeriod  string     `json:"year_period"`
	SubmitDate  string     `json:"submit_date"`
	Document    []Document `json:"document"`
}

type Document struct {
	DocumentName string   `json:"document_name"`
	OriginalPath string   `json:"original_path"`
	FileType     string   `json:"file_type"`
	Pages        []string `json:"pages"`
	Analysis     string   `json:"analysis,omitempty"`
}

func main() {
	// URL to crawl
	url := "https://market.sec.or.th/public/idisc/th/Viewmore/fs-r561"

	// Fetch the URL
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Error: status code %d", res.StatusCode)
	}

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Data structure to hold the extracted row data
	var data []Data

	// Counter for the number of records processed
	recordCount := 0
	const maxRecords = 100

	// Find all <tr> elements
	doc.Find("tr").Each(func(rowIndex int, row *goquery.Selection) {
		if recordCount >= maxRecords {
			return // Stop processing if we've reached the limit
		}

		// Create a slice to hold cell data for this row
		tds := row.Find("td")
		if tds.Length() >= 4 {
			companyName := tds.Eq(0).Text()
			yearPeriod := tds.Eq(1).Text()
			submitDate := tds.Eq(2).Text()

			link := tds.Eq(3).Find("a[href]")
			var documents []Document
			if link.Length() > 0 {
				href, exists := link.Attr("href")
				if exists {
					// Download and Extract file
					documents, err = downloadAndExtractZip(href)
					if err != nil {
						log.Printf("Error Download or Extract a File: %v", err)
					}
				}
			}

			// Convert yearPeriod and submitDate to ISO8601 format
			yearPeriodISO, err := convertToISO8601(yearPeriod)
			if err != nil {
				log.Printf("Error converting yearPeriod to ISO8601: %v", err)
			}
			submitDateISO, err := convertSubmitDateToISO8601(submitDate)
			if err != nil {
				log.Printf("Error converting submitDate to ISO8601: %v", err)
			}

			// Assign data to the Data struct
			data = append(data, Data{
				CompanyName: companyName,
				YearPeriod:  yearPeriodISO,
				SubmitDate:  submitDateISO,
				Document:    documents,
			})

			recordCount++
		}

		// Stop processing if we've reached the limit
		if recordCount >= maxRecords {
			return
		}
	})
	// Convert the Data to JSON
	// jsonData, err := json.MarshalIndent(data, "", "  ")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Save the data to JSON File
	err = saveResiltsToFile("result.json", data)
	if err != nil {
		log.Fatalf("Error saving resilts to file %v", err)
	}

	// // Print the JSON data
	// fmt.Println(string(jsonData))
}

func convertToISO8601(dateStr string) (string, error) {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return "", fmt.Errorf("date string is empty")
	}

	if len(dateStr) == 4 && isYear(dateStr) {
		year, _ := strconv.Atoi(dateStr)
		gregorianYear := year - 543
		return fmt.Sprintf("%d", gregorianYear), nil
	}

	return "", fmt.Errorf("invalid date format: %s", dateStr)
}

func convertSubmitDateToISO8601(dateStr string) (string, error) {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return "", fmt.Errorf("date string is empty")
	}

	inputFormat := "02/01/2006"

	parts := strings.Split(dateStr, "/")
	if len(parts) == 3 {
		year, err := strconv.Atoi(parts[2])
		if err == nil {
			adjustYear := year - 543
			parts[2] = strconv.Itoa(adjustYear)
			adjustedDateStr := strings.Join(parts, "/")

			parsedDate, err := time.Parse(inputFormat, adjustedDateStr)
			if err == nil {
				return parsedDate.Format(time.RFC3339), nil
			}
		}
	}

	return "", fmt.Errorf("invalid date format: %s", dateStr)
}

func isYear(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil && len(s) == 4
}

func downloadAndExtractZip(url string) ([]Document, error) {
	// Download ZIP file
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error: status code %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check if the response body is actually a ZIP file
	if !isZipFile(body) {
		return nil, fmt.Errorf("the downloaded file is not a valid ZIP archive")
	}

	// Create ZIP file reader
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, err
	}

	var documents []Document
	for _, file := range zipReader.File {
		f, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer f.Close()

		// Read file content
		content, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}

		fileType := getFileType(file.Name)
		var pages []string
		var analysis string

		switch fileType {
		case "PDF":
			pages, analysis, err = extractPDFPages(content)
		case "DOC":
			pages, err = extractDocPages(content)
		case "XLSX":
			pages, err = extractExcelPages(content)
		case "PPTX":
			pages, err = extractPPTXPages(content)
		default:
			pages = []string{string(content)}
		}

		if err != nil {
			log.Printf("Error processing file %s: %v", file.Name, err)
			continue
		}

		doc := Document{
			DocumentName: getFileName(file.Name),
			OriginalPath: file.Name,
			FileType:     fileType,
			Pages:        pages,
			Analysis:     analysis,
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

func getFileName(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func isZipFile(body []byte) bool {
	r := bytes.NewReader(body)
	_, err := zip.NewReader(r, int64(len(body)))
	return err == nil
}

func getFileType(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) > 1 {
		return strings.ToUpper(parts[len(parts)-1])
	}
	return ""
}

func extractPDFPages(content []byte) ([]string, string, error) {
	// Create a temporary file
	pdfFile, err := os.CreateTemp("", "*.pdf")
	if err != nil {
		return nil, "", err
	}
	defer os.Remove(pdfFile.Name())

	_, err = pdfFile.Write(content)
	if err != nil {
		return nil, "", err
	}
	pdfFile.Close()

	// Open the PDF File
	f, r, err := pdf.Open(pdfFile.Name())
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	var pages []string
	var analysis string

	for pageNum := 1; pageNum <= r.NumPage(); pageNum++ {
		page := r.Page(pageNum)
		if page.V.IsNull() {
			continue
		}
		text, err := extractTextFromPageHandlingErrors(page)
		if err != nil {
			log.Printf("Error extracting text from page %d of file %s: %v", pageNum, pdfFile.Name(), err)
			continue
		}
		if isImagePDF(text) {
			// If the page contains an image, send it to OpenAI for text extraction
			analysis, err = sendPDFToOpenAI(content)
			if err != nil {
				return nil, "", err
			}
			break
		} else if containsTable(page) {
			// If the page contains a table, format it as TSV
			pages = append(pages, formatAsTSV(text))
		} else {
			pages = append(pages, text)
		}
	}

	return pages, analysis, nil
}

func isImagePDF(text string) bool {
	return len(text) < 1000
}

func sendPDFToOpenAI(content []byte) (string, error) {
	return "Extracted Text from image", nil
}

func containsTable(page pdf.Page) bool {
	pageContent, err := page.GetPlainText(nil)
	if err != nil {
		log.Printf("Error extracting text from page: %v", err)
		return false
	}

	lines := strings.Split(pageContent, "\n")
	tabCounts := make(map[int]int)
	for _, line := range lines {
		tabCount := strings.Count(line, "\t")
		if tabCount > 0 {
			tabCounts[tabCount]++
		}
	}

	for _, count := range tabCounts {
		if count > 1 {
			return true
		}
	}

	return false
}

func formatAsTSV(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.Join(strings.Fields(line), "\t")
	}
	return strings.Join(lines, "\n")
}

func extractTextFromPageHandlingErrors(page pdf.Page) (string, error) {
	var buf bytes.Buffer
	pageContent, err := page.GetPlainText(nil)
	if err != nil {
		log.Printf("Error extracting text from page: %v", err)
		buf.WriteString(handleTextExtractionErrors(pageContent))
	} else {
		buf.WriteString(pageContent)
	}
	return buf.String(), nil
}

func handleTextExtractionErrors(text string) string {
	cleanedText := strings.ReplaceAll(text, "\x00", "")

	if !utf8.ValidString(cleanedText) {
		validText := make([]rune, 0, len(cleanedText))
		for i, r := range cleanedText {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(cleanedText[i:])
				if size == 1 {
					continue
				}
			}
			validText = append(validText, r)
		}
		cleanedText = string(validText)
	}

	return cleanedText
}

func extractDocPages(content []byte) ([]string, error) {
	// Create a temporary file
	docFile, err := os.CreateTemp("", "*.docx")
	if err != nil {
		return nil, err
	}
	defer os.Remove(docFile.Name())

	_, err = docFile.Write(content)
	if err != nil {
		return nil, err
	}
	docFile.Close()

	// Open the DOCX file
	r, err := docx.ReadDocxFile(docFile.Name())
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// Extract text from DOCX file
	doc := r.Editable()
	contentString := doc.GetContent()

	var pages []string
	pages = append(pages, contentString)

	return pages, nil
}

func extractExcelPages(content []byte) ([]string, error) {
	// Create temporary file
	excelFile, err := os.CreateTemp("", "*.xlsx")
	if err != nil {
		return nil, err
	}
	defer os.Remove(excelFile.Name())

	_, err = excelFile.Write(content)
	if err != nil {
		return nil, err
	}
	excelFile.Close()

	// Open Excel File
	f, err := excelize.OpenFile(excelFile.Name())
	if err != nil {
		return nil, err
	}

	// Extract text from all sheets
	var pages []string
	for _, sheet := range f.GetSheetMap() {
		rows, err := f.GetRows(sheet)
		if err != nil {
			return nil, err
		}
		for _, row := range rows {
			pages = append(pages, strings.Join(row, " "))
		}
	}

	return pages, nil
}

func extractPPTXPages(content []byte) ([]string, error) {
	// Create temporary file
	pptxFile, err := os.CreateTemp("", "*.pptx")
	if err != nil {
		return nil, err
	}
	defer os.Remove(pptxFile.Name())

	_, err = pptxFile.Write(content)
	if err != nil {
		return nil, err
	}
	pptxFile.Close()

	// Open the PPTX file
	ppt, err := excelize.OpenFile(pptxFile.Name())
	if err != nil {
		return nil, err
	}
	defer ppt.Close()

	var pages []string
	for _, name := range ppt.GetSheetMap() {
		rows, err := ppt.GetRows(name)
		if err != nil {
			return nil, err
		}
		for _, row := range rows {
			pages = append(pages, strings.Join(row, " "))
		}
	}

	return pages, nil
}

func saveResiltsToFile(filename string, data []Data) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}
