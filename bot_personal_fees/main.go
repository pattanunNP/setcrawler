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

type PersonalLoanDetails struct {
	ServiceProvider string         `json:"service_provider"`
	Product         string         `json:"product"`
	GeneralFees     GeneralFees    `json:"general_fees"`
	RevolvingFees   RevolvingFees  `json:"revolving_fees"`
	PaymentFees     PaymentFees    `json:"payment_fees"`
	OtherFees       OtherFees      `json:"other_fees"`
	AdditionalInfo  AdditionalInfo `json:"additional_info"`
}

type FeeDetail struct {
	Text       string   `json:"text"`
	Fee        *float64 `json:"fee"`
	Conditions string   `json:"conditions"`
}

type GeneralFees struct {
	DefaultInterest    FeeDetail `json:"default_interest"`
	DebtCollection     FeeDetail `json:"debt_collection"`
	Prepayment         FeeDetail `json:"prepayment"`
	CreditBureauCheck  FeeDetail `json:"credit_bureau_check"`
	StampDuty          FeeDetail `json:"stamp_duty"`
	ChequeReturn       FeeDetail `json:"cheque_return"`
	InsufficientFunds  FeeDetail `json:"insufficient_funds"`
	StatementReissue   FeeDetail `json:"statement_reissue"`
	TransactionInquiry FeeDetail `json:"transaction_inquiry"`
}

type RevolvingFees struct {
	CardFee         FeeDetail `json:"card_fee"`
	CardReplacement FeeDetail `json:"card_replacement"`
	PINReplacement  FeeDetail `json:"pin_replacement"`
	CurrencyRisk    FeeDetail `json:"currency_risk"`
}

type PaymentFees struct {
	FreePaymentMethods []string  `json:"free_payment_methods"`
	DebitFromProvider  FeeDetail `json:"debit_from_provider"`
	DebitFromOthers    FeeDetail `json:"debit_from_others"`
	BankBranch         FeeDetail `json:"bank_branch"`
	OtherBranch        FeeDetail `json:"other_branch"`
	CounterService     FeeDetail `json:"counter_service"`
	OnlinePayment      FeeDetail `json:"online_payment"`
	ATM_CDM            FeeDetail `json:"atm_cdm"`
	PhonePayment       FeeDetail `json:"phone_payment"`
	ChequePayment      FeeDetail `json:"cheque_payment"`
	OtherMethods       FeeDetail `json:"other_methods"`
}

type OtherFees struct {
	OtherFeeDetails []string `json:"other_fee_details"`
}

type AdditionalInfo struct {
	FeeWebsite string `json:"fee_website"`
}

