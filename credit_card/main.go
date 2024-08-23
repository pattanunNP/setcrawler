package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ProductDetails struct {
	ProductName                string               `json:"product_name"`
	ProductFeatures            Features             `json:"product_features_conditions"`
	PrimaryCardApplicantAge    int                  `json:"primary_card_applicant_age"`
	MinimumIncomeAndConditions IncomeConditions     `json:"minimum_income_and_conditions"`
	InterestFreePeriod         int                  `json:"interest_free_period"`
	CreditLimit                int                  `json:"credit_limit"`
	GeneralFees                []GeneralFee         `json:"general_fees"`
	PaymentMethods             []PaymentMethod      `json:"payment_methods"`
	LatePaymentPenalties       []LatePaymentPenalty `json:"late_payment_penalties"`
	CashWithdrawalFees         []CashWithdrawalFee  `json:"cash_withdrawal_fees"`
	SupplementaryCard          SupplementaryCard    `json:"supplementary_card"`
	AdditionInfo               AdditionInfo         `json:"additioninfo"`
}

type Features struct {
	BenefitType string   `json:"benefit_type"`
	Details     []string `json:"details"`
}

type IncomeConditions struct {
	IncomeRequirement string   `json:"income_requirement"`
	Conditions        []string `json:"conditions"`
}

type GeneralFee struct {
	FeeType     string       `json:"fee_type"`
	Amount      int          `json:"amount,omitempty"`
	AmountRange *AmountRange `json:"amount_range"`
	Conditions  []string     `json:"conditions"`
}

type AmountRange struct {
	Min int `json:"min,omitempty"`
	Max int `json:"max,omitempty"`
}

type PaymentMethod struct {
	MethodType string   `json:"method_type"`
	Details    []string `json:"details"`
}

type LatePaymentPenalty struct {
	PenaltyType      string   `json:"penalty_type"`
	AmountPercentage int      `json:"amount_percentage,omitempty"`
	MinimumAmount    int      `json:"minimum_amount,omitempty"`
	InterestRate     int      `json:"interest_rate,omitempty"`
	Conditions       []string `json:"conditions,omitempty"`
}

type CashWithdrawalFee struct {
	FeeType          string `json:"fee_type"`
	InterestRate     int    `json:"interest_rate,omitempty"`
	AmountPercentage int    `json:"amount_percentage,omitempty"`
	ConditionsType   string `json:"conditions_type"`
	Details          string `json:"details,omitempty"`
}

type SupplementaryCard struct {
	MaxNumberOfCards int            `json:"max_number_of_cards"`
	AgeRequirement   AgeRequirement `json:"age_requirement"`
	Fees             []FeeDetails   `json:"fees"`
}

type AgeRequirement struct {
	MinAge     int      `json:"min_age"`
	MaxAge     int      `json:"max_age"`
	Conditions []string `json:"conditions"`
}

type FeeDetails struct {
	FeeType    string   `json:"fee_type"`
	Conditions []string `json:"conditions"`
}

type AdditionInfo struct {
	ProductURL string `json:"product_url"`
	FeeURL     string `json:"fee_url"`
}

