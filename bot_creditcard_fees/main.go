package main

import (
	"bytes"
	"creditcrad_fee/models"
	"creditcrad_fee/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/CreditFee/CompareProductList"
	method := "POST"
	payloadTemplate := `{"ProductIdList":"5148,5180,5114,5213,5233,5177,5161,4471,4472,4445,4482,4483,4484,4479,4480,4475,4481,4452,4453,4476,4477,4478,4454,4455,4456,4457,4458,4459,4460,4461,4462,4463,4464,4465,4466,4444,4446,4447,4448,4449,4450,4451,4467,4473,4474,4468,4470,4469,5588,5589,5593,5594,5590,5591,5586,5587,5592,5467,5478,5473,5462,2244,2255,2260,2256,2245,2248,2240,2254,2243,2253,2246,2257,2247,2250,2241,2242,2258,2249,2252,2259,2251,4654,4653,4655,4656,4657,4658,4659,4660,5578,5580,5582,749,752,750,753,751,754,3627,3631,3600,3664,3624,3601,3602,3632,3625,3603,3633,3626,3604,3661,3620,3634,3636,3635,3606,3638,3637,3640,3639,3605,3672,3668,3608,3628,3641,3662,3621,3648,3607,3642,3649,3643,3651,3613,3644,3652,3609,3645,3653,3646,3654,3647,3615,3656,3610,3650,3671,3665,3658,3655,3616,3669,3663,3670,3611,3617,3612,3618,3660,3657,3614,3666,3667,3619,3659,3629,3622,3630,3623,5551,5524,5539,5566,5540,5525,5567,5527,5554,5541,5526,5553,5552,5528,5568,5555,5542,5529,5556,5543,5557,5562,5536,5574,5570,5531,5544,5558,5547,5561,5535,5563,5571,5532,5545,5559,5560,5572,5533,5546,5573,5534,5548,5538,5575,5549,5564,5565,5576,5550,5537,5577,5429,5435,3996,3995,3997,4760,4765,1606,1605,1600,1601,1602,1603,1604,1607,5489,5500,5480,5499,5479,5492,5491,5494,5502,5482,5493,5481,5496,5497,5495,5483,5484,5486,5498,5501,5503,3802,3808,3804,3806,3797,3798,3799,3800,3801,5287,5282,5285,5193,3994,3993,5211,5173,5137,5156,5126,5222,5162,5202,5158,5147,5240,5294,5296,5323,5306,5300,5329,5266,5335,5305,5304,5246","Page":%d,"Limit":3}`
	// payloadTemplate := `{"ProductIdList":"5148,5180,5114,5213,5233,5177,5161,4471,4472,4445,4482,4483,4484,4479,4480,4475,4481,4452","Page":%d,"Limit":3}`
	initialPayload := fmt.Sprintf(payloadTemplate, 1)
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(initialPayload)))
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

	var fees []models.CreditCardFee

	// Loop through each page to gather all data
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
			log.Printf("Error sending request for page %d: %v", page, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Request for page %d failed with status: %d\n", page, resp.StatusCode)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body for page %d: %v", page, err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			log.Printf("Error parsing HTML for page %d: %v", page, err)
			continue
		}

		for i := 1; i <= 3; i++ {
			col := "col" + strconv.Itoa(i)
			provider := utils.CleanText(doc.Find(fmt.Sprintf("th.%s span", col)).Text())
			product := utils.CleanText(doc.Find(fmt.Sprintf("th.prod-col%d span", i)).Text())

			feeDetails := models.CreditCardFee{
				Provider: provider,
				Product:  product,
				GeneralFees: models.GeneralFees{
					EntranceFeeMainCard:          utils.ExtractFee(doc, "attr-primaryHolderEntranceFeeDisplay", col),
					AnnualFeeMainCard:            utils.SplitByCondition(utils.ExtractFee(doc, "attr-primaryHolderAnnualFee", col), "ปี"),
					CurrencyConversionRisk:       utils.SplitByCondition(utils.ExtractFee(doc, "attr-CostFXRisk", col), "เงื่อนไข"),
					CashAdvanceFee:               utils.ExtractFee(doc, "attr-cashAdvanceFee", col),
					ReplacementCardFee:           utils.ExtractFee(doc, "attr-replacementCardFee", col),
					EntranceFeeSupplementaryCard: utils.ExtractFee(doc, "attr-supplementaryCardHolderEntranceFeeDisplay", col),
					AnnualFeeSupplementaryCard:   utils.SplitByCondition(utils.ExtractFee(doc, "attr-supplementaryCardHolderAnnualFeeFirstYear", col), "ปี"),
					NewPINRequestFee:             utils.ExtractFee(doc, "attr-replacementCardFPinFee", col),
					StatementCopyFee:             utils.ExtractFee(doc, "attr-copyStatementFee", col),
					TransactionVerificationFee:   utils.ExtractFee(doc, "attr-TransactionVerifyFee", col),
					SalesSlipCopyFee:             utils.ExtractFee(doc, "attr-copySaleSlipFee", col),
					ReturnedChequeFee:            utils.ExtractFee(doc, "attr-fineChequeReturn", col),
					TaxPaymentFee:                utils.ExtractFee(doc, "attr-GovernmentAgencyRelatedPaymentFee", col),
					DebtCollectionFee:            utils.SplitByCondition(utils.ExtractFee(doc, "attr-debtCollectionFee", col), "-"),
				},
				PaymentFees: models.PaymentFees{
					FeeFreeChannels:        utils.SplitByCondition(utils.ExtractFee(doc, "attr-freePaymentChannel", col), "-"),
					DirectDebitServiceFee:  utils.ExtractFee(doc, "attr-directDebitFromAccountFee", col),
					DirectDebitOtherFee:    utils.ExtractFee(doc, "attr-directDebitFromAccountFeeOther", col),
					BankCounterFee:         utils.ExtractFee(doc, "attr-BankCounterServiceFee", col),
					OtherBankCounterFee:    utils.ExtractFee(doc, "attr-BankCounterServiceFeeOther", col),
					PaymentServicePointFee: utils.SplitByCondition(utils.ExtractFee(doc, "attr-CounterServiceFeeOther", col), "-"),
					OnlinePaymentFee:       utils.ExtractFee(doc, "attr-paymentOnlineFee", col),
					ATMPaymentFee:          utils.ExtractFee(doc, "attr-paymentCDMATMFee", col),
					PhonePaymentFee:        utils.ExtractFee(doc, "attr-paymentPhoneFee", col),
					ChequeOrMoneyOrderFee:  utils.ExtractFee(doc, "attr-paymentChequeOrMoneyOrderFee", col),
					OtherPaymentChannels:   utils.SplitByCondition(utils.ExtractFee(doc, "attr-paymentOtherChannelFee", col), "-"),
				},
				OtherFees: models.OthersFees{
					OtherFees: utils.ExtractFeePtr(doc, "attr-other", col),
				},
				AdditionalInfo: models.AdditionalInfo{
					WebsiteFeeLink: utils.ExtractLink(doc, "attr-Feeurl", col),
				},
			}

			fees = append(fees, feeDetails)
		}
	}

	jsonData, err := json.MarshalIndent(fees, "", " ")
	if err != nil {
		log.Fatalf("Error converting to JSON: %v", err)
	}

	err = os.WriteFile("creditcard_fees.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	fmt.Println("Data saved to creditcard_fees.json")
}
