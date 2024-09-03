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

type Product struct {
	Provider              *string                `json:"provider"`
	Product               *string                `json:"product"`
	FeaturesAndConditions *FeaturesAndConditions `json:"features_and_conditions"`
	TopUp                 *TopUp                 `json:"top_up"`
	Fees                  *Fees                  `json:"fees"`
	SpendingFees          *SpendingFees          `json:"spending_fees"`
	CancellationFees      *CancellationFees      `json:"cancellation_fees"`
	AdditionalInfo        *AdditionalInfo        `json:"additional_info"`
}

type FeaturesAndConditions struct {
	HighlightFeatures      *string         `json:"highlight_features"`
	AgeRequirement         *string         `json:"age_requirement"`
	ApplicantQualification []string        `json:"applicant_qualification"`
	UsageConditions        UsageConditions `json:"usage_conditions"`
}

type UsageConditions struct {
	UsageMethod          *string `json:"usage_method"`
	Lifetime             *string `json:"lifetime"`
	PaymentMethod        *string `json:"payment_method"`
	SupportedMerchants   *string `json:"supported_merchants"`
	WebsiteList          *string `json:"website_list"`
	InternationalService *string `json:"international_service"`
}

type TopUp struct {
	TopUpFrequency          *string  `json:"top_up_frequency"`
	FirstTopUpValue         *string  `json:"first_top_up_value"`
	FirstTopUpCondition     *string  `json:"first_top_up_condition"`
	NextTopUpValue          *string  `json:"next_top_up_value"`
	NextTopUpCondition      *string  `json:"next_top_up_condition"`
	MaxBalance              *int     `json:"max_balance"`
	MaxBalanceCondition     *string  `json:"max_balance_condition"`
	FreeTopUpChannels       []string `json:"free_top_up_channels"`
	FeeTopUpChannels        []string `json:"fee_top_up_channels"`
	FeeTopUpChannelsNumeric []int    `json:"fee_top_up_channels_numeric"`
}

type Fees struct {
	InitialFee            *string `json:"initial_fee"`
	InitialFeeNumeric     *int    `json:"initial_fee_numeric"`
	AnnualFee             *string `json:"annual_fee"`
	AnnualFeeNumeric      *int    `json:"annual_fee_numeric"`
	CardReissueFee        *string `json:"card_reissue_fee"`
	CardReissueFeeNumeric *int    `json:"card_reissue_fee_numeric"`
	MaintenanceFee        *string `json:"maintenance_fee"`
	MaintenanceFeeNumeric *int    `json:"maintenance_fee_numeric"`
	OtherFees             *string `json:"other_fees"`
}

type SpendingFees struct {
	SpendingFee                    *int    `json:"spending_fee"`
	SpendingNotificationFee        *string `json:"spending_notification_fee"`
	SpendingNotificationFeeNumeric *int    `json:"spending_notification_fee_numeric"`
	InternationalWithdrawalFee     *string `json:"international_withdrawal_fee"`
	CurrencyConversionFee          *string `json:"currency_conversion_fee"`
}

type CancellationFees struct {
	CashRedemptionFee             *string `json:"cash_redemption_fee"`
	CashRedemptionFeeNumeric      *int    `json:"cash_redemption_fee_numeric"`
	EarlyTerminationFee           *string `json:"early_termination_fee"`
	EarlyTerminationFeePercentage *int    `json:"early_termination_fee_percentage"`
}

type AdditionalInfo struct {
	ProductWebsite *string `json:"product_website"`
	FeeWebsite     *string `json:"fee_website"`
}

func main() {
	url := "https://app.bot.or.th/1213/MCPD/ProductApp/EMoney/CompareProductList"
	payloadTemplate := `{"ProductIdList":"93,10150,10151,10397,10485,10107,10106,10105,10108,10220,10604,10675,10669,10670,10664,10677,10605,10600,10601,29,10603,92,57,54,56,104,10523,10653,10228,10227,10602,85,10474,10392,10294,10292,10658,10644,10239,10219,10584,10260,10152,10153,28,10674,31,10660,10235,10234,10662,30,10398,10484,10486,55,60,10291,10668,10606,53,10607,10632,10609,10619,10643,10681,10290,10296,10637,10616,58","Page":%d,"Limit":3}`

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

	var allProducts []Product
	// Loop through each page
	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages) // Log the pagination value

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

		products := extractProductsFromPage(doc)
		allProducts = append(allProducts, products...)

		// Wait for 2 seconds before making the next request to avoid overloading the server
		time.Sleep(2 * time.Second)
	}

	// Convert the combined products to JSON and save to a file
	jsonData, err := json.MarshalIndent(allProducts, "", "  ")
	if err != nil {
		log.Fatalf("Failed to convert struct to JSON: %v", err)
	}

	// Save JSON to a file
	err = os.WriteFile("e_money.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write JSON to file: %v", err)
	}

	fmt.Println("Product details saved to e_money.json")
}