func main() {
	url := "https://app.bot.or.th/1213/MCPD/ProductApp/Credit/CompareProductList"
	payloadTemplate := `{"ProductIdList":"3629,3622,3630,4633,4632,4634,4655,4656,1604,1607,2259,5573,5161,2251,5534,5537,5177,5491,5522,5523,3627,3631,3600,3664,3624,3601,3625,3603,3633,3626,3604,3661,3620,3634,3636,3635,3606,3638,3637,3640,3639,3605,3672,3668,3628,3662,3621,3648,3607,3642,3649,3643,3651,3613,3644,3652,3609,3645,3653,3646,3647,3615,3656,3610,3650,3671,3658,3655,3616,3670,3612,4653,4657,4658,4659,4660,5568,5528,2256,5570,2245,5531,4471,4472,4467,5479,5494,5484,5497,5483,5503,1603,4760,4765,5473,3804,3800,4445,4482,4483,4479,4480,4452,4453,4476,4477,4454,4456,4457,4458,4460,4461,4462,4463,4464,4465,4444,4446,4447,4449,4450,4473,4470,5486,5498,5501,5540,5555,2242,5560,5496,1605,1600,1601,1602,4475,2246,5563,5565,5539,5567,5525,5556,2244,5562,2255,5536,2260,5574,2252,5546,5548,5538,5577,5518,5514,5521,5517,5527,3802,5529,5542,2258,5572,2249,5533,5516,4459,4474,4469,5499,5492,5481,5148,5435,3608,3641,3669,3611,3617,3660,5114,5287,5282,5285,5193,5137,5211,5173,5162,5156,5222,5126,5202,5158,5147,5240,5294,5296,5323,5306,5300,5329,749,753,752,751,750,754,5489,5500,5480,5495,4484,4481,4478,4455,4466,4451,3808,3797,3798,5502,5482,5566,2257,5571,2247,5532,5575,5576,5541,5554,2248,5544,2240,5558,5429,3996,3993,5520,3994,3995,3602,3632,3654,3665,3663,3618,3657,3623,1606,5526,5553,5552,3806,3799,3801,2253,5535,5564,3614,3666,3667,3619,3659,2250,5545,2241,5559,5549,2254,5547,2243,5561,5550,4448,4468,5519,5493,5266,5305,5304,5246,4654,5462,5213,5233,5478,5467,3997,5180,5557,5543,5551,5524,5335","Page":%d,"Limit":3}`

	// First, get the total number of pages from the initial request
	initialPage := 1
	payload := fmt.Sprintf(payloadTemplate, initialPage)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	setHeaders(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to retrieve data: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Error loading HTML: %v", err)
	}

	// Extract total number of pages
	totalPages := 1
	doc.Find("ul.pagination li.MoveLast a").Each(func(i int, s *goquery.Selection) {
		dataPage, exists := s.Attr("data-page")
		if exists {
			totalPages, err = strconv.Atoi(dataPage)
			if err != nil {
				log.Fatalf("Error converting data-page to int: %v", err)
			}
		}
	})

	var allProducts []ProductDetails

	// Loop through each page
	for page := 1; page <= totalPages; page++ {
		payload := fmt.Sprintf(payloadTemplate, page)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}

		setHeaders(req)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error sending request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Failed to retrieve data: %v", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response: %v", err)
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			log.Fatalf("Error loading HTML: %v", err)
		}

		colCount := 1
		doc.Find("th.font-black.text-center").Each(func(i int, s *goquery.Selection) {
			productName := strings.TrimSpace(s.Find("span.text-bold").Text() + " " + s.Find("span.txt-normal").Text())
			if productName != "" {
				product := ProductDetails{
					ProductName:                productName,
					ProductFeatures:            extractFeatures(doc, colCount),
					PrimaryCardApplicantAge:    extractAgeRequirement(doc, colCount),
					MinimumIncomeAndConditions: extractIncomeConditions(doc, colCount),
					InterestFreePeriod:         extractInterestFreePeriod(doc, colCount),
					CreditLimit:                extractCreditLimit(doc, colCount),
					GeneralFees:                extractGeneralFees(doc, colCount),
					PaymentMethods:             extractPaymentMethods(doc, colCount),
					LatePaymentPenalties:       extractLatePaymentPenalties(doc, colCount),
					CashWithdrawalFees:         extractCashWithdrawalFees(doc, colCount),
					SupplementaryCard:          extractSupplementaryCard(doc, colCount),
					AdditionInfo:               extractAdditionInfo(doc, colCount),
				}
				allProducts = append(allProducts, product)
				colCount++
			}
		})

		// Stop for 5 seconds before making the next request to avoid overloading the server
		time.Sleep(5 * time.Second)
	}

	// Convert the combined products to JSON and save to a file
	jsonData, err := json.MarshalIndent(allProducts, "", "  ")
	if err != nil {
		log.Fatalf("Failed to convert struct to JSON: %v", err)
	}

	// Save JSON to a file
	err = os.WriteFile("credit_card_compare.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write JSON to file: %v", err)
	}

	fmt.Println("Product details saved to credit_card_compare.json")
}

