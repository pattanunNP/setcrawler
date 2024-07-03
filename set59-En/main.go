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
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Define the Transaction struct
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

// Define the Report struct
type Report struct {
	BatchNo          string        `json:"BatchNo"`
	Company          string        `json:"Company"`
	Reporter         string        `json:"Reporter"`
	Position         string        `json:"Position"`
	SubmitDate       string        `json:"SubmitDate"`
	BusinessTypeCode string        `json:"BusinessTypeCode"`
	PrintDate        string        `json:"PrintDate"`
	TransactionList  []Transaction `json:"TransactionList"`
}

// Define the ResponseStatus struct
type ResponseStatus struct {
	Seq    int     `json:"Seq"`
	Value  string  `json:"Value"`
	TextTh *string `json:"TextTh"`
	TextEn *string `json:"TextEn"`
}

// Define the ApiResponse struct
type ApiResponse struct {
	ResponseStatus ResponseStatus `json:"ResponseStatus"`
	Report         Report         `json:"Report"`
}

// Define the Record struct to store desired fields
type Record struct {
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
}

func main() {
	initialURL := "https://market.sec.or.th/public/idisc/th/r59"

	// Request the HTML page.
	res, err := http.Get(initialURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var allRecords []Record

	// Find and process all links
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Find("td.RgCol_Center a").Attr("href")
		if exists {
			// Extract batchNo from the link
			batchNo, err := extractQueryParam(link, "batchNo")
			if err != nil {
				log.Fatal(err)
			}
			// Extract transId and reporter from the link
			transId, err := extractQueryParam(link, "transId")
			if err != nil {
				log.Fatal(err)
			}
			reporterParam, err := extractQueryParam(link, "reporter")
			if err != nil {
				log.Fatal(err)
			}

			// Collect and process content from the link
			records := collectContent(batchNo, transId, reporterParam)
			allRecords = append(allRecords, records...)
		}
	})
	printRecordsAsJSON(allRecords)

	// Save all records to a JSON file
	//saveRecordsToFile(allRecords, "set59th.json")
}

func collectContent(batchNo, transId, reporterParam string) []Record {
	postURL := "https://market.sec.or.th/r59/publicapi/report"
	formData := map[string]string{
		"BatchNo": batchNo,
		"Lang":    "En",
	}

	jsonData, err := json.Marshal(formData)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", postURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	// Send the POST request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Check the response status
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the JSON response
	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Fatal("Failed to unmarshal JSON:", err)
	}

	// Create a list to hold all records
	var records []Record

	// Iterate over all transactions and create a record for each
	for _, transaction := range apiResponse.Report.TransactionList {

		// Convert TransDate to ISO 8601 format
		transDateISO, err := convertDateToISO8601(transaction.TransDate)
		if err != nil {
			log.Printf("Failed to convert date: %s\n", transaction.TransDate)
			continue
		}

		// Convert AvgPrice to float64
		avgPriceFloat, err := strconv.ParseFloat(transaction.AvgPrice, 64)
		if err != nil {
			log.Printf("Failed to convert avg price: %s\n", transaction.AvgPrice)
			continue
		}

		record := Record{
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
		}
		records = append(records, record)
	}
	return records
}

func convertDateToISO8601(dateStr string) (string, error) {
	// Assume the input date is in the format "DD/MM/YYYY"
	const inputFormat = "02/01/2006"
	parsedDate, err := time.Parse(inputFormat, dateStr)
	if err != nil {
		return "", err
	}
	return parsedDate.Format(time.RFC3339), nil
}

func extractQueryParam(urlStr, param string) (string, error) {
	// Parse the URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	// Extract query parameters
	queryParams := parsedURL.Query()

	// Retrieve the value of the specified query parameter
	value := queryParams.Get(param)
	if value == "" {
		return "", fmt.Errorf("query parameter '%s' not found", param)
	}

	return value, nil
}

// func saveRecordsToFile(records []Record, filename string) {
//  jsonData, err := json.MarshalIndent(records, "", "  ")
//  if err != nil {
//      log.Fatal(err)
//  }

//  err = os.WriteFile(filename, jsonData, 0644)
//  if err != nil {
//      log.Fatal(err)
//  }

//  fmt.Printf("Records saved to %s\n", filename)
// }

func printRecordsAsJSON(records []Record) {
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonData))
}
