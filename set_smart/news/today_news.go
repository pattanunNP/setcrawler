package news

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ledongthuc/pdf"
	"github.com/nguyenthenguyen/docx"
	"github.com/xuri/excelize/v2"
)

type NewsItem struct {
	DateTime string `json:"date_time"`
	Symbol   string `json:"symbol"`
	Source   string `json:"source"`
	Subject  string `json:"suject"`
	Detail   string `json:"detail"`
}

type NewsDetail struct {
	DateTime            string `json:"date_time"`
	Headline            string `json:"headline"`
	Symbol              string `json:"symbol"`
	FulldetailedNews    string `json:"full_detailed_news"`
	AnnouncementDetails string `json:"announcement_details"`
	FileContent         string `json:"file_content"`
}

func FetchNews(cookieStr, symbol, locale string) ([]NewsItem, error) {
	url := "https://www.setsmart.com/ism/searchTodayNews.html"
	payload := strings.NewReader(fmt.Sprintf("companyNews=on&exchangeNews=on&lstView=bySymbol&symbol=%s&regulatorSymbol=&lstSecType=&lstSector=A_0_99_0_M&lstFavorite=0&txtSubject=&newsType=&submit.x=0&submit.y=0", symbol))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("cookie", cookieStr)
	req.Header.Set("origin", "https://www.setsmart.com")
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("referer", "https://www.setsmart.com/ism/searchTodayNews.html")
	req.Header.Set("sec-ch-ua", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return extractNewsItem(doc, cookieStr)
}

func extractNewsItem(doc *goquery.Document, cookieStr string) ([]NewsItem, error) {
	var newsItems []NewsItem

	doc.Find("#item tbody tr").Each(func(i int, row *goquery.Selection) {
		var item NewsItem
		row.Find("td").Each(func(i int, cell *goquery.Selection) {
			switch i {
			case 0:
				item.DateTime = convertToISO8601(cell.Text())
			case 2:
				item.Symbol = cell.Text()
			case 3:
				item.Source = cell.Text()
			case 4:
				item.Subject = cell.Text()
			case 5:
				detailURL := cell.Find("a").AttrOr("href", "")
				if detailURL != "" {
					detail, err := fetchDetailHTML("https://www.setsmart.com"+detailURL, cookieStr)
					if err == nil {
						detailJSON, err := json.Marshal(detail)
						if err == nil {
							item.Detail = string(detailJSON)
						}
					}
				}
			}
		})
		newsItems = append(newsItems, item)
	})
	return newsItems, nil
}

func fetchDetailHTML(url, cookieStr string) (*NewsDetail, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("cookie", cookieStr)
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error: non-200 status code %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	newsDetail := &NewsDetail{}
	details := make(map[string]string)

	doc.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() >= 2 {
			key := strings.TrimSpace(cells.Eq(0).Text())
			value := strings.TrimSpace(cells.Eq(1).Text())

			fmt.Printf("Key: %s, Value: %s\n", key, value)

			switch key {
			case "Date/Time:":
				newsDetail.DateTime = convertToISO8601(value)
			case "Headline:":
				newsDetail.Headline = value
			case "Symbol:":
				newsDetail.Symbol = value
			case "Full Detailed News:":
				pdfURL := cells.Eq(1).Find("a").AttrOr("href", "")
				if pdfURL != "" {
					pdfContent, err := downloadAndExtractFile("https://www.setsmart.com" + pdfURL)
					if err == nil {
						newsDetail.FileContent = pdfContent
					}
				}
			}
		}
	})

	announcementDetails := doc.Find(".newsstory-body").Text()
	details["Announcement Details"] = announcementDetails

	detailsStr, err := json.MarshalIndent(details, "", " ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling details to JSON: %v", err)
	}
	newsDetail.AnnouncementDetails = string(detailsStr)

	return newsDetail, nil
}

