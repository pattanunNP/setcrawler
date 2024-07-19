package MajorShareHolder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ShareHolderDetail struct {
	Rank          string         `json:"rank"`
	Shareholder   string         `json:"shareholder"`
	Shares        string         `json:"shares"`
	PercentShares string         `json:"percent_shares"`
	LinkContent   DesiredContent `json:"link_content"`
}

type MajorShareHolder struct {
	Date   string            `json:"date"`
	Values map[string]string `json:"values"`
}

type DesiredContent struct {
	Symbol          string `json:"symbol"`
	Shareholder     string `json:"shareholder"`
	Shares          string `json:"shares"`
	PercenShares    string `json:"percent_shares"`
	ShareholderAsof string `json:"shareholder_as_of"`
}

type MajorShareHolders struct {
	FreeFloat   []MajorShareHolder `json:"free_float"`
	Shareholder []MajorShareHolder `json:"shareholder"`
}

type CombinedShareHolderData struct {
	Symbol      string              `json:"symbol"`
	FreeFloat   []MajorShareHolder  `json:"free_float"`
	Shareholder []MajorShareHolder  `json:"shareholder"`
	Details     []ShareHolderDetail `json:"details"`
}

func FetchDetailedShareholderData(postURL, cookieStr, formData string) (string, error) {
	req, err := http.NewRequest("POST", postURL, bytes.NewBufferString(formData))
	if err != nil {
		return "", fmt.Errorf("error creating POST request: %w", err)
	}

	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Add("cache-control", "max-age=0")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cookie", cookieStr)
	req.Header.Add("origin", "https://www.setsmart.com")
	req.Header.Add("priority", "u=0, i")
	req.Header.Add("referer", "https://www.setsmart.com/ism/majorshareholder.html?locale=en_US")
	req.Header.Add("sec-ch-ua", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "macOS")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error maing POST requesy: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}

func FetchLinkData(link, cookieStr string) (string, error) {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Cookie", cookieStr)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}

func ExtractOptionValues(htmlStr, selectName string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return nil, fmt.Errorf("error loading HTML document: %w", err)
	}

	var values []string

	doc.Find(fmt.Sprintf("select[name=%s] option", selectName)).Each(func(index int, item *goquery.Selection) {
		value, exists := item.Attr("value")
		if exists {
			values = append(values, value)
		}
	})
	return values, nil
}

func ExtractDataBlock(doc *goquery.Document, selectName string) (map[string]string, error) {
	data := make(map[string]string)

	table := doc.Find("table.tfont")
	if table.Length() == 0 {
		return nil, fmt.Errorf("table with class 'tfont' not found")
	}

	table = table.FilterFunction(func(i int, s *goquery.Selection) bool {
		return s.Find(fmt.Sprintf("select[name='%s']", selectName)).Length() > 0
	})

	if table.Length() == 0 {
		return nil, fmt.Errorf("specific table block not found")
	}

	table.Find("tr").Each(func(index int, row *goquery.Selection) {
		tds := row.Find("td")
		if tds.Length() >= 4 {
			key1 := strings.TrimSpace(tds.Eq(0).Text())
			value1 := strings.TrimSpace(tds.Eq(1).Text())
			key2 := strings.TrimSpace(tds.Eq(2).Text())
			value2 := strings.TrimSpace(tds.Eq(3).Text())
			if key1 != "" && value1 != "" {
				data[key1] = value1
			}
			if key2 != "" && value2 != "" {
				data[key2] = value2
			}
		}
	})

	return data, nil
}

func ExtractShareHolderDetails(doc *goquery.Document, cookieStr string) ([]ShareHolderDetail, error) {
	var details []ShareHolderDetail

	doc.Find("table.tfont").Each(func(index int, table *goquery.Selection) {
		table.Find("tr").Each(func(index int, row *goquery.Selection) {
			tds := row.Find("td")
			if tds.Length() == 4 {
				rank := strings.TrimSpace(tds.Eq(0).Text())
				shareholder := strings.TrimSpace(tds.Eq(1).Text())
				shares := strings.TrimSpace(tds.Eq(2).Text())
				percentShares := strings.TrimSpace(tds.Eq(3).Text())

				if rank == "Minor Shareholders (Free float)" || rank == "Total Shareholders" || rank == "Rank" {
					return
				}

				var linkContent DesiredContent
				if rank != "" && shareholder != "" && shares != "" && percentShares != "" {

					link, exists := tds.Eq(1).Find("a.olink").Attr("href")
					if exists {
						fullLink := fmt.Sprintf("https://www.setsmart.com%s", link)
						fmt.Printf("processing link: %s\n", fullLink)
						linkData, err := FetchLinkData(fullLink, cookieStr)
						if err != nil {
							fmt.Printf("Error fetching link data for %s: %v\n", shareholder, err)
						} else {
							linkDoc, err := goquery.NewDocumentFromReader(strings.NewReader(linkData))
							if err != nil {
								fmt.Printf("Error parsing link data for %s: %v\n", shareholder, err)
							} else {
								linkContent = extractDesiredContent(linkDoc)
							}
						}
					}

					details = append(details, ShareHolderDetail{
						Rank:          rank,
						Shareholder:   shareholder,
						Shares:        shares,
						PercentShares: percentShares,
						LinkContent:   linkContent,
					})
				}
			}
		})
	})

	if len(details) == 0 {
		return nil, fmt.Errorf("no shareholder details found")
	}

	return details, nil
}