func setHeaders(req *http.Request) {
	req.Header.Set("Accept", "text/plain, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Cookie", `verify=test; verify=test; verify=test; mycookie=\u0021IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; _uid6672=16B5DEBD.6; visit_time=2171; _ga_NLQFGWVNXN=GS1.1.1723798115.10.1.1723798115.60.0.0`)
	req.Header.Set("Origin", "https://app.bot.or.th")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://app.bot.or.th/1213/MCPD/ProductApp/Credit/CompareProduct")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	req.Header.Set("Verificationtoken", "jZKLX6IRz_o2TRzkK3Muo-MQJJyTKCCqYH_nqRqFPSlMqzPukrPg--T1qwu7lBF6ikqALmGzObfK1bEHefY_iHAfKTD0-PqTh6CzSbnPS4M1,qZsSwtss8Ueiv1fJBVwqoyWusJv_BnjoVmkYsXcNn3E6JAWC5FBaJ1jlFqgXFD_9nIEjon7NPJ-AfhnaK11irhNoXh9jfWaN2j8YIaZp-3g1")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
}

func splitAndClean(text string) []string {
	// Initialize an empty slice to store the parts
	parts := []string{}

	// Split the text by dash and trim whitespace
	for _, part := range strings.Split(text, "-") {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}

	// Return nil if no valid parts are found to avoid empty arrays in JSON
	if len(parts) == 0 {
		return nil
	}

	return parts
}

func extractFeatures(doc *goquery.Document, colCount int) Features {
	benefitType := cleanText(doc.Find(fmt.Sprintf("tr.attr-header.attr-productBenefitType.trbox-shadow .cmpr-col.col%d", colCount)).Text())
	detailsText := cleanText(doc.Find(fmt.Sprintf("tr.attr-header.attr-productBenefitMain.trbox-shadow .cmpr-col.col%d span", colCount)).Text())

	details := []string{}
	for _, part := range strings.Split(detailsText, "-") {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			details = append(details, trimmed)
		}
	}

	return Features{
		BenefitType: benefitType,
		Details:     details,
	}
}

func extractAgeRequirement(doc *goquery.Document, colCount int) int {
	var age int
	ageText := cleanText(doc.Find(fmt.Sprintf("tr.attr-header.attr-primaryAgeOfHolder.trbox-shadow .cmpr-col.col%d span", colCount)).Text())
	fmt.Sscanf(ageText, "%d", &age)
	return age
}

func extractIncomeConditions(doc *goquery.Document, colCount int) IncomeConditions {
	incomeRequirement := cleanText(doc.Find(fmt.Sprintf("tr.attr-header.attr-minIncomePerMonthDisplay.trbox-shadow .cmpr-col.col%d", colCount)).Text())
	conditions := []string{cleanText(doc.Find(fmt.Sprintf("tr.attr-header.attr-minIncomePerMonthDisplay.trbox-shadow .cmpr-col.col%d .text-primary", colCount)).Text())}

	return IncomeConditions{
		IncomeRequirement: incomeRequirement,
		Conditions:        conditions,
	}
}

func extractInterestFreePeriod(doc *goquery.Document, colCount int) int {
	var days int
	interestText := cleanText(doc.Find(fmt.Sprintf("tr.attr-header.attr-interestFreePeriodDisplay.trbox-shadow .cmpr-col.col%d span", colCount)).Text())
	fmt.Sscanf(interestText, "%d", &days)
	return days
}

func extractCreditLimit(doc *goquery.Document, colCount int) int {
	var limit int
	limitText := cleanText(doc.Find(fmt.Sprintf("tr.attr-header.attr-crditLineMax.trbox-shadow .cmpr-col.col%d span", colCount)).Text())
	fmt.Sscanf(limitText, "%d", &limit)
	return limit
}

