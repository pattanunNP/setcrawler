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

type RecordEN struct {
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

type FormData struct {
	Company  string `json:"ctl00$CPH$ddlCompany"`
	DateType string `json:"ctl00$CPH$rblDateType"`
	DateFrom string `json:"ctl00$CPH$BSDateFrom"`
	DateTo   string `json:"ctl00$CPH$BSDateTo"`
}

func main() {
	initialURL := "https://market.sec.or.th/public/idisc/th/r59"

	dropdownValues, err := getDropDownValues(initialURL)
	if err != nil {
		log.Fatal(err)
	}

	radioValue, err := getRadioValue(initialURL)
	if err != nil {
		log.Fatal(err)
	}

	// Simulate date
	startDate := time.Date(1975, 4, 30, 0, 0, 0, 0, time.UTC)
	endDate := time.Now()

	step := 30 * 24 * time.Hour
	dateRanges := generateDateRanges(startDate, endDate, step)

	resultCount := 0
	const maxResults = 30

	for _, dropdownValue := range dropdownValues {
		for _, dateRange := range dateRanges {
			formData := FormData{
				Company:  dropdownValue,
				DateType: radioValue,
				DateFrom: dateRange.DateFrom,
				DateTo:   dateRange.DateTo,
			}

			responseBody, err := postFormData(initialURL, formData)
			if err != nil {
				log.Printf("Failed to post form data for company %s, range %s to %s: %v\n", dropdownValue, dateRange.DateFrom, dateRange.DateTo, err)
				continue
			}

			fmt.Printf("Successfully posted form data for company %s, range %s to %s\n", dropdownValue, dateRange.DateFrom, dateRange.DateTo)

			records, recordsEN := collectContent(responseBody)
			if len(records) == 0 && len(recordsEN) == 0 {
				fmt.Printf("No relevant data found for company %s, range %s to %s. Skipping...\n", dropdownValue, dateRange.DateFrom, dateRange.DateTo)
				continue
			}

			printRecordsAsJSON(records, recordsEN)

			resultCount++
			if resultCount >= maxResults {
				fmt.Println("Reached the maximum number of results. Stopping..")
				return
			}
		}
	}

	//fmt.Println(dropdownValues)
}

func collectContent(responseBody string) ([]Record, []RecordEN) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(responseBody))
	if err != nil {
		log.Fatal("failed to parse response body:", err)
	}

	var allRecords []Record
	var allRecordEN []RecordEN

	//fmt.Println(responseBody)

	// check if the document contain any data
	if doc.Find("td.RgCol_Center a").Length() == 0 {
		fmt.Println("No 'td.RgCol_Center a' elements found in response.")
		return allRecords, allRecordEN
	}

	// Extract batchNo, transId, reporter from the HTML
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Find("td.RgCol_Center a").Attr("href")
		if exists {
			fmt.Println(link) // debug
			batchNo, err := extractQueryParam(link, "batchNo")
			if err != nil {
				log.Fatal(err)
				fmt.Println(batchNo) // debug
			}
			transId, err := extractQueryParam(link, "transId")
			if err != nil {
				log.Fatal(err)
				fmt.Println(transId) //debug
			}
			reporterParam, err := extractQueryParam(link, "reporter")
			if err != nil {
				log.Fatal(err)
				fmt.Print(reporterParam) //debug
			}

			// Collect and process content from the link
			records, recordsEN := collectContentFromLink(batchNo, transId, reporterParam)
			if len(records) > 0 || len(recordsEN) > 0 {
				allRecords = append(allRecords, records...)
				allRecordEN = append(allRecordEN, recordsEN...)
			}
		}
	})

	return allRecords, allRecordEN
}

