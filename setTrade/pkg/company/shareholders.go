package company

import (
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Shareholder struct {
	No      int     `json:"no"`
	Name    string  `json:"name"`
	Percent float64 `json:"percent"`
}

func ParseShareholders(doc *goquery.Document) []Shareholder {
	var shareholders []Shareholder

	// Handle both types of tables
	doc.Find("table[role='table']").First().Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		var shareholder Shareholder
		shareholder.No = i + 1
		shareholder.Name = strings.TrimSpace(s.Find("td").Eq(1).Text())
		percentStr := strings.TrimSpace(s.Find("td").Eq(2).Text())
		percentStr = strings.Replace(percentStr, ",", "", -1)
		percentStr = strings.TrimSuffix(percentStr, "%")
		shareholder.Percent, _ = strconv.ParseFloat(percentStr, 64)
		shareholders = append(shareholders, shareholder)
	})

	doc.Find("div[data-element='element_text_editor'] table").First().Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		var shareholder Shareholder
		shareholder.No = i + 1
		shareholder.Name = strings.TrimSpace(s.Find("td").Eq(1).Text())
		percentStr := strings.TrimSpace(s.Find("td").Eq(2).Text())
		percentStr = strings.Replace(percentStr, ",", "", -1)
		percentStr = strings.TrimSuffix(percentStr, "%")
		shareholder.Percent, _ = strconv.ParseFloat(percentStr, 64)
		shareholders = append(shareholders, shareholder)
	})

	if len(shareholders) == 0 {
		log.Println("No shareholders found")
	}

	return shareholders
}
