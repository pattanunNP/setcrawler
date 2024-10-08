package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Define structs with additional fields for extracted numeric values
type Product struct {
	Provider                     string                       `json:"provider"`
	Product                      string                       `json:"product"`
	FeaturesAndConditions        FeaturesAndConditions        `json:"features_and_conditions"`
	GeneralFees                  GeneralFees                  `json:"general_fees"`
	TransactionFeesDomestic      TransactionFeesDomestic      `json:"transaction_fees_domestic"`
	TransactionFeesInternational TransactionFeesInternational `json:"transaction_fees_international"`
	Insurance                    Insurance                    `json:"insurance"`
	AdditionalInfo               AdditionalInfo               `json:"additional_info"`
}

type FeaturesAndConditions struct {
	ProductType             string            `json:"product_type"`
	Network                 string            `json:"network"`
	Highlights              []string          `json:"highlights"`
	AgeRequirement          int               `json:"age_requirement"`
	ApplicantQualifications []string          `json:"applicant_qualificaitons"`
	UsageConditions         []string          `json:"usage_conditions"`
	CardExpiry              string            `json:"card_expiry"`
	PaymentOptions          []string          `json:"payment_options"`
	SupplementaryCard       SupplementaryCard `json:"supplementary_card"`
}

type SupplementaryCard struct {
	Available  []string `json:"available"`
	Conditions []string `json:"conditions"`
}

type GeneralFees struct {
	EntranceFee                 FeeDetail               `json:"entrance_fee"`
	AnnualFee                   FeeDetail               `json:"annual_fee"`
	CardReplacementFee          FeeDetail               `json:"card_replacement_fee"`
	PinReplacementFee           FeeDetail               `json:"pin_replacement_fee"`
	StatementCopyFee            FeeDetail               `json:"statement_copy_fee"`
	SlipCopyFee                 FeeDetailWithConditions `json:"slip_copy_fee"`
	TransactionInvestigationFee FeeDetail               `json:"transaction_investigation_fee"`
	OtherFees                   []string                `json:"other_fees"`
}

type FeeDetailWithConditions struct {
	OriginalText []string `json:"original_text"`
	Amount       []int    `json:"amount"`
}

type FeeDetail struct {
	OriginalText string `json:"original_text"`
	Amount       int    `json:"amount"`
}

type TransactionFeesDomestic struct {
	FreeTransactionsPerMonth        int       `json:"free_transactions_per_month"`
	CashWithdrawal                  []string  `json:"cash_withdrawal"`
	InServiceAreaBalanceInquiryFee  FeeDetail `json:"in_service_area_balance_inquiry_fee"`
	OutServiceAreaBalanceInquiryFee FeeDetail `json:"out_service_area_balance_inquiry_fee"`
	InServiceAreaCashWithdrawalFee  FeeDetail `json:"in_service_area_cash_withdrawal_fee"`
	OutServiceAreaCashWithdrawalFee FeeDetail `json:"out_service_area_cash_withdrawal_fee"`
	InServiceAreaTransferFee        FeeDetail `json:"in_service_area_transfer_fee"`
	OutServiceAreaTransferFee       FeeDetail `json:"out_service_area_transfer_fee"`
	TransferBetweenProvidersFee     FeeDetail `json:"transfer_between_providers_fee"`
	Under10000Fee                   FeeDetail `json:"under_10000_fee"`
	Between10001And50000Fee         FeeDetail `json:"between_10001_and_50000_fee"`
	AdditionalFee                   FeeDetail `json:"additional_fee"`
	OtherConditions                 string    `json:"other_conditions"`
}

type TransactionFeesInternational struct {
	WithdrawalFee         FeeDetail `json:"withdrawal_fee"`
	BalanceInquiryFee     FeeDetail `json:"balance_inquiry_fee"`
	CurrencyConversionFee float64   `json:"currency_conversion_fee"`
}

