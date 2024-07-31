package main

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/nfnt/resize"
	"github.com/nguyenthenguyen/docx"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/xuri/excelize/v2"
)

type Data struct {
	CompanyName string     `json:"company_name"`
	YearPeriod  string     `json:"year_period"`
	SubmitDate  string     `json:"submit_date"`
	Document    []Document `json:"document"`
}

type Document struct {
	DocumentName string `json:"document_name"`
	OriginalPath string `json:"original_path"`
	FileType     string `json:"file_type"`
	Pages        []Page `json:"pages"`
	Analysis     string `json:"analysis,omitempty"`
}

type Page struct {
	PageNumber int    `json:"page"`
	Text       string `json:"text"`
}

type OpenAIRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content []struct {
		Type     string `json:"type"`
		Text     string `json:"text,omitempty"`
		ImageURL struct {
			URL string `json:"url,omitempty"`
		} `json:"image_url,omitempty"`
	} `json:"content"`
}

type OpenAiResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func main() {
	url := "https://market.sec.or.th/public/idisc/th/Viewmore/fs-r561"
	apiKey := "apiKey"
	modelName := "gpt-4o"

	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Error: status code %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Failed to parse HTML document: %v", err)
	}

	var data []Data
	recordCount := 0
	const maxRecords = 4

	doc.Find("tr").Each(func(rowIndex int, row *goquery.Selection) {
		if recordCount >= maxRecords {
			return
		}

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
					log.Printf("Extract href link: %s", href)
					documents, err = downloadAndProcessFile(href, apiKey, modelName)
					if err != nil {
						log.Printf("Error downloading or extracting file: %v", err)
					}
				}
			}

			yearPeriodISO, err := convertToISO8601(yearPeriod)
			if err != nil {
				log.Printf("Error converting yearPeriod to ISO8601: %v", err)
			}
			submitDateISO, err := convertSubmitDateToISO8601(submitDate)
			if err != nil {
				log.Printf("Error converting submitDate to ISO8601: %v", err)
			}

			data = append(data, Data{
				CompanyName: companyName,
				YearPeriod:  yearPeriodISO,
				SubmitDate:  submitDateISO,
				Document:    documents,
			})

			recordCount++
		}

		if recordCount >= maxRecords {
			return
		}
	})

	err = saveResultsToFile("result.json", data)
	if err != nil {
		log.Fatalf("Error saving results to file: %v", err)
	}
}

func downloadAndProcessFile(url, apiKey, modelName string) ([]Document, error) {
	url = strings.Replace(url, "as_of=0000-00-00 00:00:00", "as_of=0000-00-00&00:00:00", 1)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	setRequestHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if strings.Contains(string(body), "<script") {
		redirectURL, err := extractRedirectURL(body)
		if err != nil {
			return nil, fmt.Errorf("failed to extract redirect URL: %w", err)
		}

		req, err = http.NewRequest("GET", redirectURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create redirect request: %w", err)
		}
		req.Header = resp.Request.Header

		resp, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to execute redirect request: %w", err)
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read redirect response body: %w", err)
		}
	}

	contentDisposition := resp.Header.Get("Content-Disposition")
	contentType := resp.Header.Get("Content-Type")

	fileName := getFileNameFromDisposition(contentDisposition, url)
	fileType := getFileTypeFromContentType(contentType)
	if fileType == "" {
		fileType = getFileType(fileName)
	}

	log.Printf("File Type: %s", fileType)

	if strings.Contains(contentType, "text/html") {
		log.Printf("HTML response detected: %s", string(body[:100]))
		return nil, fmt.Errorf("server returned an HTML page indicating an error: %s", string(body[:100]))
	}

	if fileType == "ZIP" || fileType == "ASPX" {
		return extractZipFiles(body, apiKey, modelName)
	}

	return processFile(body, fileType, fileName, apiKey, modelName)
}

func shouldConvertToImage(text string) bool {
	return strings.Contains(text, "à¾×!Í¡ÒÃÅ§·Ø¹")
}

func setRequestHeaders(req *http.Request) {
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", "_ga=GA1.1.1504185108.1719892671; TS01703d9d028=01f20e05678a331b2ff9ae1991862aadbe27d4c9f52f2a4b53e65262a893586401a36dab66350d3a58f20c0d785c99d9db2aed93b7; _ga_3NH0QL72D6=GS1.1.1721623445.33.1.1721623576.60.0.0; TS01703d9d=012c1f76db3fb8c2be621423f63d2a297d80aa4e25e794a5754319f9d4ee1378b9449c4888e3a1931d8efa4747e36f77a8a76feac6; TS023e49ee027=08f2067569ab2000f05a908732303db7cc41b0036b90e7ec856654067dfc60b58e4c3a9a4feb8f84087a14ac201130002d0b89c82386b4400d37b703987824bf59b8611a44677822f6efdab1e59aa6f6b8d4f49e5dce66cd3bb8d72210ca1334")
}

