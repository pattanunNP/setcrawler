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
	AdditionInfo             AdditionInfo             `json:"addition_info"`
}

type GeneralFee struct {
	CoinCountingFee                string   `json:"coin_counting_fee"`
	CrossBankDepositWithdrawalFee  string   `json:"cross_bank_deposit_withdrawal_fee"`
	OtherProviderDepositFeeCDMATM  []string `json:"other_provider_deposit_fee_cdm_atm"`
	SameProviderDepositFeeCDMATM   string   `json:"same_provider_deposit_fee_cdm_atm"`
	DepositWithdrawalAgentFee      string   `json:"deposit_withdrawal_agent_fee"`
	AutoTransferFeeSavingsChecking []string `json:"auto_transfer_fee_savings_checking"`
	CrossBankTransferFee           []string `json:"cross_bank_transfer_fee"`
	OtherFees                      string   `json:"other_fees"`
}

type AdditionInfo struct {
	ProductWebsiteLink string `json:"product_website_link"`
	FeeWebsiteLink     string `json:"fee_website_link"`
}

type ProductFees struct {
	AccountMaintenanceFee  string   `json:"account_maintenance_fee"`
	SMSNotificationFee     string   `json:"sms_notification_fee"`
	PassbookReplacementFee *int     `json:"passbook_replacement_fee"`
	TransactionHistoryFee  []string `json:"transaction_history_fee"`
	AccountClosureFee      string   `json:"account_closure_fee"`
}

type ProductUsageConditions struct {
	MinimumDepositPerTransaction      *int     `json:"minimum_deposit_per_transaction"`
	AdditionalDepositsAllowed         string   `json:"additional_deposits_allowed"`
	PartialWithdrawalsAllowed         string   `json:"partial_withdrawals_allowed"`
	DepositWithdrawTransferConditions []string `json:"deposit_withdraw_transfer_conditions"`
	AccountRenewalWhenDue             string   `json:"account_renewal_when_due"`
}

type Interest struct {
	MinRate                   float64  `json:"min_rate"`
	MaxRate                   float64  `json:"max_rate"`
	ConditionalRate           []string `json:"conditional_rate"`
	InterestCalculationMethod string   `json:"interest_calculation_method"`
	TaxFree                   string   `json:"tax_free"`
	InterestPaymentPeriod     string   `json:"interest_payment_period"`
	InterestPaymentMethod     string   `json:"interest_payment_method"`
	InterestPenaltyConditions string   `json:"interest_penalty_conditions"`
}

type AccountOpeningConditions struct {
	FixedTerm                 string   `json:"fixed_term"`
	MinimumOpeningBalance     *int     `json:"minimum_opening_balance,omitempty"`
	MaximumDepositLimit       *int     `json:"maximum_deposit_limit,omitempty"`
	OtherProductRequirements  string   `json:"other_product_requirements"`
	MinAge                    *int     `json:"min_age,omitempty"`
	MaxAge                    *int     `json:"max_age,omitempty"`
	SpecificOpeningConditions []string `json:"specific_opening_conditions"`
}

type Insurance struct {
	Insurance              string `json:"insurance"`
	InsuranceCompany       string `json:"insurance_company"`
	InsuranceCoverageLimit *int   `json:"insurance_coverage_limit"`
	InsuranceConditions    string `json:"insurance_conditions"`
}

