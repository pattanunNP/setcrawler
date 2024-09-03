package loanproducts

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CleanText(s string) string {
	return strings.TrimSpace(s)
}

// SplitText splits a string by a delimiter and trims spaces
func SplitText(s, delimiter string) []string {
	parts := strings.Split(s, delimiter)
	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// ParseConditions splits conditions by "-" and trims spaces
func ParseConditions(raw string) []string {
	conditions := strings.Split(raw, "-")
	var result []string
	for _, condition := range conditions {
		trimmed := strings.TrimSpace(condition)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// WriteJSONToFile writes the given data to a JSON file
func WriteJSONToFile(data interface{}, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, jsonData, 0644)
}

// Helper function to extract min and max values from a string
func extractMinMax(text string) (int, int) {
	// Remove non-numeric characters except for spaces, hyphens, and commas
	re := regexp.MustCompile(`[^\d\s\-]`)
	cleanedText := re.ReplaceAllString(text, "")

	// Split the cleaned string into parts based on spaces, hyphens, and commas
	parts := regexp.MustCompile(`[,\s\-]+`).Split(cleanedText, -1)

	// Parse the numbers into min and max
	var numbers []int
	for _, part := range parts {
		if part != "" {
			if number, err := strconv.Atoi(part); err == nil {
				numbers = append(numbers, number)
			}
		}
	}

	// Determine min and max from the extracted numbers
	var min, max int
	if len(numbers) >= 2 {
		min = numbers[0]
		max = numbers[1]
	} else if len(numbers) == 1 {
		min = numbers[0]
		max = min
	} else {
		min = 0
		max = 0
	}
	return min, max
}

// Helper function to get Borrower Type by column index
func getBorrowerType(doc *goquery.Document, index int) []string {
	selector := fmt.Sprintf("tr.attr-BorrowerType td.cmpr-col.col%d span", index)
	text := doc.Find(selector).Text()
	return ParseConditions(text)
}

// Helper function to get Age Limit by column index
func getAgeLimit(doc *goquery.Document, index int) AgeLimit {
	selector := fmt.Sprintf("tr.attr-BorrowerAge td.cmpr-col.col%d span", index)
	text := doc.Find(selector).Text()
	min, max := extractMinMax(text)
	return AgeLimit{
		MinAge: min,
		MaxAge: max,
	}
}

// Helper function to get Other Borrower Conditions by column index
func getOtherBorrowerConditions(doc *goquery.Document, index int) *string {
	selector := fmt.Sprintf("tr.attr-ConditionOfBorrower td.cmpr-col.col%d span", index)
	text := CleanText(doc.Find(selector).Text())
	if text == "" {
		return nil
	}
	return &text
}

// Helper function to get Minimum Income by column index
func getMinimumIncome(doc *goquery.Document, index int) []string {
	selector := fmt.Sprintf("tr.attr-MinimumMonthlyIncome td.cmpr-col.col%d span", index)
	text := doc.Find(selector).Text()
	return ParseConditions(text)
}

// Helper function to get Income Conditions by column index
func getIncomeConditions(doc *goquery.Document, index int) *string {
	selector := fmt.Sprintf("tr.attr-MonthlyIncomeCondition td.cmpr-col.col%d span", index)
	text := CleanText(doc.Find(selector).Text())
	if text == "" {
		return nil
	}
	return &text
}

// Helper function to get VehicleType by column index
func getVehicleType(doc *goquery.Document, index int) []string {
	selector := fmt.Sprintf("tr.attr-TypeOfCollateral td.cmpr-col.col%d span", index)
	text := doc.Find(selector).Text()
	return SplitText(text, "/")
}

// Helper function to get VehicleCondition by column index
func getVehicleCondition(doc *goquery.Document, index int) []string {
	selector := fmt.Sprintf("tr.attr-TypeOfCollateralCondition td.cmpr-col.col%d span", index)
	text := doc.Find(selector).Text()
	return ParseConditions(text)
}

// Helper function to get LoanType by column index
func getLoanType(doc *goquery.Document, index int) string {
	selector := fmt.Sprintf("tr.attr-CreditLineType td.cmpr-col.col%d span", index)
	return CleanText(doc.Find(selector).Text())
}

func getInterestRate(doc *goquery.Document, index int) InterestRate {
	annualSelector := fmt.Sprintf("tr.attr-InterestWithServiceFee td.cmpr-col.col%d span", index)
	conditionSelector := fmt.Sprintf("tr.attr-InterestWithServiceFeeCondition td.cmpr-col.col%d span", index)
	penaltySelector := fmt.Sprintf("tr.attr-DefaultInterestRate td.cmpr-col.col%d span", index)

	annualInterest := CleanText(doc.Find(annualSelector).Text())
	conditions := ParseConditions(doc.Find(conditionSelector).Text())
	penaltyInterest := CleanText(doc.Find(penaltySelector).Text())

	return InterestRate{
		AnnualInterestRate:     annualInterest,
		InterestRateConditions: conditions,
		PenaltyInterestRate:    penaltyInterest,
	}
}

func getCreditLimitAndInstallment(doc *goquery.Document, index int) CreditLimitAndInstallment {
	creditLimit := getCreditLimit(doc, index)
	creditLimitConditions := getCreditLimitConditions(doc, index)
	installmentPeriod := getInstallmentPeriod(doc, index)
	loanReceivingChannel := getLoanReceivingChannel(doc, index)

	return CreditLimitAndInstallment{
		CreditLimit:           creditLimit,
		CreditLimitConditions: creditLimitConditions,
		InstallmentPeriod:     installmentPeriod,
		LoanReceivingChannel:  loanReceivingChannel,
	}
}

// Extract and parse credit limit
func getCreditLimit(doc *goquery.Document, index int) CreditLimit {
	selector := fmt.Sprintf("tr.attr-CreditLimit td.cmpr-col.col%d span", index)
	text := doc.Find(selector).Text()
	min, max := extractMinMax(text)
	return CreditLimit{
		MinLimit: min,
		MaxLimit: max,
	}
}

// Extract and parse credit limit conditions
func getCreditLimitConditions(doc *goquery.Document, index int) []string {
	selector := fmt.Sprintf("tr.attr-CreditLimitCondition td.cmpr-col.col%d span", index)
	text := doc.Find(selector).Text()
	return ParseConditions(text)
}

// Extract and parse installment period
func getInstallmentPeriod(doc *goquery.Document, index int) InstallmentPeriod {
	selector := fmt.Sprintf("tr.attr-InstallmentPeriod td.cmpr-col.col%d span", index)
	text := doc.Find(selector).Text()
	min, max := extractMinMax(text)
	return InstallmentPeriod{
		MinMonth: min,
		MaxMonth: max,
	}
}

// Extract and parse loan receiving channel
func getLoanReceivingChannel(doc *goquery.Document, index int) []string {
	selector := fmt.Sprintf("tr.attr-LoanReceivingChannel td.cmpr-col.col%d span", index)
	text := doc.Find(selector).Text()
	return ParseConditions(text)
}

func getBorrowerQualifications(doc *goquery.Document, index int) BorrowerQualifications {
	borrowerType := getBorrowerType(doc, index)
	ageLimit := getAgeLimit(doc, index)
	otherConditions := getOtherBorrowerConditions(doc, index)
	minimumIncome := getMinimumIncome(doc, index)
	incomeConditions := getIncomeConditions(doc, index)

	return BorrowerQualifications{
		BorrowerType:     borrowerType,
		AgeLimit:         ageLimit,
		OtherConditions:  otherConditions,
		MinimumIncome:    minimumIncome,
		IncomeConditions: incomeConditions,
	}
}

func getGeneralFees(doc *goquery.Document, index int) GeneralFees {
	stampDuty := getStampDutyFee(doc, index)
	returnedCheque := getReturnedChequeFee(doc, index)
	debtCollectionFee := getDebtCollectionFee(doc, index)

	return GeneralFees{
		StampDuty:         stampDuty,
		ReturnCheque:      returnedCheque,
		DebtCollectionFee: debtCollectionFee,
	}
}

// Extract Stamp Duty Fee by column index
func getStampDutyFee(doc *goquery.Document, index int) []string {
	selector := fmt.Sprintf("tr.attr-StampDutyFee td.cmpr-col.col%d span", index)
	htmlContent, err := doc.Find(selector).Html()
	if err != nil {
		fmt.Printf("Error extracting Stamp Duty Fee content: %v\n", err)
		return []string{}
	}
	return ParseConditions(htmlContent)
}

// Extract Returned Cheque Fee by column index
func getReturnedChequeFee(doc *goquery.Document, index int) string {
	selector := fmt.Sprintf("tr.attr-ReturnedCheque td.cmpr-col.col%d span", index)
	text := CleanText(doc.Find(selector).Text())
	return text
}

// Extract Debt Collection Fee by column index
func getDebtCollectionFee(doc *goquery.Document, index int) []string {
	selector := fmt.Sprintf("tr.attr-DebtCollectionFee td.cmpr-col.col%d span", index)
	htmlContent, err := doc.Find(selector).Html()
	if err != nil {
		fmt.Printf("Error extracting Debt Collection Fee content: %v\n", err)
		return []string{}
	}
	return ParseConditions(htmlContent)
}

func getCardFees(doc *goquery.Document, index int) CardFees {
	cardFee := getCardFee(doc, index)
	cardReplcementFee := getCardReplacementFee(doc, index)
	creditWithdrawalFee := getCreditWithdrawalFee(doc, index)

	return CardFees{
		CardFees:            cardFee,
		CardReplacementFee:  cardReplcementFee,
		CreditWithdrawalFee: creditWithdrawalFee,
	}
}

// Extract Card Fee by column index
func getCardFee(doc *goquery.Document, index int) string {
	selector := fmt.Sprintf("tr.attr-CardFee td.cmpr-col.col%d span", index)
	text := CleanText(doc.Find(selector).Text())
	return text
}

// Extract Card Replacement Fee by column index
func getCardReplacementFee(doc *goquery.Document, index int) string {
	selector := fmt.Sprintf("tr.attr-CardReplacementFee td.cmpr-col.col%d span", index)
	text := CleanText(doc.Find(selector).Text())
	return text
}

// Extract Credit Withdrawal Fee by column index
func getCreditWithdrawalFee(doc *goquery.Document, index int) *string {
	selector := fmt.Sprintf("tr.attr-CreditWithdrawalFee td.cmpr-col.col%d span", index)
	text := CleanText(doc.Find(selector).Text())
	if text == "" {
		return nil
	}
	return &text
}

func getPaymentFees(doc *goquery.Document, index int) PaymentFees {
	freePaymentChannel := getFreePaymentChannel(doc, index)
	deductingFromServiceProvider := getDeductingFromServiceProvider(doc, index)
	deductingFromOtherServiceProvider := getDeductingFromOtherServiceProvider(doc, index)
	serviceProviderCounter := getServiceProviderCounter(doc, index)
	otherProviderCounter := getOtherProviderCounter(doc, index)
	paymentServicePoints := cleanPaymentServicePoints(getTextOrEmpty(doc, fmt.Sprintf("tr.attr-OthersPaymentCounter td.cmpr-col.col%d span", index)))
	onlinePayment := getOnlinePayment(doc, index)
	cdmAtmPayment := getCDMATMPayment(doc, index)
	phonePayment := getPhonePayment(doc, index)
	chequeMoneyOrderPayment := getChequeMoneyOrderPayment(doc, index)
	otherChannelPayment := getOtherChannelPayment(doc, index)

	return PaymentFees{
		FreePaymentChannel:                freePaymentChannel,
		DeductingFromServiceProvider:      deductingFromServiceProvider,
		DeductingFromOtherServiceProvider: deductingFromOtherServiceProvider,
		ServiceProviderCounter:            serviceProviderCounter,
		OtherProviderCounter:              otherProviderCounter,
		PaymentServicePoints:              paymentServicePoints,
		OnlinePayment:                     onlinePayment,
		CDMATMPayment:                     cdmAtmPayment,
		PhonePayment:                      phonePayment,
		ChequeMoneyOrderPayment:           chequeMoneyOrderPayment,
		OtherChannelPayment:               otherChannelPayment,
	}
}

func cleanPaymentServicePoints(raw string) []string {
	raw = strings.ReplaceAll(raw, "\u003cbr/\u003e", "")
	raw = strings.ReplaceAll(raw, "\n", "")
	// Split by "-"
	points := strings.Split(raw, "-")
	var cleaned []string
	for _, point := range points {
		trimmed := strings.TrimSpace(point)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return cleaned
}

func getFreePaymentChannel(doc *goquery.Document, index int) []string {
	selector := fmt.Sprintf("tr.attr-FreePaymentChannel td.cmpr-col.col%d span", index)
	text := CleanText(doc.Find(selector).Text())
	if text == "" {
		return []string{}
	}
	return strings.Split(text, ",")
}

func getDeductingFromServiceProvider(doc *goquery.Document, index int) string {
	selector := fmt.Sprintf("tr.attr-DeductingFromBankACFee td.cmpr-col.col%d span", index)
	return CleanText(doc.Find(selector).Text())
}

func getDeductingFromOtherServiceProvider(doc *goquery.Document, index int) string {
	selector := fmt.Sprintf("tr.attr-DeductingFromOtherBankACFee td.cmpr-col.col%d span", index)
	return CleanText(doc.Find(selector).Text())
}

func getServiceProviderCounter(doc *goquery.Document, index int) string {
	selector := fmt.Sprintf("tr.attr-ServiceProviderCounter td.cmpr-col.col%d span", index)
	return CleanText(doc.Find(selector).Text())
}

func getOtherProviderCounter(doc *goquery.Document, index int) string {
	selector := fmt.Sprintf("tr.attr-OtherProviderCounter td.cmpr-col.col%d span", index)
	return CleanText(doc.Find(selector).Text())
}

func getOnlinePayment(doc *goquery.Document, index int) string {
	selector := fmt.Sprintf("tr.attr-OnlinePaymentFee td.cmpr-col.col%d span", index)
	return CleanText(doc.Find(selector).Text())
}

func getCDMATMPayment(doc *goquery.Document, index int) string {
	selector := fmt.Sprintf("tr.attr-CDMATMPaymentFee td.cmpr-col.col%d span", index)
	return CleanText(doc.Find(selector).Text())
}

func getPhonePayment(doc *goquery.Document, index int) string {
	selector := fmt.Sprintf("tr.attr-PhonePaymentFee td.cmpr-col.col%d span", index)
	return CleanText(doc.Find(selector).Text())
}

func getChequeMoneyOrderPayment(doc *goquery.Document, index int) string {
	selector := fmt.Sprintf("tr.attr-ChequeMoneyOrderPaymentFee td.cmpr-col.col%d span", index)
	return CleanText(doc.Find(selector).Text())
}

func getOtherChannelPayment(doc *goquery.Document, index int) *string {
	selector := fmt.Sprintf("tr.attr-OtherChannelPaymentFee td.cmpr-col.col%d span", index)
	text := CleanText(doc.Find(selector).Text())
	if text == "" {
		return nil
	}
	return &text
}

func getOtherFees(doc *goquery.Document, index int) OtherFees {
	litigationLawyerFee := getLitigationLawyerFee(doc, index)
	otherFeesDetails := getOtherFeesDetails(doc, index)

	return OtherFees{
		LitigationLawyerFee: litigationLawyerFee,
		OtherFeesDetails:    otherFeesDetails,
	}
}

func getLitigationLawyerFee(doc *goquery.Document, index int) []string {
	selector := fmt.Sprintf("tr.attr-lawyerFeeInCaseOfLitigation td.cmpr-col.col%d span", index)
	htmlContent, err := doc.Find(selector).Html()
	if err != nil {
		fmt.Printf("Error extracting Litigation Lawyer Fee content: %v\n", err)
		return []string{}
	}

	// Split by "เงื่อนไข:" and further split by "-"
	conditions := strings.Split(htmlContent, "เงื่อนไข:")
	if len(conditions) > 1 {
		return ParseConditions(conditions[1])
	}
	return []string{CleanText(htmlContent)}
}

func getOtherFeesDetails(doc *goquery.Document, index int) []string {
	selector := fmt.Sprintf("tr.attr-OtherFee td.cmpr-col.col%d span", index)
	text := CleanText(doc.Find(selector).Text())
	if text == "" {
		return []string{}
	}

	// Adjusted regex to match "number. description" without using lookahead
	re := regexp.MustCompile(`\d+\..*?`)
	matches := re.FindAllString(text, -1)

	// Trim spaces and add to the result
	var result []string
	for _, match := range matches {
		result = append(result, strings.TrimSpace(match))
	}
	return result
}

func getAdditionInfo(doc *goquery.Document, index int) *AdditionalInfo {
	productWebsite := getProductWebsite(doc, index)
	feeWebsite := getFeeWebsite(doc, index)

	if productWebsite == nil && feeWebsite == nil {
		return nil
	}

	return &AdditionalInfo{
		ProductWebsite: productWebsite,
		FeeWebsite:     feeWebsite,
	}
}

func getProductWebsite(doc *goquery.Document, index int) *string {
	selector := fmt.Sprintf("tr.attr-URL td.cmpr-col.col%d a", index)
	link, exists := doc.Find(selector).Attr("href")
	if !exists || CleanText(link) == "" {
		return nil
	}
	return &link
}

func getFeeWebsite(doc *goquery.Document, index int) *string {
	selector := fmt.Sprintf("tr.attr-FeeURL td.cmpr-col.col%d a", index)
	link, exists := doc.Find(selector).Attr("href")
	if !exists || CleanText(link) == "" {
		return nil
	}
	return &link
}

func getTextOrEmpty(doc *goquery.Document, selector string) string {
	return CleanText(doc.Find(selector).Text())
}
