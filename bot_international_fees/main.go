package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"internationla_fees/models"
	"internationla_fees/utils"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/InternationalTransactionFee/CompareProductList"
	payloadTemplate := `{"ProductIdList":"122,107,108,97,81,116,13,34,39,37,99,72,20,30,38,118,12,4,52,24,19,85,50,91,104,73,15,109,100,21,66,65,51","Page":%d,"Limit":3}`

	// Initialize first request to determine total pages
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

	// Check if response status is OK
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	totalPages := utils.DetermineTotalPage(doc)
	if totalPages == 0 {
		log.Fatal("Could not determine the total number of pages")
	}

	var international_fees []models.Fees

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

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Printf("Error parsing HTML for page %d: %v", page, err)
			continue
		}

		for i := 1; i <= 3; i++ {
			col := "col" + strconv.Itoa(i)

			var feeDetails models.Fees

			feeDetails.Provider = utils.CleanText(doc.Find(fmt.Sprintf("th.%s span", col)).Text())

			// Extract and populate InternationalTransferFees
			inwardRemittanceText := utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-InwardRemittanceFee td.%s", col))
			outwardRemittanceText := utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-OutwardRemittanceFee td.%s", col))

			feeDetails.InternationalTransferFees.InwardRemittance = models.InwardRemittance{
				Fee:                     utils.ExtractFee(inwardRemittanceText),
				ExchangeCompensationFee: utils.ExtractCompensationFee(inwardRemittanceText),
			}

			feeDetails.InternationalTransferFees.OutwardRemittance = models.OutwardRemittance{
				FeeType:                 utils.ExtractFeeType(outwardRemittanceText),
				Conditions:              utils.ProcessTextWithPattern(outwardRemittanceText),
				ExchangeCompensationFee: utils.ExtractCompensationFee(outwardRemittanceText),
			}

			// Extract and populate CheckAndDraftFees
			feeDetails.CheckAndDraftFees.TravelerChequeBuyingFee = utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-TravelerChequeBuyingFee td.%s", col))
			feeDetails.CheckAndDraftFees.TravelerChequeSellingFee = utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-TravelerChequeSellingFee td.%s", col))
			feeDetails.CheckAndDraftFees.DraftBuyingFee = utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-DraftBuyingFee td.%s", col))
			feeDetails.CheckAndDraftFees.DraftSellingFee = utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-DraftSellingFee td.%s", col))
			feeDetails.CheckAndDraftFees.ForeignBillBuyingFee = utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-ForeignBillBuyingFee td.%s", col))
			feeDetails.CheckAndDraftFees.ForeignBillSellingFee = utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-ForeignBillSellingFee td.%s", col))
			feeDetails.CheckAndDraftFees.ExchangeCompensationFee = utils.ExtractCompensationFeeFromTd(doc, fmt.Sprintf(".attr-ExchangeCompensationFee td.%s", col))

			// Extract and populate LetterOfCreditFees
			feeDetails.LetterOfCreditFees.ForeignLC = utils.CleanText(utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-ForeignLetterOfCreditFee td.%s", col)))
			feeDetails.LetterOfCreditFees.DomesticLC = utils.CleanText(utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-DomesticLetterOfCreditFee td.%s", col)))

			// Extract and populate BillCollectionFees
			feeDetails.BillCollectionFees.InwardBillFee = utils.CleanText(utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-InwardBillFee td.%s", col)))
			feeDetails.BillCollectionFees.OutwardBillFeeExporter = utils.CleanText(utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-OutwardBillFeeFromExporter td.%s", col)))
			feeDetails.BillCollectionFees.OutwardBillFeeImporter = utils.CleanText(utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-OutwardBillFeeFromImporter td.%s", col)))
			feeDetails.BillCollectionFees.ImportBillFee = utils.CleanText(utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-ImportBillFee td.%s", col)))
			feeDetails.BillCollectionFees.ExportBillFeeSeller = models.InvoiceFees{
				FirstInvoice:       utils.CleanText(utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-FirstInvoice td.%s", col))),
				SubsequentInvoices: utils.CleanText(utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-SubsequentInvoices td.%s", col))),
			}
			feeDetails.BillCollectionFees.ExportBillFeeBuyer = utils.CleanText(utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-ExportBillFeeFromBuyer td.%s", col)))

			// Extract and populate OtherFees
			otherFeeText := utils.CleanText(utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-OtherFees td.%s", col)))
			if otherFeeText != "" {
				otherFeesArray := utils.SplitAndCleanText(otherFeeText, "-")
				feeDetails.OtherFees.OtherFee = otherFeesArray
			}

			// Extract and populate AdditionalInformation (only one URL field)
			feeURL := doc.Find(fmt.Sprintf(".attr-FeeURL td.%s a.prod-url", col)).AttrOr("href", "")
			if feeURL != "" {
				feeDetails.AdditionalInformation.FeeURL = &feeURL
			}

			international_fees = append(international_fees, feeDetails)
		}
		time.Sleep(2 * time.Second)
	}

	file, err := json.MarshalIndent(international_fees, "", " ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	err = os.WriteFile("international_fee.json", file, 0644)
	if err != nil {
		log.Fatalf("Failed to write JSON to file: %v", err)
	}

	fmt.Println("Data successfully saved to international_fee.json")
}
