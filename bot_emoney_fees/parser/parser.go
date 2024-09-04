package parser

import (
	"emoney_fees/models"
	"emoney_fees/utils"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseTopUpDetails(doc *goquery.Document, colIndex int) models.TopUpDetails {
	noFeeChannels := []string{}
	feeChannels := ""

	doc.Find("tr.attr-TopUpChannelsWithoutFee td.cmpr-col").Each(func(i int, s *goquery.Selection) {
		text := utils.CleanText(s.Text())
		noFeeChannels = append(noFeeChannels, strings.Split(text, "-")...)
	})
	noFeeChannels = utils.CleanAndFilter(noFeeChannels)

	doc.Find("tr.attr-TopUpChannelsWithFee td.cmpr-col").Each(func(i int, s *goquery.Selection) {
		feeChannels = utils.CleanText(s.Text())
	})

	return models.TopUpDetails{
		NoFeeChannels: noFeeChannels,
		FeeChannels:   feeChannels,
	}
}

func ParseGeneralFees(doc *goquery.Document, colIndex int) models.GeneralFees {
	// Extract fee details as raw text
	cardReplacementFee := utils.CleanText(doc.Find("tr.attr-CardReplacementFee td.cmpr-col").First().Text())

	// Parse amount and condition from raw text
	cardReplacementAmount, cardReplacementCond := extractAmountAndCondition(cardReplacementFee)

	return models.GeneralFees{
		EntranceFee:           utils.CleanText(doc.Find("tr.attr-EntranceFeeAmount td.cmpr-col").First().Text()),
		AnnualFee:             utils.CleanText(doc.Find("tr.attr-AnnualFee td.cmpr-col").First().Text()),
		CardReplacementFee:    cardReplacementFee,
		CardReplacementAmount: cardReplacementAmount,
		CardReplacementCond:   cardReplacementCond,
		MaintenanceFee:        utils.CleanText(doc.Find("tr.attr-ProductMaintenanceFee td.cmpr-col").First().Text()),
	}
}

func ParseSpendingFees(doc *goquery.Document, colIndex int) models.SpendingFees {
	// Extract fee details as raw text
	overseasWithdrawalFee := utils.CleanText(doc.Find("tr.attr-OverseasCashWithdrawalFee td.cmpr-col").First().Text())
	currencyConversionFee := utils.CleanText(doc.Find("tr.attr-CurrencyConversionRiskFeeRate td.cmpr-col").First().Text())

	// Parse amounts and conditions from raw text
	overseasWithdrawalAmount, overseasWithdrawalCond := extractAmountAndCondition(overseasWithdrawalFee)
	currencyConversionRate, currencyConversionCond := extractRateAndCondition(currencyConversionFee)

	return models.SpendingFees{
		SpendingFee:              utils.CleanText(doc.Find("tr.attr-SpendingFee td.cmpr-col").First().Text()),
		SpendingAlertFee:         utils.CleanText(doc.Find("tr.attr-SpendingAlertFee td.cmpr-col").First().Text()),
		OverseasWithdrawalFee:    overseasWithdrawalFee,
		OverseasWithdrawalAmount: overseasWithdrawalAmount,
		OverseasWithdrawalCond:   overseasWithdrawalCond,
		CurrencyConversionFee:    currencyConversionFee,
		CurrencyConversionRate:   currencyConversionRate,
		CurrencyConversionCond:   currencyConversionCond,
	}
}

func ParseTerminationFees(doc *goquery.Document, colIndex int) models.TerminationFees {
	cashRefundFee := utils.CleanText(doc.Find("tr.attr-CashRefundFee td.cmpr-col").First().Text())
	terminationFeeStr := utils.CleanText(doc.Find("tr.attr-TerminationFee td.cmpr-col").First().Text())

	terminationFee := 0
	if terminationFeeStr != "ไม่มีค่าธรรมเนียม" {
		terminationFee, _ = strconv.Atoi(strings.TrimSpace(strings.ReplaceAll(terminationFeeStr, " บาท", "")))
	}

	return models.TerminationFees{
		CashRefundFee:   cashRefundFee,
		TerminationFees: terminationFee,
	}
}

func ParseOtherFes(doc *goquery.Document, colIndex int) models.OtherFees {
	otherFeeDetails := []string{}

	// Extract the raw text from the document
	doc.Find("tr.attr-OtherFees td.cmpr-col").Each(func(i int, s *goquery.Selection) {
		text := utils.CleanText(s.Text())
		otherFeeDetails = append(otherFeeDetails, text)
	})

	// Combine all extracted text into a single string for easier processing
	rawDetails := strings.Join(otherFeeDetails, " ")

	// Format the extracted text properly
	formattedDetails := formatFeeDetails(rawDetails)

	return models.OtherFees{
		OtherFeesDetails: []string{formattedDetails},
	}
}

func formatFeeDetails(details string) string {
	// Replace specific patterns with newlines for better formatting
	details = strings.ReplaceAll(details, "1. ", "\n1. ")
	details = strings.ReplaceAll(details, "2. ", "\n\n2. ")
	details = strings.ReplaceAll(details, "- ", "\n- ")
	details = strings.ReplaceAll(details, "-\n", "\n- ")
	details = strings.ReplaceAll(details, " (", "\n(")
	details = strings.ReplaceAll(details, " ต่อรายการ", " ต่อรายการ\n")
	details = strings.ReplaceAll(details, " คิด", "\nคิด")
	details = strings.ReplaceAll(details, " โดย", "\nโดย")
	details = strings.ReplaceAll(details, "-วง", "\n-วง")
	details = strings.ReplaceAll(details, " จำกัด", "\nจำกัด")
	details = strings.ReplaceAll(details, " 1.", "\n1.")
	details = strings.ReplaceAll(details, " -", "\n-")
	details = strings.ReplaceAll(details, " 2.", "\n2.")
	details = strings.ReplaceAll(details, " -", "\n-")
	details = strings.ReplaceAll(details, " 3.", "\n3.")
	details = strings.ReplaceAll(details, " ,", ",")

	// Trim any leading/trailing whitespace and return
	return strings.TrimSpace(details)
}

func AdditionalInfo(doc *goquery.Document, colIndex int) models.AdditionalInfo {
	feeURL := ""

	doc.Find("tr.attr-FeeURL td.cmpr-col a.prod-url").Each(func(i int, s *goquery.Selection) {
		feeURL, _ = s.Attr("href")
	})

	return models.AdditionalInfo{
		FeeURL: feeURL,
	}
}

// Extract amount and conditions from fee description
func extractAmountAndCondition(fee string) (float64, string) {
	amountRegex := regexp.MustCompile(`(\d+\.?\d*) บาท`)
	conditionRegex := regexp.MustCompile(`เงื่อนไข: (.+)`)

	amountMatch := amountRegex.FindStringSubmatch(fee)
	conditionMatch := conditionRegex.FindStringSubmatch(fee)

	var amount float64
	var condition string

	if len(amountMatch) > 1 {
		amount, _ = strconv.ParseFloat(amountMatch[1], 64)
	}

	if len(conditionMatch) > 1 {
		condition = conditionMatch[1]
	}

	return amount, condition
}

// Extract percentage rate and conditions from fee description
func extractRateAndCondition(fee string) (float64, string) {
	rateRegex := regexp.MustCompile(`(\d+\.?\d*)%`)
	conditionRegex := regexp.MustCompile(`เงื่อนไข: (.+)`)

	rateMatch := rateRegex.FindStringSubmatch(fee)
	conditionMatch := conditionRegex.FindStringSubmatch(fee)

	var rate float64
	var condition string

	if len(rateMatch) > 1 {
		rate, _ = strconv.ParseFloat(rateMatch[1], 64)
	}

	if len(conditionMatch) > 1 {
		condition = conditionMatch[1]
	}

	return rate, condition
}
