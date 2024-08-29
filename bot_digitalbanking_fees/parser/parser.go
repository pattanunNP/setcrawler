package parser

import (
	"digitalbanking_fees/models"
	"digitalbanking_fees/utils"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseServiceDetails(doc *goquery.Document, col string, index int) models.ServiceDetails {
	// Always check if selection is found before accessing
	serviceTypeSelection := doc.Find(fmt.Sprintf("tr.attr-ServiceTypeId .cmpr-col.%s span", col))
	serviceType := ""
	if serviceTypeSelection != nil {
		serviceType = utils.CleanText(serviceTypeSelection.Text())
	}

	mainFeatureSelection := doc.Find(fmt.Sprintf("tr.attr-ServiceMainCharacteristic .cmpr-col.%s span", col))
	mainFeature := ""
	if mainFeatureSelection != nil {
		mainFeature = utils.CleanText(mainFeatureSelection.Text())
	}

	customerTextSelection := doc.Find(fmt.Sprintf("tr.attr-CustomerCharacterApplyCondition .cmpr-col.%s span", col))
	customerText := ""
	if customerTextSelection != nil {
		customerText = utils.CleanText(customerTextSelection.Text())
	}

	customerGroups := parseCustomerGroups(customerText)

	return models.ServiceDetails{
		Type:           serviceType,
		MainFeature:    mainFeature,
		CustomerGroups: customerGroups,
	}
}

func ParseFeeDetails(doc *goquery.Document, col string, index int) models.FeeDetails {
	// Parse PromptPay Fee
	promptPayFeeText := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-PromptPayTransferFee .cmpr-col.%s span", col)).Text())
	promptPayFee := parseFee(promptPayFeeText)

	// Parse Interbank Transfer Fees
	interbankText := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-InterbankTransferFee .cmpr-col.%s span", col)).Text())
	interbankFees := parseInterbankFees(interbankText)

	// Parse Intrabank Transfer Fee
	intrabankFeeText := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-IntrabankTransferFee .cmpr-col.%s span", col)).Text())
	intrabankFee := parseFee(intrabankFeeText)

	// Parse Cardless Withdrawal Fees
	cardlessText := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-CardlessCashWithdrawalFee .cmpr-col.%s span", col)).Text())
	cardlessFees := parseCardlessFees(cardlessText)

	// Parse Entrance Fee
	entranceFeeText := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-EntranceFee .cmpr-col.%s span", col)).Text())
	entranceFee := parseFee(entranceFeeText)

	// Parse Annual Fee
	annualFeeText := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-AnnualFee .cmpr-col.%s span", col)).Text())
	annualFee := parseFee(annualFeeText)

	// Parse Other Fees
	otherFeesText := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-OtherFee .cmpr-col.%s span", col)).Text())
	otherFees := parseOtherFees(otherFeesText)

	return models.FeeDetails{
		PromptPayTransfer:  promptPayFee,
		InterbankTransfer:  interbankFees,
		IntrabankTransfer:  intrabankFee,
		CardlessWithdrawal: cardlessFees,
		EntranceFee:        entranceFee,
		AnnualFee:          annualFee,
		OtherFees:          otherFees,
	}
}

func ParseAdditionalDetails(doc *goquery.Document, col string, index int) models.AdditionalDetails {
	serviceWebsite := doc.Find(fmt.Sprintf("tr.attr-Url .cmpr-col.%s a", col)).AttrOr("href", "")
	feeWebsite := doc.Find(fmt.Sprintf("tr.attr-Feeurl .cmpr-col.%s a", col)).AttrOr("href", "")

	var serviceWebsitePtr, feeWebsitePtr *string

	if strings.TrimSpace(serviceWebsite) != "" {
		serviceWebsitePtr = &serviceWebsite
	}
	if strings.TrimSpace(feeWebsite) != "" {
		feeWebsitePtr = &feeWebsite
	}

	return models.AdditionalDetails{
		ServiceWebsite: serviceWebsitePtr,
		FeeWebsite:     feeWebsitePtr,
	}
}

// Helper function to parse customer group details dynamically
func parseCustomerGroups(text string) []models.CustomerGroup {
	var customerGroups []models.CustomerGroup

	// Example logic to extract customer details from text
	// This assumes customer details are separated by semicolons or commas
	parts := utils.SplitTextByDelimeter(text)
	for _, part := range parts {
		description := "N/A"
		ageRequirement := "N/A"
		accountRequirements := []string{}

		// Parse age requirement if mentioned in the text
		if strings.Contains(part, "อายุตั้งแต่") {
			ageRequirement = part
		}

		// Parse description (assuming it might contain specific words)
		if strings.Contains(part, "บุคคลธรรมดา") || strings.Contains(part, "นิติบุคคล") {
			description = part
		}

		// Example logic to extract account requirements
		if strings.Contains(part, "บัญชี") || strings.Contains(part, "บัตร") {
			accountRequirements = append(accountRequirements, part)
		}

		customerGroups = append(customerGroups, models.CustomerGroup{
			Description:         description,
			AgeRequirement:      ageRequirement,
			AccountRequirements: accountRequirements,
		})
	}

	return customerGroups
}

func parseFee(text string) models.Fee {
	feeAmount := extractNumericValue(text)
	amount := *feeAmount
	return models.Fee{
		FeeText:   text,
		FeeAmount: amount,
	}
}

func parseInterbankFees(text string) []models.TransferCondition {
	cleanText := strings.ReplaceAll(text, ",", "")
	parts := utils.SplitTextByDelimeter(cleanText)
	var fees []models.TransferCondition

	for _, part := range parts {
		description := "Interbank Transfer"
		conditionText := extractCondition(part)
		feeAmount := extractNumericValue(part)
		amount := *feeAmount

		// Initialize the ConditionRange correctly, avoid nil dereference
		conditionRange := determineRangeFromCondition(conditionText)

		fees = append(fees, models.TransferCondition{
			Description:    description,
			ConditionText:  conditionText,
			ConditionRange: conditionRange, // Make sure conditionRange is properly handled
			FeeText:        part,
			FeeAmount:      amount, // using int instead of *int
		})
	}

	return fees
}

// Helper function to parse cardless withdrawal fees
func parseCardlessFees(text string) models.Fee {
	feeAmount := extractNumericValue(text)
	return models.Fee{
		Description: "Cardless Withdrawal",
		FeeText:     text,
		FeeAmount:   *feeAmount,
	}
}

func parseOtherFees(text string) []models.OtherFee {
	parts := utils.SplitTextByDelimeterAndNumber(text)
	var otherFees []models.OtherFee

	for _, part := range parts {
		conditionText := extractCondition(part)
		feeAmount := extractNumericValue(part)
		amount := 0 // default value
		if feeAmount != nil {
			amount = *feeAmount
		}
		otherFees = append(otherFees, models.OtherFee{
			Description: "Other Fees",
			Conditions: []models.OtherFeeCondition{
				{
					ConditionText: conditionText,
					Currency: map[string]models.CurrencyFee{
						"THB": {
							FeeText:   part,
							FeeAmount: amount, // using int instead of *int
						},
					},
				},
			},
		})
	}

	return otherFees
}

func extractNumericValue(text string) *int {
	cleanedText := strings.ReplaceAll(text, ",", "")

	if strings.Contains(cleanedText, "ไม่มีค่าธรรมเนียม") || strings.Contains(cleanedText, "ไม่มีค่าบริการ") {
		zero := 0
		return &zero
	}

	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(cleanedText, -1)

	if len(matches) == 0 {
		zero := 0
		return &zero
	}

	val, err := strconv.Atoi(matches[0])
	if err != nil {
		zero := 0
		return &zero
	}

	return &val
}

func extractCondition(text string) string {

	re := regexp.MustCompile(`(ไม่เกิน|เกิน)\s?(\d{1,3}(,\d{3})*)(\s?บาท)?`)
	matches := re.FindStringSubmatch(text)

	if len(matches) > 0 {
		conditionType := matches[1] // ไม่เกิน or เกิน
		amount := matches[2]        // Amount like 100,000
		condition := fmt.Sprintf("โอนเงิน%s %s บาท", conditionType, amount)
		return condition
	}

	// If no specific pattern is found, return the original text
	return text
}

// Utility function to determine range from condition text
func determineRangeFromCondition(condition string) models.Range {
	cleanedCondition := strings.ReplaceAll(condition, ",", "")

	if strings.Contains(cleanedCondition, "ไม่เกิน") {
		re := regexp.MustCompile(`\d+`)
		matches := re.FindAllString(cleanedCondition, -1)
		if len(matches) > 0 {
			max, err := strconv.Atoi(matches[0])
			if err == nil {
				return models.Range{Min: 0, Max: max}
			}
		}
	} else if strings.Contains(cleanedCondition, "เกิน") {
		re := regexp.MustCompile(`\d+`)
		matches := re.FindAllString(cleanedCondition, -1)
		if len(matches) > 0 {
			min, err := strconv.Atoi(matches[0])
			if err == nil {
				return models.Range{Min: min, Max: 1000000} // Assume upper bound
			}
		}
	}

	return models.Range{Min: 0, Max: 0} // Default to zero range if no condition found

}