func extractGeneralFees(doc *goquery.Document, colCount int) []GeneralFee {
	var generalFees []GeneralFee

	feeTypes := []struct {
		FeeTypeSelector string
		AmountSelector  string
		IsAmountRange   bool
	}{
		{"attr-primaryHolderEntranceFeeDisplay", "attr-primaryHolderEntranceFeeDisplay", false},
		{"attr-primaryHolderAnnualFee", "attr-primaryHolderAnnualFee", true},
		{"attr-replacementCardFee", "attr-replacementCardFee", false},
		{"attr-CostFXRisk", "attr-CostFXRisk", true},
		{"attr-replacementCardFPinFee", "attr-replacementCardFPinFee", false},
		{"attr-copyStatementFee", "attr-copyStatementFee", false},
		{"attr-TransactionVerifyFee", "attr-TransactionVerifyFee", true},
		{"attr-copySaleSlipFee", "attr-copySaleSlipFee", false},
		{"attr-fineChequeReturn", "attr-fineChequeReturn", false},
		{"attr-GovernmentAgencyRelatedPaymentFee", "attr-GovernmentAgencyRelatedPaymentFee", true},
		{"attr-otherFee", "attr-otherFee", false},
	}

	for _, feeType := range feeTypes {
		fee := GeneralFee{
			FeeType: cleanText(doc.Find(fmt.Sprintf("tr.%s .text-center.frst-col", feeType.FeeTypeSelector)).Text()),
		}
		conditionsText := cleanText(doc.Find(fmt.Sprintf("tr.%s .text-primary", feeType.FeeTypeSelector)).Text())
		fee.Conditions = splitAndClean(conditionsText)

		if feeType.IsAmountRange {
			amountRange := extractAmountRange(doc, colCount, feeType.AmountSelector)
			// Set AmountRange only if it has non-zero values
			if amountRange.Min > 0 || amountRange.Max > 0 {
				fee.AmountRange = &amountRange
			} else {
				fee.AmountRange = nil
			}
		} else {
			amount := extractAmount(doc, colCount, feeType.AmountSelector)
			fee.Amount = amount
			// Set AmountRange to nil since it's not a range type
			if amount > 0 {
				fee.AmountRange = &AmountRange{Min: amount, Max: amount}
			} else {
				fee.AmountRange = nil
			}
		}

		// Append the fee to the list
		generalFees = append(generalFees, fee)
	}

	return generalFees
}

func extractAmount(doc *goquery.Document, colCount int, selector string) int {
	var amount int
	amounText := cleanText(doc.Find(fmt.Sprintf("tr.%s .cmpr-col.col%d span", selector, colCount)).Text())
	fmt.Sscanf(amounText, "%d", &amount)
	return amount
}

func extractAmountRange(doc *goquery.Document, colCount int, selector string) AmountRange {
	var min, max int
	amountText := cleanText(doc.Find(fmt.Sprintf("tr.%s .cmpr-col.col%d span", selector, colCount)).Text())

	if strings.Contains(amountText, "-") {
		fmt.Sscanf(amountText, "%d-%d", &min, &max)
	} else {
		fmt.Sscanf(amountText, "%d", &min)
		max = min
	}
	return AmountRange{Min: min, Max: max}
}

func extractPaymentMethods(doc *goquery.Document, colCount int) []PaymentMethod {
	var paymentMethods []PaymentMethod

	paymentTypes := []struct {
		MethodTypeSelector string
		DetailsSelector    string
	}{
		{"attr-freePaymentChannel", "attr-freePaymentChannel"},
		{"attr-directDebitFromAccountFee", "attr-directDebitFromAccountFee"},
		{"attr-directDebitFromAccountFeeOther", "attr-directDebitFromAccountFeeOther"},
		{"attr-BankCounterServiceFee", "attr-BankCounterServiceFee"},
		{"attr-BankCounterServiceFeeOther", "attr-BankCounterServiceFeeOther"},
		{"attr-CounterServiceFeeOther", "attr-CounterServiceFeeOther"},
		{"attr-paymentOnlineFee", "attr-paymentOnlineFee"},
		{"attr-paymentCDMATMFee", "attr-paymentCDMATMFee"},
		{"attr-paymentPhoneFee", "attr-paymentPhoneFee"},
		{"attr-paymentChequeOrMoneyOrderFee", "attr-paymentChequeOrMoneyOrderFee"},
		{"attr-paymentOtherChannelFee", "attr-paymentOtherChannelFee"},
	}

	for _, pt := range paymentTypes {
		methodType := cleanText(doc.Find(fmt.Sprintf("tr.%s .text-center.frst-col", pt.MethodTypeSelector)).Text())
		detailsText := cleanText(doc.Find(fmt.Sprintf("tr.%s .cmpr-col.col%d span", pt.DetailsSelector, colCount)).Text())

		details := splitAndClean(detailsText)

		if methodType != "" && len(details) > 0 {
			paymentMethods = append(paymentMethods, PaymentMethod{
				MethodType: methodType,
				Details:    details,
			})
		}
	}

	return paymentMethods
}

