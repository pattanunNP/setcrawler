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

type ProductInfo struct {
	Provider                 string                   `json:"provider"`
	Product                  string                   `json:"product"`
	Interest                 Interest                 `json:"interest"`
	AccountOpeningConditions AccountOpeningConditions `json:"account_opening_conditions"`
	ProductUsageConditions   ProductUsageConditions   `json:"product_usage_conditions"`
	Insurance                Insurance                `json:"insurance"`
	ProductFees              ProductFees              `json:"product_fees"`
	GeneralFee               GeneralFee               `json:"general_fee"`
	AdditionInfo             AdditionInfo             `json:"additional_info"`
}

type GeneralFee struct {
	CoinCountingFee                string   `json:"coin_counting_fee"`
	CrossBankDepositWithdrawalFee  []string `json:"cross_bank_deposit_withdrawal_fee"`
	OtherProviderDepositFeeCDMATM  []string `json:"other_provider_deposit_fee_cdm_atm"`
	SameProviderDepositFeeCDMATM   []string `json:"same_provider_deposit_fee_cdm_atm"`
	DepositWithdrawalAgentFee      []string `json:"deposit_withdrawal_agent_fee"`
	AutoTransferFeeSavingsChecking []string `json:"auto_transfer_fee_savings_checking"`
	CrossBankTransferFee           []string `json:"cross_bank_transfer_fee"`
	OtherFees                      string   `json:"other_fees"`
}

type AdditionInfo struct {
	ProductWebsiteLink string `json:"product_website_link"`
	FeeWebsiteLink     string `json:"fee_website_link"`
}

type ProductFees struct {
	AccountMaintenanceFee  string `json:"account_maintenance_fee"`
	SMSNotificationFee     string `json:"sms_notification_fee"`
	PassbookReplacementFee string `json:"passbook_replacement_fee"`
	TransactionHistoryFee  string `json:"transaction_history_fee"`
	AccountClosureFee      string `json:"account_closure_fee"`
}

type ProductUsageConditions struct {
	MinimumDepositPerTransaction      int      `json:"minimum_deposit_per_transaction"`
	AdditionalDepositsAllowed         string   `json:"additional_deposits_allowed"`
	PartialWithdrawalsAllowed         string   `json:"partial_withdrawals_allowed"`
	DepositWithdrawTransferConditions []string `json:"deposit_withdraw_transfer_conditions"`
	AccountRenewalWhenDue             string   `json:"account_renewal_when_due"`
}

type Interest struct {
	MinRate                   float64 `json:"min_rate"`
	MaxRate                   float64 `json:"max_rate"`
	InterestCalculationMethod string  `json:"interest_calculation_method"`
	TaxFree                   string  `json:"tax_free"`
	InterestPaymentPeriod     string  `json:"interest_payment_period"`
	MinAmount                 int     `json:"min_amount,omitempty"`
	MaxAmount                 int     `json:"max_amount,omitempty"`
	InterestPaymentMethod     string  `json:"interest_payment_method"`
	InterestPenaltyConditions string  `json:"interest_penalty_conditions"`
}

type AccountOpeningConditions struct {
	FixedTerm                 string   `json:"fixed_term"`
	MinimumOpeningBalance     int      `json:"minimum_opening_balance"`
	MaximumDepositLimit       int      `json:"maximum_deposit_limit"`
	SpecificOpeningConditions []string `json:"specific_opening_conditions"`
}

type Insurance struct {
	Insurance              string `json:"insurance"`
	InsuranceCompany       string `json:"insurance_company"`
	InsuranceCoverageLimit string `json:"insurance_coverage_limit"`
	InsuranceConditions    string `json:"insurance_conditions"`
}