type Insurance struct {
	InsuranceType           string `json:"insurance_type"`
	InsuranceCompany        string `json:"insurance_company"`
	MaxCoverageAmount       int    `json:"max_coverage_amount"`
	CoverageDetails         string `json:"coverage_details"`
	DeathDisabilityCoverage int    `json:"death_disability_coverage"`
	MedicalExpensesCoverage int    `json:"medical_expenses_coverage"`
	HospitalIncomeCoverage  int    `json:"hospital_income_coverage"`
	OtherBenefits           string `json:"other_benefits"`
	CoverageConditions      string `json:"coverage_conditions"`
	CoveragePeriod          string `json:"coverage_period"`
	Exclusions              string `json:"exclusions"`
	ClaimProcedure          string `json:"claim_procedure"`
	ContactInsuranceCompany string `json:"contact_insurance_company"`
}

type AdditionalInfo struct {
	ProductWebsite string `json:"product_website"`
	FeeWebsite     string `json:"fee_website"`
}

func extractIntFromString(input string) int {
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(input)
	if match == "" {
		return 0
	}
	value, err := strconv.Atoi(match)
	if err != nil {
		return 0
	}
	return value
}

func extractAmounts(conditions []string) []int {
	var amounts []int
	for _, condition := range conditions {
		// Use regex to find all numbers in the string
		re := regexp.MustCompile(`\d+`)
		matches := re.FindAllString(condition, -1)
		for _, match := range matches {
			amount, err := strconv.Atoi(match)
			if err == nil {
				amounts = append(amounts, amount)
			}
		}
	}
	return amounts
}

func extractFloatFromString(input string) float64 {
	re := regexp.MustCompile(`\d+(\.\d+)?`)
	match := re.FindString(input)
	if match == "" {
		return 0
	}
	value, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0
	}
	return value
}

var hyphenOrNumberRegex = regexp.MustCompile(`(?:\s*-\s*|\s*\d+\.\s*)`)

func splitByHyphenOrNumber(input string) []string {
	result := hyphenOrNumberRegex.Split(input, -1)

	for i, val := range result {
		result[i] = strings.TrimSpace(val)
	}
	return filterEmptyStrings(result)
}

func filterEmptyStrings(arr []string) []string {
	var filtered []string
	for _, str := range arr {
		if str != "" {
			filtered = append(filtered, str)
		}
	}
	return filtered
}

func splitByHyphen(input string) []string {
	// Split the text by "-" and remove newline and extra spaces
	parts := strings.Split(input, "-")
	var cleanedParts []string
	for _, part := range parts {
		cleaned := strings.TrimSpace(part)
		cleaned = strings.ReplaceAll(cleaned, "\n", "")
		cleaned = strings.ReplaceAll(cleaned, "เงื่อนไข:", "") // Remove unnecessary parts
		if cleaned != "" {
			cleanedParts = append(cleanedParts, cleaned)
		}
	}
	return cleanedParts
}

func splitAndCleanText(input string) []string {
	// Replace multiple spaces with a single space and remove newline characters
	cleaned := strings.ReplaceAll(input, "\n", "")
	cleaned = strings.Join(strings.Fields(cleaned), " ")

	// Split the cleaned text by hyphen
	parts := strings.Split(cleaned, "-")

	var cleanedParts []string
	for _, part := range parts {
		cleanedPart := strings.TrimSpace(part)
		if cleanedPart != "" {
			cleanedParts = append(cleanedParts, cleanedPart)
		}
	}
	return cleanedParts
}

