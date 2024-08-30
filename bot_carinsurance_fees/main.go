package main

import (
	"bytes"
	"carinsurance/models"
	"carinsurance/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// URL and request payload
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/TitleLoanFee/CompareProductList"
	// payloadTemplate := `{"ProductIdList":"128,375,115,166,246,278,103,102,83,374,243,167,325,134,39,135,136,235,237,402,386,395,138,139,140,351,33,277,274,408,163,327,156,160,152,154,155,185,252,250,251,227,146,82,176,376,137,77,84,256,92,91,330,105,104,380,381,201,3,379,276,286,56,230,229","Page":%d,"Limit":3}`
	payloadTemplate := `{"ProductIdList":"128,375,115,166,246,278,103,102,83,374,243,167,325,134,39,135,136,235,237,402,386,395,138,139,140,351","Page":%d,"Limit":3}`

	// Set up HTTP request
	initialPayload := fmt.Sprintf(payloadTemplate, 1)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(initialPayload)))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	utils.AddHeader(req)

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	totalPages := utils.DetermineTotalPage(doc)
	if totalPages == 0 {
		log.Fatal("Could not determine the total number of pages")
	}

	// Extract data
	var carInsuranceDetails []models.CarInsuranceDetails
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
			// Find service provider name using adjusted selector
			serviceProvider := utils.CleanText(doc.Find(fmt.Sprintf("th.col-s.col-s-%d span", i)).Text())

			// Find product description using adjusted selector
			product := utils.CleanText(doc.Find(fmt.Sprintf("th.font-black.text-center.prod-col%d span", i)).Text())

			// Extract detailed fees
			debtCollectionFeeText := utils.CleanText(doc.Find(fmt.Sprintf(".attr-DebtCollectionFee .col%d span", i)).Text())
			stampDutyFeeText := utils.CleanText(doc.Find(fmt.Sprintf(".attr-StampDutyFee .col%d span", i)).Text())

			generalFees := models.GeneralFees{
				LatePaymentInterest:     utils.CleanText(doc.Find(fmt.Sprintf(".attr-DefaultInterestRate .col%d span", i)).Text()),
				DebtCollectionFee:       utils.SplitText(debtCollectionFeeText, "-"),
				StampDutyFee:            utils.SplitText(stampDutyFeeText, "-"),
				ChequeReturnFee:         utils.CleanText(doc.Find(fmt.Sprintf(".attr-ReturnedCheque .col%d span", i)).Text()),
				DebtCollectionFeeValues: utils.ExtractNumbersFromText(debtCollectionFeeText),
				StampDutyFeeValues:      utils.ExtractFloatNumbersFromText(stampDutyFeeText),
			}

			cardFees := models.CardFees{
				CardFee:            utils.CleanText(doc.Find(fmt.Sprintf(".attr-CardFee .col%d span", i)).Text()),
				CardReplacementFee: utils.CleanText(doc.Find(fmt.Sprintf(".attr-CardReplacementFee .col%d span", i)).Text()),
			}

			creditWithdrawalFee := utils.CleanText(doc.Find(fmt.Sprintf(".attr-CreditWithdrawalFee .col%d span", i)).Text())
			if creditWithdrawalFee == "" {
				cardFees.CreditWithdrawalFee = nil
			} else {
				cardFees.CreditWithdrawalFee = &creditWithdrawalFee
			}

			paymentFees := models.PaymentFees{
				FreePaymentChannels:          utils.SplitText(utils.CleanText(doc.Find(fmt.Sprintf(".attr-FreePaymentChannel .col%d span", i)).Text()), "/"),
				ProviderAccountDeductionFee:  utils.CleanText(doc.Find(fmt.Sprintf(".attr-DeductingFromBankACFee .col%d span", i)).Text()),
				OtherProviderAccountFee:      utils.CleanText(doc.Find(fmt.Sprintf(".attr-DeductingFromOtherBankACFee .col%d span", i)).Text()),
				ServiceProviderBranchFee:     utils.CleanText(doc.Find(fmt.Sprintf(".attr-ServiceProviderCounter .col%d span", i)).Text()),
				OtherBranchFee:               utils.CleanText(doc.Find(fmt.Sprintf(".attr-OtherProviderCounter .col%d span", i)).Text()),
				ServiceCounterFee:            utils.CleanText(doc.Find(fmt.Sprintf(".attr-OthersPaymentCounter .col%d span", i)).Text()),
				OnlinePaymentFee:             utils.SplitText(utils.CleanText(doc.Find(fmt.Sprintf(".attr-OnlinePaymentFee .col%d span", i)).Text()), "-"),
				CDMATMPaymentFee:             utils.CleanText(doc.Find(fmt.Sprintf(".attr-CDMATMPaymentFee .col%d span", i)).Text()),
				TelephonePaymentFee:          utils.CleanText(doc.Find(fmt.Sprintf(".attr-PhonePaymentFee .col%d span", i)).Text()),
				ChequeOrMoneyOrderPaymentFee: utils.CleanText(doc.Find(fmt.Sprintf(".attr-ChequeMoneyOrderPaymentFee .col%d span", i)).Text()),
				OtherPaymentChannelsFee:      utils.CleanText(doc.Find(fmt.Sprintf(".attr-OtherChannelPaymentFee .col%d span", i)).Text()),
			}

			otherFees := models.OtherFees{
				LawyerFeeLitigation: utils.SplitText(utils.CleanText(doc.Find(fmt.Sprintf(".attr-lawyerFeeInCaseOfLitigation .col%d span", i)).Text()), "-"),
			}

			otherFee := utils.CleanText(doc.Find(fmt.Sprintf(".attr-OtherFee .col%d span", i)).Text())
			if otherFee == "" {
				otherFees.OtherFees = nil
			} else {
				otherFees.OtherFees = &otherFee
			}

			additionalInformation := models.AdditionalInformation{}
			feeWebsiteLink := utils.CleanText(doc.Find(fmt.Sprintf(".attr-FeeURL .col%d a", i)).AttrOr("href", ""))
			if feeWebsiteLink == "" {
				additionalInformation.FeeWebsiteLink = nil
			} else {
				additionalInformation.FeeWebsiteLink = &feeWebsiteLink
			}

			// Create a new instance of CarInsuranceDetails
			detail := models.CarInsuranceDetails{
				ServiceProvider:       serviceProvider,
				Product:               product,
				GeneralFees:           generalFees,
				CardFees:              cardFees,
				PaymentFees:           paymentFees,
				OtherFees:             otherFees,
				AdditionalInformation: additionalInformation,
			}

			// Populate comparable fields for sorting/filtering
			detail.PopulateComparableFields()

			// Append to the carInsuranceDetails slice
			carInsuranceDetails = append(carInsuranceDetails, detail)
		}
	}

	// Sort carInsuranceDetails by ServiceProvider (example sorting logic)
	models.SortByServiceProvider(carInsuranceDetails)

	// Filter carInsuranceDetails by max late payment interest (example filtering logic)
	filteredDetails := models.FilterByMaxLatePaymentInterest(carInsuranceDetails, 200.0)

	// Convert to JSON
	jsonData, err := json.MarshalIndent(filteredDetails, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Save to file
	if err := os.WriteFile("carinsurance_fees.json", jsonData, 0644); err != nil {
		log.Fatalf("Failed to write JSON to file: %v", err)
	}

	fmt.Println("Data saved to carinsurance_fees.json")
}
