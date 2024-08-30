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

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/AvalAndAcceptanceServiceFee/CompareProductList"
	payload := `{"ProductIdList":"60,46,58,52,53,33,48,7,30,61,16,5,41,11,27,21,56,34,57,39,20","Page":1,"Limit":3}`

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
		log.Fatalf("Error reading response body: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Error parsing HTML document: %v", err)
	}

	var moneyTicketList []models.MoneyTicketFees
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

		// Add to the list
		moneyTicketList = append(moneyTicketList, moneyTicket)
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