func downloadAndExtractFile(url string) (string, error) {
	fmt.Printf("Download file from URL: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("error: non-200 status code %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	fmt.Printf("Download file content type: %s\n", contentType)
	if strings.Contains(contentType, "application/pdf") {
		return extractPDF(resp.Body)
	} else if strings.Contains(contentType, "application/zip") {
		return extractZIP(resp.Body)
	} else {
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}
}

func extractPDF(reader io.Reader) (string, error) {
	tempFile, err := os.CreateTemp("", "*.pdf")
	if err != nil {
		return "", fmt.Errorf("error creating temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, reader); err != nil {
		return "", fmt.Errorf("error writing to temporary file: %v", err)
	}

	pdfFile, pdfReader, err := pdf.Open(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("error opening PDF file: %v", err)
	}
	defer pdfFile.Close()

	var buf bytes.Buffer
	textReader, err := pdfReader.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("error extracting text from PDF: %v", err)
	}
	if _, err := buf.ReadFrom(textReader); err != nil {
		return "", fmt.Errorf("error reading from text reader: %v", err)
	}

	return buf.String(), nil
}

func extractZIP(reader io.Reader) (string, error) {
	body, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("error reading ZIP data: %v", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return "", fmt.Errorf("error opening ZIP archive: %v", err)
	}

	var extractedText strings.Builder
	for _, file := range zipReader.File {
		if err := extractZIPFile(file, &extractedText); err != nil {
			return "", fmt.Errorf("error extracting file from ZIP: %v", err)
		}
	}
	return extractedText.String(), nil
}

func extractZIPFile(file *zip.File, extractedText *strings.Builder) error {
	rc, err := file.Open()
	if err != nil {
		return fmt.Errorf("error opening file in ZIP: %v", err)
	}
	defer rc.Close()

	fileName := strings.ToLower(file.Name)
	if strings.HasSuffix(fileName, ".pdf") {
		content, err := extractPDF(rc)
		if err != nil {
			return err
		}
		extractedText.WriteString(content)
	} else if strings.HasSuffix(fileName, ".docx") {
		content, err := extractDOCX(rc)
		if err != nil {
			return err
		}
		extractedText.WriteString(content)
	} else if strings.HasSuffix(fileName, ".xlsx") {
		content, err := extractXLSX(rc)
		if err != nil {
			return err
		}
		extractedText.WriteString(content)
	} else {
		content, err := io.ReadAll(rc)
		if err != nil {
			return err
		}
		extractedText.Write(content)
	}
	return nil
}

func extractDOCX(reader io.Reader) (string, error) {
	tempFile, err := os.CreateTemp("", "*.docx")
	if err != nil {
		return "", fmt.Errorf("error creating temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, reader); err != nil {
		return "", fmt.Errorf("error writing to temporary file: %v", err)
	}
	return extractDOCXPages(tempFile.Name())
}

func extractDOCXPages(filePath string) (string, error) {
	r, err := docx.ReadDocxFile(filePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	doc := r.Editable()
	contentString := doc.GetContent()
	return extractTextFromDocxContent(contentString), nil
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

func extractXLSX(reader io.Reader) (string, error) {
	tempFile, err := os.CreateTemp("", "*.xlsx")
	if err != nil {
		return "", fmt.Errorf("error creating temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, reader); err != nil {
		return "", fmt.Errorf("error writing to temporary file: %v", err)
	}

	return extractExcelXlsxPages(tempFile.Name())
}

func extractExcelXlsxPages(filePath string) (string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var sb strings.Builder
	sheets := f.GetSheetList()
	for _, sheetName := range sheets {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			return "", err
		}
		for _, row := range rows {
			sb.WriteString(strings.TrimSpace(strings.Join(row, "\t")) + "\n")
		}
	}

	return sb.String(), nil
}

func convertToISO8601(dateStr string) string {
	const inputLayout = "02/01/2006 15:04"
	const outputLayout = time.RFC3339

	t, err := time.Parse(inputLayout, dateStr)
	if err != nil {
		fmt.Printf("Error parsing date: %v\n", err)
		return dateStr
	}
	return t.Format(outputLayout)
}

func SaveToFile(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("error encoding JSON to file: %v", err)
	}
	return nil
}