func collectContentFromLink(batchNo, transId, reporterParam string) ([]Record, []RecordEN) {
	postURL := "https://market.sec.or.th/r59/publicapi/report"
	formDataTh := map[string]string{
		"BatchNo": batchNo,
		"Lang":    "Th",
	}
	formDataEn := map[string]string{
		"BatchNo": batchNo,
		"Lang":    "En",
	}

	jsonDataTh, err := json.Marshal(formDataTh)
	if err != nil {
		log.Fatal(err)
	}
	jsonDataEn, err := json.Marshal(formDataEn)
	if err != nil {
		log.Fatal(err)
	}

	reqTh, err := http.NewRequest("POST", postURL, bytes.NewBuffer(jsonDataTh))
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resTh, err := client.Do(reqTh)
	if err != nil {
		log.Fatal(err)
	}
	defer resTh.Body.Close()

	if resTh.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resTh.StatusCode, resTh.Status)
	}

	bodyTh, err := io.ReadAll(resTh.Body)
	if err != nil {
		log.Fatal(err)
	}

	var apiResponseTh ApiResponse
	err = json.Unmarshal(bodyTh, &apiResponseTh)
	if err != nil {
		log.Fatal("Failed to unmarshal JSON:", err)
	}

	// Create Send POST request for Eng
	reqEn, err := http.NewRequest("POST", postURL, bytes.NewBuffer(jsonDataEn))
	if err != nil {
		log.Fatal(err)
	}
	resEn, err := client.Do(reqEn)
	if err != nil {
		log.Fatal(err)
	}
	defer resEn.Body.Close()

	if resEn.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resEn.StatusCode, resEn.Status)
	}

	bodyEn, err := io.ReadAll(resEn.Body)
	if err != nil {
		log.Fatal(err)
	}

	var apiResponseEn ApiResponse
	err = json.Unmarshal(bodyEn, &apiResponseEn)
	if err != nil {
		log.Fatal("Failed to unmarshal JSON:", err)
	}

	// Create lists to hold all records
	var records []Record
	var recordsEN []RecordEN

	// Iterate over all transactions and create a record for each in Thai
	for _, transaction := range apiResponseTh.Report.TransactionList {
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
			BatchNo:           apiResponseTh.Report.BatchNo,
			Company:           apiResponseTh.Report.Company,
			Reporter:          apiResponseTh.Report.Reporter,
			Position:          apiResponseTh.Report.Position,
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

	// Iterate over all transactions and create a record for each in English
	for _, transaction := range apiResponseEn.Report.TransactionList {
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

		recordEN := RecordEN{
			BatchNo:           apiResponseEn.Report.BatchNo,
			Company:           apiResponseEn.Report.Company,
			Reporter:          apiResponseEn.Report.Reporter,
			Position:          apiResponseEn.Report.Position,
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
		recordsEN = append(recordsEN, recordEN)
	}

	return records, recordsEN
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

func printRecordsAsJSON(records []Record, recordsEN []RecordEN) {
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))

	jsonDataEN, err := json.MarshalIndent(recordsEN, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonDataEN))
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

	// Load HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var dropdownName string
	doc.Find("select").Each(func(i int, s *goquery.Selection) {
		if name, exists := s.Attr("name"); exists {
			dropdownName = name
			return
		}
	})

	if len(dropdownName) == 0 {
		return nil, fmt.Errorf("no dropdown found")
	}

	var dropdownValues []string
	doc.Find(fmt.Sprintf("select[name='%s'] option", dropdownName)).Each(func(i int, s *goquery.Selection) {
		value, exists := s.Attr("value")
		if exists {
			dropdownValues = append(dropdownValues, value)
		}
	})
	return dropdownValues, nil
}

func getRadioValue(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	var radioValue string
	doc.Find("input[type='radio']").Each(func(i int, s *goquery.Selection) {
		if checked, exists := s.Attr("checked"); exists && checked == "checked" {
			value, _ := s.Attr("value")
			radioValue = value
			return
		}
	})

	if radioValue == "" {
		return "", fmt.Errorf("no radio button found")
	}

	return radioValue, nil
}

func formatDate(date time.Time) string {
	return date.Format("02/01/2006")
}

func generateDateRanges(startDate, endDate time.Time, step time.Duration) []FormData {
	var formDatalist []FormData
	for start := startDate; !start.After(endDate); start = start.Add(step) {
		formData := FormData{
			DateFrom: formatDate(start),
			DateTo:   formatDate(start.Add(step - time.Hour*24)),
		}
		formDatalist = append(formDatalist, formData)
	}
	return formDatalist
}

func postFormData(postURL string, formData FormData) (string, error) {
	form := url.Values{}
	form.Set("ctl00$CPH$ddlCompany", formData.Company)
	form.Set("ctl00$CPH$rblDateType", formData.DateType)
	form.Set("ctl00$CPH$BSDateFrom", formData.DateFrom)
	form.Set("ctl00$CPH$BSDateTo", formData.DateTo)

	req, err := http.NewRequest("POST", postURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Do(req)
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
