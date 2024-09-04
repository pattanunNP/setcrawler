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
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

// Helper function to extract numeric value from a string using regex
func extractNumericValue(text string) *int {
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(text)
	if match == "" {
		return nil
	}
	value, err := strconv.Atoi(match)
	if err != nil {
		log.Printf("Error parsing amount from text '%s': %v", text, err)
		return nil
	}
	return &value
}

// Function to parse amount from text
func parseAmount(text string) *int {
	return extractNumericValue(text)
}

// Function to parse percentage from text
func parsePercentage(text string) *float64 {
	re := regexp.MustCompile(`\d+(\.\d+)?`)
	match := re.FindString(text)
	if match == "" {
		return nil
	}
	percentage, err := strconv.ParseFloat(match, 64)
	if err != nil {
		log.Printf("Error parsing percentage from text '%s': %v", text, err)
		return nil
	}
	return &percentage
}

func main() {
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/CreditFee/CompareProductList"
	method := "POST"
	payloadTemplate := `{"ProductIdList":"5148,5180,5114,5213,5233,5177,5161,4471,4472,4445,4482,4483,4484,4479,4480,4475,4481,4452","Page":%d,"Limit":3}`
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

			entranceFeeText := utils.ExtractFee(doc, "attr-primaryHolderEntranceFeeDisplay", col)
			if entranceFeeText == "" {
				log.Printf("Entrance fee text not found for provider: %s, product: %s", provider, product)
			}

			feeDetails := models.CreditCardFee{
				Provider: provider,
				Product:  product,
				GeneralFees: models.GeneralFees{
					EntranceFeeMainCard: models.FeeDetail{
						Text:   entranceFeeText,
						Amount: parseAmount(entranceFeeText),
					},
					AnnualFeeMainCard: []models.AnnualFeeDetail{
						{
							Text:          utils.ExtractFee(doc, "attr-primaryHolderAnnualFee", col),
							InitialAmount: parseAmount(utils.ExtractFee(doc, "attr-primaryHolderAnnualFee", col)),
						},
					},
					CurrencyConversionRisk: models.FeeDetail{
						Text:       utils.ExtractFee(doc, "attr-CostFXRisk", col),
						Percentage: parsePercentage(utils.ExtractFee(doc, "attr-CostFXRisk", col)),
					},
					CashAdvanceFee: models.FeeDetail{
						Text:   utils.ExtractFee(doc, "attr-cashAdvanceFee", col),
						Amount: parseAmount(utils.ExtractFee(doc, "attr-cashAdvanceFee", col)),
					},
					ReplacementCardFee: models.FeeDetail{
						Text:   utils.ExtractFee(doc, "attr-replacementCardFee", col),
						Amount: parseAmount(utils.ExtractFee(doc, "attr-replacementCardFee", col)),
					},
					EntranceFeeSupplementaryCard: models.FeeDetail{
						Text:   utils.ExtractFee(doc, "attr-supplementaryCardHolderEntranceFeeDisplay", col),
						Amount: parseAmount(utils.ExtractFee(doc, "attr-supplementaryCardHolderEntranceFeeDisplay", col)),
					},
					AnnualFeeSupplementaryCard: []models.AnnualFeeDetail{
						{
							Text:          utils.ExtractFee(doc, "attr-supplementaryCardHolderAnnualFeeFirstYear", col),
							InitialAmount: parseAmount(utils.ExtractFee(doc, "attr-supplementaryCardHolderAnnualFeeFirstYear", col)),
						},
					},
					NewPINRequestFee: models.FeeDetail{
						Text:   utils.ExtractFee(doc, "attr-replacementCardFPinFee", col),
						Amount: parseAmount(utils.ExtractFee(doc, "attr-replacementCardFPinFee", col)),
					},
					StatementCopyFee: models.FeeDetail{
						Text:   utils.ExtractFee(doc, "attr-copyStatementFee", col),
						Amount: parseAmount(utils.ExtractFee(doc, "attr-copyStatementFee", col)),
					},
					TransactionVerificationFee: models.FeeDetail{
						Text:   utils.ExtractFee(doc, "attr-TransactionVerifyFee", col),
						Amount: parseAmount(utils.ExtractFee(doc, "attr-TransactionVerifyFee", col)),
					},
					SalesSlipCopyFee: models.FeeDetail{
						Text:   utils.ExtractFee(doc, "attr-copySaleSlipFee", col),
						Amount: parseAmount(utils.ExtractFee(doc, "attr-copySaleSlipFee", col)),
					},
					ReturnedChequeFee: models.FeeDetail{
						Text:   utils.ExtractFee(doc, "attr-fineChequeReturn", col),
						Amount: parseAmount(utils.ExtractFee(doc, "attr-fineChequeReturn", col)),
					},
					TaxPaymentFee: models.FeeDetail{
						Text:   utils.ExtractFee(doc, "attr-GovernmentAgencyRelatedPaymentFee", col),
						Amount: parseAmount(utils.ExtractFee(doc, "attr-GovernmentAgencyRelatedPaymentFee", col)),
					},
					DebtCollectionFee: []models.FeeDetail{
						{
							Text:   utils.ExtractFee(doc, "attr-debtCollectionFee", col),
							Amount: parseAmount(utils.ExtractFee(doc, "attr-debtCollectionFee", col)),
						},
					},
				},
				PaymentFees: models.PaymentFees{
					FeeFreeChannels: utils.SplitByCondition(utils.ExtractFee(doc, "attr-freePaymentChannel", col), "-"),
					DirectDebitServiceFee: models.FeeDetail{
						Text:   utils.ExtractFee(doc, "attr-directDebitFromAccountFee", col),
						Amount: parseAmount(utils.ExtractFee(doc, "attr-directDebitFromAccountFee", col)),
					},
					BankCounterFee: models.FeeDetail{
						Text:   utils.ExtractFee(doc, "attr-BankCounterServiceFee", col),
						Amount: parseAmount(utils.ExtractFee(doc, "attr-BankCounterServiceFee", col)),
					},
					OnlinePaymentFee: models.FeeDetail{
						Text:   utils.ExtractFee(doc, "attr-paymentOnlineFee", col),
						Amount: parseAmount(utils.ExtractFee(doc, "attr-paymentOnlineFee", col)),
					},
					ATMPaymentFee: models.FeeDetail{
						Text:   utils.ExtractFee(doc, "attr-paymentCDMATMFee", col),
						Amount: parseAmount(utils.ExtractFee(doc, "attr-paymentCDMATMFee", col)),
					},
					PhonePaymentFee: models.FeeDetail{
						Text:   utils.ExtractFee(doc, "attr-paymentPhoneFee", col),
						Amount: parseAmount(utils.ExtractFee(doc, "attr-paymentPhoneFee", col)),
					},
					OtherPaymentChannels: utils.SplitByCondition(utils.ExtractFee(doc, "attr-paymentOtherChannelFee", col), "-"),
				},
				OtherFees: models.OthersFees{
					OtherFees: func() *models.FeeDetail {
						if feeTextPtr := utils.ExtractFeePtr(doc, "attr-other", col); feeTextPtr != nil {
							return &models.FeeDetail{
								Text:   *feeTextPtr,
								Amount: parseAmount(*feeTextPtr),
							}
						}
						return nil
					}(),
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