func main() {
	url := "https://app.bot.or.th/1213/MCPD/ProductApp/Deposit/CompareProductList"
	payloadTemplate := `{"ProductIdList":"64434,63687,62863,63513,62861,62862,63726,63904,63935,63888,62859,63723,63881,63659,63523,62860,63725,63642,63644,63530,64306,63864,64002,63879,63908,64021,64158,64159,63722,64331,63300,62251,64066,64001,63450,64157,63294,64188,64177,64189,62850,64020,62507,64124,64095,63724,63997,64433,63995,64176,63392,62505,62511,63385,64004,64016,64171,64163,63366,63359,64155,62854,63443,63878,64003,64014,59174,64161,64169,63346,63340,64153,62852,64230,64087,64415,62517,64499,61132,61106,63880,63707,63635,59180,62856,63861,63448,61113,64316,64318,63617,63606,61123,61112,61147,61131,61117,61135,63909,63982,63892,63981,63925,64484,64505,63863,63271,63875,64498,61150,61095,61138,61139,61109,61126,61096,61144,61098,61128,61119,61103,62516,64487,64445,63865,59183,63860,63337,63417,63446,61108,64376,63903,64192,63613,61115,61146,61130,61116,61104,61120,63899,63922,62247,61136,61125,61110,61145,61149,61133,63887,63919,62234,63866,63643,64421,64436,64125,63585,63900,64019,63704,64191,63545,62201,62212,61107,61142,61097,61127,61151,61105,63898,62184,62227,63854,63602,62502,63962,63907,63378,63874,63658,62209,62188,62195,62510,62258,64167,63877,64018,63867,64256,59185,63327,63322,59182,62515,63936,63905,63937,64000,64497,64478,63929,63288,64037,64503,64395,64429,64459,52356,64412,64391,62199,62192,62241,62239,62259,63528,62508,62257,63876,63374,63435,64058,64193,64377,63862,63994,62858,63889,63487,63615,62513,63926,63541,64447,64496,64475,64406,64071,64050,64116,64042,64130,52354,62506,64467,63706,59175,64168,64160,63445,64022,64013,64267,64152,63902,59181,64100,64079,64072,64028,52358,62853,61137,61124,61141,62265,62193,63916,63278,64178,64179,63680,63656,22398,63356,63426,64166,64427,64451,64457,64444,64392,64476,64512,64411,64190,64182,22220,63917,63976,63918,63953,63221,64092,63688,63891,64129,64101,64123,62260,61100,64165,63896,63910,62242,64504,64404,64501,64393,22218,63518,63695,63570,63882,63491,64080,64144,62847,64174,63987,64009,62512,62504,62266,63855,62204,64328,63984,62198,63227,62501,64492,64474,64510,59184,64456,64464,64494,62851,59171,61148,61118,61102,63869,63998,64049,64094,64137,64086,64108,64057,63488,63490,52355,63449,22217,63489,64073,63951,62205,64164,64172,64024,64017,64353,64181,64183,22221,64156,62857,62213,62194,52351,64051,64007,63988,63447,63996,64180,63890,64194,63897,22391,22386,22388,22393,22394,62228,62226,62514,63930,63978,63967,64162,64170,64154,62190,63314,22219,63895,64023,64015,64308,64048,64109,64093,64136,64064,64107,64056,64078,64027,64070,64115,64041,64035,62243,62210,62187,62236,62855,64242,64128,64337,63964,63966,62264,63893,59173,64361,62237,64151,63992,62509,62503,64065,64324,63954,63920,63975,64286,63652,63709,64479,64113,64084,64121,63486,59172,64026,64138,64114,64034,64085,64122,64010,62202,61134,63650,64098,64032,64076,64047,63991,61099,63287,64455,64430,63230,63317,64011,63868,64413,64441,64446,64449,63217,63857,62208,62222,63297,63405,64356,64382,64384,64354,22216,22390,52357,52353,61140,64143,64006,64069,62252,63454,63459,59170,62186,63485,64502,64394,64490,64493,63484,64060,64097,64068,64105,64127,64142,62235,52352,63853,64103,63927,63440,63441,63989,64126,64118,64089,64054,64046,64134,64184,64173,64185,64040,64186,63438,62231,64221,62843,63708,63492,22222,62500,64091,64062,64033,64025,64096,64067,64038,64063,64055,22392,64141,23264,64112,64083,23259,22384,23263,62225,61122,63610,63646,22385,22389,64045,64133,64104,64075,64031,64119,64090,64061,23270,23273,22381,63870,64139,63444,63856,63442,64111,64082,64053,64039,23269,23276,60673,23260,64106,23274,23275,63977,63974,62842,64147,62845,63990,63437,61129,63985,64005,64150,64246,63439,64145,64365,64146,63986,61121,64148,61101,23262,23266,63993,63455,64175,64088,63859,64489,63698,22382,22380,64368,64509,59179,23261,64511,64036,64491,64488,64463,64402,64428,64485,23268,63452,63458,62254,37623,37616,64008,64149,23272,23267,63928,37615,37602,22387,64077,64044,64132,64074,64110,64081,64052,23271,64030,64059,63894,63885,64244,64349,64131,64371,64277,64120,59167,59177,59176,59168,61114,63858,63873,62848,63554,63883,63884,63886,59166,22383,59178,59169,64187,63872,64448,64290,63983,37605,23258,64135,63483,63231,64012,23265,64117,64140,22400,22397,22395,22399,22396","Page":%d,"Limit":3}`

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

	totalPages := 234

	var allProducts []ProductInfo

	// Loop through each page
	for page := 1; page <= totalPages; page++ {
		payload := fmt.Sprintf(payloadTemplate, page)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		setHeaders(req)

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
			provider := s.Find("th.col-s span").Last().Text()
			product := s.Find("th.font-black.text-center").Text()
			interestText := doc.Find(fmt.Sprintf("tbody td.cmpr-col.col%d.text-center span.text-bold", i+1)).Text()
			minRate, maxRate := parseInterestRate(interestText)

			conditionalRate := parseConditionalRate(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-intr.trbox-shadow td.cmpr-col.col%d", i+1)).Text())

			interestCalculationMethod := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-intrmthd.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())
			taxFree := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-intrwotax.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())
			interestPaymentPeriod := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-intrterm.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())
			interestPaymentMethod := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-intrch.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())
			interestPenaltyConditions := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-intrexc.trbox-shadow td.cmpr-col.col%d", i+1)).Text())

			fixedTerm := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-dpstterm.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())
			minimumOpeningBalanceText := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-blncmin.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())
			minimumOpeningBalance := parseInt(minimumOpeningBalanceText)
			maximumDepositLimitText := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-blncmax.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())
			maximumDepositLimit := parseInt(maximumDepositLimitText)
			otherProductRequirements := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-prodbuy.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())
			minAge, maxAge := parseAgeRange(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-age.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())

			specificOpeningConditions := parseSectionedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-opencond.trbox-shadow td.cmpr-col.col%d", i+1)).Text())

			minimumDepositPerTransactionText := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-mindpst.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())
			var minimumDepositPerTransaction *int
			if minimumDepositPerTransactionText == "ไม่มีกำหนด" {
				minimumDepositPerTransaction = nil
			} else {
				minimumDepositPerTransactionValue := parseInt(minimumDepositPerTransactionText)
				minimumDepositPerTransaction = minimumDepositPerTransactionValue
			}

			additionalDepositsAllowed := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-topup.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())
			partialWithdrawalsAllowed := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-wdprtprnc.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())

			depositWithdrawTransferConditions := parseSectionedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-wdprtxnblnc.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())

			accountRenewalWhenDue := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-accrenew.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())

			insurance := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-insrnc.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())
			insuranceCompany := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-insrnccompany.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())

			insuranceCoverageLimitText := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-insrnclimit.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())
			var insuranceCoverageLimit *int
			if insuranceCoverageLimitText != "" {
				limit, err := strconv.Atoi(strings.ReplaceAll(insuranceCoverageLimitText, ",", ""))
				if err == nil {
					insuranceCoverageLimit = &limit
				}
			}

			insuranceConditions := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-insrnccond.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())

			accountMaintenanceFee := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-accmtnc.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())
			smsNotificationFee := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-accmsms.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())

			passbookReplacementFeeText := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-accmopenbk.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())
			var passbookReplacementFee *int
			if passbookReplacementFeeText != "" {
				fee, err := strconv.Atoi(strings.ReplaceAll(passbookReplacementFeeText, " บาท/เล่ม", ""))
				if err == nil {
					passbookReplacementFee = &fee
				}
			}

			transactionHistoryFee := parseTranSactionHistoryFee(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-accmbranch.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())

			accountClosureFee := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-accmclosebk.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())

			// General Fee parsing
			coinCountingFee := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-feecorn.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())
			crossBankDepositWithdrawalFee := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-feedeposit.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())

			// Updated parsing for other_provider_deposit_fee_cdm_atm
			otherProviderDepositFeeCDMATM := parseSlashSeparatedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-feecdm.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())

			sameProviderDepositFeeCDMATM := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-feecdm2.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())
			depositWithdrawalAgentFee := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-feeother2 trbox-shadow td.cmpr-col.col%d span", i+1)).Text())

			autoTransferFeeSavingsChecking := parseAutoTransferFee(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-feetranfer.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())

			crossBankTransferFee := parseSlashSeparatedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-feetranfer2.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())

			otherFees := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-feeother1.trbox-shadow td.cmpr-col.col%d span", i+1)).Text())

			// Additional Info parsing
			productWebsiteLink := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-url.trbox-shadow td.cmpr-col.col%d a.prod-url", i+1)).AttrOr("href", ""))
			feeWebsiteLink := cleanedText(doc.Find(fmt.Sprintf("tbody tr.attr-header.attr-feeurl.trbox-shadow td.cmpr-col.col%d a.prod-url", i+1)).AttrOr("href", ""))

			if strings.TrimSpace(provider) != "" && strings.TrimSpace(product) != "" {
				productInfo := ProductInfo{
					Provider: cleanedText(provider),
					Product:  cleanedText(product),
					Interest: Interest{
						MinRate:                   minRate,
						MaxRate:                   maxRate,
						ConditionalRate:           conditionalRate,
						InterestCalculationMethod: interestCalculationMethod,
						TaxFree:                   taxFree,
						InterestPaymentPeriod:     interestPaymentPeriod,
						InterestPaymentMethod:     interestPaymentMethod,
						InterestPenaltyConditions: interestPenaltyConditions,
					},
					AccountOpeningConditions: AccountOpeningConditions{
						FixedTerm:                 fixedTerm,
						MinimumOpeningBalance:     minimumOpeningBalance,
						MaximumDepositLimit:       maximumDepositLimit,
						OtherProductRequirements:  otherProductRequirements,
						MinAge:                    minAge,
						MaxAge:                    maxAge,
						SpecificOpeningConditions: specificOpeningConditions,
					},
					ProductUsageConditions: ProductUsageConditions{
						MinimumDepositPerTransaction:      minimumDepositPerTransaction,
						AdditionalDepositsAllowed:         additionalDepositsAllowed,
						PartialWithdrawalsAllowed:         partialWithdrawalsAllowed,
						DepositWithdrawTransferConditions: depositWithdrawTransferConditions,
						AccountRenewalWhenDue:             accountRenewalWhenDue,
					},
					Insurance: Insurance{
						Insurance:              insurance,
						InsuranceCompany:       insuranceCompany,
						InsuranceCoverageLimit: insuranceCoverageLimit,
						InsuranceConditions:    insuranceConditions,
					},
					ProductFees: ProductFees{
						AccountMaintenanceFee:  accountMaintenanceFee,
						SMSNotificationFee:     smsNotificationFee,
						PassbookReplacementFee: passbookReplacementFee,
						TransactionHistoryFee:  transactionHistoryFee,
						AccountClosureFee:      accountClosureFee,
					},

					GeneralFee: GeneralFee{
						CoinCountingFee:                coinCountingFee,
						CrossBankDepositWithdrawalFee:  crossBankDepositWithdrawalFee,
						OtherProviderDepositFeeCDMATM:  otherProviderDepositFeeCDMATM,
						SameProviderDepositFeeCDMATM:   sameProviderDepositFeeCDMATM,
						DepositWithdrawalAgentFee:      depositWithdrawalAgentFee,
						AutoTransferFeeSavingsChecking: autoTransferFeeSavingsChecking,
						CrossBankTransferFee:           crossBankTransferFee,
						OtherFees:                      otherFees,
					},

					AdditionInfo: AdditionInfo{
						ProductWebsiteLink: productWebsiteLink,
						FeeWebsiteLink:     feeWebsiteLink,
					},
				}
				allProducts = append(allProducts, productInfo)
			}
		})

		// Stop for 5 seconds before making the next request to avoid overloading the server
		time.Sleep(5 * time.Second)
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