func main() {
	url := "https://app.bot.or.th/1213/MCPD/ProductApp/Deposit/CompareProductList"
	payloadTemplate := `{"ProductIdList":"64434,63687,62863,63513,62861,62862,63726,63904,63935,63888,62859,63723,63881,63659,63523,62860,63725,63642,63644,63530,64306,63864,64002,63879,63908,64021,64158,64159,63722,64331,63300,62251,64066,64001,63450,64157,63294,64188,64177,64189,62850,64020,62507,64124,64095,63724,63997,64433,63995,64176,63392,62505,62511,63385,64004,64016,64171,64163,63366,63359,64155,62854,63443,63878,64003,64014,59174,64161,64169,63346,63340,64153,62852,64230,64087,64415,62517,64499,61132,61106,63880,63707,63635,59180,62856,63861,63448,61113,64316,64318,63617,63606,61123,61112,61147,61131,61117,61135,63909,63982,63892,63981,63925,64484,64505,63863,63271,63875,64498,61150,61095,61138,61139,61109,61126,61096,61144,61098,61128,61119,61103,62516,64487,64445,63865,59183,63860,63337,63417,63446,61108,64376,63903,64192,63613,61115,61146,61130,61116,61104,61120,63899,63922,62247,61136,61125,61110,61145,61149,61133,63887,63919,62234,63866,63643,64421,64436,64125,63585,63900,64019,63704,64191,63545,62201,62212,61107,61142,61097,61127,61151,61105,63898,62184,62227,63854,63602,62502,63962,63907,63378,63874,63658,62209,62188,62195,62510,62258,64167,63877,64018,63867,64256,59185,63327,63322,59182,62515,63936,63905,63937,64000,64497,64478,63929,63288,64037,64503,64395,64429,64459,52356,64412,64391,62199,62192,62241,62239,62259,63528,62508,62257,63876,63374,63435,64058,64193,64377,63862,63994,62858,63889,63487,63615,62513,63926,63541,64447,64496,64475,64406,64071,64050,64116,64042,64130,52354,62506,64467,63706,59175,64168,64160,63445,64022,64013,64267,64152,63902,59181,64100,64079,64072,64028,52358,62853,61137,61124,61141,62265,62193,63916,63278,64178,64179,63680,63656,22398,63356,63426,64166,64427,64451,64457,64444,64392,64476,64512,64411,64190,64182,22220,63917,63976,63918,63953,63221,64092,63688,63891,64129,64101,64123,62260,61100,64165,63896,63910,62242,64504,64404,64501,64393,22218,63518,63695,63570,63882,63491,64080,64144,62847,64174,63987,64009,62512,62504,62266,63855,62204,64328,63984,62198,63227,62501,64492,64474,64510,59184,64456,64464,64494,62851,59171,61148,61118,61102,63869,63998,64049,64094,64137,64086,64108,64057,63488,63490,52355,63449,22217,63489,64073,63951,62205,64164,64172,64024,64017,64353,64181,64183,22221,64156,62857,62213,62194,52351,64051,64007,63988,63447,63996,64180,63890,64194,63897,22391,22386,22388,22393,22394,62228,62226,62514,63930,63978,63967,64162,64170,64154,62190,63314,22219,63895,64023,64015,64308,64048,64109,64093,64136,64064,64107,64056,64078,64027,64070,64115,64041,64035,62243,62210,62187,62236,62855,64242,64128,64337,63964,63966,62264,63893,59173,64361,62237,64151,63992,62509,62503,64065,64324,63954,63920,63975,64286,63652,63709,64479,64113,64084,64121,63486,59172,64026,64138,64114,64034,64085,64122,64010,62202,61134,63650,64098,64032,64076,64047,63991,61099,63287,64455,64430,63230,63317,64011,63868,64413,64441,64446,64449,63217,63857,62208,62222,63297,63405,64356,64382,64384,64354,22216,22390,52357,52353,61140,64143,64006,64069,62252,63454,63459,59170,62186,63485,64502,64394,64490,64493,63484,64060,64097,64068,64105,64127,64142,62235,52352,63853,64103,63927,63440,63441,63989,64126,64118,64089,64054,64046,64134,64184,64173,64185,64040,64186,63438,62231,64221,62843,63708,63492,22222,62500,64091,64062,64033,64025,64096,64067,64038,64063,64055,22392,64141,23264,64112,64083,23259,22384,23263,62225,61122,63610,63646,22385,22389,64045,64133,64104,64075,64031,64119,64090,64061,23270,23273,22381,63870,64139,63444,63856,63442,64111,64082,64053,64039,23269,23276,60673,23260,64106,23274,23275,63977,63974,62842,64147,62845,63990,63437,61129,63985,64005,64150,64246,63439,64145,64365,64146,63986,61121,64148,61101,23262,23266,63993,63455,64175,64088,63859,64489,63698,22382,22380,64368,64509,59179,23261,64511,64036,64491,64488,64463,64402,64428,64485,23268,63452,63458,62254,37623,37616,64008,64149,23272,23267,63928,37615,37602,22387,64077,64044,64132,64074,64110,64081,64052,23271,64030,64059,63894,63885,64244,64349,64131,64371,64277,64120,59167,59177,59176,59168,61114,63858,63873,62848,63554,63883,63884,63886,59166,22383,59178,59169,64187,63872,64448,64290,63983,37605,23258,64135,63483,63231,64012,23265,64117,64140,22400,22397,22395,22399,22396","Page":%d,"Limit":3}`
	// payloadTemplate := `{"ProductIdList":"64434,63687,62863,63513,62861,62862,63726,63904,63935,63888,62859,63723,63881,63659,63523,62860,63725,63642,63644,63530,64306,63864,64002,63879,63908,64021,64158,64159,63722,64331,63300,62251,64066,64001,63450","Page":%d,"Limit":3}`

	initialPage := 1
	payload := fmt.Sprintf(payloadTemplate, initialPage)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	totalPages := 5

	var allProducts []ProductInfo

	// Loop through each page
	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages)
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
			fmt.Println("Error reading response:", err)
			return
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			fmt.Println("Error loading HTML:", err)
			return
		}

		doc.Find("th.cmpr-col").Each(func(i int, s *goquery.Selection) {
			// Get provider and product
			provider := cleanedText(getTextSafely(s.Find("th.col-s span").Last()))
			product := cleanedText(getTextSafely(s.Find("th.font-black.text-center")))

			if provider == "" || product == "" {
				fmt.Printf("Skipping column %d due to empty provider or product\n", i+1)
				return
			}
			// Get interest rates
			interestText := cleanedText(getTextSafely(doc.Find(fmt.Sprintf("tbody td.cmpr-col.col%d.text-center span.text-bold", i+1))))
			minRate, maxRate := parseInterestRate(interestText)

			interest := Interest{
				MinRate:                   minRate,
				MaxRate:                   maxRate,
				InterestCalculationMethod: getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-intrmthd.trbox-shadow td.cmpr-col.col%d span", i+1))),
				TaxFree:                   getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-intrwotax.trbox-shadow td.cmpr-col.col%d span", i+1))),
				InterestPaymentPeriod:     getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-intrterm.trbox-shadow td.cmpr-col.col%d span", i+1))),
				InterestPaymentMethod:     getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-intrch.trbox-shadow td.cmpr-col.col%d span", i+1))),
				InterestPenaltyConditions: cleanedText(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-intrexc.trbox-shadow td.cmpr-col.col%d", i+1)))),
			}

			// Parse amounts from the interest payment period text
			interest.MinAmount, interest.MaxAmount = parseInterestPaymentPeriodAmounts(interest.InterestPaymentPeriod)

			// Account Opening Conditions
			accountOpeningConditions := AccountOpeningConditions{
				FixedTerm:                 getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-dpstterm.trbox-shadow td.cmpr-col.col%d span", i+1))),
				MinimumOpeningBalance:     parseInt(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-blncmin.trbox-shadow td.cmpr-col.col%d span", i+1)))),
				MaximumDepositLimit:       parseInt(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-blncmax.trbox-shadow td.cmpr-col.col%d span", i+1)))),
				SpecificOpeningConditions: parseConditions(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-opencond.trbox-shadow td.cmpr-col.col%d", i+1)))),
			}

			// Product Usage Conditions
			productUsageConditions := ProductUsageConditions{
				MinimumDepositPerTransaction:      parseInt(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-mindpst.trbox-shadow td.cmpr-col.col%d span", i+1)))),
				AdditionalDepositsAllowed:         getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-topup.trbox-shadow td.cmpr-col.col%d span", i+1))),
				PartialWithdrawalsAllowed:         getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-wdprtprnc.trbox-shadow td.cmpr-col.col%d span", i+1))),
				DepositWithdrawTransferConditions: parseConditions(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-wdprtxnblnc.trbox-shadow td.cmpr-col.col%d span", i+1)))),
				AccountRenewalWhenDue:             getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-accrenew.trbox-shadow td.cmpr-col.col%d span", i+1))),
			}

			// Insurance Information
			insurance := Insurance{
				Insurance:              cleanedText(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-insrnc.trbox-shadow td.cmpr-col.col%d span", i+1)))),
				InsuranceCompany:       cleanedText(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-insrnccompany.trbox-shadow td.cmpr-col.col%d span", i+1)))),
				InsuranceCoverageLimit: cleanedText(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-insrnclimit.trbox-shadow td.cmpr-col.col%d span", i+1)))),
				InsuranceConditions:    cleanedText(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-insrnccond.trbox-shadow td.cmpr-col.col%d span", i+1)))),
			}

			// Product Fees
			productFees := ProductFees{
				AccountMaintenanceFee:  cleanedText(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-accmtnc.trbox-shadow td.cmpr-col.col%d span", i+1)))),
				SMSNotificationFee:     cleanedText(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-accmsms.trbox-shadow td.cmpr-col.col%d span", i+1)))),
				PassbookReplacementFee: cleanedText(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-accmopenbk.trbox-shadow td.cmpr-col.col%d span", i+1)))),
				TransactionHistoryFee:  cleanedText(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-accmbranch.trbox-shadow td.cmpr-col.col%d span", i+1)))),
				AccountClosureFee:      cleanedText(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-accmclosebk.trbox-shadow td.cmpr-col.col%d span", i+1)))),
			}

			// General Fee
			generalFee := GeneralFee{
				CoinCountingFee:                cleanedText(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-feecorn.trbox-shadow td.cmpr-col.col%d span", i+1)))),
				CrossBankDepositWithdrawalFee:  parseFeeConditions(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-feedeposit.trbox-shadow td.cmpr-col.col%d span", i+1)))),
				OtherProviderDepositFeeCDMATM:  parseFeeConditions(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-feecdm.trbox-shadow td.cmpr-col.col%d span", i+1)))),
				SameProviderDepositFeeCDMATM:   parseFeeConditions(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-feecdm2.trbox-shadow td.cmpr-col.col%d span", i+1)))),
				DepositWithdrawalAgentFee:      parseFeeConditions(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-feeother2.trbox-shadow td.cmpr-col.col%d span", i+1)))),
				AutoTransferFeeSavingsChecking: parseFeeConditions(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-feetranfer.trbox-shadow td.cmpr-col.col%d span", i+1)))),
				CrossBankTransferFee:           parseFeeConditions(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-feetranfer2.trbox-shadow td.cmpr-col.col%d span", i+1)))),
				OtherFees:                      cleanedText(getTextSafely(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-feeother1.trbox-shadow td.cmpr-col.col%d span", i+1)))),
			}

			// Additional Info (handle href attribute extraction separately)
			productWebsiteLink := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-url.trbox-shadow td.cmpr-col.col%d a.prod-url", i+1)).AttrOr("href", ""))
			feeWebsiteLink := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-feeurl.trbox-shadow td.cmpr-col.col%d a.prod-url", i+1)).AttrOr("href", ""))

			additionalInfo := AdditionInfo{
				ProductWebsiteLink: productWebsiteLink,
				FeeWebsiteLink:     feeWebsiteLink,
			}

			// Create ProductInfo struct
			productInfo := ProductInfo{
				Provider:                 provider,
				Product:                  product,
				Interest:                 interest,
				AccountOpeningConditions: accountOpeningConditions,
				ProductUsageConditions:   productUsageConditions,
				Insurance:                insurance,
				ProductFees:              productFees,
				GeneralFee:               generalFee,
				AdditionInfo:             additionalInfo,
			}

			allProducts = append(allProducts, productInfo)
		})

		time.Sleep(2 * time.Second)
	}

	// Convert the combined products to JSON and save to a file
	jsonData, err := json.MarshalIndent(allProducts, "", "  ")
	if err != nil {
		fmt.Println("Failed to convert struct to JSON:", err)
		return
	}

	// Save JSON to a file
	err = os.WriteFile("deposit.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Failed to write JSON to file:", err)
		return
	}

	fmt.Println("Product details saved to deposit.json")
}

func parseInterestPaymentPeriodAmounts(text string) (int, int) {
	var minAmount, maxAmount int

	// Regular expressions to find numeric values
	amountRegex := regexp.MustCompile(`\d{1,3}(?:,\d{3})*`)
	amounts := amountRegex.FindAllString(text, -1)

	if len(amounts) == 0 {
		return 0, 0
	}

	// Parse the first amount found as minAmount
	if len(amounts) >= 1 {
		minAmountStr := strings.ReplaceAll(amounts[0], ",", "")
		minAmount, _ = strconv.Atoi(minAmountStr)
	}

	// Parse the second amount found as maxAmount (if exists)
	if len(amounts) >= 2 {
		maxAmountStr := strings.ReplaceAll(amounts[1], ",", "")
		maxAmount, _ = strconv.Atoi(maxAmountStr)
	}

	return minAmount, maxAmount
}

func parseFeeConditions(text string) []string {
	// Split by the "-" character
	parts := strings.Split(text, "-")

	var conditions []string
	for _, part := range parts {
		// Clean and trim each part
		cleanedPart := cleanedText(part)
		if cleanedPart != "" {
			conditions = append(conditions, cleanedPart)
		}
	}
	return conditions
}

func getTextSafely(s *goquery.Selection) string {
	if s.Length() == 0 {
		return ""
	}
	return strings.TrimSpace(s.Text())
}

func cleanedText(input string) string {
	// Replace \u003c with < and \u003e with >
	cleaned := strings.ReplaceAll(input, `\u003c`, "<")
	cleaned = strings.ReplaceAll(cleaned, `\u003e`, ">")

	// Remove any other unexpected escape sequences
	escapeRegex := regexp.MustCompile(`\\u[0-9a-fA-F]{4}`)
	cleaned = escapeRegex.ReplaceAllStringFunc(cleaned, func(r string) string {
		switch r {
		case `\u003c`:
			return "<"
		case `\u003e`:
			return ">"
		default:
			return ""
		}
	})

	// Remove newline characters
	noNewLines := strings.ReplaceAll(cleaned, "\n", "")

	// Replace multiple spaces with a single space
	spaceRegex := regexp.MustCompile(`\s+`)
	singleSpaced := spaceRegex.ReplaceAllString(noNewLines, " ")

	return strings.TrimSpace(singleSpaced)
}

func parseInterestRate(rateText string) (float64, float64) {
	rateText = cleanedText(rateText)

	if strings.Contains(rateText, "-") {
		rates := strings.Split(rateText, "-")
		minRateStr := strings.TrimSpace(strings.Trim(strings.TrimSpace(rates[0]), "%"))
		maxRateStr := strings.TrimSpace(strings.Trim(strings.TrimSpace(rates[1]), "%"))

		minRate, errMin := strconv.ParseFloat(minRateStr, 64)
		maxRate, errMax := strconv.ParseFloat(maxRateStr, 64)

		if errMin != nil || errMax != nil {
			fmt.Println("Error parsing rates:", errMin, errMax)
			return 0, 0
		}

		return minRate, maxRate
	}

	rateStr := strings.TrimSpace(strings.Trim(rateText, "%"))
	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		fmt.Println("Error parsing rate:", err)
		return 0, 0
	}
	return rate, rate
}

func parseConditions(text string) []string {
	// Split by the "-" character
	parts := strings.Split(text, "-")

	var conditions []string
	for _, part := range parts {
		// Clean and trim each part
		cleanedPart := cleanedText(part)
		if cleanedPart != "" {
			conditions = append(conditions, cleanedPart)
		}
	}
	return conditions
}

// func parseConditionalRateDetails(text string) []ConditionalRateDetail {
// 	var details []ConditionalRateDetail
// 	lines := strings.Split(text, "\n")
// 	for _, line := range lines {
// 		fields := strings.Fields(line)
// 		if len(fields) >= 5 {
// 			rate, _ := strconv.ParseFloat(strings.Trim(fields[0], "%"), 64)
// 			amountMin, _ := strconv.Atoi(fields[1])
// 			amountMax, _ := strconv.Atoi(fields[3])
// 			details = append(details, ConditionalRateDetail{
// 				Range:         fmt.Sprintf("%d - %d บาท", amountMin, amountMax),
// 				Rate:          fmt.Sprintf("%.2f%%", rate),
// 				RateFloat:     rate,
// 				PaymentPeriod: fields[4],
// 				AmountMin:     amountMin,
// 				AmountMax:     amountMax,
// 			})
// 		}
// 	}
// 	return details
// }

// func parseSpecificOpeningConditions(text string) []SpecificOpeningCondition {
// 	var conditions []SpecificOpeningCondition
// 	lines := strings.Split(text, "\n")
// 	for _, line := range lines {
// 		parts := strings.Split(line, " ")
// 		if len(parts) >= 3 {
// 			amountMin, _ := strconv.Atoi(parts[0])
// 			amountMax, _ := strconv.Atoi(parts[2])
// 			condition := strings.Join(parts[3:], " ")
// 			conditions = append(conditions, SpecificOpeningCondition{
// 				Condition: condition,
// 				AmountMin: amountMin,
// 				AmountMax: amountMax,
// 			})
// 		}
// 	}
// 	return conditions
// }

// func parseDepositWithdrawTransferConditions(text string) []DepositWithdrawTransferCondition {
// 	var conditions []DepositWithdrawTransferCondition
// 	lines := strings.Split(text, "\n")
// 	for _, line := range lines {
// 		condition := strings.TrimSpace(line)
// 		if condition != "" {
// 			conditions = append(conditions, DepositWithdrawTransferCondition{Condition: condition})
// 		}
// 	}
// 	return conditions
// }

func parseInt(text string) int {
	// Return 0 for "ไม่กำหนด" (meaning "not specified") and similar cases
	if text == "" || strings.Contains(text, "ไม่กำหนด") {
		return 0
	}

	// Remove any non-digit characters (e.g., " บาท", " /เดือน", etc.)
	cleanedText := regexp.MustCompile(`\D`).ReplaceAllString(text, "")

	// Parse the cleaned string into an integer
	num, err := strconv.Atoi(cleanedText)
	if err != nil {
		fmt.Printf("Error parsing integer: %v, input: %s\n", err, text)
		return 0
	}
	return num
}

// func parseAgeRange(ageText string) (*int, *int) {
// 	ages := strings.Split(ageText, "-")
// 	if len(ages) != 2 {
// 		return nil, nil
// 	}

// 	minAgeStr := strings.TrimSpace(ages[0])
// 	maxAgeStr := strings.TrimSpace(ages[1])

// 	minAge, errMin := strconv.Atoi(minAgeStr)
// 	maxAge, errMax := strconv.Atoi(maxAgeStr)

// 	if errMin != nil || errMax != nil {
// 		return nil, nil
// 	}

// 	return &minAge, &maxAge
// }

// func parseSectionedText(text string) []string {
// 	sectionPattern := regexp.MustCompile(`(\d+\.)|(-)`)
// 	matches := sectionPattern.FindAllStringIndex(text, -1)

// 	var result []string
// 	startIndex := 0

// 	for _, match := range matches {
// 		if match[0] > startIndex {
// 			section := strings.TrimSpace(text[startIndex:match[0]])
// 			if section != "" {
// 				result = append(result, section)
// 			}
// 		}
// 		startIndex = match[0]
// 	}

// 	// Append the last section if any
// 	if startIndex < len(text) {
// 		section := strings.TrimSpace(text[startIndex:])
// 		if section != "" {
// 			result = append(result, section)
// 		}
// 	}

// 	return result
// }

// func parseSlashSeparatedText(text string) []string {
// 	// Remove newline characters and unnecessary spaces
// 	cleanedText := strings.ReplaceAll(text, "\n", "")
// 	cleanedText = strings.Join(strings.Fields(cleanedText), " ")

// 	sections := strings.Split(cleanedText, "/")
// 	var result []string

// 	for _, section := range sections {
// 		cleanSection := strings.TrimSpace(section)
// 		if cleanSection != "" {
// 			result = append(result, cleanSection)
// 		}
// 	}
// 	return result
// }

// func parseTranSactionHistoryFee(text string) []string {
// 	sections := []string{}
// 	parts := strings.Split(text, "เงื่อนไข:")
// 	if len(parts) > 0 {
// 		mainPart := strings.TrimSpace(parts[0])
// 		sections = append(sections, strings.Split(mainPart, "ขอใบแสดงรายการย้อนหลัง")...)
// 	}

// 	if len(parts) > 1 {
// 		conditionPart := strings.TrimSpace(parts[1])
// 		conditionSections := strings.Split(conditionPart, "-")
// 		for _, section := range conditionSections {
// 			cleanedSection := strings.TrimSpace(section)
// 			if cleanedSection != "" {
// 				sections = append(sections, cleanedSection)
// 			}
// 		}
// 	}

// 	for i, sec := range sections {
// 		sections[i] = strings.TrimSpace(sec)
// 		if strings.HasPrefix(sec, "น้อยกว่า") || strings.HasPrefix(sec, "ตั้งแต่") || strings.HasPrefix(sec, "มากกว่า") {
// 			sections[i] = "- " + sec
// 		} else if sec != "" && !strings.HasPrefix(sec, "-") {
// 			sections[i] = "ขอใบแสดงรานการย้อนหลัง" + sec
// 		}
// 	}
// 	return sections
// }

// func parseAutoTransferFee(text string) []string {
// 	cleanedText := strings.ReplaceAll(text, "\n", "")
// 	cleanedText = strings.Join(strings.Fields(cleanedText), " ")

// 	sections := []string{}
// 	parts := strings.Split(cleanedText, "เงื่อนไข")
// 	if len(parts) > 0 {
// 		mainPart := strings.TrimSpace(parts[0])
// 		if mainPart != "" && mainPart != "-" {
// 			sections = append(sections, mainPart)
// 		}
// 	}

// 	if len(parts) > 1 {
// 		conditionsPart := strings.TrimSpace(parts[1])
// 		conditionSections := strings.Split(conditionsPart, "-")
// 		for _, section := range conditionSections {
// 			cleanSection := strings.TrimSpace(section)
// 			if cleanSection != "" {
// 				sections = append(sections, cleanSection)
// 			}
// 		}
// 	}

// 	// Remove colon if present
// 	for i, sec := range sections {
// 		sections[i] = strings.ReplaceAll(sec, ":", "")
// 	}

// 	return sections
// }

// func parseConditionalRate(text string) []string {
// 	text = cleanedText(text)
// 	return strings.Split(text, "เงินฝากส่วนที่เกิน")
// }

func setHeaders(req *http.Request) {
	req.Header.Set("Accept", "text/plain, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Cookie", "verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; _uid6672=16B5DEBD.63; _ctout6672=1; RT=\"z=1&dm=app.bot.or.th&si=0d61fb4b-0525-401c-af19-c7a1013eb434&ss=m0lwtfky&sl=5&tt=30c&obo=4&rl=1\"; _ga_NLQFGWVNXN=GS1.1.1725336521.68.1.1725337171.28.0.0")
	req.Header.Set("Origin", "https://app.bot.or.th")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://app.bot.or.th/1213/MCPD/ProductApp/Deposit/CompareProduct")
	req.Header.Set("Sec-CH-UA", `"Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome";v="128"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	req.Header.Set("VerificationToken", "3Bwwcx-u3Mgqt-1mnda_kJ6shdGICCn8Qffjx3G1MzU-fu6OrzCUUIsyV2Rre-vLMfyWKNa4fZCZrRbiWgWJvOSLZUPHNHcDS3T1j_LKc441,ITUVxn0RLt2U8y_QKPeU-dpTKHXO3d4CAG2GvEmG9bfJcVDnh2sH2gVhReNysIQldG_BLu6khulVkGlpPCGyKmGzhNTg6DV6754QWDtlKtk1")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
}