func extractRedirectURL(body []byte) (string, error) {
	re := regexp.MustCompile(`document.location\s*=\s*'(.*?)';`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) >= 2 {
		return "http://capital.sec.or.th" + matches[1], nil
	}
	return "", fmt.Errorf("redirect URL not found")
}

func processFile(body []byte, fileType, fileName, apiKey, modelName string) ([]Document, error) {
	tempFile, err := os.CreateTemp("", fmt.Sprintf("tempfile_*%s", filepath.Ext(fileName)))
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(body); err != nil {
		return nil, fmt.Errorf("failed to write to temp file: %v", err)
	}
	tempFile.Close()

	if fileType == "PDF" {
		log.Printf("Processing PDF File: %s", tempFile.Name())
		pageFiles, err := splitPDFIntoPages(tempFile.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to split PDF into pages: %v", err)
		}

		var documentPages []Page
		for i, pageFile := range pageFiles {
			defer os.Remove(pageFile)

			output, err := extractTextWithPDFBOX(pageFile)
			if err != nil {
				return nil, fmt.Errorf("failed to extract text with PDFBox for page %s: %v", pageFile, err)
			}

			if shouldConvertToImage(output) || len(output) < 1000 {
				log.Printf("Special text pattern detected on page %d, converting to image", i+1)
				base64Image, err := convertPageToBase64Image(tempFile.Name(), i+1)
				if err != nil {
					return nil, fmt.Errorf("failed to convert page %d to image: %v", i+1, err)
				}
				output, err = sendImageToOpenAI(base64Image, apiKey, modelName)
				if err != nil {
					return nil, fmt.Errorf("failed to send image to OpenAI for page %d: %v", i+1, err)
				}
				documentPages = append(documentPages, Page{
					PageNumber: i + 1,
					Text:       output,
				})
			} else {
				cleanedText := CleanText(output)
				documentPages = append(documentPages, Page{
					PageNumber: i + 1,
					Text:       cleanedText,
				})
			}
		}

		doc := Document{
			DocumentName: filepath.Base(tempFile.Name()),
			OriginalPath: filepath.Base(tempFile.Name()),
			FileType:     "PDF",
			Pages:        documentPages,
		}
		return []Document{doc}, nil
	}

	return extractTextFromFile(tempFile.Name(), apiKey, modelName, nil)
}

func splitPDFIntoPages(pdfPath string) ([]string, error) {
	outputDir, err := os.MkdirTemp("", "pdf_pages")
	if err != nil {
		return nil, fmt.Errorf("error creating temp output directory: %v", err)
	}

	cmd := exec.Command("pdfcpu", "split", pdfPath, outputDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error splitting PDF into pages: %v, output: %s", err, output)
	}

	pageFiles, err := filepath.Glob(filepath.Join(outputDir, "*.pdf"))
	if err != nil {
		return nil, fmt.Errorf("error listing split page files: %v", err)
	}
	return pageFiles, nil
}

func extractTextWithPDFBOX(filePath string) (string, error) {
	log.Printf("Executing pdfbox-app-3.0.2.jar on file: %s", filePath)

	output, err := exec.Command("java", "-jar", "pdfbox-app-3.0.2.jar", "export:text",
		"-encoding", "UTF-8", "-i", filePath, "-console").Output()
	if err != nil {
		log.Printf("Failed to execute pdfbox-app-3.0.2.jar: %v", err)
		return "", fmt.Errorf("failed to execute pdfbox-app-3.0.2.jar: %v", err)
	}

	return string(output), nil
}

func cleanExtractedText(text string) string {
	text = strings.ReplaceAll(text, "The encoding parameter is ignored when writing to the console.", "")
	return text
}

func getTotalPages(pdfPath string) (int, error) {
	ctx, err := api.ReadContextFile(pdfPath)
	if err != nil {
		return 0, err
	}
	return ctx.PageCount, nil
}