func parseInt(text string) *int {
	if text == "ไม่กำหนด" || text == "" {
		return nil
	}

	// Remove commas from the text
	cleanedText := strings.ReplaceAll(text, ",", "")

	// Extract the numeric part of the string
	numStr := regexp.MustCompile(`\d+`).FindString(cleanedText)
	if numStr == "" {
		return nil
	}

	// Convert the extracted numeric string to an integer
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return nil
	}
	return &num
}

func parseAgeRange(ageText string) (*int, *int) {
	ages := strings.Split(ageText, "-")
	if len(ages) != 2 {
		return nil, nil
	}

	minAgeStr := strings.TrimSpace(ages[0])
	maxAgeStr := strings.TrimSpace(ages[1])

	minAge, errMin := strconv.Atoi(minAgeStr)
	maxAge, errMax := strconv.Atoi(maxAgeStr)

	if errMin != nil || errMax != nil {
		return nil, nil
	}

	return &minAge, &maxAge
}

func parseSectionedText(text string) []string {
	sectionPattern := regexp.MustCompile(`(\d+\.)|(-)`)
	matches := sectionPattern.FindAllStringIndex(text, -1)

	var result []string
	startIndex := 0

	for _, match := range matches {
		if match[0] > startIndex {
			section := strings.TrimSpace(text[startIndex:match[0]])
			if section != "" {
				result = append(result, section)
			}
		}
		startIndex = match[0]
	}

	// Append the last section if any
	if startIndex < len(text) {
		section := strings.TrimSpace(text[startIndex:])
		if section != "" {
			result = append(result, section)
		}
	}

	return result
}

