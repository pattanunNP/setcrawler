package management

import (
	"fmt"
	"log"
	"login_token/utils"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ManagementData struct {
	Name       string       `json:"name"`
	Position   string       `json:"position"`
	StartDate  string       `json:"start_date"`
	EndDate    string       `json:"end_date"`
	DetailPage []DetailPage `json:"detail_page"`
}

type DetailPage struct {
	Symbol    string `json:"symbol"`
	Name      string `json:"name"`
	Position  string `json:"position"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func FetchManagementHTML(cookieStr, locale, symbol string) ([]ManagementData, error) {
	baseURL := "https://www.setsmart.com/ism/management.html"

	// Get the current date
	now := time.Now()

	// Determine the correct date format based on the locale
	date := now.Format("02/01/2006")
	if locale == "th_TH" {
		date = convertToBuddhistDate(now)
	}

	// Manually construct the URL with parameters in the correct order
	fullURL := fmt.Sprintf(
		"%s?date=%s&symbol=%s&period=D&lstSort=P&txtLastName=&locale=%s&txtFirstName=&showDate=%s&submit.x=19&submit.y=11",
		baseURL, date, symbol, locale, date,
	)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", locale)
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("cookie", cookieStr)
	req.Header.Set("origin", "https://www.setsmart.com")
	req.Header.Set("referer", "https://www.setsmart.com/ism/management.html")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	log.Printf("Requesting management data from URL: %s\n", fullURL)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var managementData []ManagementData

	doc.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		if row.HasClass("heading-set") || row.Find("td").Length() != 4 {
			return
		}
		var data ManagementData
		row.Find("td").Each(func(j int, cell *goquery.Selection) {
			switch j {
			case 0:
				data.Name = strings.TrimSpace(cell.Text())
			case 1:
				data.Position = strings.TrimSpace(cell.Text())
			case 2:
				data.StartDate = utils.ConvertToISO8601(strings.TrimSpace(cell.Text()))
			case 3:
				data.EndDate = utils.ConvertToISO8601(strings.TrimSpace(cell.Text()))
			}
		})

		link, exists := row.Find("td.table a.glink").Attr("href")
		if exists {
			detailURL := "https://www.setsmart.com" + link
			log.Printf("Fetching detail page: %s", detailURL) // Log the detail page link
			detailPages, err := FetchDetailPage(detailURL, cookieStr)
			if err != nil {
				log.Printf("Error fetching detail page: %v", err)
			} else {
				data.DetailPage = append(data.DetailPage, detailPages...)
			}
		}

		if data.Name != "" && !strings.Contains(data.Name, "function openFavoriteQuery") && data.Name != "คณะกรรมการ / ผู้บริหาร" && data.Name != "Symbol" && data.Name != "Name" && data.Name != "Group by" && data.Name != "As of date" && data.Name != "ชื่อ" {
			managementData = append(managementData, data)
		}
	})

	return managementData, nil
}

func FetchDetailPage(url, cookieStr string) ([]DetailPage, error) {
	var detailPages []DetailPage

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return detailPages, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("cookie", cookieStr)
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return detailPages, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return detailPages, fmt.Errorf("error reading response body: %v", err)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return detailPages, fmt.Errorf("error reading response body: %v", err)
	}

	doc.Find("table.tfont tbody tr").Each(func(i int, row *goquery.Selection) {
		var detailPage DetailPage
		skip := false
		row.Find("td").Each(func(i int, cell *goquery.Selection) {
			switch i {
			case 0:
				detailPage.Symbol = strings.TrimSpace(cell.Text())
			case 1:
				detailPage.Name = strings.TrimSpace(cell.Text())
			case 2:
				detailPage.Position = strings.TrimSpace(cell.Text())
			case 3:
				detailPage.StartDate = utils.ConvertToISO8601(strings.TrimSpace(cell.Text()))
			case 4:
				detailPage.EndDate = utils.ConvertToISO8601(strings.TrimSpace(cell.Text()))
			}

			if detailPage.Symbol == "คณะกรรมการ / ผู้บริหาร" || detailPage.Symbol == "หลักทรัพย์" || detailPage.Name == "ชื่อ" || detailPage.Position == "ตำแหน่ง" || detailPage.StartDate == "วันที่เริ่มต้น" || detailPage.EndDate == "วันที่สื้นสุด" || detailPage.Symbol == "" && detailPage.Name == "" && detailPage.Position == "" && detailPage.StartDate == "" && detailPage.EndDate == "" {
				skip = true
			}
		})
		if !skip {
			detailPages = append(detailPages, detailPage)
		}
	})
	return detailPages, nil
}

func convertToBuddhistDate(t time.Time) string {
	buddhistYear := t.Year() + 543
	return fmt.Sprintf("%02d/%02d/%04d", t.Day(), t.Month(), buddhistYear)
}
