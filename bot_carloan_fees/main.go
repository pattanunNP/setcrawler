package main

import (
	"bytes"
	"carloan_fees/models"
	"carloan_fees/parser"
	"carloan_fees/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// Define the URL and the initial payload with the page parameter
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/HirePurchaseFee/CompareProductList"
	payloadTemplate := `{"ProductIdList":"33,337,246,357,331,327,223,259,258,335,301,299,193,260,277,273,34,341,261,241,237,233,232,318,358,226,317,275,270,249","Page":%d,"Limit":3}`
	initialPayload := fmt.Sprintf(payloadTemplate, 1)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(initialPayload)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	utils.AddHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error performing request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Error loading HTTP response body into GoQuery: %v", err)
	}

	// Determine the total number of pages
	totalPages := utils.DetermineTotalPage(doc)
	if totalPages == 0 {
		log.Fatal("Could not determine the total number of pages")
	}

	var carLoanFeesDetails []models.CarLoanFeeDetail

	// Iterate over each page to collect all data
	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages)
		payload := fmt.Sprintf(payloadTemplate, page)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			log.Printf("Error creating request for page %d: %v", page, err)
			continue
		}

		utils.AddHeader(req)

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error making request for page %d: %v", page, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Request for page %d failed with status: %d\n", page, resp.StatusCode)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response for page %d: %v", page, err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			log.Printf("Error parsing HTML for page %d: %v", page, err)
			continue
		}

		// Process each column in the page
		for i := 1; i <= 3; i++ {
			col := "col" + strconv.Itoa(i)
			provider := utils.CleanText(doc.Find(fmt.Sprintf("th.%s span", col)).Text())

			contractEffectiveDate := parser.ParseContractEffectiveDate(doc, i)
			generalFees := parser.ParseGeneralFees(doc, i)
			paymentFees := parser.ParsePaymentFees(doc, i)
			otherFees := parser.ParseOtherFees(doc, i)
			additionalInfo := parser.ParseAdditinoalInfo(doc, i)

			carLoanFeesDetails = append(carLoanFeesDetails, models.CarLoanFeeDetail{
				Provider:              provider,
				ContractEffectiveDate: contractEffectiveDate,
				GeneralFees:           generalFees,
				PaymentFees:           paymentFees,
				OtherFees:             otherFees,
				AdditionalInfo:        additionalInfo,
			})
		}
	}

	// Marshal the data to JSON and write it to a file
	jsonData, err := json.MarshalIndent(carLoanFeesDetails, "", " ")
	if err != nil {
		fmt.Println("Error marshaling data to JSON:", err)
		return
	}

	err = os.WriteFile("carloan_fees_details.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing JSON to file: %v", err)
	}

	fmt.Println("Data successfully saved to carloan_fees_details.json")
}
