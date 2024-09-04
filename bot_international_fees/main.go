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
	"regexp"
	"strconv"
	"strings"
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

	var internationalFees []models.Fees

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

			inwardFee := utils.SplitAndCleanText(utils.ExtractFee(inwardRemittanceText), "-")
			inwardFeeNumeric := extractNumericValues(inwardFee)

			if numSlice, ok := inwardFeeNumeric.([]float64); ok && len(numSlice) == 0 {
				inwardFeeNumeric = nil
			}

			feeDetails.InternationalTransferFees.InwardRemittance = models.InwardRemittance{
				Fee:                     inwardFee,
				FeeNumeric:              inwardFeeNumeric,
				ExchangeCompensationFee: utils.ExtractCompensationFee(inwardRemittanceText),
			}

			outwardConditions := utils.ProcessTextWithPattern(outwardRemittanceText)
			outwardConditionsNumeric := extractConditionsNumeric(outwardConditions)

			feeDetails.InternationalTransferFees.OutwardRemittance = models.OutwardRemittance{
				FeeType:                 utils.ExtractFeeType(outwardRemittanceText),
				Conditions:              outwardConditions,
				ConditionsNumeric:       outwardConditionsNumeric,
				ExchangeCompensationFee: utils.ExtractCompensationFee(outwardRemittanceText),
			}

			// Extract and populate CheckAndDraftFees
			feeDetails.CheckAndDraftFees.TravelerChequeBuyingFee = utils.SplitAndCleanText(utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-TravelerChequeBuyingFee td.%s", col)), "-")
			feeDetails.CheckAndDraftFees.TravelerChequeBuyingFeeNumeric = extractCheckAndDraftFeeNumeric(feeDetails.CheckAndDraftFees.TravelerChequeBuyingFee)
			feeDetails.CheckAndDraftFees.TravelerChequeSellingFee = utils.SplitAndCleanText(utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-TravelerChequeSellingFee td.%s", col)), "-")
			feeDetails.CheckAndDraftFees.DraftBuyingFee = utils.SplitAndCleanText(utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-DraftBuyingFee td.%s", col)), "-")
			feeDetails.CheckAndDraftFees.DraftBuyingFeeNumeric = extractCheckAndDraftFeeNumeric(feeDetails.CheckAndDraftFees.DraftBuyingFee)
			feeDetails.CheckAndDraftFees.DraftSellingFee = utils.SplitAndCleanText(utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-DraftSellingFee td.%s", col)), "-")
			feeDetails.CheckAndDraftFees.DraftSellingFeeNumeric = extractCheckAndDraftFeeNumeric(feeDetails.CheckAndDraftFees.DraftSellingFee)
			feeDetails.CheckAndDraftFees.ForeignBillBuyingFee = utils.SplitAndCleanText(utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-ForeignBillBuyingFee td.%s", col)), "-")
			feeDetails.CheckAndDraftFees.ForeignBillBuyingFeeNumeric = extractCheckAndDraftFeeNumeric(feeDetails.CheckAndDraftFees.ForeignBillBuyingFee)
			feeDetails.CheckAndDraftFees.ForeignBillSellingFee = utils.SplitAndCleanText(utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-ForeignBillSellingFee td.%s", col)), "-")
			feeDetails.CheckAndDraftFees.ExchangeCompensationFee = utils.SplitAndCleanText(utils.ExtractCompensationFeeFromTd(doc, fmt.Sprintf(".attr-ExchangeCompensationFee td.%s", col)), "-")

			// Extract and populate LetterOfCreditFees
			feeDetails.LetterOfCreditFees.ForeignLC = utils.CleanText(utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-ForeignLetterOfCreditFee td.%s", col)))
			feeDetails.LetterOfCreditFees.ForeignLCNumeric = extractLCFeeNumeric(feeDetails.LetterOfCreditFees.ForeignLC)
			feeDetails.LetterOfCreditFees.DomesticLC = utils.CleanText(utils.ExtractTextInsideElement(doc, fmt.Sprintf(".attr-DomesticLetterOfCreditFee td.%s", col)))
			feeDetails.LetterOfCreditFees.DomesticLCNumeric = extractLCFeeNumeric(feeDetails.LetterOfCreditFees.DomesticLC)

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

			internationalFees = append(internationalFees, feeDetails)
		}
		time.Sleep(2 * time.Second)
	}

	file, err := json.MarshalIndent(internationalFees, "", " ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	err = os.WriteFile("international_fee.json", file, 0644)
	if err != nil {
		log.Fatalf("Failed to write JSON to file: %v", err)
	}

	fmt.Println("Data successfully saved to international_fee.json")
}

// Helper function to extract numeric values from fee text
func extractNumericValues(text []string) interface{} {
	var numbers []float64
	for _, part := range text {
		re := regexp.MustCompile(`\d+(\.\d+)?`)
		matches := re.FindAllString(part, -1)
		for _, match := range matches {
			num, err := strconv.ParseFloat(match, 64)
			if err == nil {
				numbers = append(numbers, num)
			}
		}
	}
	return numbers
}

// Helper function to extract numeric values from conditions text into a structured ConditionsNumeric
func extractConditionsNumeric(conditions []string) models.ConditionsNumeric {
	var conditionsNumeric models.ConditionsNumeric
	for _, condition := range conditions {
		if strings.Contains(condition, "ตั้งแต่") && strings.Contains(condition, "บาท") {
			re := regexp.MustCompile(`\d+`)
			numbers := re.FindAllString(condition, -1)
			if len(numbers) >= 2 {
				conditionsNumeric.TransactionRange = []float64{
					parseStringToFloat(numbers[0]),
					parseStringToFloat(numbers[1]),
				}
			}
		}
		if strings.Contains(condition, "เรียกเก็บที่") {
			re := regexp.MustCompile(`\d+`)
			number := re.FindString(condition)
			fee := parseStringToFloat(number)
			conditionsNumeric.FeePerTransaction = &fee
		}
		if strings.Contains(condition, "ตั้งแต่") && strings.Contains(condition, "ถึงวันที่") {
			dates := extractDatesFromCondition(condition)
			if len(dates) == 2 {
				startDate, _ := time.Parse("2 January 2006", dates[0])
				endDate, _ := time.Parse("2 January 2006", dates[1])
				conditionsNumeric.PromotionStartDate = &startDate
				conditionsNumeric.PromotionEndDate = &endDate
			}
		}
		if strings.Contains(condition, "ยกเลิก") && strings.Contains(condition, "เรียกเก็บ") {
			re := regexp.MustCompile(`\d+`)
			number := re.FindString(condition)
			cancellationFee := parseStringToFloat(number)
			conditionsNumeric.CancellationFee = &cancellationFee
		}
	}
	return conditionsNumeric
}

// Helper function to extract numeric values from CheckAndDraft fee text
func extractCheckAndDraftFeeNumeric(text []string) models.CheckAndDraftFeeNumeric {
	var feeNumeric models.CheckAndDraftFeeNumeric
	for _, line := range text {
		if strings.Contains(line, "บาท") && strings.Contains(line, "ฉบับ") {
			re := regexp.MustCompile(`\d+`)
			numbers := re.FindAllString(line, -1)
			if len(numbers) >= 1 {
				baseFee := parseStringToFloat(numbers[0])
				feeNumeric.BaseFee = &baseFee
			}
			if len(numbers) >= 2 {
				stampDuty := parseStringToFloat(numbers[1])
				feeNumeric.StampDuty = &stampDuty
			}
		}
		if strings.Contains(line, "เช็คคืน") {
			re := regexp.MustCompile(`\d+`)
			number := re.FindString(line)
			returnFee := parseStringToFloat(number)
			feeNumeric.ReturnFee = &returnFee
		}
		if strings.Contains(line, "Stop Payment") {
			re := regexp.MustCompile(`\d+`)
			number := re.FindString(line)
			stopPaymentFee := parseStringToFloat(number)
			feeNumeric.StopPaymentFee = &stopPaymentFee
		}
	}
	return feeNumeric
}

// Helper function to extract numeric values from LC fees text
func extractLCFeeNumeric(text string) models.LCFeesNumeric {
	var lcNumeric models.LCFeesNumeric
	re := regexp.MustCompile(`(\d+(\.\d+)?%)|(\d+(\.\d+)? บาท)`)
	matches := re.FindAllString(text, -1)
	for _, match := range matches {
		if strings.Contains(match, "%") {
			percent := parseStringToFloat(strings.TrimSuffix(match, "%"))
			lcNumeric.PercentFee = &percent
		} else if strings.Contains(match, "บาท") {
			amount := parseStringToFloat(strings.TrimSuffix(match, " บาท"))
			lcNumeric.MinFee = &amount
		} else {
			otherFee := parseStringToFloat(match)
			lcNumeric.OtherFees = append(lcNumeric.OtherFees, otherFee)
		}
	}
	return lcNumeric
}

// Helper function to parse string to float
func parseStringToFloat(s string) float64 {
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return num
}

// Helper function to extract dates from condition text
func extractDatesFromCondition(condition string) []string {
	re := regexp.MustCompile(`\d{1,2} \w+ \d{4}`)
	matches := re.FindAllString(condition, 2)
	return matches
}
