package company

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type CompanyData struct {
	No              int           `json:"no"`
	Name            string        `json:"name"`
	Website         string        `json:"website"`
	NoSharehold     []Shareholder `json:"no_sharehold"`
	BoardOfDirector []BoardMember `json:"board_of_director"`
	Branches        []Branch      `json:"branches"`
}

func FetchCompanyData(url string) CompanyData {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Received non-200 response code: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	var companyData CompanyData
	companyData.No = ParseCompanyNo(doc)
	companyData.Name = ParseCompanyName(doc)
	companyData.Website = ParseCompanyWebsite(doc)
	companyData.NoSharehold = ParseShareholders(doc)
	companyData.BoardOfDirector = ParseBoardMembers(doc)

	branhURL, exists := doc.Find("div[data-element='element_button_image'] a").Attr("href")
	if exists && branhURL != "" {
		companyData.Branches = FetchBranchData(branhURL)
	}
	return companyData
}

func ParseCompanyNo(doc *goquery.Document) int {
	var noText string

	noText = doc.Find("div.title-member span.text-primary").Text()
	if noText == "" {
		noText = doc.Find("span:contains('Broker')").Parent().Find("span:contains('Number') span").Text()
	}

	if noText == "" {
		doc.Find("span").Each(func(i int, s *goquery.Selection) {
			text := s.Text()
			if strings.Contains(text, "หมายเลข") {
				noText = s.Find("span").Text()
			} else if strings.Contains(text, "Number") {
				noText = s.Find("span").Eq(1).Text()
			}
		})
	}

	noText = strings.TrimSpace(noText)
	if noText == "" {
		log.Println("Company number not found")
		return 0
	}

	re := regexp.MustCompile(`\d+`)
	noStr := re.FindString(noText)

	no, err := strconv.Atoi(noStr)
	if err != nil {
		log.Printf("Error conveting company number: %v, raw text %s", err, noText)
		return 0
	}
	return no
}

func ParseCompanyName(doc *goquery.Document) string {
	name := strings.TrimSpace(doc.Find("h2[data-element='element_heading']").Text())
	if name == "" {
		name = strings.TrimSpace(doc.Find("h2").Text())
	}
	if name == "" {
		log.Println("Company name not found")
	}
	return name
}

func ParseCompanyWebsite(doc *goquery.Document) string {
	website := doc.Find("div[data-element='element_text_editor'] a").AttrOr("href", "")
	if website == "" {
		website = doc.Find("a[href^='http']").AttrOr("href", "")
	}
	if website == "" {
		log.Println("Company website not found")
	}
	return website
}
