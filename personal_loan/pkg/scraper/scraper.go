package scraper

import (
	"bytes"
	"fmt"
	"personal_loan/helpers"
	"personal_loan/pkg/model"
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

func ParseHTML(body []byte) ([]model.LoanProduct, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var loanProducts []model.LoanProduct

	for i := 1; i <= 3; i++ {
		col := fmt.Sprintf("col%d", i)
		serviceProvider := helpers.CleanString(doc.Find("th." + col + " .compare-header span").Eq(1).Text())

		product := helpers.CleanString(doc.Find("th.prod-" + col + " span").First().Text())

		salaryInterest := helpers.GetNullableText(doc.Find("tr.attr-interestForSalaryIncomeDisplay td." + col + " span").First().Text())
		businessInterest := helpers.GetNullableText(doc.Find("tr.attr-interestForSalaryIncomeDisplay td." + col + " span").Last().Text())

		interestConditions := helpers.GetNullableText(doc.Find("tr.attr-interestCondition td." + col + " span").First().Text())
		interestPromotionsRaw := doc.Find(fmt.Sprintf("tr.attr-promotionOrCampaignOfInterest td.%s span", col)).First().Text()
		interestPromotions := helpers.SplitAndFormatArray(interestPromotionsRaw)

		minimumPayment := helpers.GetNullableText(doc.Find(fmt.Sprintf("tr.attr-minimumPayment td.%s span", col)).First().Text())

		creditLimitRaw := helpers.GetNullableText(doc.Find(fmt.Sprintf("tr.attr-maximumTimesOfIncomeForCreditLines td.%s span", col)).First().Text())
		var creditLimit []string
		if creditLimitRaw != nil {
			creditLimit = helpers.FormatCreditLimit(*creditLimitRaw)
		}

		loanAmount := helpers.ParseLoanAmount(doc.Find(fmt.Sprintf("tr.attr-creditLineAmoutDisplay td.%s span", col)).First().Text())

		loanDuration := helpers.ParseLoanDuration(doc.Find(fmt.Sprintf("tr.attr-termDisplay td.%s span", col)).First().Text())

		moneyTransferMethod := helpers.GetNullableText(doc.Find(fmt.Sprintf("tr.attr-channelCreditDelivery td.%s span", col)).First().Text())

		moneyTransferConditionsRaw := helpers.GetNullableText(doc.Find(fmt.Sprintf("tr.attr-conditionOfCreditDelivery td.%s span", col)).First().Text())
		var moneyTransferConditions []string
		if moneyTransferConditionsRaw != nil {
			moneyTransferConditions = helpers.SplitMoneyTransferConditions(*moneyTransferConditionsRaw)
		}

		// Parsing applicant requirements for salary employees
		ageText := doc.Find(fmt.Sprintf("tr.attr-header.attr-ageForSalaryIncome td.%s span", col)).First().Text()
		salaryAge := helpers.ParseInt(ageText)

		minIncomeText := doc.Find(fmt.Sprintf("tr:contains('รายได้ขั้นต่ำ') td.%s span", col)).First().Text()
		salaryMinIncome := helpers.ParseInt(minIncomeText)

		workExperienceText := doc.Find(fmt.Sprintf("tr.attr-header.attr-ageForSalaryIncome td.cmpr-col.%s span", col)).Last().Text()
		workExperience := helpers.CleanString(workExperienceText)

		businessAgeText := doc.Find(fmt.Sprintf("tr.attr-header.attr-ageForSelfEmployed td.%s span", col)).First().Text()
		businessOwnerAge := helpers.ParseInt(businessAgeText)

		businessMinIncomeText := doc.Find(fmt.Sprintf("tr:contains('รายได้ขั้นต่ำ') td.%s span", col)).Last().Text()
		businessOwnerMinIncome := helpers.ParseInt(businessMinIncomeText)

		businessDurationText := doc.Find(fmt.Sprintf("tr:contains('อายุงานขั้นต่ำ') td.%s span", col)).Last().Text()
		businessOwnerDuration := helpers.CleanString(businessDurationText)

		applicantReq := model.ApplicantRequirements{
			SalaryEmployee: model.SalaryEmployeeRequirements{
				Age:            salaryAge,
				MinimumIncome:  salaryMinIncome,
				WorkExperience: workExperience,
			},
			BusinessOwnerRequirements: model.BusinessOwnerRequirements{
				Age:              businessOwnerAge,
				MinimumIncome:    businessOwnerMinIncome,
				BusinessDuration: businessOwnerDuration,
			},
		}

		// Penalty rates extraction
		var penaltyRates []model.PenaltyRate
		doc.Find("tr.attr-penaltyFee td." + col).Each(func(index int, item *goquery.Selection) {
			rate := helpers.GetNullableText(item.Find("span").First().Text())
			conditionRaw := item.Find("div span.text-primary").Text()
			conditionArray := helpers.SplitAndFormatArray(helpers.CleanString(conditionRaw))
			if rate != nil && len(conditionArray) > 0 {
				penaltyRates = append(penaltyRates, model.PenaltyRate{
					Rate:      rate,
					Condition: conditionArray,
				})
			}
		})

		applicationConditionsRaw := doc.Find(fmt.Sprintf("tr.attr-header.attr-conditionToApply td.%s span", col)).First().Text()
		applicationConditions := helpers.SplitAndFormatArray(applicationConditionsRaw)

		installmentService := helpers.ParseOptionalField(doc, fmt.Sprintf("tr.attr-header.attr-cardUsedCondition td.%s span", col))
		benefits := helpers.ParseOptionalField(doc, fmt.Sprintf("tr.attr-header.attr-benefit td.%s span", col))
		cardFee := helpers.ParseOptionalField(doc, fmt.Sprintf("tr.attr-header.attr-cardAnnuallyFee td.%s span", col))
		cardReplacementFee := helpers.ParseOptionalField(doc, fmt.Sprintf("tr.attr-header.attr-cardReplacementFee td.%s span", col))
		pinReissuingFee := helpers.ParseOptionalField(doc, fmt.Sprintf("tr.attr-header.attr-pINReissuingFee td.%s span", col))
		fxRiskConversionFee := helpers.ParseOptionalField(doc, fmt.Sprintf("tr.attr-header.attr-costOfFXRisk td.%s span", col))

		revolvingCreditFeeIngo := model.RevolvingCreditFeeInfo{
			InstallmentService:  installmentService,
			Benefits:            benefits,
			CardFee:             cardFee,
			CardReplacementFee:  cardReplacementFee,
			PinReissuingFee:     pinReissuingFee,
			FxRiskConversionFee: fxRiskConversionFee,
		}

		// Parsing ServiceFee details
		creditCheckFee := helpers.CleanString(doc.Find(fmt.Sprintf("tr.attr-header.attr-creditBureauFee td.%s span", col)).First().Text())
		stampDuty := helpers.CleanString(doc.Find(fmt.Sprintf("tr.attr-header.attr-dutyStampFee td.%s span", col)).First().Text())
		earlyRepaymentFee := helpers.CleanString(doc.Find(fmt.Sprintf("tr.attr-header.attr-prepaymentFee td.%s span", col)).First().Text())
		chequeReturnedFee := helpers.CleanString(doc.Find(fmt.Sprintf("tr.attr-header.attr-chequeReturnedFee td.%s span", col)).First().Text())
		insufficientDirectDebitCharge := helpers.CleanString(doc.Find(fmt.Sprintf("tr.attr-header.attr-insufficientDirectDebitCharge td.%s span", col)).First().Text())
		statementReissuingFee := helpers.CleanString(doc.Find(fmt.Sprintf("tr.attr-header.attr-statementReissuingFee td.%s span", col)).First().Text())
		transactionVerificationFee := helpers.CleanString(doc.Find(fmt.Sprintf("tr.attr-header.attr-transactionVerificationFee td.%s span", col)).First().Text())
		collectionFee := helpers.CleanString(doc.Find(fmt.Sprintf("tr.attr-header.attr-collectionFee td.%s span", col)).First().Text())
		otherFees := helpers.CleanString(doc.Find(fmt.Sprintf("tr.attr-header.attr-otherFee td.%s span", col)).First().Text())

		serviceFee := model.ServiceFee{
			CreditCheckFee:                helpers.GetNullableText(creditCheckFee),
			StampDuty:                     helpers.GetNullableText(stampDuty),
			EarlyRepaymentFee:             helpers.GetNullableText(earlyRepaymentFee),
			ChequeReturnedFee:             helpers.GetNullableText(chequeReturnedFee),
			InsufficientDirectDebitCharge: helpers.GetNullableText(insufficientDirectDebitCharge),
			StatementReissuingFee:         helpers.GetNullableText(statementReissuingFee),
			TransactionVerificationFee:    helpers.GetNullableText(transactionVerificationFee),
			CollectionFee:                 helpers.GetNullableText(collectionFee),
			OtherFees:                     helpers.GetNullableText(otherFees),
		}

		noFeeRaw := doc.Find(fmt.Sprintf("tr.attr-header.attr-freePaymentChannel td.%s span", col)).First().Text()
		noFee := helpers.SplitAndFormatArray(noFeeRaw)

		branchService := helpers.ParseOptionalField(doc, fmt.Sprintf("tr.attr-header.attr-deductingFromBankAcFee td.%s span", col))
		deductingFromOtherBank := helpers.ParseOptionalField(doc, fmt.Sprintf("tr.attr-header.attr-deductingFromOtherBankAcFee td.%s span", col))
		providerBranchService := helpers.ParseOptionalField(doc, fmt.Sprintf("tr.attr-header.attr-bankCounterServiceFee td.%s span", col))

		otherProviderBranchRaw := doc.Find(fmt.Sprintf("tr.attr-header.attr-otherBankCounterServiceFee td.%s span", col)).First().Text()
		otherProviderBranch := helpers.SplitAndFormatArray(otherProviderBranchRaw)

		counterServiceRaw := doc.Find(fmt.Sprintf("tr.attr-header.attr-otherCounterServiceFee td.%s span", col)).First().Text()
		counterService := helpers.SplitAndFormatArray(counterServiceRaw)

		onlinePaymentRaw := doc.Find(fmt.Sprintf("tr.attr-header.attr-onlinePaymentFee td.%s span", col)).First().Text()
		onlinePayment := helpers.SplitAndFormatArray(onlinePaymentRaw)

		cDMATMPaymentRaw := doc.Find(fmt.Sprintf("tr.attr-header.attr-cDMATMPaymentFee td.%s span", col)).First().Text()
		cDMATMPayment := helpers.SplitAndFormatArray(cDMATMPaymentRaw)

		phonePayment := helpers.ParseOptionalField(doc, fmt.Sprintf("tr.attr-header.attr-phonePaymentFee td.%s span", col))
		chequeMoneyOrderPayment := helpers.ParseOptionalField(doc, fmt.Sprintf("tr.attr-header.attr-chequeOrMoneyOrderPaymentFee td.%s span", col))
		otherChannelPayment := helpers.ParseOptionalField(doc, fmt.Sprintf("tr.attr-header.attr-otherChannelPaymentFee td.%s span", col))

		productWebsite := helpers.ParseLinkField(doc, fmt.Sprintf("tr.attr-header.attr-uRL td.%s", col))
		feeWebsite := helpers.ParseLinkField(doc, fmt.Sprintf("tr.attr-header.attr-feeuRL td.%s", col))

		paymentMethods := model.PaymentMethods{
			NoFee:                   noFee,
			BranchService:           branchService,
			DeductingFromOtherBank:  deductingFromOtherBank,
			ProviderBranchService:   providerBranchService,
			OtherProviderBranch:     otherProviderBranch,
			CounterService:          counterService,
			OnlinePayment:           onlinePayment,
			CDMATMPayment:           cDMATMPayment,
			PhonePayment:            phonePayment,
			ChequeMoneyOrderPayment: chequeMoneyOrderPayment,
			OtherChannelPayment:     otherChannelPayment,
		}

		additionalInfo := model.AdditionalInfo{
			ProductWebsite: productWebsite,
			FeeWebsite:     feeWebsite,
		}

		paymentInfo := model.PaymentFeeInfo{
			PaymentMethods: paymentMethods,
			AdditionalInfo: additionalInfo,
		}

		loanProducts = append(loanProducts, model.LoanProduct{
			ServiceProvider: serviceProvider,
			Product:         product,
			InterestRate: model.InterestRate{
				SalaryEmployee: salaryInterest,
				BusinessOwner:  businessInterest,
			},
			InterestConditions:      interestConditions,
			InterestPromotions:      interestPromotions,
			PenaltyRates:            penaltyRates,
			MinimumPayment:          minimumPayment,
			CreditLimit:             creditLimit,
			LoanAmount:              loanAmount,
			LoanDurationMonths:      loanDuration,
			MoneyTransferMethod:     moneyTransferMethod,
			MoneyTransferConditions: moneyTransferConditions,
			ApplicantRequirements:   applicantReq,
			ApplicationConditions:   applicationConditions,
			RevolvingCreditFeeInfo:  revolvingCreditFeeIngo,
			ServiceFee:              serviceFee,
			PaymentFeeInfo:          paymentInfo,
		})
	}

	return loanProducts, nil
}
