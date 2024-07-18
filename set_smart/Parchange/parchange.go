package parchange

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type TableData struct {
	ParChange []ParChange `json:"par_change"`
}

type ParChange struct {
	Symbol        string `json:"symbol"`
	SecurityType  string `json:"security_type"`
	EffectiveDate string `json:"effective_date"`
	BoardDate     string `json:"board_date"`
	AnnounceDate  string `json:"announce_date"`
	OldPar        string `json:"old_par"`
	NewPar        string `json:"new_par"`
	ChangeParType string `json:"change_par_type"`
}

func MakeParChangeRequest(cookieStr, url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Cookie", cookieStr)
	req.Header.Set("Origin", "https://www.setsmart.com")
	req.Header.Set("Referer", "https://www.setsmart.com/ism/parchange.html")
	req.Header.Set("Sec-CH-UA", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", "macOS")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil

}

func extractParChange(doc *goquery.Document) *TableData {
	var data TableData
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		var parChange ParChange
		validRow := true
		s.Find("td").Each(func(j int, td *goquery.Selection) {
			switch j {
			case 0:
				parChange.Symbol = strings.TrimSpace(td.Text())
			case 1:
				parChange.SecurityType = strings.TrimSpace(td.Text())
			case 2:
				parChange.EffectiveDate = parseDateToISO8601(strings.TrimSpace(td.Text()))
			case 3:
				parChange.BoardDate = parseDateToISO8601(strings.TrimSpace(td.Text()))
			case 4:
				parChange.AnnounceDate = parseDateToISO8601(strings.TrimSpace(td.Text()))
			case 5:
				parChange.OldPar = strings.TrimSpace(td.Text())
			case 6:
				parChange.NewPar = strings.TrimSpace(td.Text())
			case 7:
				parChange.ChangeParType = strings.TrimSpace(td.Text())
			}
		})
		if parChange.Symbol == "" || parChange.SecurityType == "" || parChange.EffectiveDate == "" ||
			parChange.BoardDate == "" || parChange.AnnounceDate == "" || parChange.OldPar == "" ||
			parChange.NewPar == "" || parChange.ChangeParType == "" {
			validRow = false
		}
		if validRow {
			data.ParChange = append(data.ParChange, parChange)
		}
	})
	return &data
}

func parseDateToISO8601(dateStr string) string {
	layout := "02/01/2006 15:04"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		layout = "02/01/2006"
		t, err = time.Parse(layout, dateStr)
		if err != nil {
			return dateStr
		}
	}
	return t.Format(time.RFC3339)
}

func collectPageLinks(doc *goquery.Document) []string {
	var links []string
	doc.Find("span.pagelinks a.olink").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			links = append(links, "https://www.setsmart.com"+href)
		}
	})
	return links
}

func FetchAllPages(cookieStr, url string) (*TableData, error) {
	var allParChanges TableData
	htmlContent, err := MakeParChangeRequest(cookieStr, url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	parChanges := extractParChange(doc)
	allParChanges.ParChange = append(allParChanges.ParChange, parChanges.ParChange...)

	links := collectPageLinks(doc)

	for _, link := range links {
		htmlContent, err := MakeParChangeRequest(cookieStr, link)
		if err != nil {
			return nil, err
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
		if err != nil {
			return nil, fmt.Errorf("error parsing HTML: %w", err)
		}

		parChanges := extractParChange(doc)
		allParChanges.ParChange = append(allParChanges.ParChange, parChanges.ParChange...)
	}
	return &allParChanges, nil
}

func GetParChange(cookieStr, market, condition, quickPeriod, beginDate, endDate string) (*TableData, error) {
	url := fmt.Sprintf("https://www.setsmart.com/ism/parchange.html?market=%s&lstCondition=%s&quickPeriod=%s&lstPeriod=D&showBeginDate=%s&beginDate=%s&showEndDate=%s&endDate=%s", market, condition, quickPeriod, beginDate, beginDate, endDate, endDate)
	return FetchAllPages(cookieStr, url)
}

func SaveToFile(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling data: %w", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}
