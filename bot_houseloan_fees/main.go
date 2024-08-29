package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"houseLoan_fees/models"
	"houseLoan_fees/utils"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/homeloanFee/CompareProductList"
	payloadTemplate := `{"ProductIdList":"21510,17709,20310,20242,20275,20317,20274,20315,20271,20312,20308,20299,20259,20314,20287,20276,20260,20270,20289,20311,20269,20328,20272,20313,20257,20330,20307,20306,20258,20268,20321,20281,20331,20241,20316,20283,20262,20261,20318,20319,20280,20253,20326,20294,20298,20255,20336,20250,20302,20277,20266,20252,20245,20278,20325,20290,20323,20322,20332,20247,20285,20320,20244,20263,20291,20333,20264,20324,20293,20246,20288,20295,20243,20303,20334,20248,20297,20273,20335,20249,20286,20301,20305,20300,20256,20265,20309,20304,20279,20267,20329,20296,20251,20284,20327,20292,20282,20254,21476,21504,21472,21514,21471,21462,21505,21488,21506,21515,21478,21498,21473,21487,21477,21481,21467,21492,21469,21475,21466,21493,21517,21497,21519,21512,21500,21483,21502,21474,21460,21461,21465,21489,21494,21513,21499,21480,21485,21518,21491,21486,21457,21463,21520,21496,21501,21459,21479,16469,16461,20220,20219,20222,20221,20229,20224,20223,20226,20225,20228,20227,15219,15221,15220,15222,20843,20822,20833,20852,20854,20856,1740,20939,20932,20936,20927,20938,20925,20923,20924,20942,20928,20940,20921,20922,20929,20941,20920,20931,20926,20937,20934,20933,20935,20919,20930,1738,1739,21210,21233,21133,21113,21244,21142,21053,21060,21227,21115,21095,21134,21226,21114,21242,21135,21453,21455,21454,21456,20794,20770,20789,20790,20779,20780,20792,20772,20773,20788,20778,20781,20791,20771,20796,20793,20795,20776,20783,20798,20768,20787,20782,20785,20775,20777,20774,20784,20786,20797,20769,17220,17221,17226,17227,17222,17223,17224,17225,17216,17217,17218,17219,16203,16176,16184,16193,16183,16171,16204,16233,16214,16223,16194,16213,16254,16234,16244,16249,16243,16224,14965,14981,14980,14982,21085,21170,21173,21194,21195,21218,21221,21084,21124,21030,21116,21193,21185,21182,21223,21126,21092,21093,21127,21032,21035,21086,21087,21196,21197,21172,21175,21220,21091,21187,21222,21198,21199,21045,21089,21184,21129,21096,21206,21241,21219,21048,21051,21140,21141,21243,21240,21132,21213,21112,21061,21139,21040,21217,21214,21215,21038,21041,21130,21131,21238,21118,21209,21103,21049,21039,21201,21059,21101,21212,21192,21239,21119,21047,21036,21098,21099,21224,21225,21189,21128,21037,21208,21211,21110,21111,21058,21121,21123,21100,21216,21043,21056,21138,700,706,707,701,702,708,703,709,710,705,704,21004,21005,21006,21007,21008,21009,21010,21011,21016,21017,21018,21019,21012,21013,21014,21015,21020,21021,21022,21023,21245,21247,21257,21256,21255,21248,21246,21249,21250,21260,21258,21259,21251,21253,21262,21254,21261,21263,21252,12618,17712,17710,17711,15456,15455","Page":%d,"Limit":3}`
	// payloadTemplate := `{"ProductIdList":"21510,17709,20310,20242,20275,20317,20274,20315,20271,20312,20308,20299,20259,20314,20287,20309,20304,20279,20267,20329,20296,20251,20284,20327,20292,20282,20254,21476","Page":%d,"Limit":3}`

	initialPayload := fmt.Sprintf(payloadTemplate, 1)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(initialPayload)))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	utils.AddHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check if response status is OK
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	totalPages := utils.DetermineTotalPage(doc)
	if totalPages == 0 {
		log.Fatal("Could not determine the total number of pages")
	}

	houseloanFeesDetailsList := []models.HouseLoanFeesDetails{}

	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages)
		payload := fmt.Sprintf(payloadTemplate, page)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			log.Printf("Error creating request for page %d: %v", page, err)
			continue
		}

		utils.AddHeader(req)

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error making request for page %d: %v", page, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Request for page %d failed with status: %d\n", page, err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Printf("Error parsing HTML for page %d: %v", page, err)
			continue
		}

		for i := 1; i <= 3; i++ {
			col := "col" + strconv.Itoa(i)
			houseloanFeesDetails := models.HouseLoanFeesDetails{
				Provider: utils.CleanText(doc.Find(fmt.Sprintf("th.%s span", col)).Eq(1).Text()),
				Product:  utils.CleanText(doc.Find(fmt.Sprintf("th.prod-col%d span", i)).First().Text()),
			}

			// Extracting fee details using ParseFeeDetail
			houseloanFeesDetails.InterestRates.DefaultInterestRate = utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-DfltInterestRate td.%s span", col)).Text())
			houseloanFeesDetails.InterestRates.SurveyAndAppraisalFee = []models.FeeDetail{
				utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-SurveyAndAppraisalFee td.%s span", col)).Text()),
			}
			houseloanFeesDetails.InterestRates.MRTA = utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-MortgageReducingTermAssuranceCancelled td.%s span", col)).Text())
			houseloanFeesDetails.InterestRates.StampDuty = utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-DutyStampFee td.%s span", col)).Text())
			houseloanFeesDetails.InterestRates.MortgageFee = utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-MortgateFee td.%s span", col)).Text())
			houseloanFeesDetails.InterestRates.TransferOwnershipFee = utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-TranfersOwnerFee td.%s span", col)).Text())
			houseloanFeesDetails.InterestRates.CreditBureauFee = utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-CreditBureauFee td.%s span", col)).Text())
			houseloanFeesDetails.InterestRates.FireInsurancePremium = utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-FireInsuracePremiumsFee td.%s span", col)).Text())
			houseloanFeesDetails.InterestRates.OtherChequeReturnedFee = utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-OtherChequeReturnedFee td.%s span", col)).Text())
			houseloanFeesDetails.InterestRates.InsufficientDirectDebitFee = utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-InsufficientDirectDebitFee td.%s span", col)).Text())
			houseloanFeesDetails.InterestRates.CopyStatementReissuingFee = utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-CopyStatementReissuingFee td.%s span", col)).Text())
			houseloanFeesDetails.InterestRates.ChequeReturnedFee = utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-ChequeReturnedFee td.%s span", col)).Text())
			houseloanFeesDetails.InterestRates.DebtCollectionFee = []models.FeeDetail{
				utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-DebtCollectionFee td.%s span", col)).Text()),
			}
			houseloanFeesDetails.InterestRates.ChangingInterestRateFee = utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-ChangingInterestRateFee td.%s span", col)).Text())
			houseloanFeesDetails.InterestRates.RefinanceFee = utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-RefinanceFee td.%s span", col)).Text())

			// Payment fees
			houseloanFeesDetails.PaymentsFees.DirectDebitFromProvider = utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-DirectDebitFromAccountFee td.%s span", col)).Text())
			houseloanFeesDetails.PaymentsFees.DirectDebitFromOtherProvider = []models.FeeDetail{
				utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-DirectDebitFromAccountFeeOther td.%s span", col)).Text()),
			}
			houseloanFeesDetails.PaymentsFees.AtProviderBranch = utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-BankCounterServiceFee td.%s span", col)).Text())
			houseloanFeesDetails.PaymentsFees.AtOtherProviderBranch = utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-BankCounterServiceFeeOther td.%s span", col)).Text())
			houseloanFeesDetails.PaymentsFees.AtPaymentServicePoint = []models.FeeDetail{
				utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-CounterServiceFee td.%s span", col)).Text()),
			}
			houseloanFeesDetails.PaymentsFees.OnlinePayment = []models.FeeDetail{
				utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-PaymentOnlineFee td.%s span", col)).Text()),
			}
			houseloanFeesDetails.PaymentsFees.CDMOrATM = utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-PaymentCDMATMFee td.%s span", col)).Text())
			houseloanFeesDetails.PaymentsFees.PhonePayment = utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-PaymentPhoneFee td.%s span", col)).Text())
			houseloanFeesDetails.PaymentsFees.ChequeOrMoneyOrderPayment = utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-PaymentChequeOrMoneyOrderFee td.%s span", col)).Text())
			houseloanFeesDetails.PaymentsFees.OtherPaymentChannels = utils.ParseFeeDetail(doc.Find(fmt.Sprintf(".attr-PaymentOtherChannelFee td.%s span", col)).Text())

			houseloanFeesDetails.OtherFees.OtherFees = utils.NullEmpty(doc.Find(fmt.Sprintf(".attr-other td.%s span", col)).Text())

			houseloanFeesDetails.AdditionalInfo.FeewebsiteLink = utils.CleanText(doc.Find(fmt.Sprintf(".attr-Feeurl td.%s a", col)).AttrOr("href", ""))

			houseloanFeesDetailsList = append(houseloanFeesDetailsList, houseloanFeesDetails)
		}

		time.Sleep(2 * time.Second)
	}

	jsonOutput, err := json.MarshalIndent(houseloanFeesDetailsList, "", " ")
	if err != nil {
		log.Fatalf("Failed to marshaling JSON: %v", err)
	}

	filename := "house_loan_details.json"
	err = os.WriteFile(filename, jsonOutput, 0644)
	if err != nil {
		log.Fatalf("Failed to write JSON to file: %v", err)
	}

	fmt.Printf("JSON data svaed to %s\n", filename)
}
