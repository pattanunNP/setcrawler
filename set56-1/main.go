package main

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/PuerkitoBio/goquery"
	"google.golang.org/api/option"
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
}

func main() {
	url := "https://market.sec.or.th/public/idisc/th/Viewmore/fs-r561"
	bucketName := "dbd-crawler"
	folderPath := "set56"
	credentialsFile := "credentail"

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
	const maxRecords = 20

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
					documents, err = downloadAndProcessFile(href)
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

	err = uploadFolderToGCS(bucketName, folderPath, credentialsFile)
	if err != nil {
		log.Fatalf("Failed to upload folder to GCS: %v", err)
	}
}

func downloadAndProcessFile(url string) ([]Document, error) {
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

	if fileType != "ZIP" && strings.Contains(fileType, "PHP") {
		fileName += ".zip"
	}

	destFolder := filepath.Join("set56", "lodash", strings.TrimSuffix(fileName, filepath.Ext(fileName)))

	if fileType == "ZIP" {
		return extractZipFiles(body, destFolder)
	}

	// Handle non-ZIP files
	err = saveNonZipFile(body, fileName, destFolder)
	if err != nil {
		return nil, fmt.Errorf("failed to save non-ZIP file: %w", err)
	}

	document := Document{
		DocumentName: fileName,
		OriginalPath: filepath.Join(destFolder, fileName),
		FileType:     fileType,
	}

	return []Document{document}, nil
}

func setRequestHeaders(req *http.Request) {
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
}

func extractRedirectURL(body []byte) (string, error) {
	re := regexp.MustCompile(`document.location\s*=\s*'(.*?)';`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) >= 2 {
		return "http://capital.sec.or.th" + matches[1], nil
	}
	return "", fmt.Errorf("redirect URL not found")
}

func extractZipFiles(body []byte, destFolder string) ([]Document, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to create zip reader: %w", err)
	}

	if err := os.MkdirAll(destFolder, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	var documents []Document

	for _, file := range zipReader.File {
		f, err := file.Open()
		if err != nil {
			log.Printf("Error opening file in ZIP: %v", err)
			continue
		}
		defer f.Close()

		filePath := filepath.Join(destFolder, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
		} else {
			err = saveFileContent(f, filePath, file.Mode())
			if err != nil {
				log.Printf("Failed to save file %s: %v", filePath, err)
				continue
			}

			doc := Document{
				DocumentName: file.Name,
				OriginalPath: filePath,
				FileType:     "ZIP",
			}
			documents = append(documents, doc)
		}
	}
	return documents, nil
}

func saveFileContent(reader io.Reader, filePath string, mode os.FileMode) error {
	outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, reader)
	return err
}

func saveNonZipFile(body []byte, fileName, destFolder string) error {
	if err := os.MkdirAll(destFolder, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(destFolder, fileName)
	return os.WriteFile(filePath, body, 0644)
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

func uploadFolderToGCS(bucketName, folderPath, credentialsFile string) error {
	// Create a context and a new storage client
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return fmt.Errorf("failed to create GCS client: %w", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)

	// Walk through the folder and upload files
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Open the file
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", path, err)
		}
		defer file.Close()

		// Define the object name in GCS
		objectName := filepath.Join(filepath.Base(folderPath), path[len(folderPath)+1:])

		// Upload the file to GCS
		wc := bucket.Object(objectName).NewWriter(ctx)
		if _, err = io.Copy(wc, file); err != nil {
			return fmt.Errorf("failed to upload file %s: %w", path, err)
		}

		if err := wc.Close(); err != nil {
			return fmt.Errorf("failed to close writer for file %s: %w", path, err)
		}

		log.Printf("Uploaded %s to %s/%s", path, bucketName, objectName)
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to upload folder %s to GCS: %w", folderPath, err)
	}

	return nil
}