func parseSlashSeparatedText(text string) []string {
	// Remove newline characters and unnecessary spaces
	cleanedText := strings.ReplaceAll(text, "\n", "")
	cleanedText = strings.Join(strings.Fields(cleanedText), " ")

	sections := strings.Split(cleanedText, "/")
	var result []string

	for _, section := range sections {
		cleanSection := strings.TrimSpace(section)
		if cleanSection != "" {
			result = append(result, cleanSection)
		}
	}
	return result
}

func parseTranSactionHistoryFee(text string) []string {
	sections := []string{}
	parts := strings.Split(text, "เงื่อนไข:")
	if len(parts) > 0 {
		mainPart := strings.TrimSpace(parts[0])
		sections = append(sections, strings.Split(mainPart, "ขอใบแสดงรายการย้อนหลัง")...)
	}

	if len(parts) > 1 {
		conditionPart := strings.TrimSpace(parts[1])
		conditionSections := strings.Split(conditionPart, "-")
		for _, section := range conditionSections {
			cleanedSection := strings.TrimSpace(section)
			if cleanedSection != "" {
				sections = append(sections, cleanedSection)
			}
		}
	}

	for i, sec := range sections {
		sections[i] = strings.TrimSpace(sec)
		if strings.HasPrefix(sec, "น้อยกว่า") || strings.HasPrefix(sec, "ตั้งแต่") || strings.HasPrefix(sec, "มากกว่า") {
			sections[i] = "- " + sec
		} else if sec != "" && !strings.HasPrefix(sec, "-") {
			sections[i] = "ขอใบแสดงรานการย้อนหลัง" + sec
		}
	}
	return sections
}

