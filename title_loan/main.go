package main

import (
	"fmt"
	"log"
	"title_loan/loanproducts"
)

func main() {
	url := "https://app.bot.or.th/1213/MCPD/ProductApp/TitleLoan/CompareProductList"
	initialpayload := `{"ProductIdList":"84,167,128,146,139,140,227,327,379,138,277,3,185,166,137,351,77,330,135,136,386,395,105,104,92,103,102,91,402,276,82,278,381,380,375,39,325,155,156,160,152,154,274,235,376,230,229,33,115,176,243,286,56,83,201,256,246,250,374,251,163,134,237,252","Page":1,"Limit":3}`

	firstPageBody, err := loanproducts.SendHTTPRequest(url, initialpayload)
	if err != nil {
		log.Fatalf("Error fetching first page: %v", err)
	}

	totalPages := loanproducts.DetermineTotalPages(firstPageBody)
	if totalPages == 0 {
		log.Fatalf("Could not determine the total number of pages")
	}

	var allTitleLoans []loanproducts.TitleLoan

	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Procssing page: %d/%d\n", page, totalPages)
		payload := fmt.Sprintf(`{"ProductIdList":"84,167,128,146,139,140,227,327,379,138,277,3,185,166,137,351,77,330,135,136,386,395,105,104,92,103,102,91,402,276,82,278,381,380,375,39,325,155,156,160,152,154,274,235,376,230,229,33,115,176,243,286,56,83,201,256,246,250,374,251,163,134,237,252","Page":%d,"Limit":3}`, page)
		titleLoans, err := loanproducts.FetchLoanProducts(url, payload)
		if err != nil {
			log.Printf("Error fetching loan products for page %d: %v", page, err)
			continue
		}
		allTitleLoans = append(allTitleLoans, titleLoans...)
	}

	err = loanproducts.WriteJSONToFile(allTitleLoans, "title_loan.json")
	if err != nil {
		log.Fatalf("Error writing JSON to file: %v", err)
	}

	fmt.Println("Data successfully written to title_loan.json")
}
