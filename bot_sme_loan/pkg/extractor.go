package pkg

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FetchHTML(url, payload string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "text/plain, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Cookie", "verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT=\"z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0\"; _uid6672=16B5DEBD.32; _ctout6672=1; _ga_NLQFGWVNXN=GS1.1.1724735144.39.1.1724735219.59.0.0")
	req.Header.Set("Origin", "https://app.bot.or.th")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://app.bot.or.th/1213/MCPD/ProductApp/SME/CompareProduct")
	req.Header.Set("Sec-CH-UA", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	req.Header.Set("VerificationToken", "7vph-1q3GpjGrBZRY41Pp71PlOYj1qPmVifHjS3waHWhWkbZgkV3XHy9QyQmHY69AmCP8rFlPTDozoA6hRM9Gu1EwPK3iT-3t5HcnQk-W-41,UYKf7DDVmHWVgCNfGsB1wxlbCmIaLrDenx5PxGy1-ZWJEQW1NPMZuhL5QH_0JzwotfFobL6XNkHjYLSUlfszH39nFMeqh7NB9yjL6nrvkZ81")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func ExtractProducts(doc *goquery.Document, count int) []Product {
	var products []Product

	for i := 1; i <= count; i++ {
		col := "col" + strconv.Itoa(i)

		serviceProvider := cleanString(doc.Find("th.col-s.col-s-" + strconv.Itoa(i) + " span").Last().Text())
		product := cleanString(doc.Find("th.font-black.text-center.prod-" + col + " span").Text())

		interestRatePerYear := extractInterestRates(doc, col)
		defaultInterest := cleanString(doc.Find(".attr-DefaultInterestRate ." + col + " span").Text())

		creditLineType := cleanString(doc.Find(".attr-CreditLineType ." + col + " span").Text())
		collateral := extractList(doc.Find(".attr-Collateral ."+col+" span").Text(), "/")
		productConditions := extractList(doc.Find(".attr-ProductCondition ."+col+" span").Text(), "1.")
		borrowerAge := cleanString(doc.Find(".attr-BorrowerAge ." + col + " span").Text())
		applicationConditions := extractList(doc.Find(".attr-ApplicationCondition ."+col+" span").Text(), "-")

		creditLimit := cleanString(doc.Find(".attr-CreditLimit ." + col + " span").Text())
		creditLimitConditions := extractList(doc.Find(".attr-CreditLimitCondition ."+col+" span").Text(), "-")
		borrowingPeriod := extractList(doc.Find(".attr-BorrowingPeriod ."+col+" span").Text(), "-")
		borrowingPeriodConditions := cleanString(doc.Find(".attr-ConditionOfBorrowingPeriod ." + col + " span").Text())

		var borrowingPeriodConditionsInterface interface{}
		if borrowingPeriodConditions == "" {
			borrowingPeriodConditionsInterface = nil
		} else {
			borrowingPeriodConditionsInterface = borrowingPeriodConditions
		}

		products = append(products, Product{
			ServiceProvider: serviceProvider,
			Product:         product,
			InterestRates: InterestRates{
				InterestRatesPerYear: interestRatePerYear,
				DefaultInterest:      defaultInterest,
			},
			ProductDetails: ProductCondition{
				CreditLineType:        creditLineType,
				Collateral:            collateral,
				ProductConditions:     productConditions,
				BorrowerAge:           borrowerAge,
				ApplicationConditions: applicationConditions,
			},
			CreditTerms: CreditAndLoanTerms{
				CreditLimit:               creditLimit,
				CreditLimitConditions:     creditLimitConditions,
				BorrowingPeriod:           borrowingPeriod,
				BorrowingPeriodConditions: borrowingPeriodConditionsInterface,
			},
			Fees:           extractFees(doc, col),
			AdditionalInfo: extractAdditionalInfo(doc, col),
		})
	}

	return products
}

func cleanString(text string) string {
	return strings.TrimSpace(strings.ReplaceAll(text, "\n", ""))
}

func extractInterestRates(doc *goquery.Document, col string) []string {
	interestRatesText := cleanString(doc.Find(".attr-InterestRatePerYear ." + col + " span").Text())
	interestRates := strings.Split(interestRatesText, "ปี")
	var interestRatePerYear []string
	for _, rate := range interestRates {
		rate = cleanString("ปี" + rate) // Reconstruct the string to include "ปี"
		if rate != "ปี" && rate != "" { // Filter out unwanted entries
			interestRatePerYear = append(interestRatePerYear, rate)
		}
	}
	return interestRatePerYear
}

func extractList(text string, delimiter string) []string {
	parts := strings.Split(text, delimiter)
	var cleanedParts []string
	for _, part := range parts {
		part = cleanString(part)
		if part != "" {
			cleanedParts = append(cleanedParts, part)
		}
	}
	return cleanedParts
}

func extractFees(doc *goquery.Document, col string) Fees {
	prepaymentFeeText := doc.Find(fmt.Sprintf(".attr-PrepaymentFeeRate .%s span", col)).Text()
	prepaymentFee := strings.Split(prepaymentFeeText, "-")
	spaceReplacer := regexp.MustCompile(`\s+`)

	for i := range prepaymentFee {
		prepaymentFee[i] = strings.ReplaceAll(prepaymentFee[i], "\n", "")
		prepaymentFee[i] = spaceReplacer.ReplaceAllString(prepaymentFee[i], " ")
		prepaymentFee[i] = strings.TrimSpace(prepaymentFee[i])
	}

	externalAppraisalFeeText := doc.Find(fmt.Sprintf(".attr-SurveyAndAppraisalFeeByExternal .%s span", col)).Text()
	externalAppraisalFee := strings.Split(externalAppraisalFeeText, "เงื่อนไข:")
	for i := range externalAppraisalFee {
		externalAppraisalFee[i] = strings.ReplaceAll(externalAppraisalFee[i], "\n", "")
		externalAppraisalFee[i] = strings.TrimSpace(externalAppraisalFee[i])
	}

	otherFeesText := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-OtherFee .%s span", col)).Text())
	var otherFees interface{}
	if otherFeesText == "" {
		otherFees = nil
	} else {
		otherFees = strings.ReplaceAll(otherFeesText, "\n", "")
		otherFees = otherFeesText
	}

	fees := Fees{
		FrontEndFee:           strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-FrontEndFeeRate .%s span", col)).Text()),
		ManagementFee:         strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-ManagementFeeRate .%s span", col)).Text()),
		CommitmentFee:         strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-CommitmentFeeRate .%s span", col)).Text()),
		CancellationFee:       strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-CancellationFeeRate .%s span", col)).Text()),
		PrepaymentFee:         prepaymentFee,
		ExtensionFee:          strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-ExtensionFeeRate .%s span", col)).Text()),
		AnnualFee:             strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-AnnualFeeRate .%s span", col)).Text()),
		InternalAppraisalFee:  strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-SurveyAndAppraisalFeeByInternal .%s span", col)).Text()),
		ExternalAppraisalFee:  externalAppraisalFee,
		DebtCollectionFee:     strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-DebtCollectionFee .%s span", col)).Text()),
		CreditCheckFee:        strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-CreditBureauFee .%s span", col)).Text()),
		StatementReissuingFee: strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-StatementReIssuingFee .%s span", col)).Text()),
		OtherFees:             otherFees, // Assume other fees are null for now
	}
	return fees
}

func extractAdditionalInfo(doc *goquery.Document, col string) AdditionalInfo {
	productWebsite := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-uRL .%s a", col)).AttrOr("href", ""))
	feeWebsite := strings.TrimSpace(doc.Find(fmt.Sprintf(".attr-FeeURL .%s a", col)).AttrOr("href", ""))

	var productsWebsitePtr, feeWensitePtr *string

	if productWebsite == "" {
		productsWebsitePtr = nil
	} else {
		productsWebsitePtr = &productWebsite
	}

	if feeWebsite == "" {
		feeWensitePtr = nil
	} else {
		feeWensitePtr = &feeWebsite
	}

	additionaInfo := AdditionalInfo{
		ProductWebsite: productsWebsitePtr,
		FeeWebsite:     feeWensitePtr,
	}

	return additionaInfo
}

func DetermineTotalPage(body []byte) int {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Printf("Error loading HTML to determine total pages: %v\n", err)
		return 0
	}

	totalPages := 1

	// Log the entire HTML for debugging
	fmt.Println("HTML content received:")
	doc.Find("html").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})

	// Check pagination
	doc.Find("ul.pagination li a").Each(func(i int, s *goquery.Selection) {
		pageNum, exists := s.Attr("data-page")
		if exists {
			fmt.Printf("Found page link: %s\n", pageNum)
			page, err := strconv.Atoi(pageNum)
			if err == nil && page > totalPages {
				totalPages = page
			}
		}
	})

	return totalPages
}
