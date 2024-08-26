package parser

import (
	"bytes"
	"fmt"
	"nano_finance/pkg/formatter"
	model "nano_finance/pkg/models"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func DetermineTotalPages(body []byte) int {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Printf("Error loading HTML to determine total pages: %v\n", err)
		return 0
	}

	totalPages := 1
	doc.Find("ul.pagination li.MoveLast a").Each(func(i int, s *goquery.Selection) {
		dataPage, exists := s.Attr("data-page")
		if exists {
			totalPages, err = strconv.Atoi(dataPage)
			if err != nil {
				fmt.Printf("Error converting data-page to int: %v\n", err)
				totalPages = 0
			}
		}
	})

	return totalPages
}

// ParseProductData parses the product data from the HTML body
func ParseProductData(body []byte) ([]model.ProductInfo, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %v", err)
	}

	var products []model.ProductInfo

	for i := 1; i <= 3; i++ {
		col := fmt.Sprintf("col%d", i)
		serviceProvider := formatter.CleanString(doc.Find("th." + col + " .col-s span").Eq(1).Text())
		product := formatter.CleanString(doc.Find("th.prod-" + col + " span").First().Text())

		// Extract details
		detailsRaw := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-ProductCondition td.%s span", col)).First().Text())
		details := formatter.SplitAndFormatArray(detailsRaw)

		// Extract credit line and merge all its details into a single array
		creditLineMain := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-CreditLine td.%s span", col)).First().Text())
		creditLineDetails := formatter.SplitAndFormatArray(formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-CreditLine td.%s div span", col)).Text()))
		creditLine := append([]string{creditLineMain}, creditLineDetails...)

		// Extract applicant conditions
		ageText := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-CustomerAge td.%s span", col)).First().Text())
		minAge, maxAge := parseAgeRange(ageText)

		qualificationsRaw := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-ConditionToApply td.%s span", col)).First().Text())
		qualifications := formatter.SplitAndFormatArray(qualificationsRaw)

		applicantConditions := model.ApplicantConditions{
			MinAge:                  minAge,
			MaxAge:                  maxAge,
			ApplicantQualifications: qualifications,
		}

		// Extract loan amount
		loanAmountText := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-CreditLimit td.%s span", col)).First().Text())
		loanAmount := parseLoanAmountRange(loanAmountText)

		// Extract loan duration
		loanDurationText := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-InstallmentPeriod td.%s span", col)).First().Text())
		minMonth, maxMonth := parseLoanDuration(loanDurationText)

		// Extract credit approval conditions
		approvalConditions := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-CreditLimitCondition td.%s span", col)).First().Text())
		if approvalConditions == "" {
			approvalConditions = "null"
		}

		// Extract repayment conditions
		repaymentConditionsRaw := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-InstallmentPeriodCondition td.%s span", col)).First().Text())
		repaymentConditions := formatter.SplitAndFormatArray(repaymentConditionsRaw)

		creditApprovalDetails := model.CreditApprovalConditions{
			LoanAmount:         loanAmount,
			ApprovalConditions: &approvalConditions,
			LoanDuration: model.LoanDuration{
				MinMonth: minMonth,
				MaxMonth: maxMonth,
			},
			RepaymentConditions: model.RepaymentConditions{
				Conditions: repaymentConditions,
			},
		}

		// Extract interest rate details
		interestWithServiceFee := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-InterestWithServiceFee td.%s span", col)).First().Text())
		interestWithServiceFeeConditionRaw := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-InterestWithServiceFeeCondition td.%s span", col)).First().Text())
		interestWithServiceFeeCondition := formatter.SplitAndFormatArray(interestWithServiceFeeConditionRaw)

		defaultInterestRate := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-DefaultInterestRate td.%s span", col)).First().Text())
		defaultInterestRateConditionRaw := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-DefaultInterestRateCondition td.%s span", col)).First().Text())
		defaultInterestRateCondition := formatter.SplitAndFormatArray(defaultInterestRateConditionRaw)

		interestRateDetails := model.InterestRateDetails{
			InterestWithServiceFee:          interestWithServiceFee,
			InterestWithServiceFeeCondition: interestWithServiceFeeCondition,
			DefaultInterestRate:             defaultInterestRate,
			DefaultInterestRateCondition:    defaultInterestRateCondition,
		}

		// Extract payment fees
		freePaymentChannelsRaw := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-FreePaymentChannel td.%s span", col)).First().Text())
		freePaymentChannels := formatter.SplitAndFormatArrayCustom(freePaymentChannelsRaw)

		deductFromProviderAccount := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-DeductingFromBankACFee td.%s span", col)).First().Text())
		deductFromOtherProviderAccount := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-DeductingFromOtherBankACFee td.%s span", col)).First().Text())
		payAtProviderBranch := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-ServiceProviderCounter td.%s span", col)).First().Text())
		payAtOtherBranch := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-OtherProviderCounter td.%s span", col)).First().Text())
		payAtPaymentCounters := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-OthersPaymentCounter td.%s span", col)).First().Text())

		// Updated extraction for online payment channels
		onlinePaymentChannelsRaw := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-OnlinePaymentFee td.%s span", col)).First().Text())
		onlinePaymentChannels := formatter.SplitAndFormatArray(onlinePaymentChannelsRaw)

		atmCdmPaymentChannelsRaw := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-CDMATMPaymentFee td.%s span", col)).First().Text())
		atmCdmPaymentChannels := formatter.SplitAndFormatArrayBy(atmCdmPaymentChannelsRaw, "/")

		phonePaymentChannels := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-PhonePaymentFee td.%s span", col)).First().Text())
		chequeMoneyOrderChannels := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-ChequeMoneyOrderPaymentFee td.%s span", col)).First().Text())
		otherPaymentChannels := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-OtherChannelPaymentFee td.%s span", col)).First().Text())

		paymentFees := model.PaymentFees{
			FreeChannels:              freePaymentChannels,
			DeductFromServiceProvider: deductFromProviderAccount,
			DeductFromOtherProvider:   deductFromOtherProviderAccount,
			ServiceProviderBranch:     payAtProviderBranch,
			OtherProviderBranch:       payAtOtherBranch,
			PaymentCounters:           payAtPaymentCounters,
			OnlineChannels:            onlinePaymentChannels,
			ATMCDMChannels:            atmCdmPaymentChannels,
			TelephoneChannels:         phonePaymentChannels,
			ChequeMoneyOrderChannels:  chequeMoneyOrderChannels,
			OtherChannels:             otherPaymentChannels,
		}

		productWebsiteLink := formatter.CleanString(doc.Find(fmt.Sprintf("tr.attr-URL td.%s a.prod-url", col)).AttrOr("href", ""))
		additionalInfo := model.AdditionalInfo{
			ProductWebste: productWebsiteLink,
		}

		products = append(products, model.ProductInfo{
			ServiceProvider:       serviceProvider,
			Product:               product,
			ProductDetails:        model.ProductDetails{Details: details, CreditLine: creditLine},
			ApplicantConditions:   applicantConditions,
			CreditApprovalDetails: creditApprovalDetails,
			InterestRateDetails:   interestRateDetails,
			PaymentFees:           paymentFees,
			AdditionalInfo:        additionalInfo,
		})
	}

	return products, nil
}