func main() {
	url := "https://app.bot.or.th/1213/MCPD/ProductApp/Debit/CompareProductList"
	payloadTemplate := `{"ProductIdList":"1234,1459,1466,1606,1460,1468,1463,976,13,14,1378,1379,1380,1381,1365,1366,733,731,1492,721,722,723,1502,720,67,1377,719,718,726,727,725,724,1256,1382,1474,954,950,961,1467,1237,1634,1490,246,1642,1618,16,1499,138,11,1587,682,472,1585,1239,1236,1593,1504,1,946,958,730,732,960,1473,1476,1475,12,15,1457,1462,1469,1456,474,744,1627,1461,1458,17,1503,1478,1477,1482,1481,1484,1480,1483,473,1235,1491,1496,1240,1494,752,749,1641,1488,1485,1487,477,972,750,746,1592,2,1501,1472,1464,748,751,1594,475,964,1465,1612,962,1471,1470,1649,139,1500,1497,1493,1498,1601,1489,1369,1367,1238,1479,1486,728,729,1495","Page":%d,"Limit":3}`
	// payloadTemplate := `{"ProductIdList":"1234,1459,1466,1606,1460,1468,1463,976,13,14,1378,1379,1380,1381,1365,1366,733,731,1492,721,722,723,1502,720","Page":%d,"Limit":3}`

	var allProducts []Product

	totalPages := 1 // Assume at least 1 page
	for page := 1; page <= totalPages; page++ {
		payload := fmt.Sprintf(payloadTemplate, page)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		setHeaders(req)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Failed to retrieve data:", resp.Status)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			fmt.Println("Error parsing HTML:", err)
			return
		}

		doc.Find("th.cmpr-col").Each(func(i int, s *goquery.Selection) {
			provider := cleanedText(s.Find("th.col-s span").Last().Text())
			product := cleanedText(s.Find("th.font-black.text-center").Text())

			productType := cleanedText(doc.Find(fmt.Sprintf(".attr-productTypeName .col%d", i+1)).Text())
			network := cleanedText(doc.Find(fmt.Sprintf(".attr-networkTypeName .col%d", i+1)).Text())

			var highlights []string
			doc.Find(fmt.Sprintf(".attr-productBenefitMain .col%d span", i+1)).Each(func(index int, item *goquery.Selection) {
				highlight := strings.TrimSpace(item.Text())
				if highlight != "" {
					highlights = append(highlights, highlight)
				}
			})

			ageStr := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-cardholderAge .col%d span", i+1)).Text())
			ageRequirement, err := strconv.Atoi(strings.ReplaceAll(ageStr, " ปีขึ้นไป", ""))
			if err != nil {
				ageRequirement = 0
			}

			applicantQualifications := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-conditionToApply .col%d span", i+1)).Text())
			applicantQualificationsArray := splitByHyphen(applicantQualifications)

			usageConditions := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-conditionToUse .col%d span", i+1)).Text())
			usageConditionsArray := splitByHyphen(usageConditions)

			cardExpiry := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-usagePeriod .col%d span", i+1)).Text())

			var paymentOptions []string
			doc.Find(fmt.Sprintf(".attr-payment .col%d span", i+1)).Each(func(index int, item *goquery.Selection) {
				option := strings.TrimSpace(item.Text())
				if option != "" {
					paymentOptions = append(paymentOptions, splitByHyphen(option)...)
				}
			})

			supplementaryAvailable := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-supplementaryCard .col%d span", i+1)).Text())
			supplementaryAvailableArray := splitByHyphen(supplementaryAvailable)

			var supplementaryConditionsArray []string

			var conditionsText string
			doc.Find(fmt.Sprintf(".attr-otherCondition .col%d span", i+1)).Each(func(index int, item *goquery.Selection) {
				condition := strings.TrimSpace(item.Text())
				if condition != "" {
					conditionsText += condition
				}
			})
			supplementaryConditionsArray = splitByHyphen(conditionsText)

			entranceFeeText := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-cardHolderEntranceFeeDisplay .col%d span", i+1)).Text())
			entranceFeeAmount := extractIntFromString(entranceFeeText)
			annualFeeText := cleanedText(doc.Find(fmt.Sprintf(".attr-annualFeeDisplay .col%d span", i+1)).Text())
			annualFeeAmount := extractIntFromString(annualFeeText)
			cardReplacementFeeText := cleanedText(doc.Find(fmt.Sprintf(".attr-replacementCardFee .col%d span", i+1)).Text())
			cardReplacementFeeAmount := extractIntFromString(cardReplacementFeeText)
			pinReplacementFeeText := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-replacementOfCardPINFee .col%d span", i+1)).Text())
			pinReplacementFeeAmount := extractIntFromString(pinReplacementFeeText)
			statementCopyFeeText := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-copyofStatementFee .col%d span", i+1)).Text())
			statementCopyFeeAmount := extractIntFromString(statementCopyFeeText)
			slipCopyFeeText := doc.Find(fmt.Sprintf(".attr-copyOfSalesSlipFee .col%d span", i+1)).Text()
			slipCopyFeeOriginalText := splitAndCleanText(slipCopyFeeText)     // Clean and split the text
			slipCopyFeeConditions := splitAndCleanConditions(slipCopyFeeText) // Further split conditions
			slipCopyFeeAmounts := extractAmounts(slipCopyFeeConditions)
			transactionInvestigationFeeText := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-transactionverificationFee .col%d span", i+1)).Text())
			transactionInvestigationFeeAmount := extractIntFromString(transactionInvestigationFeeText)
			otherFees := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-otherFee .col%d span", i+1)).Text())
			otherFeesArray := splitByHyphen(otherFees)

			freeTransactionsPerMonthStr := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-11 .col%d span", i+1)).Text())
			freeTransactionsPerMonth := extractIntFromString(freeTransactionsPerMonthStr)

			cashWithdrawal := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-11 .col%d span", i+1)).Text())
			cashWithdrawalArray := splitByHyphenOrNumber(cashWithdrawal)

			inServiceAreaBalanceInquiryFeeText := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-31 .col%d span", i+1)).Text())
			inServiceAreaBalanceInquiryFeeAmount := extractIntFromString(inServiceAreaBalanceInquiryFeeText)

			outServiceAreaBalanceInquiryFeeText := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-32 .col%d span", i+1)).Text())
			outServiceAreaBalanceInquiryFeeAmount := extractIntFromString(outServiceAreaBalanceInquiryFeeText)

			inServiceAreaCashWithdrawalFeeText := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-33 .col%d span", i+1)).Text())
			inServiceAreaCashWithdrawalFeeAmount := extractIntFromString(inServiceAreaCashWithdrawalFeeText)

			outServiceAreaCashWithdrawalFeeText := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-34 .col%d span", i+1)).Text())
			outServiceAreaCashWithdrawalFeeAmount := extractIntFromString(outServiceAreaCashWithdrawalFeeText)

			inServiceAreaTransferFeeText := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-35 .col%d span", i+1)).Text())
			inServiceAreaTransferFeeAmount := extractIntFromString(inServiceAreaTransferFeeText)

			outServiceAreaTransferFeeText := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-36 .col%d span", i+1)).Text())
			outServiceAreaTransferFeeAmount := extractIntFromString(outServiceAreaTransferFeeText)

			transferBetweenProvidersFeeText := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-FeeTranferDiffProvider .col%d span", i+1)).Text())
			transferBetweenProvidersFeeAmount := extractIntFromString(transferBetweenProvidersFeeText)

			under10000FeeText := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-41 .col%d span", i+1)).Text())
			under10000FeeAmount := extractIntFromString(under10000FeeText)

			between10001And50000FeeText := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-42 .col%d span", i+1)).Text())
			between10001And50000FeeAmount := extractIntFromString(between10001And50000FeeText)

			additionalFeeText := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-FeeAdditional .col%d span", i+1)).Text())
			additionalFeeAmount := extractIntFromString(additionalFeeText)

			otherConditions := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-FeeOtherCondition .col%d span", i+1)).Text())

			withdrawalFeeText := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-51 .col%d span", i+1)).Text())
			withdrawalFeeAmount := extractIntFromString(withdrawalFeeText)

			balanceInquiryFeeText := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-52 .col%d span", i+1)).Text())
			balanceInquiryFeeAmount := extractIntFromString(balanceInquiryFeeText)

			currencyConversionFeeStr := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-53 .col%d span", i+1)).Text())
			currencyConversionFee := extractFloatFromString(currencyConversionFeeStr)

			insuranceType := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-insuranceTypeName .col%d span", i+1)).Text())
			insuranceCompany := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-insuranceCompanyName .col%d span", i+1)).Text())

			maxCoverageAmountStr := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-maxCoverageDisplay .col%d span", i+1)).Text())
			maxCoverageAmount := extractIntFromString(maxCoverageAmountStr)

			coverageDetails := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-OtherBenefits .col%d span", i+1)).Text())

			deathDisabilityCoverageStr := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-OtherBenefits .col%d span", i+1)).Text())
			deathDisabilityCoverage := extractIntFromString(deathDisabilityCoverageStr)

			medicalExpensesCoverageStr := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-OtherBenefits .col%d span", i+1)).Text())
			medicalExpensesCoverage := extractIntFromString(medicalExpensesCoverageStr)

			hospitalIncomeCoverageStr := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-OtherBenefits .col%d span", i+1)).Text())
			hospitalIncomeCoverage := extractIntFromString(hospitalIncomeCoverageStr)

			otherBenefits := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-OtherBenefits .col%d span", i+1)).Text())

			coverageConditions := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-CoveragePeriod .col%d span", i+1)).Text())
			coveragePeriod := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-CoveragePeriod .col%d span", i+1)).Text())
			exclusions := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-CoveragePeriod .col%d span", i+1)).Text())
			claimProcedure := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-CoveragePeriod .col%d span", i+1)).Text())
			contactInsuranceCompany := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-CoveragePeriod .col%d span", i+1)).Text())

			productWebsite := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-uRL .col%d a", i+1)).AttrOr("href", ""))
			feeWebsite := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-uRLFee .col%d span", i+1)).Text())

			allProducts = append(allProducts, Product{
				Provider: provider,
				Product:  product,
				FeaturesAndConditions: FeaturesAndConditions{
					ProductType:             productType,
					Network:                 network,
					Highlights:              highlights,
					AgeRequirement:          ageRequirement,
					ApplicantQualifications: applicantQualificationsArray,
					UsageConditions:         usageConditionsArray,
					CardExpiry:              cardExpiry,
					PaymentOptions:          paymentOptions,
					SupplementaryCard: SupplementaryCard{
						Available:  supplementaryAvailableArray,
						Conditions: supplementaryConditionsArray,
					},
				},
				GeneralFees: GeneralFees{
					EntranceFee: FeeDetail{
						OriginalText: entranceFeeText,
						Amount:       entranceFeeAmount,
					},
					AnnualFee: FeeDetail{
						OriginalText: annualFeeText,
						Amount:       annualFeeAmount,
					},
					CardReplacementFee: FeeDetail{
						OriginalText: cardReplacementFeeText,
						Amount:       cardReplacementFeeAmount,
					},
					PinReplacementFee: FeeDetail{
						OriginalText: pinReplacementFeeText,
						Amount:       pinReplacementFeeAmount,
					},
					StatementCopyFee: FeeDetail{
						OriginalText: statementCopyFeeText,
						Amount:       statementCopyFeeAmount,
					},
					SlipCopyFee: FeeDetailWithConditions{
						OriginalText: slipCopyFeeOriginalText,
						Amount:       slipCopyFeeAmounts,
					},
					TransactionInvestigationFee: FeeDetail{
						OriginalText: transactionInvestigationFeeText,
						Amount:       transactionInvestigationFeeAmount,
					},
					OtherFees: otherFeesArray,
				},
				TransactionFeesDomestic: TransactionFeesDomestic{
					FreeTransactionsPerMonth: freeTransactionsPerMonth,
					CashWithdrawal:           cashWithdrawalArray,
					InServiceAreaBalanceInquiryFee: FeeDetail{
						OriginalText: inServiceAreaBalanceInquiryFeeText,
						Amount:       inServiceAreaBalanceInquiryFeeAmount,
					},
					OutServiceAreaBalanceInquiryFee: FeeDetail{
						OriginalText: outServiceAreaBalanceInquiryFeeText,
						Amount:       outServiceAreaBalanceInquiryFeeAmount,
					},
					InServiceAreaCashWithdrawalFee: FeeDetail{
						OriginalText: inServiceAreaCashWithdrawalFeeText,
						Amount:       inServiceAreaCashWithdrawalFeeAmount,
					},
					OutServiceAreaCashWithdrawalFee: FeeDetail{
						OriginalText: outServiceAreaCashWithdrawalFeeText,
						Amount:       outServiceAreaCashWithdrawalFeeAmount,
					},
					InServiceAreaTransferFee: FeeDetail{
						OriginalText: inServiceAreaTransferFeeText,
						Amount:       inServiceAreaTransferFeeAmount,
					},
					OutServiceAreaTransferFee: FeeDetail{
						OriginalText: outServiceAreaTransferFeeText,
						Amount:       outServiceAreaTransferFeeAmount,
					},
					TransferBetweenProvidersFee: FeeDetail{
						OriginalText: transferBetweenProvidersFeeText,
						Amount:       transferBetweenProvidersFeeAmount,
					},
					Under10000Fee: FeeDetail{
						OriginalText: under10000FeeText,
						Amount:       under10000FeeAmount,
					},
					Between10001And50000Fee: FeeDetail{
						OriginalText: between10001And50000FeeText,
						Amount:       between10001And50000FeeAmount,
					},
					AdditionalFee: FeeDetail{
						OriginalText: additionalFeeText,
						Amount:       additionalFeeAmount,
					},
					OtherConditions: otherConditions,
				},
				TransactionFeesInternational: TransactionFeesInternational{
					WithdrawalFee: FeeDetail{
						OriginalText: withdrawalFeeText,
						Amount:       withdrawalFeeAmount,
					},
					BalanceInquiryFee: FeeDetail{
						OriginalText: balanceInquiryFeeText,
						Amount:       balanceInquiryFeeAmount,
					},
					CurrencyConversionFee: currencyConversionFee,
				},
				Insurance: Insurance{
					InsuranceType:           insuranceType,
					InsuranceCompany:        insuranceCompany,
					MaxCoverageAmount:       maxCoverageAmount,
					CoverageDetails:         coverageDetails,
					DeathDisabilityCoverage: deathDisabilityCoverage,
					MedicalExpensesCoverage: medicalExpensesCoverage,
					HospitalIncomeCoverage:  hospitalIncomeCoverage,
					OtherBenefits:           otherBenefits,
					CoverageConditions:      coverageConditions,
					CoveragePeriod:          coveragePeriod,
					Exclusions:              exclusions,
					ClaimProcedure:          claimProcedure,
					ContactInsuranceCompany: contactInsuranceCompany,
				},
				AdditionalInfo: AdditionalInfo{
					ProductWebsite: productWebsite,
					FeeWebsite:     feeWebsite,
				},
			})
		})

		// Wait before making the next request
		time.Sleep(2 * time.Second)
	}

	jsonData, err := json.MarshalIndent(allProducts, "", " ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	file, err := os.Create("debit.json")
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing to JSON file:", err)
		return
	}

	fmt.Println("JSON file 'debit.json' created successfully")
}

