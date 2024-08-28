package parser

import (
	"emoney_fees/models"
	"emoney_fees/utils"
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
	return models.GeneralFees{
		EntranceFee:        utils.CleanText(doc.Find("tr.attr-EntranceFeeAmount td.cmpr-col").First().Text()),
		AnnualFee:          utils.CleanText(doc.Find("tr.attr-AnnualFee td.cmpr-col").First().Text()),
		CardReplacementFee: utils.CleanText(doc.Find("tr.attr-CardReplacementFee td.cmpr-col").First().Text()),
		MaintenaceFee:      utils.CleanText(doc.Find("tr.attr-ProductMaintenanceFee td.cmpr-col").First().Text()),
	}
}

func ParseSpendingFees(doc *goquery.Document, colIndex int) models.SpendingFees {
	return models.SpendingFees{
		SpendingFee:           utils.CleanText(doc.Find("tr.attr-SpendingFee td.cmpr-col").First().Text()),
		SpendingAlertFee:      utils.CleanText(doc.Find("tr.attr-SpendingAlertFee td.cmpr-col").First().Text()),
		OverseasWithdrawalFee: utils.CleanText(doc.Find("tr.attr-OverseasCashWithdrawalFee td.cmpr-col").First().Text()),
		CurrencyConversionFee: utils.CleanText(doc.Find("tr.attr-CurrencyConversionRiskFeeRate td.cmpr-col").First().Text()),
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

	doc.Find("tr.attr-OtherFees td.cmpr-col").Each(func(i int, s *goquery.Selection) {
		text := utils.CleanText(s.Text())
		// Split by numbers like "1., 2." and "-"
		parts := strings.FieldsFunc(text, func(r rune) bool {
			return r == '1' || r == '2' || r == '-' || r == ' ' || r == '.'
		})
		if len(parts) > 0 {
			otherFeeDetails = append(otherFeeDetails, parts...)
		} else {
			otherFeeDetails = append(otherFeeDetails, "null")
		}
	})

	otherFeeDetails = utils.CleanAndFilter(otherFeeDetails)

	return models.OtherFees{
		OtherFeesDetails: otherFeeDetails,
	}
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