// parseAgeRange parses the min and max age from a string like "20-70 ปี"
func parseAgeRange(ageText string) (*int, *int) {
	parts := strings.Split(ageText, "-")
	if len(parts) == 2 {
		minAge, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		maxAge, err2 := strconv.Atoi(strings.Split(strings.TrimSpace(parts[1]), " ")[0]) // splitting to remove the 'ปี' text
		if err1 == nil && err2 == nil {
			return &minAge, &maxAge
		}
	}
	return nil, nil // default values if parsing fails
}

func parseLoanAmountRange(amountText string) *model.LoanAmount {
	amountText = strings.ReplaceAll(amountText, ",", "")
	parts := strings.Split(amountText, "-")

	if len(parts) == 2 {
		// Attempt to parse minimum and maximum loan amounts
		minAmount, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		maxAmount, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))

		if err1 == nil && err2 == nil {
			// Return struct if parsing is successful
			return &model.LoanAmount{
				MinLoanAmount: minAmount,
				MaxLoanAmount: maxAmount,
			}
		}
	}
	// Return null if parsing fails or no numeric value
	return nil
}

func parseLoanDuration(durationText string) (*int, *int) {
	if durationText == "ไม่กำหนด" {
		return nil, nil
	}
	parts := strings.Split(durationText, "-")
	if len(parts) == 2 {
		minMonth, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		maxMonth, err2 := strconv.Atoi(strings.TrimSpace(strings.Split(parts[1], " ")[0]))
		if err1 == nil && err2 == nil {
			return &minMonth, &maxMonth
		}
	}
	return nil, nil // default values if parsing fails
}
