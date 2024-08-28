package main

import (
	"bytes"
	"carloan_fees/models"
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
	// Define the URL
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/HirePurchaseFee/CompareProductList"
	// Define the JSON payload
	payload := `{"ProductIdList":"33,337,246,357,331,327,223,259,258,335,301,299,193,260,277,273,34,341,261,241,237,233,232,318,358,226,317,275,270,249","Page":1,"Limit":3}`

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
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
		log.Fatalf("Error reading response: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Error loading HTML into goquery: %v", err)
	}

	var carLoanFees []models.CarLoanFee

	for i := 1; i <= 3; i++ {
		col := "col" + strconv.Itoa(i)
		provider := doc.Find("th.attr-header.attr-prod.font-black.text-center.cmpr-col." + col + " span").Eq(1).Text()
		productType := doc.Find("th.attr-header.attr-prod.font-black.text-center.cmpr-col." + col + " span").Eq(2).Text()

		carLoanFees = append(carLoanFees, models.CarLoanFee{
			Provider:    provider,
			ProductType: productType,
		})
	}

	jsonData, err := json.MarshalIndent(carLoanFees, "", " ")
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}

	err = os.WriteFile("carloan_fees.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	fmt.Println("Car loan fees have been successfully saved to carloan_fees.json")
}