func cleanedText(input string) string {
	cleaned := strings.ReplaceAll(input, `\u003c`, "<")
	cleaned = strings.ReplaceAll(cleaned, `\u003e`, ">")
	cleaned = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(cleaned, "")
	cleaned = strings.ReplaceAll(cleaned, "\n", "")
	cleaned = strings.ReplaceAll(cleaned, "\t", "")
	cleaned = strings.TrimSpace(cleaned)
	spaceRegex := regexp.MustCompile(`\s+`)
	cleaned = spaceRegex.ReplaceAllString(cleaned, " ")

	return cleaned
}

func setHeaders(req *http.Request) {
	req.Header.Set("accept", "text/plain, */*; q=0.01")
	req.Header.Set("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Header.Set("cookie", `verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m04nlc8h&sl=0&tt=0"; _uid6672=16B5DEBD.21; _ctout6672=1; _ga_NLQFGWVNXN=GS1.1.1724320535.24.1.1724323292.2.0.0`)
	req.Header.Set("origin", "https://app.bot.or.th")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://app.bot.or.th/1213/MCPD/ProductApp/Debit/CompareProduct")
	req.Header.Set("sec-ch-ua", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	req.Header.Set("verificationtoken", `ToQWsd6JpdywXWWiHC8F7T8eQZkBkBxMij7tw9cmR-ustXzjzA5kmlRIalkuj-0WblKIrki2wYe-iFBJdeGpAsL5UDE7ix8yTesristz_WY1,9R2bjBfVukm3UcFSumGCpsGpB097wGQ0InKyeYA45PZanVekI2TPT-Jc9AOGVGhWT16oGo44ZKOAzFhfM1Y8uDiDI3hm5n6jnKVf5IlbPL01`)
	req.Header.Set("x-requested-with", "XMLHttpRequest")
}

func splitAndCleanConditions(input string) []string {
	// Split the text by "-" and remove newline and extra spaces
	parts := strings.Split(input, "-")
	var cleanedParts []string
	for _, part := range parts {
		cleaned := strings.TrimSpace(part)
		cleaned = strings.ReplaceAll(cleaned, "\n", "")
		if cleaned != "" {
			cleanedParts = append(cleanedParts, cleaned)
		}
	}
	return cleanedParts
}