func main() {
	// Load initial translations (can be from a file or predefined)
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/personalloanFee/CompareProductList"
	// payloadTemplate := `{"ProductIdList":"249-1,244-1,245-1,247-1,246-1,251-1,10819-1,487-1,489-1,10478-2,376-1,272-2,11701-1,10978-1,10537-2,11678-1,278-2,279-2,375-1,11917-1,121-1,273-1,10796-1,10797-1,11452-1,11096-1,11699-1,11698-1,10818-1,10574-1,11625-1,10913-1,11700-1,10548-1,10613-1,11036-1,10701-1,11177-1,11928-1,11930-1,11019-1,10653-1,10655-1,442-1,443-1,312-2,10702-1,10411-2,11378-1,11377-1,11340-1,441-1,11601-1,11604-1,11338-1,11914-1,348-1,303-1,10736-2,129-1,10713-2,10337-2,10524-1,10336-2,10525-1,11562-2,11697-1,11535-1,11534-1,11536-1,11926-1,11715-1,11718-1,11720-1,389-1,11173-2,11550-2,10743-2,10744-2,11039-1,11041-1,10541-1,11358-2,11913-1,11231-2,11230-2,11232-2,11621-1,11566-2,63-2,11543-2,11544-2,126-1,11932-1,11542-2,11912-1,11808-1,11816-1,11535-2,11903-1,11517-2,11887-1,11872-1,11496-2,11873-1,11866-1,11519-2,11514-2,11863-1,11882-1,11488-2,11862-1,11510-2,292-2,231-1,10720-1,382-1,386-1,293-2,10817-1,10542-1,229-2,11561-2,230-1,140-1,103-1,102-1,11323-2,11324-2,11666-1,11335-2,11336-2,10722-1,10719-1,10718-1,10721-1,385-1,10816-1,10910-2,10911-2,10916-2,10915-2,11274-1,11325-2,11444-2,164-2,250-1,243-1,163-2,248-1,159-1","Page":%d,"Limit":3}`
	payloadTemplate := `{"ProductIdList":"249-1,244-1,245-1,247-1,246-1,251-1,10819-1,487-1,489-1,10478-2,376-1,272-2,11701-1,10978-1,10537-2,11678-1,278-2,279-2,375-1,11917-1,121-1,273-1,10796-1,10797-1,11452-1,11096-1,11699-1,11698-1,10818-1,10574-1,11625-1,10913-1,11700-1,10548-1,10613-1,11036-1","Page":%d,"Limit":3}`

	// Make the HTTP request
	initialPayload := fmt.Sprintf(payloadTemplate, 1)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(initialPayload)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set the headers
	setHeaders(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	totalPages := DetermineTotalPage(doc)
	if totalPages == 0 {
		log.Fatal("Could not determine the total number of pages")
	}

	// Extract data
	var details []PersonalLoanDetails
	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages)
		payload := fmt.Sprintf(payloadTemplate, page)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			log.Printf("Error creating request for page %d: %v", page, err)
			continue
		}

		setHeaders(req)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Request failed with status: %v", resp.Status)
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			log.Fatalf("Error parsing HTML: %v", err)
		}

		for i := 1; i <= 3; i++ {
			col := fmt.Sprintf("col%d", i)
			serviceProvider := cleanText(doc.Find(fmt.Sprintf("th.col-s.col-s-%d span", i)).Last().Text())
			product := cleanText(doc.Find(fmt.Sprintf("th.font-black.text-center.prod-%s span", col)).Text())

			generalFees := GeneralFees{
				DefaultInterest:    extractFeeDetail(doc, fmt.Sprintf(".attr-DfltIntRateFee .cmpr-col.%s span", col)),
				DebtCollection:     extractFeeDetail(doc, fmt.Sprintf(".attr-DebtCollectionFee .cmpr-col.%s span", col)),
				Prepayment:         extractFeeDetail(doc, fmt.Sprintf(".attr-PrepaymentFee .cmpr-col.%s span", col)),
				CreditBureauCheck:  extractFeeDetail(doc, fmt.Sprintf(".attr-CreditBureauFee .cmpr-col.%s span", col)),
				StampDuty:          extractFeeDetail(doc, fmt.Sprintf(".attr-DutyStampFee .cmpr-col.%s span", col)),
				ChequeReturn:       extractFeeDetail(doc, fmt.Sprintf(".attr-ChequeReturnedFee .cmpr-col.%s span", col)),
				InsufficientFunds:  extractFeeDetail(doc, fmt.Sprintf(".attr-InsufficientDirectDebitFee .cmpr-col.%s span", col)),
				StatementReissue:   extractFeeDetail(doc, fmt.Sprintf(".attr-CopyStatementReissuing .cmpr-col.%s span", col)),
				TransactionInquiry: extractFeeDetail(doc, fmt.Sprintf(".attr-TransactionVerificationFee .cmpr-col.%s span", col)),
			}

			revolvingFees := RevolvingFees{
				CardFee:         extractFeeDetail(doc, fmt.Sprintf(".attr-CardHolderAnnualFee .cmpr-col.%s span", col)),
				CardReplacement: extractFeeDetail(doc, fmt.Sprintf(".attr-CardReplacementFee .cmpr-col.%s span", col)),
				PINReplacement:  extractFeeDetail(doc, fmt.Sprintf(".attr-CardPINReplacement .cmpr-col.%s span", col)),
				CurrencyRisk:    extractFeeDetail(doc, fmt.Sprintf(".attr-FXRiskCost .cmpr-col.%s span", col)),
			}

			freePaymentMethods := extractFreePaymentMethods(doc, fmt.Sprintf(".attr-FreePaymentMethod .cmpr-col.%s span", col))

			paymentFees := PaymentFees{
				FreePaymentMethods: freePaymentMethods,
				DebitFromProvider:  extractFeeDetail(doc, fmt.Sprintf(".attr-DirectDebitFromAccountFee .cmpr-col.%s span", col)),
				DebitFromOthers:    extractFeeDetail(doc, fmt.Sprintf(".attr-DirectDebitFromAccountFeeOther .cmpr-col.%s span", col)),
				BankBranch:         extractFeeDetail(doc, fmt.Sprintf(".attr-BankCounterServiceFee .cmpr-col.%s span", col)),
				OtherBranch:        extractFeeDetail(doc, fmt.Sprintf(".attr-BankCounterServiceFeeOther .cmpr-col.%s span", col)),
				CounterService:     extractFeeDetail(doc, fmt.Sprintf(".attr-CounterServiceFee .cmpr-col.%s span", col)),
				OnlinePayment:      extractFeeDetail(doc, fmt.Sprintf(".attr-PaymentOnlineFee .cmpr-col.%s span", col)),
				ATM_CDM:            extractFeeDetail(doc, fmt.Sprintf(".attr-PaymentCDMATMFee .cmpr-col.%s span", col)),
				PhonePayment:       extractFeeDetail(doc, fmt.Sprintf(".attr-PaymentPhoneFee .cmpr-col.%s span", col)),
				ChequePayment:      extractFeeDetail(doc, fmt.Sprintf(".attr-PaymentChequeOrMoneyOrderFee .cmpr-col.%s span", col)),
				OtherMethods:       extractFeeDetail(doc, fmt.Sprintf(".attr-PaymentOtherChannelFee .cmpr-col.%s span", col)),
			}

			otherFeesDetails := extractOtherFeeDetails(doc, fmt.Sprintf(".attr-other .cmpr-col.%s span", col))

			otherFees := OtherFees{
				OtherFeeDetails: otherFeesDetails,
			}

			additionalInfo := AdditionalInfo{
				FeeWebsite: doc.Find(fmt.Sprintf(".attr-Feeurl .cmpr-col.%s a", col)).First().AttrOr("href", ""),
			}

			details = append(details, PersonalLoanDetails{
				ServiceProvider: serviceProvider,
				Product:         product,
				GeneralFees:     generalFees,
				RevolvingFees:   revolvingFees,
				PaymentFees:     paymentFees,
				OtherFees:       otherFees,
				AdditionalInfo:  additionalInfo,
			})
		}
		time.Sleep(2 * time.Second)
	}

	// Save to JSON file
	file, err := os.Create("personalLoanFees.json")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(details); err != nil {
		log.Fatalf("Error encoding JSON to file: %v", err)
	}

	fmt.Println("Data saved to personalLoanFees.json")
}

