package company

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Branch struct {
	BranchName  string   `json:"branch_name"`
	Address     string   `json:"address"`
	PhoneNumber []string `json:"phone_number"`
	FaxNumber   []string `json:"fax_number"`
}

func FetchBranchData(url string) []Branch {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making request to branch URL: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Received non-200 response code from branch URL: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Error parsing branch HTML: %v", err)
	}

	var branches []Branch

	doc.Find("table.rgMasterTable").Each(func(tableIdx int, table *goquery.Selection) {
		table.Find("tbody tr").Each(func(rowIdx int, row *goquery.Selection) {
			if row.HasClass("rgNoRecords") {
				return
			}
			var branch Branch
			branch.BranchName = strings.TrimSpace(row.Find("td").Eq(1).Text())
			branch.Address = strings.TrimSpace(row.Find("td").Eq(2).Text())

			phoneNumbers := strings.Split(strings.TrimSpace(row.Find("td").Eq(3).Text()), ",")
			for i := range phoneNumbers {
				phoneNumbers[i] = strings.TrimSpace(phoneNumbers[i])
			}
			branch.PhoneNumber = phoneNumbers

			faxNumbers := strings.Split(strings.TrimSpace(row.Find("td").Eq(4).Text()), ",")
			for i := range faxNumbers {
				faxNumbers[i] = strings.TrimSpace(faxNumbers[i])
			}
			branch.FaxNumber = faxNumbers

			if branch.BranchName != "" && branch.Address != "" {
				branches = append(branches, branch)
			}
		})
	})
	return branches
}
