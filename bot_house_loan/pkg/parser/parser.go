package parser

import (
	"bytes"
	"errors"
	"fmt"
	"house_loan/pkg/models"
	"house_loan/pkg/utils"
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

func ParseHouseLoanDetails(body []byte) ([]models.HouseLoan, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var loans []models.HouseLoan
	for i := 1; i <= 3; i++ {
		col := fmt.Sprintf("col%d", i)
		averageInterestRateStr := utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-percentAverageInterestRateForThreeYearsDisplay td.%s span.text-bold", col)).Text())
		averageInterestRate, _ := strconv.ParseFloat(strings.ReplaceAll(averageInterestRateStr, "%", ""), 64)

		interestRateConditions := utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-conditionOfYearInterestRateDisplay td.%s span", col)).Text())
		interestRateConditionList := utils.ParseInterestRateConditions(interestRateConditions)

		maximumNormalInterestRateStr := utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-maximumNormalInterestRateDisplay td.%s span", col)).Text())
		maximumNormalInterestRate, _ := strconv.ParseFloat(strings.ReplaceAll(maximumNormalInterestRateStr, "%", ""), 64)

		borrowerAgeStr := utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-borrowerAge td.%s span", col)).Text())
		borrowerAge, _ := utils.ExtractNumericValue(strings.Join(strings.Fields(borrowerAgeStr), "")) // Extract numeric value

		combinedLoanConditions := utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-conditionOfProduct td.%s span", col)).Text())
		if combinedLoanConditions == "-" {
			combinedLoanConditions = ""
		}

		productSpecificConditions := utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-conditionOfProduct td.%s span", col)).Text())
		if productSpecificConditions == "-" {
			productSpecificConditions = ""
		}

		// Extract and split LTV ratio
		ltvRatioStr := utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-loanToValueRatio td.%s span", col)).Text())
		ltvRatioList := strings.Split(ltvRatioStr, "/")

		// Extract survey and appraisal fee and split it by "-"
		surveyAppraisalFeeStr := utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-surveyAndAppraisalFeeDisplay td.%s span", col)).Text())
		surveyAppraisalFeeList := strings.Split(surveyAppraisalFeeStr, "-")

		deductingFromOtherBankAC := utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-deductingFromOtherBankACFee td.%s span", col)).Text())
		deductingFromOtherBankACList := strings.Split(deductingFromOtherBankAC, "-")

		otherCounterService := utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-otherCounterServiceFee td.%s span", col)).Text())
		otherCounterServiceList := strings.Split(otherCounterService, "-")

		onlinePayment := utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-onlinePaymentFee td.%s span", col)).Text())
		onlinePaymentList := strings.Split(onlinePayment, "-")

		loan := models.HouseLoan{
			ServiceProvider: utils.CleanString(doc.Find("th.col-s-" + fmt.Sprint(i) + " span").Eq(1).Text()),
			Product:         utils.CleanString(doc.Find("th.font-black.text-center." + col + " span").Last().Text()),
			InterestRate: models.InterestRate{
				AverageInterestRateThreeYears: averageInterestRate,
				InterestRateConditions:        interestRateConditionList,
				EffectiveInterestRate:         utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-effectiveInterestRateDisplay td.%s span", col)).Text()),
				MaximumNormalInterestRate:     maximumNormalInterestRate,
				DefaultInterestRate:           utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-defaultInterestRateAndRelatedCondition td.%s span", col)).Text()),
			},
			ProductDetails: models.ProductDetails{
				LoanType:                  utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-loanTypeName td.%s span", col)).Text()),
				CollateralType:            utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-typeOfCollateraDisplay td.%s span", col)).Text()),
				BorrowerQualifications:    utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-characterOfBorrowerToLoanInterestRate td.%s span", col)).Text()),
				LoanConditions:            utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-conditionOfLoanWithOtherProducts td.%s span", col)).Text()),
				CombinedLoanConditions:    utils.ParseOptionalCondition(combinedLoanConditions),
				ProductSpecificConditions: utils.ParseOptionalCondition(productSpecificConditions),
				BorrowerAge:               borrowerAge,
				MinimumIncome:             utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-minimumMonthlyIncomeDisplay td.%s span", col)).Text()),
				ApplicationConditions:     utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-conditionToApply td.%s span", col)).Text()),
			},
			LoanCreditRepayment: models.LoanCreditRepayment{
				CreditLimitRange:      utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-creditLimitDisplay td.%s span", col)).Text()),
				LTVRatio:              ltvRatioList,
				CreditLimitConditions: utils.ParseOptionalCondition(utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-conditionOfCreditLimit td.%s span", col)).Text())),
				LoanTerm:              utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-periodOfBorrowing td.%s span", col)).Text()),
				RepaymentConditions:   utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-conditionOfInstallment td.%s span", col)).Text()),
			},
			InsuranceDetails: models.InsuranceDetails{
				MRTAConditions:      utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-mortgageReducingTermAssuranceDisplay td.%s span", col)).Text()),
				MRTACancellationFee: utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-mRTACancellationFee td.%s span", col)).Text()),
			},
			GeneralFees: models.GeneralFees{
				SurveyAndAppraisalFee:     surveyAppraisalFeeList,
				StampDuty:                 utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-stampDuty td.%s span", col)).Text()),
				MortgageFee:               utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-mortgageFee td.%s span", col)).Text()),
				TransferFee:               utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-transferFee td.%s span", col)).Text()),
				CreditInfoVerificationFee: utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-creditInformationVerificationFee td.%s span", col)).Text()),
				FireInsurancePremium:      utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-fireInsurancePremiums td.%s span", col)).Text()),
				ChequeReturnFee:           utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-feeForChequeReturned td.%s span", col)).Text()),
				DeficiencyBalanceFee:      utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-feeForDeficiencyBalanceAC td.%s span", col)).Text()),
				StatementCopyFee:          utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-copyOfStatementFee td.%s span", col)).Text()),
				ChequeReturnFine:          strings.Split(utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-finesForChequeReturned td.%s span", col)).Text()), "-"),
				DebtCollectionFee:         strings.Split(utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-debtCollectionFee td.%s span", col)).Text()), "-"),
				InterestRateChangeFee:     utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-feeforChangingInterestRate td.%s span", col)).Text()),
				RefinanceFee:              strings.Split(utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-refinanceFee td.%s span", col)).Text()), "-"),
				OtherFees:                 utils.ParseOptionalCondition(utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-otherFees td.%s span", col)).Text())),
			},
			PaymentFees: models.PaymentFees{
				DeductingFromBankAccount:  utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-deductingFromBankACFee td.%s span", col)).Text()),
				DeductingFromOtherBankAC:  deductingFromOtherBankACList,
				BankCounterService:        utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-bankCounterServiceFee td.%s span", col)).Text()),
				OtherBankCounterService:   utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-otherBankCounterServiceFee td.%s span", col)).Text()),
				OtherCounterService:       otherCounterServiceList,
				OnlinePayment:             onlinePaymentList,
				CDMATMPayment:             utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-cDMATMPaymentFee td.%s span", col)).Text()),
				PhonePayment:              utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-phonePaymentFee td.%s span", col)).Text()),
				ChequeOrMoneyOrderPayment: utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-chequeOrMoneyOrderPaymentFee td.%s span", col)).Text()),
				OtherChannelPayment:       utils.CleanString(doc.Find(fmt.Sprintf("tr.attr-otherChannelPaymentFee td.%s span", col)).Text()),
			},
			ProductWebsite: doc.Find(fmt.Sprintf("tr.attr-uRL td.%s a.prod-url", col)).AttrOr("href", ""),
			FeeWebsite:     doc.Find(fmt.Sprintf("tr.attr-feeuRL td.%s a.prod-url", col)).AttrOr("href", ""),
		}
		loans = append(loans, loan)
	}

	if len(loans) == 0 {
		return nil, errors.New("no loan details found")
	}
	return loans, nil
}
