package parser

import (
	"digitalbanking_fees/models"
	"digitalbanking_fees/utils"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseServiceDetails(doc *goquery.Document, col string, index int) models.ServiceDetails {
	serviceType := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-ServiceTypeId .cmpr-col.%s span", col)).Text())
	mainFeature := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-ServiceMainCharacteristic .cmpr-col.%s span", col)).Text())

	customerText := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-CustomerCharacterApplyCondition .cmpr-col.%s span", col)).Text())
	customerGroup := utils.SplitTextByDelimeterAndNumber(customerText)

	return models.ServiceDetails{
		ServiceType: serviceType,
		MainFeature: mainFeature,
		CustomGroup: customerGroup,
	}
}

func ParseFeeDetails(doc *goquery.Document, col string, index int) models.FeeDetails {
	promptPayFee := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-PromptPayTransferFee .cmpr-col.%s span", col)).Text())

	interbankText := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-InterbankTransferFee .cmpr-col.%s span", col)).Text())
	interbankFees := utils.SplitTextByDelimeter(interbankText)

	intrabankFee := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-IntrabankTransferFee .cmpr-col.%s span", col)).Text())

	cardlessText := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-CardlessCashWithdrawalFee .cmpr-col.%s span", col)).Text())
	cardlessFees := utils.SplitTextByDelimeter(cardlessText)

	entranceFee := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-EntranceFee .cmpr-col.%s span", col)).Text())
	annualFee := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-AnnualFee .cmpr-col.%s span", col)).Text())

	otherFeesText := utils.CleanText(doc.Find(fmt.Sprintf("tr.attr-OtherFee .cmpr-col.%s span", col)).Text())
	otherFees := utils.SplitTextByDelimeterAndNumber(otherFeesText)

	return models.FeeDetails{
		PromptPayTransferFee:  promptPayFee,
		InterbankTransferFee:  interbankFees,
		IntrabankTransferFee:  intrabankFee,
		CardlessWithdrawalFee: cardlessFees,
		EntranceFee:           entranceFee,
		AnnualFee:             annualFee,
		OtherFees:             otherFees,
	}
}

func ParseAdditionalDetails(doc *goquery.Document, col string, index int) models.AdditionalDetails {
	serviceWebsite := doc.Find(fmt.Sprintf("tr.attr-Url .cmpr-col.%s a", col)).AttrOr("href", "")
	feeWebsite := doc.Find(fmt.Sprintf("tr.attr-Feeurl .cmpr-col.%s a", col)).AttrOr("href", "")

	var serviceWebsitePtr, feeWebsitrPtr *string

	if strings.TrimSpace(serviceWebsite) != "" {
		serviceWebsitePtr = &serviceWebsite
	}
	if strings.TrimSpace(feeWebsite) != "" {
		feeWebsitrPtr = &feeWebsite
	}

	return models.AdditionalDetails{
		ServiceWebsite: serviceWebsitePtr,
		FeeWebsite:     feeWebsitrPtr,
	}
}