func extractTextFromFile(filePath, apiKey, modelName string, documents []Document) ([]Document, error) {
	fileType := getFileType(filePath)
	var pages []Page
	var err error

	switch fileType {
	case "PDF", "pdf":
		pages, _, err = extractPDFPages(filePath, apiKey, modelName)
	case "DOCX":
		pages, _, err = extractDocxPages(filePath)
	case "DOC":
		var docxFilePath string
		docxFilePath, err = convertDocToDocx(filePath)
		if err == nil {
			defer os.Remove(docxFilePath)
			pages, _, err = extractDocxPages(docxFilePath)
		}
	case "XLS":
		var xlsxFilePath string
		xlsxFilePath, err = convertXlsToXlsx(filePath)
		if err == nil {
			pages, err = extractExcelXlsxPages(xlsxFilePath)
			fileType = "XLSX"
		}
	case "XLSX":
		pages, err = extractExcelXlsxPages(filePath)
	case "PPTX":
		var content []byte
		content, err = os.ReadFile(filePath)
		if err != nil {
			pages, err = extractPPTXPages(content)
		}
	default:
		err = fmt.Errorf("unsupported file type: %s", fileType)
	}

	if err != nil {
		return documents, err
	}
	fileName := filepath.Base(filePath)

	doc := Document{
		DocumentName: fileName,
		OriginalPath: fileName,
		FileType:     fileType,
		Pages:        pages,
	}

	documents = append(documents, doc)
	return documents, nil
}

func convertXlsToXlsx(filePath string) (string, error) {
	xlsxFilePath := strings.TrimSuffix(filePath, filepath.Ext(filePath)) + ".xlsx"
	cmd := exec.Command("soffice", "--headless", "--convert-to", "xlsx", "--outdir", filepath.Dir(filePath), filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to convert XLS to XLSX: %v, output: %s", err, output)
	}

	if _, err := os.Stat(xlsxFilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("failed to convert XLS to XLSX: %v", xlsxFilePath)
	}
	return xlsxFilePath, nil
}

func extractExcelXlsxPages(filePath string) ([]Page, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var pages []Page
	for _, sheetName := range f.GetSheetMap() {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			return nil, err
		}
		var sb strings.Builder
		for _, row := range rows {
			sb.WriteString(strings.TrimSpace(strings.Join(row, "\t")) + "\n")
		}
		pages = append(pages, Page{
			PageNumber: len(pages) + 1,
			Text:       sb.String(),
		})
	}

	return pages, nil
}

func convertDocToDocx(docFilePath string) (string, error) {
	docxFilePath := strings.TrimSuffix(docFilePath, filepath.Ext(docFilePath)) + ".docx"
	cmd := exec.Command("unoconv", "-f", "docx", "-o", docxFilePath, docFilePath)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to convert .doc to .docx: %w", err)
	}
	return docxFilePath, nil
}

func extractDocxPages(filePath string) ([]Page, string, error) {
	r, err := docx.ReadDocxFile(filePath)
	if err != nil {
		return nil, "", err
	}
	defer r.Close()

	doc := r.Editable()
	contentString := doc.GetContent()
	text := extractTextFromDocxContent(contentString)
	pages := []Page{{
		PageNumber: 1,
		Text:       text,
	}}
	return pages, "DOCX", nil
}

func extractTextFromDocxContent(content string) string {
	var sb strings.Builder
	inTag := false

	for _, rune := range content {
		switch rune {
		case '<':
			inTag = true
		case '>':
			inTag = false
		default:
			if !inTag {
				sb.WriteRune(rune)
			}
		}
	}
	return sb.String()
}

func extractPDFPages(filePath, apiKey, modelName string) ([]Page, string, error) {
	log.Printf("Starting to extract PDF pages from file: %s", filePath)
	var processedTexts []Page
	totalPages, err := getTotalPages(filePath)
	if err != nil {
		log.Printf("Error getting total pages: %v", err)
		return nil, "", fmt.Errorf("error getting total pages: %v", err)
	}

	for i := 1; i <= totalPages; i++ {
		log.Printf("Processing page %d of %d", i, totalPages)
		text, err := extractTextByPage(filePath, i)
		if err != nil {
			log.Printf("Error extracting text from page %d: %v", i, err)
			processedTexts = append(processedTexts, Page{
				PageNumber: i,
				Text:       text,
			})
		} else {
			if shouldConvertToImage(text) || len(text) < 1000 {
				log.Printf("Special text pattern on page %d, converting to image", i)
				base64Image, err := convertPageToBase64Image(filePath, i)
				if err != nil {
					log.Printf("Failed to convert page %d to image: %v", i, err)
					continue
				}
				// Log the base64 image data to verify if it's sending the right page
				log.Printf("Base64 image data for page %d: %s", i, base64Image)
				ocrText, err := sendImageToOpenAI(base64Image, apiKey, modelName)
				if err != nil {
					return nil, "", fmt.Errorf("failed to send image to OpeAI for page %d: %v", i, err)
				}
				processedTexts = append(processedTexts, Page{
					PageNumber: i,
					Text:       ocrText,
				})
			} else {
				cleanedText := CleanText(text)
				processedTexts = append(processedTexts, Page{
					PageNumber: i,
					Text:       cleanedText,
				})
			}
		}
	}
	log.Printf("Finished extracting PDF pages from file: %s", filePath)
	return processedTexts, "PDF", nil
}