// Helper functions

func setHeaders(req *http.Request) {
	req.Header.Set("Accept", "text/plain, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Cookie", `verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95@AdobeOrg=1; AMCV_F915091E62ED182D0A495F95@AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0"; _uid6672=16B5DEBD.54; _ctout6672=1; visit_time=13; _ga_NLQFGWVNXN=GS1.1.1725005721.61.1.1725008703.41.0.0`)
	req.Header.Set("Origin", "https://app.bot.or.th")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://app.bot.or.th/1213/MCPD/FeeApp/OtherFee/SafeDepositBoxServiceFee/CompareProduct")
	req.Header.Set("Sec-CH-UA", `"Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome";v="128"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	req.Header.Set("VerificationToken", "wOtMHHCbxVnBmyxoDil8c2CVVE-xPABkb54TQGyS_sTucouzC9oXom7guar3k56sz2tupUOoXZcWMLRLVnBgkTMo92r-UGVvnIRHI4OOjuQ1,d2wxcDXHhYSdYYlu31OZKUYs2IJ4Ig8xWtqVIvH7YsOJK8Bd4hfrA27l98vYRdp4aKPL1AUUxsHmUnz0KZ3zCi2vbuow7BuTpWHRLGye8Sg1")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
}

func cleanText(text string) string {
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, ",", "")
	text = strings.Join(strings.Fields(text), " ")
	return text
}

func extractFeeDetail(doc *goquery.Document, selector string) FeeDetail {
	text := cleanText(doc.Find(selector).Text())
	fee, conditions := parseFeeText(text)
	return FeeDetail{Text: text, Fee: fee, Conditions: conditions}
}

func parseFeeText(text string) (*float64, string) {
	if strings.Contains(text, "ไม่มีค่าธรรมเนียม") {
		fee := 0.0
		return &fee, "None"
	}
	if strings.Contains(text, "ไม่มีบริการ") {
		return nil, "No service"
	}
	// Extract numerical fee
	feeStr := strings.Fields(text)[0]
	fee, err := strconv.ParseFloat(feeStr, 64)
	if err != nil {
		return nil, "Error"
	}
	return &fee, ""
}

func extractFreePaymentMethods(doc *goquery.Document, selector string) []string {
	text := cleanText(doc.Find(selector).Text())
	methods := filterEmptyStrings(strings.Split(text, "-"))
	return methods
}

func extractOtherFeeDetails(doc *goquery.Document, selector string) []string {
	text := cleanText(doc.Find(selector).Text())
	parts := strings.Split(text, "1.")
	return filterEmptyStrings(parts)
}

func filterEmptyStrings(input []string) []string {
	var output []string
	for _, str := range input {
		trimmed := strings.TrimSpace(str)
		if trimmed != "" {
			output = append(output, trimmed)
		}
	}
	return output
}

func DetermineTotalPage(doc *goquery.Document) int {
	totalPages := 1

	doc.Find("ul.pagination li a").Each(func(i int, s *goquery.Selection) {
		pagenum, exists := s.Attr("data-page")
		if exists {
			page, err := strconv.Atoi(pagenum)
			if err == nil && page > totalPages {
				totalPages = page
			}
		}
	})
	return totalPages
}
