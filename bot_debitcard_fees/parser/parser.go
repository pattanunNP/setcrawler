package parser

import (
	"debitcard_fees/models"
	"debitcard_fees/utils"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseGeneralFees(doc *goquery.Document, index int) models.Fee {
	col := "col" + strconv.Itoa(index)
	annualFeeText := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-CardHolderAnnualFee td.cmpr-col.%s span", col)).Text())
	annualFeeConditions := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-CardHolderAnnualFee td.cmpr-col.%s .text-primary", col)).Text())

	var annualFee *models.AnnualFeeDetail

	// Check if the text indicates no fee
	if strings.Contains(annualFeeText, "ไม่มีค่าธรรมเนียม") {
		annualFee = &models.AnnualFeeDetail{
			Amount:     0,
			Conditions: nil, // Set to nil or empty string if there are no conditions
		}
	} else {
		// Attempt to extract a numerical value from the text
		re := regexp.MustCompile(`\d+(?:,\d+)*`)
		amountStr := re.FindString(annualFeeText)
		amountStr = strings.ReplaceAll(amountStr, ",", "")
		amount := 0

		if amountStr != "" {
			var err error
			amount, err = strconv.Atoi(amountStr)
			if err != nil {
				log.Printf("Error parsing annual fee amount: %v", err)
				amount = 0
			}
		}

		// Only set conditions if they exist
		var conditions *string
		if annualFeeConditions != "" {
			conditions = &annualFeeConditions
		}
		annualFee = &models.AnnualFeeDetail{
			Amount:     amount,
			Conditions: conditions,
		}
	}

	return models.Fee{
		EntranceFee:                utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-CardHolderEntranceFee td.cmpr-col.%s span", col)).Text()),
		AnnualFee:                  annualFee,
		CardReplacementFee:         parseOptionalArray(doc, fmt.Sprintf("tr.attr-CardReplacementFee td.cmpr-col.%s span", col)),
		PINReplacementFee:          parseOptionalString(doc, fmt.Sprintf("tr.attr-CardPINReplacement td.cmpr-col.%s span", col)),
		StatementRequestFee:        parseOptionalArray(doc, fmt.Sprintf("tr.attr-CopyStatementFee td.cmpr-col.%s span", col)),
		TransactionSlipRequestFee:  parseOptionalString(doc, fmt.Sprintf("tr.attr-CopySaleSlipFee td.cmpr-col.%s span", col)),
		TransactionVerificationFee: parseOptionalString(doc, fmt.Sprintf("tr.attr-TransactionVerification td.cmpr-col.%s span", col)),
	}
}

func ParseDomesticFees(doc *goquery.Document, index int) models.DomesticTransaction {
	col := "col" + strconv.Itoa(index)
	freeTransactionText := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-FeeInternal td.cmpr-col.%s span", col)).Text())
	freeTransactionArray := parseOptionalArrayFromString(freeTransactionText, "-")

	return models.DomesticTransaction{
		FreeTransactionCount:    freeTransactionArray,
		BalanceInquiryFeeOut:    parseOptionalString(doc, fmt.Sprintf("tr.attr-KioskCheckBalanceFee td.cmpr-col.%s span", col)),
		WithdrawFeeOut:          parseOptionalString(doc, fmt.Sprintf("tr.attr-KioskWitddrawFee td.cmpr-col.%s span", col)),
		TransferFeeOut:          parseOptionalString(doc, fmt.Sprintf("tr.attr-KiosTransferFee td.cmpr-col.%s span", col)),
		BalanceInquiryFeeIn:     parseOptionalString(doc, fmt.Sprintf("tr.attr-KioskBalanceInFee td.cmpr-col.%s span", col)),
		WithdrawFeeIn:           parseOptionalString(doc, fmt.Sprintf("tr.attr-KioskWithdrawInFee td.cmpr-col.%s span", col)),
		TransferFeeIn:           parseOptionalString(doc, fmt.Sprintf("tr.attr-KioskTransferInFee td.cmpr-col.%s span", col)),
		BalanceInquiryFeeOutAlt: parseOptionalString(doc, fmt.Sprintf("tr.attr-KioskBalanceOutFee td.cmpr-col.%s span", col)),
		WithdrawFeeOutAlt:       parseOptionalString(doc, fmt.Sprintf("tr.attr-KioskWithdrawOutFee td.cmpr-col.%s span", col)),
		TransferFeeOutAlt:       parseOptionalString(doc, fmt.Sprintf("tr.attr-KioskTransferOutFee td.cmpr-col.%s span", col)),
		TransferLimit10k:        parseOptionalString(doc, fmt.Sprintf("tr.attr-KioskTransfer10kFee td.cmpr-col.%s span", col)),
		TransferLimit50k:        parseOptionalString(doc, fmt.Sprintf("tr.attr-KioskTransfer50kFee td.cmpr-col.%s span", col)),
		AdditionalFee:           parseOptionalString(doc, fmt.Sprintf("tr.attr-KioskOtherFee td.cmpr-col.%s span", col)),
	}
}

func ParseInternationalFees(doc *goquery.Document, index int) *models.InternationalTransaction {
	col := "col" + strconv.Itoa(index)
	withdrawalFee := parseOptionalString(doc, fmt.Sprintf("tr.attr-InteralWithdrawFee td.cmpr-col.%s span", col))
	balanceInquiryFee := parseOptionalString(doc, fmt.Sprintf("tr.attr-InteralBalnace td.cmpr-col.%s span", col))
	currencyExchangeFee := parseOptionalString(doc, fmt.Sprintf("tr.attr-FXRiskCost td.cmpr-col.%s span", col))

	if withdrawalFee == nil && balanceInquiryFee == nil && currencyExchangeFee == nil {
		return nil
	}

	return &models.InternationalTransaction{
		WithdrawalFee:       withdrawalFee,
		BalanceInquiryFee:   balanceInquiryFee,
		CurrencyExchangeFee: currencyExchangeFee,
	}
}

func ParseOtherFees(doc *goquery.Document, index int) *models.OtherFees {
	col := "col" + strconv.Itoa(index)
	otherFeesText := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-OtherFee td.cmpr-col.%s span", col)).Text())
	otherFeesArray := parseOptionalArrayFromString(otherFeesText, "-")

	if otherFeesArray == nil {
		return nil
	}

	return &models.OtherFees{
		OtherFees: otherFeesArray,
	}
}

func ParseAdditionalInfo(doc *goquery.Document, index int) *models.AdditionalInfo {
	col := "col" + strconv.Itoa(index)
	link := doc.Find(fmt.Sprintf("tr.attr-feeurl td.cmpr-col.%s a", col)).AttrOr("href", "")

	if link == "" {
		return nil
	}

	return &models.AdditionalInfo{
		FeeWebsite: &link,
	}
}

// Helper function to parse optional string
func parseOptionalString(doc *goquery.Document, selector string) *string {
	text := utils.CleanText(doc.Find(selector).Text())
	if text == "" {
		return nil
	}
	return &text
}

func parseOptionalArray(doc *goquery.Document, selector string) *[]string {
	text := utils.CleanText(doc.Find(selector).Text())
	return parseOptionalArrayFromString(text, "-")
}

// Helper function to parse optional array from string using a delimiter
func parseOptionalArrayFromString(text, delimiter string) *[]string {
	if text == "" {
		return nil
	}

	parts := utils.SplitText(text, delimiter)
	cleanedParts := []string{}
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			cleanedParts = append(cleanedParts, part)
		}
	}

	if len(cleanedParts) == 0 {
		return nil
	}

	return &cleanedParts
}
