package capitalmovement

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

type CapitalMovement struct {
	AuthorizedShare  string            `json:"Authorized Share (Shares)"`
	PaidUpShares     string            `json:"Paid-up Shares(Shares)"`
	ListedShare      string            `json:"Listed Share (Shares)"`
	ParValue         string            `json:"Par (Baht)"`
	CorporateActions []CorporateAction `json:"Corporate Actions"`
}

type CorporateAction struct {
	Date                string `json:"Date"`
	CorporateActionType string `json:"Corporate Action Type"`
	AuthorizedShare     string `json:"Authorized Share (Shares)"`
	ChangedAuthShare    string `json:"Changed Authorized Share(Shares)"`
	PaidUpShare         string `json:"Paid-Up Share(Shares)"`
	ChangedPaidUpShare  string `json:"Changed Paid-Up Shares(Shares)"`
	ListedShare         string `json:"Listed Share (Shares)"`
	ChangedListedShared string `json:"Changed Listed Share (Shares)"`
	ParValue            string `json:"Par (Baht)"`
}

func MakeCapitalMovementRequest(cookieStr, symbol, locale string) (string, error) {
	data := fmt.Sprintf("symbol=%s&locale=%s&submit.x=0&submit.y=0", symbol, locale)
	requestURL := "https://www.setsmart.com/ism/capitalmovement.html"

	req, err := http.NewRequest("POST", requestURL, strings.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Cookie", cookieStr)
	req.Header.Set("Origin", "https://www.setsmart.com")
	req.Header.Set("Referer", fmt.Sprintf("https://www.setsmart.com/ism/capitalmovement.html?locale=%s", locale))
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

	// Check the response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}

func ExtractCapitalMovement(doc *goquery.Document) (*CapitalMovement, error) {
	var movement CapitalMovement

	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		if i == 6 {
			s.Find("td").Each(func(j int, td *goquery.Selection) {
				switch j {
				case 0:
					movement.AuthorizedShare = strings.TrimSpace(td.Text())
				case 1:
					movement.PaidUpShares = strings.TrimSpace(td.Text())
				case 2:
					movement.ListedShare = strings.TrimSpace(td.Text())
				case 3:
					movement.ParValue = strings.TrimSpace(td.Text())
				}
			})
		}
	})

	doc.Find("#movement tbody tr").Each(func(i int, s *goquery.Selection) {
		var action CorporateAction
		s.Find("td").Each(func(j int, t *goquery.Selection) {
			switch j {
			case 0:
				action.Date = parseDateToISO8601(strings.TrimSpace(t.Text()))
			case 1:
				action.CorporateActionType = strings.TrimSpace(t.Text())
			case 2:
				action.AuthorizedShare = strings.TrimSpace(t.Text())
			case 3:
				action.ChangedAuthShare = strings.TrimSpace(t.Text())
			case 4:
				action.PaidUpShare = strings.TrimSpace(t.Text())
			case 5:
				action.ChangedPaidUpShare = strings.TrimSpace(t.Text())
			case 6:
				action.ListedShare = strings.TrimSpace(t.Text())
			case 7:
				action.ChangedListedShared = strings.TrimSpace(t.Text())
			case 8:
				action.ParValue = strings.TrimSpace(t.Text())
			}
		})
		movement.CorporateActions = append(movement.CorporateActions, action)
	})

	return &movement, nil
}

func GetCapitalMovement(cookieStr, symbol, locale string) (*CapitalMovement, error) {
	htmlContent, err := MakeCapitalMovementRequest(cookieStr, symbol, locale)
	if err != nil {
		return nil, fmt.Errorf("capital movement request error for symbol %s: %v", symbol, err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("error extracting capital movement for symbol %s: %w", symbol, err)
	}

	movement, err := ExtractCapitalMovement(doc)
	if err != nil {
		return nil, fmt.Errorf("error extracting capital movement for symbol %s: %w", symbol, err)
	}

	return movement, nil
}

func SaveToFile(filename string, data map[string]map[string]CapitalMovement) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating JSON file: %w", err)
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error writing to JSON file: %w", err)
	}

	return nil
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