func extractTextByPage(pdfPath string, page int) (string, error) {
	log.Printf("Extracting text from page %d of file: %s", page, pdfPath)
	cmd := exec.Command("java", "-jar", "pdfbox-app-3.0.2.jar", "export:text",
		"-encoding", "UTF-8", "-startPage", strconv.Itoa(page), "-endPage", strconv.Itoa(page),
		"-i", pdfPath, "-console")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running PDFBox CLI for page %d: %v", page, err)
	}

	log.Printf("Extracted text from page %d: %s", page, string(output))
	return string(output), nil
}

func CleanText(s string) string {
	s = strings.ReplaceAll(s, "\\n", "\n")
	s = strings.TrimSpace(s)
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	s = regexp.MustCompile(`\s([,;.!?])`).ReplaceAllString(s, "$1")
	unnecessaryChars := regexp.MustCompile(`[\x00-\x1F\x7F-\x9F]`)
	s = unnecessaryChars.ReplaceAllString(s, "")
	s = strings.Map(func(r rune) rune {
		if r == 0x200B || r == 0x200C || r == 0x200D || r == 0xFEFF {
			return -1
		}
		return r
	}, s)
	s = strings.ReplaceAll(s, "“", "\"")
	s = strings.ReplaceAll(s, "”", "\"")
	s = strings.ReplaceAll(s, "‘", "'")
	s = strings.ReplaceAll(s, "’", "'")
	return s
}

func convertPageToBase64Image(pdfPath string, pageNumber int) (string, error) {
	imagePath, err := convertPDFToImage(pdfPath, pageNumber)
	if err != nil {
		return "", err
	}
	defer os.Remove(imagePath)

	resizedImagePath, err := resizeImage(imagePath)
	if err != nil {
		return "", err
	}
	defer os.Remove(resizedImagePath)

	return imageToBase64(resizedImagePath)
}

func resizeImage(imagePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %v", err)
	}

	resizedImage := resize.Resize(400, 0, img, resize.Lanczos2)

	outputPath := fmt.Sprintf("resized_%s", filepath.Base(imagePath))
	out, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create resized image file: %v", err)
	}
	defer out.Close()

	err = png.Encode(out, resizedImage)
	if err != nil {
		return "", fmt.Errorf("failed to encode resized image: %v", err)
	}

	return outputPath, nil
}

func convertPDFToImage(pdfPath string, pageNumber int) (string, error) {
	outputPrefix := fmt.Sprintf("page_%d", pageNumber)
	cmd := exec.Command("pdftoppm", "-f", strconv.Itoa(pageNumber), "-l", strconv.Itoa(pageNumber), "-png", pdfPath, outputPrefix)
	output, err := cmd.CombinedOutput()

	log.Printf("Executing command: pdftoppm -f %d -l %d -png %s %s", pageNumber, pageNumber, pdfPath, outputPrefix)
	log.Printf("Command output: %s", string(output))

	if err != nil {
		log.Printf("Command error: %v", err)
		return "", fmt.Errorf("failed to convert PDF page to image: %v, output: %s", err, output)
	}

	imagePath := fmt.Sprintf("%s-1.png", outputPrefix)
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		imagePath = fmt.Sprintf("%s-01.png", outputPrefix)
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			imagePath = fmt.Sprintf("%s-001.png", outputPrefix)
			if _, err := os.Stat(imagePath); os.IsNotExist(err) {
				log.Printf("None of the expected image files %s-1.png, %s-01.png, or %s-001.png exist", outputPrefix, outputPrefix, outputPrefix)
				return "", fmt.Errorf("none of the expected image files %s-1.png, %s-01.png, or %s-001.png exist", outputPrefix, outputPrefix, outputPrefix)
			}
		}
	}

	newImagePath := fmt.Sprintf("%s-converted.png", outputPrefix)
	input, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open generated image file: %v", err)
	}
	defer input.Close()

	outputFile, err := os.Create(newImagePath)
	if err != nil {
		return "", fmt.Errorf("failed to create new image file: %v", err)
	}
	defer outputFile.Close()

	if _, err := io.Copy(outputFile, input); err != nil {
		return "", fmt.Errorf("failed to copy image data to new file: %v", err)
	}

	return newImagePath, nil
}