func parseAutoTransferFee(text string) []string {
	cleanedText := strings.ReplaceAll(text, "\n", "")
	cleanedText = strings.Join(strings.Fields(cleanedText), " ")

	sections := []string{}
	parts := strings.Split(cleanedText, "เงื่อนไข")
	if len(parts) > 0 {
		mainPart := strings.TrimSpace(parts[0])
		if mainPart != "" && mainPart != "-" {
			sections = append(sections, mainPart)
		}
	}

	if len(parts) > 1 {
		conditionsPart := strings.TrimSpace(parts[1])
		conditionSections := strings.Split(conditionsPart, "-")
		for _, section := range conditionSections {
			cleanSection := strings.TrimSpace(section)
			if cleanSection != "" {
				sections = append(sections, cleanSection)
			}
		}
	}

	// Remove colon if present
	for i, sec := range sections {
		sections[i] = strings.ReplaceAll(sec, ":", "")
	}

	return sections
}

func parseConditionalRate(text string) []string {
	text = cleanedText(text)
	return strings.Split(text, "เงินฝากส่วนที่เกิน")
}

func setHeaders(req *http.Request) {
	req.Header.Set("accept", "text/plain, */*; q=0.01")
	req.Header.Set("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Header.Set("cookie", `verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m01s4s5a&sl=0&tt=0"; _uid6672=16B5DEBD.18; _ctout6672=1; visit_time=287; _ga_NLQFGWVNXN=GS1.1.1724292807.22.1.1724293102.60.0.0`)
	req.Header.Set("origin", "https://app.bot.or.th")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://app.bot.or.th/1213/MCPD/ProductApp/Deposit/CompareProduct")
	req.Header.Set("sec-ch-ua", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "macOS")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	req.Header.Set("verificationtoken", "bbRZjCUNPATqay7Tj6HdKx1qQ-tcSOFS3_P4itlr8FBNruy-3MAPLQPnO9gdpzapGwEyizGCY5dtcRdp4PECDFTZGEVqkdWGR7QemV7y_RA1,s0LfE_Sesi9xWyLJErtHVN4k5o4DsrW_b1MhwcdDFa6ez66W4GPFb4GRkmVucv7_Q0neSgiGyEhqwv9U0iJqn4gdbIJ_p0jOYlQEUHuX-_U1")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
}
