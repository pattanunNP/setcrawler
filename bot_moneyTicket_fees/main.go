package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"moneyticket_fees/models"
	"moneyticket_fees/utils"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/AvalAndAcceptanceServiceFee/CompareProductList"
	payloadTemplate := `{"ProductIdList":"60,46,58,52,53,33,48,7,30,61,16,5,41,11,27,21,56,34,57,39,20","Page":%d,"Limit":3}`

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
		log.Fatalf("Error parsing HTML document: %v", err)
	}

	totalPages := utils.DetermineTotalPage(doc)
	if totalPages == 0 {
		log.Fatal("Could not determine the total number of pages")
	}

	var moneyTicketList []models.MoneyTicketFees
	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages)
		payload := fmt.Sprintf(payloadTemplate, page)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			log.Printf("Error creating request for page %d: %v", page, err)
			continue
		}

		utils.AddHeader(req)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			log.Fatalf("Error parsing HTML document: %v", err)
		}

		for i := 1; i <= 3; i++ {
			col := "col" + strconv.Itoa(i)
			var moneyTicket models.MoneyTicketFees

			// Extract the provider name
			moneyTicket.Provider = utils.CleanText(doc.Find(fmt.Sprintf("th.%s span", col)).Text())

			// Extract acceptance fees
			acceptanceFeeSelector := fmt.Sprintf("tr.attr-Acceptance td.%s", col)
			moneyTicket.AcceptanceFee = utils.ParseFeeDetailsAsArray(doc, acceptanceFeeSelector)

			// Extract aval fees
			avalFeeSelector := fmt.Sprintf("tr.attr-Aval td.%s", col)
			moneyTicket.AvalFee = utils.ParseFeeDetailsAsArray(doc, avalFeeSelector)

			// Extract other fees (if any)
			otherFeeText := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-other td.%s span", col)).Text())
			if otherFeeText == "" {
				moneyTicket.OtherFees = nil
			} else {
				moneyTicket.OtherFees = &otherFeeText
			}

			// Extract additional information
			var additionalInfo models.AdditionalInfo
			doc.Find(fmt.Sprintf("tr.attr-Feeurl td.%s a.prod-url", col)).Each(func(index int, item *goquery.Selection) {
				href, exists := item.Attr("href")
				if exists {
					additionalInfo.FeeLinks = append(additionalInfo.FeeLinks, href)
				}
			})
			moneyTicket.AdditionalInfo = additionalInfo

			// Extract and parse numeric values for sorting/filtering
			moneyTicket.ExtractedInfo = utils.ExtractNumericInfo(moneyTicket.AcceptanceFee, moneyTicket.AvalFee)

			// Add to the list
			moneyTicketList = append(moneyTicketList, moneyTicket)
		}
		time.Sleep(2 * time.Second)
	}

	// Marshal into JSON and write to file
	jsonData, err := json.MarshalIndent(moneyTicketList, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	err = os.WriteFile("moneyticket_fees.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing JSON to file: %v", err)
	}

	fmt.Println("Data successfully saved to moneyticket_fees.json")
}