func extractDesiredContent(doc *goquery.Document) DesiredContent {
	var content DesiredContent

	doc.Find("table#holder").Each(func(index int, table *goquery.Selection) {
		table.Find("tbody tr").Each(func(index int, row *goquery.Selection) {
			tds := row.Find("td")
			if tds.Length() == 5 {
				content.Symbol = strings.TrimSpace(tds.Eq(0).Text())
				content.Shareholder = strings.TrimSpace(tds.Eq(1).Text())
				content.Shares = strings.TrimSpace(tds.Eq(2).Text())
				content.PercenShares = strings.TrimSpace(tds.Eq(3).Text())
				content.ShareholderAsof = strings.TrimSpace(tds.Eq(4).Text())
			}
		})
	})
	return content
}

func ConvertToISO8601(dateStr string) string {
	layout := "2006-01-02"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return dateStr
	}
	return t.Format(time.RFC3339)
}

func GetMajorShareHoldersAndDetails(cookieStr, symbol, locale string) (CombinedShareHolderData, error) {
	initialURL := "https://www.setsmart.com/ism/majorshareholder.html"
	formData := fmt.Sprintf("txtSymbol=%s&locale=%s&submit.x=0&submit.y=0", symbol, locale)

	htmlStr, err := FetchDetailedShareholderData(initialURL, cookieStr, formData)
	if err != nil {
		return CombinedShareHolderData{}, fmt.Errorf("error fetching initial shareholder data: %w", err)
	}

	freeFloatDates, err := ExtractOptionValues(htmlStr, "lstFreeFloatDate")
	if err != nil {
		return CombinedShareHolderData{}, fmt.Errorf("error extracting option values: %w", err)
	}

	shareHolderDates, err := ExtractOptionValues(htmlStr, "lstDate")
	if err != nil {
		return CombinedShareHolderData{}, fmt.Errorf("error extracting option values: %w", err)
	}

	var combinedData CombinedShareHolderData
	combinedData.Symbol = symbol

	for _, value := range freeFloatDates {
		formData = fmt.Sprintf("radChoice=1&txtSymbol=%s&radShow=2&hidAction=&hidLastContentType=&lstFreeFloatDate=%s", symbol, value)
		response, err := FetchDetailedShareholderData(initialURL, cookieStr, formData)
		if err != nil {
			return CombinedShareHolderData{}, fmt.Errorf("error fetching detailed shareholder data: %w", err)
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(response))
		if err != nil {
			return CombinedShareHolderData{}, fmt.Errorf("error loading HTML document: %w", err)
		}

		shareHolders, err := ExtractDataBlock(doc, "lstFreeFloatDate")
		if err != nil {
			return CombinedShareHolderData{}, fmt.Errorf("error extracting shareholders data: %w", err)
		}

		combinedData.FreeFloat = append(combinedData.FreeFloat, MajorShareHolder{
			Date:   ConvertToISO8601(value),
			Values: shareHolders,
		})
	}

	for _, value := range shareHolderDates {
		formData = fmt.Sprintf("radChoice=1&txtSymbol=%s&radShow=2&hidAction=&hidLastContentType=&lstDate=%s", symbol, value)
		response, err := FetchDetailedShareholderData(initialURL, cookieStr, formData)
		if err != nil {
			return CombinedShareHolderData{}, fmt.Errorf("error fetching detailed shareholder data: %w", err)
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(response))
		if err != nil {
			return CombinedShareHolderData{}, fmt.Errorf("error loading HTML document: %w", err)
		}

		shareHolders, err := ExtractDataBlock(doc, "lstDate")
		if err != nil {
			return CombinedShareHolderData{}, fmt.Errorf("error extracting shareholders data: %w", err)
		}

		combinedData.Shareholder = append(combinedData.Shareholder, MajorShareHolder{
			Date:   ConvertToISO8601(value),
			Values: shareHolders,
		})

		details, err := ExtractShareHolderDetails(doc, cookieStr)
		if err != nil {
			return CombinedShareHolderData{}, fmt.Errorf("error extracting shareholder details: %w", err)
		}

		combinedData.Details = append(combinedData.Details, details...)
	}

	return combinedData, nil
}

func SaveFile(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return fmt.Errorf("error marshalling data: %w", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}