func setHeaders(req *http.Request) {
	req.Header.Set("Accept", "text/plain, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Cookie", `verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m04nlc8h&sl=0&tt=0"; _uid6672=16B5DEBD.22; _ctout6672=1; _ga_NLQFGWVNXN=GS1.1.1724333049.25.1.1724333066.43.0.0`)
	req.Header.Set("Origin", "https://app.bot.or.th")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://app.bot.or.th/1213/MCPD/ProductApp/EMoney/CompareProduct")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	req.Header.Set("Verificationtoken", "66E3SgnhJHakEQ0_6fAzoiDB6UZo6oYDgI47vsjQYrGkzwZGDWiXyannTm8BODrH4rwr4zPthXTjY4lLYIaaXYQ1_gvstIlt7M8Y13USPS41,8qffyaAOxeFXfU_brP3gEBudpQWNaZABXrAPP5_nwl5816HNwWct2afbKxISRDGiv8WjSxcsNRVAPu-hgUmRS7PI8PervZzgPsrDfo_UVEo1")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
}

func extractProductsFromPage(doc *goquery.Document) []Product {
	var products []Product

	for i := 1; i <= 3; i++ {
		colClass := "col" + strconv.Itoa(i)

		// Extract data for each product and trim spaces
		provider := cleanText(doc.Find(fmt.Sprintf("th.%s span", colClass)).Text())
		product := cleanText(doc.Find(fmt.Sprintf("th.prod-%s span", colClass)).Text())
		highlightFeatures := cleanText(doc.Find(fmt.Sprintf("tr.attr-KeyProduct td.%s span", colClass)).Text())
		ageRequirement := cleanText(doc.Find(fmt.Sprintf("tr.attr-CustomerAge td.%s span", colClass)).Text())
		applicantQualification := doc.Find(fmt.Sprintf("tr.attr-ConditionToApply td.%s span", colClass)).Text()
		usageMethod := cleanText(doc.Find(fmt.Sprintf("tr.attr-UsageCharacteristic td.%s span", colClass)).Text())
		lifetime := cleanText(doc.Find(fmt.Sprintf("tr.attr-UsagePeriod td.%s span", colClass)).Text())
		paymentMethod := cleanText(doc.Find(fmt.Sprintf("tr.attr-Payment td.%s span", colClass)).Text())
		supportedMerchants := cleanText(doc.Find(fmt.Sprintf("tr.attr-ParticipatedShops td.%s span", colClass)).Text())
		websiteList := cleanText(doc.Find(fmt.Sprintf("tr.attr-ParticipatedShopsURL td.%s span", colClass)).Text())
		internationalService := cleanText(doc.Find(fmt.Sprintf("tr.attr-OverseasSpending td.%s span", colClass)).Text())
		topUpFrequency := getStringOrNil(cleanText(doc.Find(fmt.Sprintf("tr.attr-TopUpFrequency td.%s span", colClass)).Text()))
		firstTopUpValue := getStringOrNil(cleanText(doc.Find(fmt.Sprintf("tr.attr-FirstTopUpAmount td.%s span", colClass)).Text()))
		firstTopUpCondition := getStringOrNil(cleanText(doc.Find(fmt.Sprintf("tr.attr-FirstTopUpAmountCondition td.%s span", colClass)).Text()))
		nextTopUpValue := getStringOrNil(cleanText(doc.Find(fmt.Sprintf("tr.attr-FollowingTopUpAmount td.%s span", colClass)).Text()))
		nextTopUpCondition := getStringOrNil(cleanText(doc.Find(fmt.Sprintf("tr.attr-FollowingTopUpAmountCondition td.%s span", colClass)).Text()))
		maxBalanceStr := cleanText(doc.Find(fmt.Sprintf("tr.attr-RemainingBalance td.%s span", colClass)).Text())

		var maxBalance *int
		if maxBalanceStr != "ไม่มีกำหนด" && maxBalanceStr != "" {
			if mb, err := strconv.Atoi(strings.ReplaceAll(maxBalanceStr, ",", "")); err == nil {
				maxBalance = &mb
			}
		} else {
			maxBalance = nil
		}
		maxBalanceCondition := getStringOrNil(cleanText(doc.Find(fmt.Sprintf("tr.attr-RemainingBalanceCondition td.%s span", colClass)).Text()))

		freeTopUpChannels := splitAndTrim(doc.Find(fmt.Sprintf("tr.attr-TopUpChannelsWithoutFee td.%s span", colClass)).Text(), "-")
		feeTopUpChannels := splitAndTrim(doc.Find(fmt.Sprintf("tr.attr-TopUpChannelsWithFee td.%s span", colClass)).Text(), "-")

		// Extract numeric values for fee top-up channels
		var feeTopUpChannelsNumeric []int
		for _, channel := range feeTopUpChannels {
			if fee := extractNumericFee(channel); fee != nil {
				feeTopUpChannelsNumeric = append(feeTopUpChannelsNumeric, *fee)
			}
		}

		// Extract Fees information
		initialFee := getStringOrNil(cleanText(doc.Find(fmt.Sprintf("tr.attr-EntranceFeeAmount td.%s span", colClass)).Text()))
		initialFeeNumeric := extractNumericFee(cleanText(doc.Find(fmt.Sprintf("tr.attr-EntranceFeeAmount td.%s span", colClass)).Text()))
		annualFee := getStringOrNil(cleanText(doc.Find(fmt.Sprintf("tr.attr-AnnualFee td.%s span", colClass)).Text()))
		annualFeeNumeric := extractNumericFee(cleanText(doc.Find(fmt.Sprintf("tr.attr-AnnualFee td.%s span", colClass)).Text()))
		cardReissueFee := getStringOrNil(cleanText(doc.Find(fmt.Sprintf("tr.attr-CardReplacementFee td.%s span", colClass)).Text()))
		cardReissueFeeNumeric := extractNumericFee(cleanText(doc.Find(fmt.Sprintf("tr.attr-CardReplacementFee td.%s span", colClass)).Text()))
		maintenanceFee := getStringOrNil(cleanText(doc.Find(fmt.Sprintf("tr.attr-ProductMaintenanceFee td.%s span", colClass)).Text()))
		maintenanceFeeNumeric := extractNumericFee(cleanText(doc.Find(fmt.Sprintf("tr.attr-ProductMaintenanceFee td.%s span", colClass)).Text()))
		otherFees := getStringOrNil(cleanText(doc.Find(fmt.Sprintf("tr.attr-OtherFees td.%s span", colClass)).Text()))

		// Extract Spending Fees information
		spendingFeeText := cleanText(doc.Find(fmt.Sprintf("tr.attr-SpendingFee td.%s span", colClass)).Text())
		var spendingFee *int
		if spendingFeeText != "ไม่มีค่าธรรมเนียม" && spendingFeeText != "" {
			if sf, err := strconv.Atoi(spendingFeeText); err == nil {
				spendingFee = &sf
			}
		}
		spendingNotificationFee := getStringOrNil(cleanText(doc.Find(fmt.Sprintf("tr.attr-SpendingAlertFee td.%s span", colClass)).Text()))
		spendingNotificationFeeNumeric := extractNumericFee(cleanText(doc.Find(fmt.Sprintf("tr.attr-SpendingAlertFee td.%s span", colClass)).Text()))
		internationalWithdrawalFee := getStringOrNil(cleanText(doc.Find(fmt.Sprintf("tr.attr-OverseasCashWithdrawalFee td.%s span", colClass)).Text()))
		currencyConversionFee := getStringOrNil(cleanText(doc.Find(fmt.Sprintf("tr.attr-CurrencyConversionRiskFeeRate td.%s span", colClass)).Text()))

		// Extract Cancellation Fees information
		cashRedemptionFee := getStringOrNil(cleanText(doc.Find(fmt.Sprintf("tr.attr-CashRefundFee td.%s span", colClass)).Text()))
		cashRedemptionFeeNumeric := extractNumericFee(cleanText(doc.Find(fmt.Sprintf("tr.attr-CashRefundFee td.%s span", colClass)).Text()))
		earlyTerminationFee := getStringOrNil(cleanText(doc.Find(fmt.Sprintf("tr.attr-TerminationFee td.%s span", colClass)).Text()))
		earlyTerminationFeePercentage := extractPercentage(cleanText(doc.Find(fmt.Sprintf("tr.attr-TerminationFee td.%s span", colClass)).Text()))

		// Extract Additional Information
		productWebsite := cleanText(doc.Find(fmt.Sprintf("tr.attr-URL td.%s a.prod-url", colClass)).AttrOr("href", ""))
		feeWebsite := cleanText(doc.Find(fmt.Sprintf("tr.attr-FeeURL td.%s a.prod-url", colClass)).AttrOr("href", ""))

		// Split applicantQualification by "-" and trim spaces
		qualificationArray := splitAndTrim(applicantQualification, "-")

		// Append the product to the products slice
		products = append(products, Product{
			Provider: getStringOrNil(provider),
			Product:  getStringOrNil(product),
			FeaturesAndConditions: &FeaturesAndConditions{
				HighlightFeatures:      getStringOrNil(highlightFeatures),
				AgeRequirement:         getStringOrNil(ageRequirement),
				ApplicantQualification: qualificationArray,
				UsageConditions: UsageConditions{
					UsageMethod:          getStringOrNil(usageMethod),
					Lifetime:             getStringOrNil(lifetime),
					PaymentMethod:        getStringOrNil(paymentMethod),
					SupportedMerchants:   getStringOrNil(supportedMerchants),
					WebsiteList:          getStringOrNil(websiteList),
					InternationalService: getStringOrNil(internationalService),
				},
			},
			TopUp: &TopUp{
				TopUpFrequency:          topUpFrequency,
				FirstTopUpValue:         firstTopUpValue,
				FirstTopUpCondition:     firstTopUpCondition,
				NextTopUpValue:          nextTopUpValue,
				NextTopUpCondition:      nextTopUpCondition,
				MaxBalance:              maxBalance,
				MaxBalanceCondition:     maxBalanceCondition,
				FreeTopUpChannels:       freeTopUpChannels,
				FeeTopUpChannels:        feeTopUpChannels,
				FeeTopUpChannelsNumeric: feeTopUpChannelsNumeric,
			},
			Fees: &Fees{
				InitialFee:            initialFee,
				InitialFeeNumeric:     initialFeeNumeric,
				AnnualFee:             annualFee,
				AnnualFeeNumeric:      annualFeeNumeric,
				CardReissueFee:        cardReissueFee,
				CardReissueFeeNumeric: cardReissueFeeNumeric,
				MaintenanceFee:        maintenanceFee,
				MaintenanceFeeNumeric: maintenanceFeeNumeric,
				OtherFees:             otherFees,
			},
			SpendingFees: &SpendingFees{
				SpendingFee:                    spendingFee,
				SpendingNotificationFee:        spendingNotificationFee,
				SpendingNotificationFeeNumeric: spendingNotificationFeeNumeric,
				InternationalWithdrawalFee:     internationalWithdrawalFee,
				CurrencyConversionFee:          currencyConversionFee,
			},
			CancellationFees: &CancellationFees{
				CashRedemptionFee:             cashRedemptionFee,
				CashRedemptionFeeNumeric:      cashRedemptionFeeNumeric,
				EarlyTerminationFee:           earlyTerminationFee,
				EarlyTerminationFeePercentage: earlyTerminationFeePercentage,
			},
			AdditionalInfo: &AdditionalInfo{
				ProductWebsite: getStringOrNil(productWebsite),
				FeeWebsite:     getStringOrNil(feeWebsite),
			},
		})
	}

	return products
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	var trimmed []string
	for _, p := range parts {
		trimmedText := strings.TrimSpace(p)
		if trimmedText != "" {
			trimmed = append(trimmed, trimmedText)
		}
	}
	return trimmed
}

func cleanText(s string) string {
	return strings.TrimSpace(strings.ReplaceAll(s, "\n", ""))
}

func getStringOrNil(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}

func extractNumericFee(text string) *int {
	// Extract numeric value from text, assuming the text has the format like "ค่าธรรมเนียมเริ่มต้น 5 บาท"
	parts := strings.Fields(text)
	for _, part := range parts {
		if num, err := strconv.Atoi(strings.ReplaceAll(part, ",", "")); err == nil {
			return &num
		}
	}
	return nil
}

func extractPercentage(text string) *int {
	// Extract percentage from text, assuming the text contains a format like "4% ของยอดเงินที่เอาออกจากกระเป๋า"
	if strings.Contains(text, "%") {
		parts := strings.Fields(text)
		for _, part := range parts {
			if strings.Contains(part, "%") {
				part = strings.TrimSuffix(part, "%")
				if num, err := strconv.Atoi(part); err == nil {
					return &num
				}
			}
		}
	}
	return nil
}
