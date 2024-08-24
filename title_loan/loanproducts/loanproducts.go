package loanproducts

import (
	"bytes"
	"fmt"
	"strconv"

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

func FetchLoanProducts(url, payload string) ([]TitleLoan, error) {
	// Perform the HTTP POST request
	body, err := SendHTTPRequest(url, payload)
	if err != nil {
		return nil, err
	}

	// Parse the HTML response
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %v", err)
	}

	var titleLoans []TitleLoan
	// Iterate over columns to extract provider and product info
	for i := 1; i <= 3; i++ {
		colClass := "col" + strconv.Itoa(i)

		// Extract data for each product and trim spaces
		provider := CleanText(doc.Find(fmt.Sprintf("th.%s span", colClass)).Eq(1).Text())
		product := CleanText(doc.Find(fmt.Sprintf("th.prod-%s span", colClass)).Text())

		vehicleType := getVehicleType(doc, i)
		vehicleCondition := getVehicleCondition(doc, i)
		loanType := getLoanType(doc, i)
		interestRate := getInterestRate(doc, i)
		creditLimitAndInstallment := getCreditLimitAndInstallment(doc, i)
		borroweQualifications := getBorrowerQualifications(doc, i)
		generalFees := getGeneralFees(doc, i)
		cardFee := getCardFees(doc, i)
		paymentFees := getPaymentFees(doc, i)
		otherFees := getOtherFees(doc, i)
		additionalInfo := getAdditionInfo(doc, i)

		titleLoan := TitleLoan{
			Provider:                  provider,
			Product:                   product,
			VehicleType:               vehicleType,
			VehicleCondition:          vehicleCondition,
			LoanType:                  loanType,
			InterestRate:              interestRate,
			CreditLimitAndInstallment: creditLimitAndInstallment,
			BorrowerQualifications:    borroweQualifications,
			GeneralFees:               generalFees,
			CardFees:                  cardFee,
			PaymentFees:               paymentFees,
			OtherFees:                 otherFees,
			AdditionalInfo:            additionalInfo,
		}
		titleLoans = append(titleLoans, titleLoan)
	}

	return titleLoans, nil
}
