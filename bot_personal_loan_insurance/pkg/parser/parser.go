package parser

import (
	"bot_personal_insurance/pkg/models"
	"bot_personal_insurance/pkg/utils"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func DetermineTotalPage(body []byte) int {
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

func FetchHTML(url, payload string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, err
	}

	// Add headers to the request
	req.Header.Set("Accept", "text/plain, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Cookie", `verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0"; _uid6672=16B5DEBD.30; _ctout6672=1; _ga_NLQFGWVNXN=GS1.1.1724668043.37.1.1724668084.19.0.0; visit_time=14`)
	req.Header.Set("Origin", "https://app.bot.or.th")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://app.bot.or.th/1213/MCPD/ProductApp/PLoanwithorwithoutCollateral/CompareProduct")
	req.Header.Set("Sec-CH-UA", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	req.Header.Set("VerificationToken", "OYXWhAvsZKJvMmvV8EMX3Tdo1plzeSP_3vwzhjI9OUljk9c160KtMQ4DR0ksFOWu7OliK4nG_59Z6_q7qw1r5vxFVrbId0QhTr_NOyNT32E1,8nzgL9iCG8TF7d5I9rszDVPJHVLsmm2lfxA4E3hYzXvca3cXEYqDGxqcgipt-BcgbfCD26OmIwLhD6-mfwax5XQyzUG-F-slLuL1Ny2xhvM1")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func ParsePersonalLoanDetails(body []byte) ([]models.PersonalLoan, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var loans []models.PersonalLoan

	// Correct the loop condition
	for i := 1; i <= 3; i++ {
		col := fmt.Sprintf("col%d", i)

		// Use the correct selectors with the col variable
		serviceProvider := utils.CleanString(doc.Find(fmt.Sprintf("th.col-s.col-s-%d span", i)).Last().Text())
		product := utils.CleanString(doc.Find(fmt.Sprintf("th.font-black.text-center.prod-%s span", col)).Text())

		productTypeText := utils.CleanString(doc.Find(fmt.Sprintf(".attr-CreditType .%s span", col)).Text())
		productType := strings.Split(productTypeText, "/")

		creditCharacter := utils.CleanString(doc.Find(fmt.Sprintf(".attr-CharacterOfCredit .%s span", col)).Text())

		// Split the collateral text by "-" to create an array and filter out empty strings
		collateralText := utils.CleanString(doc.Find(fmt.Sprintf(".attr-Collateral .%s span", col)).Text())
		collateral := strings.Split(collateralText, "-")
		collateral = utils.FilterEmptyStrings(collateral) // Filter out empty strings
		for i := range collateral {
			collateral[i] = strings.TrimSpace(collateral[i]) // Trim whitespace from each item
		}

		creditLineType := utils.CleanString(doc.Find(fmt.Sprintf(".attr-CreditLineType .%s span", col)).Text())
		lifeInsurance := utils.CleanString(doc.Find(fmt.Sprintf(".attr-MortgageReducingTermAssurance .%s span", col)).Text())

		interestRatePerYear := utils.CleanString(doc.Find(fmt.Sprintf(".attr-InterestRatePerYear .%s span", col)).Text())
		interestRateCondition := utils.CleanString(doc.Find(fmt.Sprintf(".attr-InterestRateCondition .%s span", col)).Text())
		defaultInterestRate := utils.CleanString(doc.Find(fmt.Sprintf(".attr-DefaultInterestRate .%s span", col)).Text())

		creditLimit := utils.CleanString(doc.Find(fmt.Sprintf(".attr-CreditLimit .%s span", col)).Text())
		creditLimitConditionText := utils.CleanString(doc.Find(fmt.Sprintf(".attr-CreditLimitCondition .%s span", col)).Text())
		creditLimitCondition := utils.ParseTextIntoArray(creditLimitConditionText, "-")

		installmentPeriod := utils.CleanString(doc.Find(fmt.Sprintf(".attr-InstallmentPeriod .%s span", col)).Text())
		installmentPeriodCondition := utils.CleanString(doc.Find(fmt.Sprintf(".attr-InstallmentPeriodCondition .%s span", col)).Text())
		installmentPlanDetail := utils.CleanString(doc.Find(fmt.Sprintf(".attr-InstallmentPlanDetail .%s span", col)).Text())

		if installmentPeriodCondition == "" {
			installmentPeriodCondition = ""
		}
		if installmentPlanDetail == "" {
			installmentPlanDetail = ""
		}

		ageText := utils.CleanString(doc.Find(fmt.Sprintf(".attr-BorrowerAge .%s span", col)).Text())
		minAge, maxAge := utils.ParseAgeRange(ageText)

		borrowerConditionsText := utils.CleanString(doc.Find(fmt.Sprintf(".attr-ApplicationCondition .%s span", col)).Text())
		borrowerConditions := utils.ParseTextIntoArray(borrowerConditionsText, "-")

		internalAppraisalFee := utils.CleanString(doc.Find(fmt.Sprintf(".attr-SurveyAndAppraisalFeeByInternal .%s span", col)).Text())
		externalAppraisalFee := utils.CleanString(doc.Find(fmt.Sprintf(".attr-SurveyAndAppraisalFeeByExternal .%s span", col)).Text())
		stampDutyFeeText := utils.CleanString(doc.Find(fmt.Sprintf(".attr-StampDutyFee .%s span", col)).Text())
		stampDutyFee := utils.ParseTextIntoArray(stampDutyFeeText, "-")

		mortgageFee := utils.CleanString(doc.Find(fmt.Sprintf(".attr-MortgageFee .%s span", col)).Text())
		creditCheckFeeText := utils.CleanString(doc.Find(fmt.Sprintf(".attr-CreditBureau .%s span", col)).Text())
		creditCheckFee := utils.ParseTextIntoArray(creditCheckFeeText, "-")

		returnedChequeFee := utils.CleanString(doc.Find(fmt.Sprintf(".attr-ReturnedCheque .%s span", col)).Text())
		insufficientFundsFee := utils.CleanString(doc.Find(fmt.Sprintf(".attr-InsufficientDirectDebitCharge .%s span", col)).Text())
		statementReIssueFee := utils.CleanString(doc.Find(fmt.Sprintf(".attr-StatementReIssuingFee .%s span", col)).Text())
		debtCollectionFeeText := utils.CleanString(doc.Find(fmt.Sprintf(".attr-DebtCollectionFee .%s span", col)).Text())
		debtCollectionFee := utils.ParseTextIntoArray(debtCollectionFeeText, "-")

		otherFeesText := utils.CleanString(doc.Find(fmt.Sprintf(".attr-OtherFee .%s span", col)).Text())
		otherFees := utils.ParseTextIntoArray(otherFeesText, "-")

		// Parse payment fees
		directDebitProvider := utils.CleanString(doc.Find(fmt.Sprintf(".attr-DirectDebitFromAccountFee .%s span", col)).Text())
		directDebitOtherProvider := utils.CleanString(doc.Find(fmt.Sprintf(".attr-DirectDebitFromAccountFeeOther .%s span", col)).Text())
		bankCounterService := utils.CleanString(doc.Find(fmt.Sprintf(".attr-BankCounterServiceFee .%s span", col)).Text())
		bankCounterOtherService := utils.CleanString(doc.Find(fmt.Sprintf(".attr-BankCounterServiceFeeOther .%s span", col)).Text())

		counterServiceFeeText := utils.CleanString(doc.Find(fmt.Sprintf(".attr-CounterServiceFeeOther .%s span", col)).Text())
		counterServiceFee := utils.ParseTextIntoArray(counterServiceFeeText, "-")

		onlinePaymentFee := utils.CleanString(doc.Find(fmt.Sprintf(".attr-paymentOnlineFee .%s span", col)).Text())
		cdmATMPaymentFee := utils.CleanString(doc.Find(fmt.Sprintf(".attr-paymentCDMATMFee .%s span", col)).Text())
		phonePaymentFee := utils.CleanString(doc.Find(fmt.Sprintf(".attr-paymentPhoneFee .%s span", col)).Text())
		chequePaymentFee := utils.CleanString(doc.Find(fmt.Sprintf(".attr-paymentChequeOrMoneyOrderFee .%s span", col)).Text())

		otherPaymentChannelsText := utils.CleanString(doc.Find(fmt.Sprintf(".attr-paymentOtherChannelFee .%s span", col)).Text())
		otherPaymentChannels := utils.ParseTextIntoArray(otherPaymentChannelsText, "-")

		website := utils.CleanString(doc.Find(fmt.Sprintf(".attr-URL .%s a", col)).AttrOr("href", ""))
		feeWebsite := utils.CleanString(doc.Find(fmt.Sprintf(".attr-FeeURL .%s a", col)).AttrOr("href", ""))

		websitePtr := utils.NullIfEmpty(website)
		FeeWebsitePtr := utils.NullIfEmpty(feeWebsite)

		personalLoan := models.PersonalLoan{
			ServiceProvider: serviceProvider,
			Product:         product,
			ProductDetails: models.ProductDetails{
				ProductType:     productType,
				CreditCharacter: creditCharacter,
				Collateral:      collateral,
				CreditLineType:  creditLineType,
				LifeInsurance:   lifeInsurance,
			},
			InterestDetails: models.InterestDetails{
				InterestRatePerYear:    interestRatePerYear,
				InterestRateConditions: interestRateCondition,
				DefaultInterestRate:    defaultInterestRate,
			},
			LoanDetails: models.LoanDetails{
				CreditLimit:                creditLimit,
				CreditLimitCondition:       creditLimitCondition,
				InstallmentPeriod:          installmentPeriod,
				InstallmentPeriodCondition: installmentPeriodCondition,
				InstallmentPlanDetail:      installmentPlanDetail,
			},
			BorrowerDetails: models.BorrowerDetails{
				MinAge:             minAge,
				MaxAge:             maxAge,
				BorrowerConditions: borrowerConditions,
			},
			FeeDetails: models.FeeDetails{
				InternalAppraisalFee: internalAppraisalFee,
				ExternalAppraisalFee: externalAppraisalFee,
				StampDutyFee:         stampDutyFee,
				MortgageFee:          mortgageFee,
				CreditCheckFee:       creditCheckFee,
				ReturnedChequeFee:    returnedChequeFee,
				InsufficientFundsFee: insufficientFundsFee,
				StatementReIssueFee:  statementReIssueFee,
				DebtCollectionFee:    debtCollectionFee,
				OtherFees:            otherFees,
			},
			PaymentFees: models.PaymentFees{
				DirectDebitProvider:      directDebitProvider,
				DirectDebitOtherProvider: directDebitOtherProvider,
				BankCounterService:       bankCounterService,
				BankCounterOtherService:  bankCounterOtherService,
				CounterServiceFee:        counterServiceFee,
				OnlinePaymentFee:         onlinePaymentFee,
				CDMATMPaymentFee:         cdmATMPaymentFee,
				PhonePaymentFee:          phonePaymentFee,
				ChequePaymentFee:         chequePaymentFee,
				OtherPaymentChannels:     otherPaymentChannels,
			},
			AdditionalInfo: models.AdditionalInfo{
				Website:    websitePtr,
				FeeWebsite: FeeWebsitePtr,
			},
		}
		loans = append(loans, personalLoan)

	}
	return loans, nil

}
