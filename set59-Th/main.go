package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Define the Table struct to store the extracted table data
type Record struct {
	CompanyName  string `json:"company_name"`
	Reporter     string `json:"reporter"`
	Relation     string `json:"relation"`
	AssetType    string `json:"asset_type"`
	TransDate    string `json:"trans_date"`
	Amount       string `json:"amount"`
	Price        string `json:"price"`
	MarketSource string `json:"market_source"`
	Note         Note   `json:"note"`
}

type Note struct {
	BatchNo           string  `json:"batch_no"`
	Company           string  `json:"company"`
	Reporter          string  `json:"reporter"`
	Position          string  `json:"position"`
	TransExecutor     string  `json:"trans_executor"`
	SecuType          string  `json:"secu_type"`
	TransDate         string  `json:"trans_date"`
	OutstandingBefore string  `json:"outstanding_before"`
	TransVolumn       string  `json:"trans_volumn"`
	AvgPrice          float64 `json:"avg_price"`
	OutstandingAfter  string  `json:"outstanding_after"`
	TransType         string  `json:"trans_type"`
	MarketSource      string  `json:"market_source"`
	TargetInfo        string  `json:"target_info"`
	RecordStatus      string  `json:"record_status"`
	TransId           string  `json:"trans_id"`
	ReporterUrl       string  `json:"reporter_url"`
	Language          string  `json:"language"`
}

type TableData struct {
	Records []Record `json:"records"`
}

type ApiResponse struct {
	Report Report `json:"Report"`
}

type Report struct {
	BatchNo         string        `json:"BatchNo"`
	Company         string        `json:"Company"`
	Reporter        string        `json:"Reporter"`
	Position        string        `json:"Position"`
	TransactionList []Transaction `json:"TransactionList"`
}

type Transaction struct {
	TransExecutor     string `json:"TransExecutor"`
	SecuType          string `json:"SecuType"`
	TransDate         string `json:"TransDate"`
	OutstandingBefore string `json:"OutstandingBefore"`
	TransVolumn       string `json:"TransVolumn"`
	AvgPrice          string `json:"AvgPrice"`
	OutstandingAfter  string `json:"OutstandingAfter"`
	TransType         string `json:"TransType"`
	MarketSource      string `json:"MarketSource"`
	TargetInfo        string `json:"TargetInfo"`
	RecordStatus      string `json:"RecordStatus"`
}

func main() {
	initialURL := "https://market.sec.or.th/public/idisc/th/r59"

	// Make the first request to get dropdown values
	dropdownValues, err := getDropDownValues(initialURL)
	if err != nil {
		log.Fatal(err)
	}

	if len(dropdownValues) == 0 {
		log.Fatal("No dropdown values found")
	}

	maxIterations := 1
	if len(dropdownValues) > maxIterations {
		dropdownValues = dropdownValues[:maxIterations]
	}

	// fmt.Printf("Dropdown values: %v\n", dropdownValues)
	// counter := 0

	// Loop through all dropdown values to make requests and print results
	for _, dropdownValue := range dropdownValues {
		// if counter >= maxIterations {
		// 	break
		// }
		startDate := convertToThaiYear("20120101")
		endDate := convertToThaiYear(time.Now().Format("20060102"))
		requestURL := fmt.Sprintf("https://market.sec.or.th/public/idisc/th/Viewmore/r59-2?UniqueIdReference=%s&DateType=1&DateFrom=%s&DateTo=%s",
			dropdownValue, startDate, endDate)

		responseBody, err := getRequest(requestURL)
		if err != nil {
			log.Printf("Failed to get data for dropdown value %s: %v\n", dropdownValue, err)
			continue
		}

		// Process the response body to collect content
		tableData := collectContent(responseBody)
		if len(tableData.Records) > 0 {
			printTablesAsJSON(tableData)
		}

		//counter++
	}
}

func getDropDownValues(url string) ([]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var dropdownValues []string
	doc.Find("select[name='ctl00$CPH$ddlCompany'] option").Each(func(i int, s *goquery.Selection) {
		value, exists := s.Attr("value")
		if exists {
			dropdownValues = append(dropdownValues, value)
		}
	})
	return dropdownValues, nil
}

func getRequest(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func collectContent(responseBody string) TableData {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(responseBody))
	if err != nil {
		log.Fatal("failed to parse response body:", err)
	}

	var tableData TableData

	// Get the table name from div.cardheading
	tableName := doc.Find(".card card-table .cardheading").Text()
	fmt.Printf("Table Name: %s\n", tableName)

	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		var record Record

		tds := s.Find("td")
		if tds.Length() < 9 {
			return
		}

		record.CompanyName = tds.Eq(0).Text()
		record.Reporter = tds.Eq(1).Text()
		record.Relation = tds.Eq(2).Text()
		record.AssetType = tds.Eq(3).Text()
		record.TransDate = tds.Eq(4).Text()
		record.Amount = tds.Eq(5).Text()
		record.Price = tds.Eq(6).Text()
		record.MarketSource = tds.Eq(7).Text()

		link, exists := tds.Eq(8).Find("a").Attr("href")
		if exists {
			fmt.Printf("Processing link: %s\n", link) // Debugging log
			noteContent := collectContentFromLink(link)
			record.Note = noteContent
		} else {
			record.Note = Note{TransVolumn: tds.Eq(8).Text()}
		}

		tableData.Records = append(tableData.Records, record)
	})

	return tableData
}

