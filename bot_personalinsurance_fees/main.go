package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"personalinsurance_fees/models"
	"personalinsurance_fees/utils"
	"time"

	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/PLoanwithorwithoutCollateralFee/CompareProductList"
	// payloadTemplate := `{"ProductIdList":"6545,6536,877,875,7912,7923,7922,7908,7914,7918,7921,7927,7915,7925,7909,7924,7920,7917,7911,7916,7913,7907,7910,7926,7919,8331,8341,8524,8571,8310,8516,8545,8336,8444,8293,8562,8472,8268,8479,8284,8258,8442,8348,8377,8387,8384,8382,8522,8379,8273,8497,8275,8356,8309,8495,8523,8578,8519,8409,8321,8431,8415,8561,8394,8412,8567,8461,8254,8360,8311,8316,8314,8484,8446,8298,8287,7257,7253,7254,7255,7256,2948,2983,3016,2982,2996,2976,2991,2968,2984,2974,2953,5808,5811,5813,5812,5807,5806,5810,6537,6543,6544,6541,6539,6538,6542,6535,6540,3006,3007,2972,3014,3011,2977,2990,3015,2956,2955,2999,3010,2964,2958,2992,3002,3000,2973,3017,2949,3008,2967,2965,2957,1436,1434,1260,1265,1261,1262,1267,3888,1002,229,231,228,230,232,874,876,7259,7258,7251,7252,8308,8417,8485,8299,8551,8393,8547,8297,8264,8504,8493,8304,8569,8447,8402,8580,8454,8549,8262,8474,8555,8457,8450,8424,8344,8564,8570,8259,8518,8464,8281,8416,8257,8276,8546,8434,8488,8398,8420,8449,8427,8539,8509,8430,8426,8330,8500,8532,8371,8397,8366,8381,8334,8411,8573,8329,8338,8490,8358,8306,8271,8300,8463,8529,1591,1592,1593","Page":%d,"Limit":3}`
	payloadTemplate := `{"ProductIdList":"6545,6536,877,875,7912,7923,7922,7908,7914,7918,7921,7927,7915,7925","Page":%d,"Limit":3}`

	initialPayload := fmt.Sprintf(payloadTemplate, 1)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(initialPayload)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	utils.AddHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Non-OK HTTP status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	totalPages := utils.DetermineTotalPage(doc)
	if totalPages == 0 {
		log.Fatal("Could not determine the total number of pages")
	}

	var personalFeeDetails []models.PersonalFeeDetails
	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages)
		payload := fmt.Sprintf(payloadTemplate, page)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			log.Printf("Error creating request for page %d: %v", page, err)
			continue
		}
		utils.AddHeader(req)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error sending request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Non-OK HTTP status: %s", resp.Status)
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			log.Fatalf("Error parsing HTML: %v", err)
		}

		for i := 1; i <= 3; i++ {
			col := fmt.Sprintf("col%d", i)
			serviceProvider := utils.CleanText(doc.Find(fmt.Sprintf("th.col-s.col-s-%d span", i)).Last().Text())
			product := utils.CleanText(doc.Find(fmt.Sprintf("th.font-black.text-center.prod-%s span", col)).Text())

			fees := models.PersonalFeeDetails{
				ServiceProvider: serviceProvider,
				Product:         product,
			}

			// Extracting General Fees with numerical values
			fees.GeneralFees.InternalEvaluation = extractFeeDetails(doc.Find(fmt.Sprintf(".attr-SurveyAndAppraisalFeeByInternal .%s span", col)).Text())
			fees.GeneralFees.ExternalEvaluation = extractFeeDetails(doc.Find(fmt.Sprintf(".attr-SurveyAndAppraisalFeeByExternal .%s span", col)).Text())
			fees.GeneralFees.StampDuty = extractFeeDetails(doc.Find(fmt.Sprintf(".attr-StampDutyFee .%s span", col)).Text())
			fees.GeneralFees.MortgageFee = extractFeeDetails(doc.Find(fmt.Sprintf(".attr-MortgageFee .%s span", col)).Text())
			fees.GeneralFees.CreditCheck = extractFeeDetails(doc.Find(fmt.Sprintf(".attr-CreditBureau .%s span", col)).Text())
			fees.GeneralFees.ReturnedChequeFee = extractFeeDetails(doc.Find(fmt.Sprintf(".attr-ReturnedCheque .%s span", col)).Text())
			fees.GeneralFees.InsufficientFundsFee = extractFeeDetails(doc.Find(fmt.Sprintf(".attr-InsufficientDirectDebitCharge .%s span", col)).Text())

			// Splitting text by "-" for array fields
			fees.GeneralFees.StatementReIssuingFee = utils.SplitAndTrim(utils.CleanText(doc.Find(fmt.Sprintf(".attr-StatementReIssuingFee .%s span", col)).Text()))
			fees.GeneralFees.DebtCollectionFee = utils.SplitAndTrim(utils.CleanText(doc.Find(fmt.Sprintf(".attr-DebtCollectionFee .%s span", col)).Text()))

			// Extracting Payment Fees
			fees.PaymentFees.DebitFromAccount = extractFeeDetails(doc.Find(fmt.Sprintf(".attr-DirectDebitFromAccountFee .%s span", col)).Text())
			fees.PaymentFees.DebitFromOtherAccount = extractFeeDetails(doc.Find(fmt.Sprintf(".attr-DirectDebitFromAccountFeeOther .%s span", col)).Text())
			fees.PaymentFees.PayAtProviderBranch = extractFeeDetails(doc.Find(fmt.Sprintf(".attr-BankCounterServiceFee .%s span", col)).Text())
			fees.PaymentFees.PayAtOtherBranch = extractFeeDetails(doc.Find(fmt.Sprintf(".attr-BankCounterServiceFeeOther .%s span", col)).Text())
			fees.PaymentFees.PayAtServicePoint = extractFeeDetails(doc.Find(fmt.Sprintf(".attr-CounterServiceFeeOther .%s span", col)).Text())
			fees.PaymentFees.PayOnline = extractFeeDetails(doc.Find(fmt.Sprintf(".attr-paymentOnlineFee .%s span", col)).Text())
			fees.PaymentFees.PayViaCDMATM = extractFeeDetails(doc.Find(fmt.Sprintf(".attr-paymentCDMATMFee .%s span", col)).Text())
			fees.PaymentFees.PayViaPhone = extractFeeDetails(doc.Find(fmt.Sprintf(".attr-paymentPhoneFee .%s span", col)).Text())
			fees.PaymentFees.PayViaChequeOrMoneyOrder = extractFeeDetails(doc.Find(fmt.Sprintf(".attr-paymentChequeOrMoneyOrderFee .%s span", col)).Text())
			fees.PaymentFees.PayViaOtherChannels = extractFeeDetails(doc.Find(fmt.Sprintf(".attr-paymentOtherChannelFee .%s span", col)).Text())

			// Extracting Other Fees
			fees.OtherFees.OtherFee = utils.CleanText(doc.Find(fmt.Sprintf(".attr-OtherFee .%s span", col)).Text())

			// Extracting Additional Info
			fees.AdditionalInfo.FeeWebsite = utils.CleanText(doc.Find(fmt.Sprintf(".attr-FeeURL .%s a", col)).AttrOr("href", ""))

			personalFeeDetails = append(personalFeeDetails, fees)
		}

		time.Sleep(2 * time.Second)
	}

	jsonData, err := json.MarshalIndent(personalFeeDetails, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	err = os.WriteFile("personalFees.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	fmt.Println("Data saved to personalFees.json")
}

// extractFeeDetails processes the fee string to include both the original text and extracted numerical values.
func extractFeeDetails(text string) models.FeeDetail {
	originalText := utils.CleanText(text)
	var feeDetail models.FeeDetail
	feeDetail.OriginalText = originalText

	// Extract percentage value if present
	if strings.Contains(originalText, "%") {
		percentageStr := strings.TrimSuffix(strings.Split(originalText, "%")[0], " ")
		percentage, err := strconv.ParseFloat(percentageStr, 64)
		if err == nil {
			feeDetail.Percentage = &percentage
		}
	}

	// Extract amount values if present
	words := strings.Fields(originalText)
	for _, word := range words {
		if amount, err := strconv.Atoi(strings.TrimSuffix(word, "บาท")); err == nil {
			feeDetail.FeeAmount = &amount
		}
	}

	return feeDetail
}