func extractLatePaymentPenalties(doc *goquery.Document, colCount int) []LatePaymentPenalty {
	var penalties []LatePaymentPenalty

	penaltyTypes := []struct {
		PenaltyTypeSelector      string
		AmountPercentageSelector string
		MinimumAmountSelector    string
		InterestRateSelector     string
		ConditionsSelector       string
	}{
		{"attr-minAmountRequiredPaymentDisplay", "attr-minAmountRequiredPaymentDisplay", "attr-minAmountRequiredPaymentDisplay", "", ""},
		{"attr-interestPenaltiyServiceFeeAndOtherChargeDisplay", "", "", "attr-interestPenaltiyServiceFeeAndOtherChargeDisplay", "attr-interestPenaltiyServiceFeeAndOtherChargeDisplay"},
		{"attr-debtCollectionFee", "", "", "", "attr-debtCollectionFee"},
	}

	for _, pt := range penaltyTypes {
		penalty := LatePaymentPenalty{
			PenaltyType:      cleanText(doc.Find(fmt.Sprintf("tr.%s .text-center.frst-col", pt.PenaltyTypeSelector)).Text()),
			AmountPercentage: extractPercentage(doc, colCount, pt.AmountPercentageSelector),
			MinimumAmount:    extractMinimumAmount(doc, colCount, pt.MinimumAmountSelector),
			InterestRate:     extractInterestRate(doc, colCount, pt.InterestRateSelector),
		}

		conditionsText := cleanText(doc.Find(fmt.Sprintf("tr.%s .text-primary", pt.ConditionsSelector)).Text())
		penalty.Conditions = splitAndClean(conditionsText)
		penalties = append(penalties, penalty)
	}

	return penalties
}

func extractCashWithdrawalFees(doc *goquery.Document, colCount int) []CashWithdrawalFee {
	var fees []CashWithdrawalFee

	feeTypes := []struct {
		FeeTypeSelector          string
		InterestRateSelector     string
		AmountPercentageSelector string
		ConditionsTypeSelector   string
		DetailsSelector          string
	}{
		{"attr-InterestPenaltyServiceFeeAndOtherChargeCashAdvance", "attr-InterestPenaltyServiceFeeAndOtherChargeCashAdvance", "", "", ""},
		{"attr-cashAdvanceFee", "", "attr-cashAdvanceFee", "", ""},
		{"attr-cashAdvanceCondition", "", "", "attr-cashAdvanceCondition", "attr-cashAdvanceCondition"},
		{"attr-cashAdvanceAmountMin", "", "", "attr-cashAdvanceAmountMin", "attr-cashAdvanceAmountMin"},
	}

	for _, ft := range feeTypes {
		fee := CashWithdrawalFee{
			FeeType:          cleanText(doc.Find(fmt.Sprintf("tr.%s .text-center.frst-col", ft.FeeTypeSelector)).Text()),
			InterestRate:     extractInterestRate(doc, colCount, ft.InterestRateSelector),
			AmountPercentage: extractPercentage(doc, colCount, ft.AmountPercentageSelector),
			ConditionsType:   cleanText(doc.Find(fmt.Sprintf("tr.%s .text-center.frst-col", ft.ConditionsTypeSelector)).Text()),
			Details:          cleanText(doc.Find(fmt.Sprintf("tr.%s .cmpr-col.col%d span", ft.DetailsSelector, colCount)).Text()),
		}
		fees = append(fees, fee)
	}

	return fees
}

