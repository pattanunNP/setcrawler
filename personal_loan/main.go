package main

import (
	"fmt"
	"log"
	"personal_loan/helpers"
	utils "personal_loan/pkg/client"
	"personal_loan/pkg/model"
	"personal_loan/pkg/scraper"
	"time"
)

func main() {
	url := "https://app.bot.or.th/1213/MCPD/ProductApp/PersonalLoan/CompareProductList"
	// payloadTemplate := `{"ProductIdList":"11699-1,11698-1,11700-1,10911-2,10916-2,11903-1,10525-1,10915-2,10910-2,11621-1,11543-2,11544-2,11041-1,11039-1,11096-1,249-1,245-1,247-1,244-1,11452-1,243-1,163-2,251-1,248-1,10736-2,246-1,250-1,11338-1,63-2,164-2,11340-1,11912-1,102-1,11562-2,11917-1,11697-1,303-1,11550-2,140-1,11913-1,11177-1,11230-2,10542-1,11535-1,11666-1,382-1,11231-2,10574-1,11930-1,11926-1,10797-1,10713-2,11536-1,11514-2,11510-2,11863-1,11862-1,11517-2,11519-2,11887-1,11873-1,11866-1,11488-2,11882-1,11496-2,11534-1,11019-1,11232-2,11678-1,278-2,11932-1,10818-1,11173-2,10702-1,129-1,11378-1,11718-1,11444-2,11566-2,11377-1,10913-1,11601-1,11808-1,11816-1,10653-1,10337-2,292-2,10722-1,10719-1,11928-1,10720-1,10744-2,10743-2,10721-1,10718-1,11625-1,293-2,273-1,375-1,231-1,10817-1,386-1,159-1,10613-1,103-1,10796-1,121-1,10411-2,11036-1,10819-1,10548-1,126-1,11336-2,11335-2,11542-2,11535-2,11323-2,11324-2,11561-2,11914-1,10336-2,10541-1,10701-1,230-1,10816-1,376-1,272-2,312-2,11274-1,442-1,443-1,487-1,11715-1,11720-1,11872-1,10524-1,489-1,279-2,229-2,389-1,385-1,10978-1,11325-2,441-1,11604-1,348-1,10655-1,10478-2,11358-2,10537-2,11701-1","Page":%d,"Limit":3}`
	payloadTemplate := `{"ProductIdList":"11699-1,11698-1,11700-1,10911-2,10916-2,11903-1,10525-1,10915-2,10910-2,11621-1,11543-2","Page":%d,"Limit":3}`

	initialPage := 1
	initialPayload := fmt.Sprintf(payloadTemplate, initialPage)
	initialBody, err := utils.MakePostRequest(url, initialPayload)
	if err != nil {
		log.Fatalf("Failed to make initial HTTP request: %v", err)
	}

	totalPages := scraper.DetermineTotalPages(initialBody)
	if totalPages == 0 {
		log.Fatalf("Failed to detemine the total number of pages")
	}

	var allLoanProducts []model.LoanProduct

	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Procssing page: %d/%d\n", page, totalPages)

		payload := fmt.Sprintf(payloadTemplate, page)
		body, err := utils.MakePostRequest(url, payload)
		if err != nil {
			log.Fatalf("Failed to make HTTP reqeust for page %d: %v", page, err)
		}

		loanProducts, err := scraper.ParseHTML(body)
		if err != nil {
			log.Fatalf("Failed to parse HTML page %d: %v", page, err)
		}

		allLoanProducts = append(allLoanProducts, loanProducts...)

		time.Sleep(2 * time.Second)
	}

	err = helpers.SaveToJSON(allLoanProducts, "personal_loan.json")
	if err != nil {
		log.Fatalf("Failed to save data to JSON file: %v", err)
	}

	fmt.Println("Product details saved to personal_loan.json")
}