func collectContentFromLink(link string) Note {
	var note Note

	// Collect data in Thai
	thContent, err := fetchContentFromLink(link, "Th")
	if err != nil {
		log.Printf("Failed to collect Thai content: %v", err)
		return note
	}

	// Collect data in English
	enContent, err := fetchContentFromLink(link, "En")
	if err != nil {
		log.Printf("Failed to collect English content: %v", err)
		return note
	}

	note = combineContent(thContent, enContent)

	return note
}

func fetchContentFromLink(link, lang string) (Note, error) {
	var note Note
	postURL := "https://market.sec.or.th/r59/publicapi/report"

	batchNo, err := extractQueryParam(link, "batchNo")
	if err != nil {
		return note, fmt.Errorf("failed to extract batchNo: %v", err)
	}
	transId, err := extractQueryParam(link, "transId")
	if err != nil {
		return note, fmt.Errorf("failed to extract transId: %v", err)
	}
	reporterParam, err := extractQueryParam(link, "reporter")
	if err != nil {
		return note, fmt.Errorf("failed to extract reporterParam: %v", err)
	}

	formData := map[string]string{
		"BatchNo": batchNo,
		"Lang":    lang,
	}

	jsonData, err := json.Marshal(formData)
	if err != nil {
		log.Fatal(err)
		return note, err
	}

	req, err := http.NewRequest("POST", postURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		return note, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return note, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return note, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return note, err
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Fatal("Failed to unmarshal JSON:", err)
		return note, err
	}

	// Assuming apiResponse has a Report field containing the relevant data
	if len(apiResponse.Report.TransactionList) > 0 {
		for _, transaction := range apiResponse.Report.TransactionList {
			transDateISO, err := convertDateToISO8601(transaction.TransDate)
			if err != nil {
				log.Printf("Failed to convert date: %s\n", transaction.TransDate)
			}
			avgPriceFloat, _ := strconv.ParseFloat(transaction.AvgPrice, 64)

			note = Note{
				BatchNo:           apiResponse.Report.BatchNo,
				Company:           apiResponse.Report.Company,
				Reporter:          apiResponse.Report.Reporter,
				Position:          apiResponse.Report.Position,
				TransExecutor:     transaction.TransExecutor,
				SecuType:          transaction.SecuType,
				TransDate:         transDateISO,
				OutstandingBefore: transaction.OutstandingBefore,
				TransVolumn:       transaction.TransVolumn,
				AvgPrice:          avgPriceFloat,
				OutstandingAfter:  transaction.OutstandingAfter,
				TransType:         transaction.TransType,
				MarketSource:      transaction.MarketSource,
				TargetInfo:        transaction.TargetInfo,
				RecordStatus:      transaction.RecordStatus,
				TransId:           transId,
				ReporterUrl:       reporterParam,
				Language:          lang,
			}
		}
	}

	return note, nil
}

func combineContent(thContent, enContent Note) Note {
	return Note{
		BatchNo:           thContent.BatchNo,
		Company:           thContent.Company,
		Reporter:          thContent.Reporter,
		Position:          thContent.Position,
		TransExecutor:     thContent.TransExecutor,
		SecuType:          thContent.SecuType,
		TransDate:         thContent.TransDate,
		OutstandingBefore: thContent.OutstandingBefore,
		TransVolumn:       thContent.TransVolumn,
		AvgPrice:          thContent.AvgPrice,
		OutstandingAfter:  thContent.OutstandingAfter,
		TransType:         thContent.TransType,
		MarketSource:      thContent.MarketSource,
		TargetInfo:        thContent.TargetInfo,
		RecordStatus:      thContent.RecordStatus,
		TransId:           thContent.TransId,
		ReporterUrl:       thContent.ReporterUrl,
		Language:          thContent.Language + " and " + enContent.Language,
	}
}

func extractQueryParam(urlStr, param string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	queryParams := parsedURL.Query()
	value := queryParams.Get(param)
	if value == "" {
		return "", fmt.Errorf("query parameter '%s' not found", param)
	}

	return value, nil
}

func convertDateToISO8601(dateStr string) (string, error) {
	// Assume the input date is in the format "DD/MM/YYYY"
	const inputFormat = "02/01/2006"
	parsedDate, err := time.Parse(inputFormat, dateStr)
	if err != nil {
		return "", nil
	}
	return parsedDate.Format(time.RFC3339), nil

}

func printTablesAsJSON(tableData TableData) {
	jsonData, err := json.MarshalIndent(tableData, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))
}

func convertToThaiYear(dateStr string) string {
	// Convert the date from ISO 8601 to Thai Buddhist calendar year
	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		log.Fatal(err)
	}
	thaiYear := date.Year() + 543
	return fmt.Sprintf("%04d%02d%02d", thaiYear, date.Month(), date.Day())

}