func extractPercentage(doc *goquery.Document, colCount int, selector string) int {
	var percentage int
	text := cleanText(doc.Find(fmt.Sprintf("tr.%s .cmpr-col.col%d span", selector, colCount)).Text())
	fmt.Sscanf(text, "%d%%", &percentage)
	return percentage
}

func extractMinimumAmount(doc *goquery.Document, colCount int, selector string) int {
	var minAmount int
	text := cleanText(doc.Find(fmt.Sprintf("tr.%s .cmpr-col.col%d span", selector, colCount)).Text())
	fmt.Sscanf(text, "ไม่น้อยกว่า %d บาท", &minAmount)
	return minAmount
}

func extractInterestRate(doc *goquery.Document, colCount int, selector string) int {
	var rate int
	text := cleanText(doc.Find(fmt.Sprintf("tr.%s .cmpr-col.col%d span", selector, colCount)).Text())
	fmt.Sscanf(text, "%d%%", &rate)
	return rate
}

func extractSupplementaryCard(doc *goquery.Document, colCount int) SupplementaryCard {
	maxNumberOfCards := extractInt(doc, colCount, "attr-supplementaryCardMax")
	minAge, maxAge := extractAgeRange(doc, colCount, "attr-supplementaryCardHolderAge")
	conditions := cleanText(doc.Find("tr.attr-supplementaryCardHolderAge .cmpr-col.col" + fmt.Sprint(colCount)).Text())

	fees := []FeeDetails{
		{
			FeeType:    cleanText(doc.Find("tr.attr-supplementaryCardHolderEntranceFeeDisplay .text-center.frst-co").Text()),
			Conditions: []string{cleanText(doc.Find("tr.attr-supplementaryCardHolderEntranceFeeDisplay .cmpr-col.col" + fmt.Sprint(colCount)).Text())},
		},
		{
			FeeType:    cleanText(doc.Find("tr.attr-supplementaryCardHolderAnnualFeeFirstYear .text-center.frst-co").Text()),
			Conditions: []string{cleanText(doc.Find("tr.attr-supplementaryCardHolderAnnualFeeFirstYear .cmpr-col.col" + fmt.Sprint(colCount)).Text())},
		},
	}

	return SupplementaryCard{
		MaxNumberOfCards: maxNumberOfCards,
		AgeRequirement: AgeRequirement{
			MinAge:     minAge,
			MaxAge:     maxAge,
			Conditions: []string{conditions},
		},
		Fees: fees,
	}
}

func extractAdditionInfo(doc *goquery.Document, colCount int) AdditionInfo {
	productURL := cleanText(doc.Find(fmt.Sprintf("tr.attr-url .cmpr-col.col%d a", colCount)).AttrOr("href", ""))
	feeURL := cleanText(doc.Find(fmt.Sprintf("tr.attr-feeurl .cmpr-col.col%d a", colCount)).AttrOr("href", ""))

	return AdditionInfo{
		ProductURL: productURL,
		FeeURL:     feeURL,
	}
}

func extractInt(doc *goquery.Document, colCount int, selector string) int {
	var result int
	text := cleanText(doc.Find(fmt.Sprintf("tr.%s .cmpr-col.col%d span", selector, colCount)).Text())
	fmt.Sscanf(text, "%d", &result)
	return result
}

func extractAgeRange(doc *goquery.Document, colCount int, selector string) (int, int) {
	var minAge, maxAge int
	text := cleanText(doc.Find(fmt.Sprintf("tr.%s .cmpr-col.col%d span", selector, colCount)).Text())
	fmt.Sscanf(text, "%d %d", &minAge, &maxAge)
	return minAge, maxAge
}

func cleanText(text string) string {
	return strings.Join(strings.Fields(strings.ReplaceAll(text, "\n", "")), " ")
}
