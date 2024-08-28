package main

import (
	"bytes"
	"digitalbanking_fees/models"
	"digitalbanking_fees/parser"
	"digitalbanking_fees/utils"
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
	// Define the URL and payload
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/DigitalBankingServiceFee/CompareProductList"
	payloadTemplate := `{"ProductIdList":"361,364,332,328,321,70,67,69,63,66,68,64,65,397,400,402,396,424,414,119,317,169,382,432,434,435,431,59,42,219,220,216,217,152,158,150,156,155,157,153,154,392,385,234,204,368,372,367,370,366,374,375,365,369,401,12,362,395,100,151,388,373,236,297,2,333,334,329,427,36,37,159,211,212,213,31,203","Page":%d,"Limit":3}`

	initialPayload := fmt.Sprintf(payloadTemplate, 1)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(initialPayload)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	utils.AddHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	totalPages := utils.DetermineTotalPage(doc)
	if totalPages == 0 {
		log.Fatal("Could not determine the total number of pages")
	}

	var result []models.DigitalBanking

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

		doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(body))
		if err != nil {
			log.Printf("Error parsing HTML for page %d: %v", page, err)
			continue
		}
		for i := 1; i <= 3; i++ {
			col := "col" + strconv.Itoa(i)
			provider := utils.CleanText(doc.Find(fmt.Sprintf("th.%s span", col)).Text())
			product := utils.CleanText(doc.Find(fmt.Sprintf("th.prod-col%d span", i)).Text())

			serviceDetails := parser.ParseServiceDetails(doc, col, i)
			fees := parser.ParseFeeDetails(doc, col, i)
			additionalInfo := parser.ParseAdditionalDetails(doc, col, i)

			result = append(result, models.DigitalBanking{
				Provider:   provider,
				Product:    product,
				Service:    serviceDetails,
				Fees:       fees,
				Additional: additionalInfo,
			})
		}
	}

	jsonData, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	err = os.WriteFile("digitalbanking_fees.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing JSON to file: %v", err)
	}

	fmt.Println("Data successfully written to digitalbanmking_fees.json")
}
