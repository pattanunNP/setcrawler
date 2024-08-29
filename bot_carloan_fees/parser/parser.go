package parser

import (
	"carloan_fees/models"
	"carloan_fees/utils"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseContractEffectiveDate(doc *goquery.Document, index int) *string {
	col := "col" + strconv.Itoa(index)
	dateText := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-ContractDate td.%s span", col)).Text())
	if dateText == "" {
		return nil
	}
	return &dateText
}

func ParseFeeDetail(doc *goquery.Document, selector string) models.FeeDetail {
	text := utils.CleanText(doc.Find(selector).Text())
	text = strings.ReplaceAll(text, ",", "")
	parts := utils.SplitText(text, "-")

	filteredParts := []string{}
	for _, part := range parts {
		if part != "" {
			filteredParts = append(filteredParts, part)
		}
	}

	min, max := ExtractMinMax(text)
	return models.FeeDetail{
		Text:      filteredParts,
		MinAmount: min,
		MaxAmount: max,
	}
}

func ExtractMinMax(text string) (min *int, max *int) {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(text, -1)
	if len(matches) == 0 {
		return nil, nil
	}

	var minVal, maxVal int
	for i, match := range matches {
		val, _ := strconv.Atoi(match)
		if i == 0 || val < minVal {
			minVal = val
		}
		if val > maxVal {
			maxVal = val
		}
	}

	return &minVal, &maxVal
}

func ParseGeneralFees(doc *goquery.Document, index int) models.GeneralFees {
	col := "col" + strconv.Itoa(index)

	return models.GeneralFees{
		NewVehicleRegistrationFee:                  ParseFeeDetail(doc, fmt.Sprintf("tr.attr-NewVehicleRegister td.%s span", col)),
		OwnershipTransferFeeOneStep:                ParseFeeDetail(doc, fmt.Sprintf("tr.attr-OneStepUponFullPayment td.%s span", col)),
		OwnershipTransferFeeTwoStep:                ParseFeeDetail(doc, fmt.Sprintf("tr.attr-TwoStepUponFullPayment td.%s span", col)),
		VehicleInspectionFee:                       ParseFeeDetail(doc, fmt.Sprintf("tr.attr-VehicleInspection td.%s span", col)),
		ServiceProviderOwnershipTransferFeeOneStep: ParseFeeDetail(doc, fmt.Sprintf("tr.attr-OneStepOwnerTransfer td.%s span", col)),
		ServiceProviderOwnershipTransferFeeTwoStep: ParseFeeDetail(doc, fmt.Sprintf("tr.attr-TwoStepOwnerTransfer td.%s span", col)),
		LeaseTransferFee:                           ParseFeeDetail(doc, fmt.Sprintf("tr.attr-HirePurchaseTransfer td.%s span", col)),
		ContractTerminationFee:                     ParseFeeDetail(doc, fmt.Sprintf("tr.attr-ContractTermination td.%s span", col)),
		LatePaymentPenalty:                         ParseFeeDetail(doc, fmt.Sprintf("tr.attr-PenaltyChargeForLate td.%s span", col)),
		DebtCollectionFee:                          ParseFeeDetail(doc, fmt.Sprintf("tr.attr-DebtCollectionFee td.%s span", col)),
		TaxRenewalFee:                              ParseFeeDetail(doc, fmt.Sprintf("tr.attr-CarTaxRenewal td.%s span", col)),
		LicensePlateProcessingFee:                  ParseFeeDetail(doc, fmt.Sprintf("tr.attr-RegistrationPlate td.%s span", col)),
		RegistrationAddressChangeFee:               ParseFeeDetail(doc, fmt.Sprintf("tr.attr-ChangeAddress td.%s span", col)),
		DocumentCopyServiceFee:                     ParseFeeDetail(doc, fmt.Sprintf("tr.attr-ContractsAndDocuments td.%s span", col)),
	}
}

func ParsePaymentFees(doc *goquery.Document, index int) models.PaymentFees {
	col := "col" + strconv.Itoa(index)

	return models.PaymentFees{
		DirectDebitFromProviderAccount:      ParseFeeDetail(doc, fmt.Sprintf("tr.attr-DirectDebitFromAccountFee td.%s span", col)),
		DirectDebitFromOtherProviderAccount: ParseFeeDetail(doc, fmt.Sprintf("tr.attr-DirectDebitFromAccountFeeOther td.%s span", col)),
		ProviderBranchPayment:               ParseFeeDetail(doc, fmt.Sprintf("tr.attr-BankCounterServiceFee td.%s span", col)),
		OtherBranchPayment:                  ParseFeeDetail(doc, fmt.Sprintf("tr.attr-BankCounterServiceFeeOther td.%s span", col)),
		PaymentServicePoints:                ParseFeeDetail(doc, fmt.Sprintf("tr.attr-CounterServiceFee td.%s span", col)),
		OnlinePayment:                       ParseFeeDetail(doc, fmt.Sprintf("tr.attr-PaymentOnlineFee td.%s span", col)),
		CDMAtmPayment:                       ParseFeeDetail(doc, fmt.Sprintf("tr.attr-PaymentCDMATMFee td.%s span", col)),
		PhonePayment:                        ParseFeeDetail(doc, fmt.Sprintf("tr.attr-PaymentPhoneFee td.%s span", col)),
		ChequeMoneyOrderPayment:             ParseFeeDetail(doc, fmt.Sprintf("tr.attr-PaymentChequeOrMoneyOrderFee td.%s span", col)),
		OtherChannelsPayment:                ParseFeeDetail(doc, fmt.Sprintf("tr.attr-PaymentOtherChannelFee td.%s span", col)),
	}
}

func ParseOtherFees(doc *goquery.Document, index int) models.OtherFees {
	col := "col" + strconv.Itoa(index)

	return models.OtherFees{
		OtherFeesAndCharges: ParseFeeDetail(doc, fmt.Sprintf("tr.attr-other td.%s span", col)),
	}
}

func ParseAdditinoalInfo(doc *goquery.Document, index int) models.AdditionalInfo {
	col := "col" + strconv.Itoa(index)
	link := doc.Find(fmt.Sprintf("tr.attr-Feeurl td.%s a", col)).AttrOr("href", "")

	if link == "" {
		return models.AdditionalInfo{
			FeeWebsiteLinks: nil,
		}
	}

	return models.AdditionalInfo{
		FeeWebsiteLinks: &link,
	}
}
