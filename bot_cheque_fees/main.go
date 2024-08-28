package main

import (
	"bytes"
	"cheque_fee/models"
	"cheque_fee/utils"
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
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/ChequeFee/CompareProductList"

	// JSON data to send in the POST request
	payloadTemplate := `{"ProductIdList":"162152,2,5,17,4,157479,15,27,34,6,194031,162151,240,155024,16,162568,237,28,241,24,163222,150920,32,9,23,37,35,30,33,26,449176,13","Page":%d,"Limit":3}`

	initialPayload := fmt.Sprintf(payloadTemplate, 1)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(initialPayload)))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	utils.AddHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	totalPages := utils.DetermineTotalPage(doc)
	if totalPages == 0 {
		log.Fatal("Could not determine the total number of pages")
	}

	var chequeFees []models.ChequeFee

	// Loop through each page to gather all data
	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages)
		payload := fmt.Sprintf(payloadTemplate, page)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			log.Printf("Error creating request for page %d: %v", page, err)
			continue
		}

		// Set headers
		utils.AddHeader(req)

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error sending request for page %d: %v", page, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Request for page %d failed with status: %d\n", page, resp.StatusCode)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body for page %d: %v", page, err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			log.Printf("Error parsing HTML for page %d: %v", page, err)
			continue
		}

		for i := 1; i <= 3; i++ {
			col := "col" + strconv.Itoa(i)
			provider := doc.Find("th.attr-header.attr-prod.font-black.text-center.cmpr-col." + col + " span").Last().Text()

			chequeFees = append(chequeFees, models.ChequeFee{
				Providers: provider,
				FeesTypes: models.FeesTypes{
					ChequeBookPurchase:                utils.SplitByCondition(utils.ExtractFee(doc, "attr-ChequeBookPurchase", col)),
					ChequeDepositAcross:               utils.SplitByCondition(utils.ExtractFee(doc, "attr-ChequeDepositAcross", col)),
					ChequeDepositInbranch:             utils.SplitByCondition(utils.ExtractFee(doc, "attr-ChequeDepositInbranch", col)),
					ChequeReturnFromInstrument:        utils.SplitByCondition(utils.ExtractFee(doc, "attr-ChequeReturnFromInstrument", col)),
					ChequeFeeReturned:                 utils.ExtractFeeArray(doc, "attr-ChequeFeeReturned", col),
					ChequeGiftPurchase:                utils.ExtractFee(doc, "attr-ChequeGiftPurchase", col),
					ChequeCashWithdrawAcross:          utils.ExtractFeeArray(doc, "attr-ChequeCashWithdrawAcross", col),
					ChequeCashWithdrawInbranch:        utils.ExtractFee(doc, "attr-ChequeCashWithdrawInbranch", col),
					CashierChequePurchase:             utils.ExtractFeeArray(doc, "attr-ChequeCashierPurchase", col),
					CashierChequeCashWithdrawAcross:   utils.SplitByCondition(utils.ExtractFee(doc, "attr-ChequeCashierCashWithdrawAcross", col)),
					CashierChequeCashWithdrawInbranch: utils.SplitByCondition(utils.ExtractFee(doc, "attr-ChequeCashierCashWithdrawInbranch", col)),
					DraftPurchaseFee:                  utils.ExtractFeeArray(doc, "attr-DraftPurchaseFee", col),
					PublicationFee:                    utils.ExtractFee(doc, "attr-PublicationFee", col),
					ChequeCancellationFee:             utils.ExtractFee(doc, "attr-ChequeCancellationFee", col),
					ChequeAdvanceDepositFee:           utils.ExtractFee(doc, "attr-ChequeAdvanceDepositFee", col),
				},
				OthersFees: models.OthersFees{
					OtherFees: utils.ExtractFee(doc, "attr-other", col),
				},
				AdditionalInfo: models.AdditionalInfo{
					WebsiteFeeLink: utils.ExtractLink(doc, "attr-Feeurl", col),
				},
			})
		}
	}

	jsonData, err := json.MarshalIndent(chequeFees, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	err = os.WriteFile("cheque_fees.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write JSON to file: %v", err)
	}

	fmt.Println("Data successfully saved to cheque_fees.json")
}