func imageToBase64(filePath string) (string, error) {
	imgData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(imgData), nil
}

func sendImageToOpenAI(base64Image, apiKey, modelName string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	content := []struct {
		Type     string `json:"type"`
		Text     string `json:"text,omitempty"`
		ImageURL struct {
			URL string `json:"url,omitempty"`
		} `json:"image_url,omitempty"`
	}{
		{
			Type: "text",
			Text: "Extract text from this image, output only text.",
		},
		{
			Type: "image_url",
			ImageURL: struct {
				URL string `json:"url,omitempty"`
			}{
				URL: "data:image/png;base64," + base64Image,
			},
		},
	}

	reqBody := OpenAIRequest{
		Model: modelName,
		Messages: []Message{
			{
				Role:    "user",
				Content: content,
			},
		},
		MaxTokens: 300,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API request failed with status %d: %s", resp.StatusCode, bodyBytes)
	}

	var response OpenAiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from OpenAI")
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

func extractZipFiles(body []byte, apiKey, modelName string) ([]Document, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to create zip reader: %w", err)
	}

	var documents []Document

	for _, file := range zipReader.File {
		f, err := file.Open()
		if err != nil {
			log.Printf("Error opening file in ZIP: %v", err)
			continue
		}
		defer f.Close()

		content, err := io.ReadAll(f)
		if err != nil {
			log.Printf("Error reading file ZIP: %v", err)
			continue
		}

		tempFile, err := os.CreateTemp("", fmt.Sprintf("tempfile_*%s", filepath.Ext(file.Name)))
		if err != nil {
			log.Printf("Failed to create temp file: %v", err)
			continue
		}
		defer os.Remove(tempFile.Name())

		if _, err := tempFile.Write(content); err != nil {
			log.Printf("Failed to write to temp file: %v", err)
			continue
		}
		tempFile.Close()

		if filepath.Ext(file.Name) == ".pdf" || filepath.Ext(file.Name) == ".PDF" {
			log.Printf("Processing PDF File: %s", tempFile.Name())
			pageFiles, err := splitPDFIntoPages(tempFile.Name())
			if err != nil {
				log.Printf("Failed to split PDF into pages: %v", err)
				continue
			}

			var documentPages []Page
			for i, pageFile := range pageFiles {
				defer os.Remove(pageFile)

				output, err := extractTextWithPDFBOX(pageFile)
				if err != nil {
					log.Printf("Failed to extract text with PDFBox for page %s: %v", pageFile, err)
					continue
				}

				output = cleanExtractedText(output)

				page := Page{
					PageNumber: i + 1,
					Text:       output,
				}

				documentPages = append(documentPages, page)
			}

			doc := Document{
				DocumentName: filepath.Base(tempFile.Name()),
				OriginalPath: filepath.Base(tempFile.Name()),
				FileType:     "PDF",
				Pages:        documentPages,
			}
			documents = append(documents, doc)
		} else {
			docs, err := extractTextFromFile(tempFile.Name(), apiKey, modelName, nil)
			if err != nil {
				log.Printf("Error processing file %s: %v", file.Name, err)
				continue
			}
			documents = append(documents, docs...)
		}
	}
	return documents, nil
}

func getFileNameFromDisposition(contentDisposition, url string) string {
	if contentDisposition != "" {
		if _, params, err := mime.ParseMediaType(contentDisposition); err == nil {
			return params["filename"]
		}
	}
	return filepath.Base(url)
}

func getFileType(fileName string) string {
	parts := strings.Split(fileName, ".")
	if len(parts) > 1 {
		return strings.ToUpper(parts[len(parts)-1])
	}
	return ""
}

func getFileTypeFromContentType(contentType string) string {
	switch contentType {
	case "application/zip":
		return "ZIP"
	case "application/pdf":
		return "PDF"
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return "DOCX"
	case "application/vnd.ms-excel":
		return "XLS"
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return "XLSX"
	default:
		return ""
	}
}

func extractPPTXPages(content []byte) ([]Page, error) {
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

	prs, err := excelize.OpenFile(pptxFile.Name())
	if err != nil {
		return nil, err
	}

	var pages []Page
	for _, name := range prs.GetSheetMap() {
		rows, err := prs.GetRows(name)
		if err != nil {
			return nil, err
		}
		var sb strings.Builder
		for _, row := range rows {
			sb.WriteString(strings.Join(row, " ") + "\n")
		}
		pages = append(pages, Page{
			PageNumber: len(pages) + 1,
			Text:       sb.String(),
		})
	}

	return pages, nil
}

func saveResultsToFile(filename string, data []Data) error {
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
