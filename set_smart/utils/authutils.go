package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type CommonStock struct {
	AuthorizedCapital             string            `json:"Authorized Capital (Common Stock)"`
	PaidUpStock                   string            `json:"Paid-up Stock (Common Stock)"`
	TreasuryStock                 string            `json:"Treasury Stock (Common Stock)"`
	VotingStockMinusTreasuryStock map[string]string `json:"Voting Stock minus Treasury Stock (Common Stock)"`
}

type PreferredStock struct {
	AuthorizedCapital             string `json:"Authorized Capital (Preferred)"`
	PaidUpCapital                 string `json:"Paid-up Capital (Preferred)"`
	PaidUpStock                   string `json:"Paid-up Stock (Preferred Stock)"`
	TreasuryStock                 string `json:"Treasury Stock (Preferred Stock)"`
	VotingStockMinusTreasuryStock string `json:"Voting Stock minus Treasury Stock (Preferred Stock)"`
}

type CompanyProfile struct {
	Name                             string         `json:"Name (Name Change)"`
	Address                          string         `json:"Address"`
	Telephone                        string         `json:"Telephone"`
	Fax                              string         `json:"Fax"`
	URL                              string         `json:"URL"`
	EstablishmentDate                string         `json:"Establishment Date"`
	JuristicPersonRegistrationNumber string         `json:"Juristic Person Registration Number"`
	CompanyType                      string         `json:"Company Type"`
	CommonStock                      CommonStock    `json:"Common Stock"`
	PreferredStock                   PreferredStock `json:"Preferred Stock"`
	Form56OneReportEng               string         `json:"Form56-1 One Report (Eng)"`
	Form56OneReportThai              string         `json:"Form56-1 One Report (Thai)"`
	ListedCompanySnapshotEng         string         `json:"Listed Company Snapshot (Eng)"`
	ListedCompanySnapshotThai        string         `json:"Listed Company Snapshot (Thai)"`
	DividendPolicy                   string         `json:"Dividend Policy"`
	AuditorAuditCompany              string         `json:"Auditor/Audit company (Effective Until 31/12/2024)"`
	FinanceResponsibility            string         `json:"The person taking the highest responsibility in finance and accounting"`
	AccountSupervision               string         `json:"The person supervising accounting"`
	ListingCondition                 string         `json:"Listing Condition"`
}

type Securities struct {
	Securities              string `json:"Securities"`
	Name                    string `json:"Name (Name Change)"`
	Market                  string `json:"Market"`
	IndustrySector          string `json:"Industry/Sector (Sector Change)"`
	SecurityType            string `json:"Security Type"`
	Status                  string `json:"Status"`
	ListedDate              string `json:"Listed Date"`
	Par                     string `json:"Par (Par Change)"`
	NoOfListedShare         string `json:"No. of Listed Share"`
	FirstTradingDate        string `json:"First Trading Date"`
	ISINNumber              string `json:"ISIN Number"`
	ForeignLimit            string `json:"Foreign Limit*"`
	ForeignAvailable        string `json:"Foreign Available*"`
	ForeignQueue            string `json:"Foreign Queue*"`
	ForeignLimitForExercise string `json:"Foreign Limit for Exercise*"`
	AccountForm             string `json:"Account Form"`
	FiscalYearEnd           string `json:"Fiscal Year End"`
	IPOPrice                string `json:"IPO Price"`
	IPOFinancialAdvisor     string `json:"IPO Financial Advisor"`
	SubscriptionPeriod      string `json:"Subscription Period"`
	IPOSilentPeriod         string `json:"IPO Silent Period"`
	Filing                  string `json:"Filing"`
	SalesReport             string `json:"Sales Report"`
	DetailOfSecurity        string `json:"Detail of Security / Information Memorandum"`
}

type OrganizedData struct {
	CompanyProfiles map[string]CompanyProfile `json:"company_profiles"`
	Securities      map[string]Securities     `json:"securities"`
}

func PerformLogin(loginURL, username, password string) (string, string, error) {
	// Create a cookie jar to store cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", "", fmt.Errorf("error creating cookie jar: %w", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	// Login credentials
	credentials := map[string]string{
		"username": username,
		"password": password,
	}

	// Convert credentials to JSON
	jsonData, err := json.Marshal(credentials)
	if err != nil {
		return "", "", fmt.Errorf("error marshalling JSON: %w", err)
	}

	// Create a POST request for login
	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("error creating login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Origin", "https://www.setsmart.com")
	req.Header.Set("Referer", "https://www.setsmart.com/ssm/login")
	req.Header.Set("Sec-CH-UA", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", "macOS")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	// Perform the login request
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("error making login request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading login response body: %w", err)
	}

	// Extract cookies from the response
	u, _ := url.Parse(loginURL)
	cookies := jar.Cookies(u)

	// Create a cookie string from the extracted cookies
	var cookieStr strings.Builder
	for _, cookie := range cookies {
		cookieStr.WriteString(fmt.Sprintf("%s=%s; ", cookie.Name, cookie.Value))
	}

	// Extract access token from response body (assuming it's in JSON format)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", "", fmt.Errorf("access token not found in response")
	}

	// Add the access token to the cookies
	cookieStr.WriteString(fmt.Sprintf("access_grant=%s; ", accessToken))

	return cookieStr.String(), accessToken, nil
}

func MakeRequestWithCookies(cookieStr, symbol, locale string) (string, CompanyProfile, Securities, error) {
	data := fmt.Sprintf("symbol=%s&locale=%s&submit.x=19&submit.y=14", symbol, locale)
	requestURL := "https://www.setsmart.com/ism/companyprofile.html"

	req, err := http.NewRequest("POST", requestURL, strings.NewReader(data))
	if err != nil {
		return "", CompanyProfile{}, Securities{}, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Cookie", cookieStr)
	req.Header.Set("Origin", "https://www.setsmart.com")
	req.Header.Set("Referer", "https://www.setsmart.com/ism/companyprofile.html?locale=en_US")
	req.Header.Set("Sec-CH-UA", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", "macOS")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Set("Priority", "u=0, i")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", CompanyProfile{}, Securities{}, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", CompanyProfile{}, Securities{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", CompanyProfile{}, Securities{}, fmt.Errorf("error reading response body: %w", err)
	}

	dataMap, companyName := parseHTML(string(body))

	// Organize data into desired structure
	organizedCompanyProfile := organizeData(dataMap)
	organizeSecurities := organizeSecurities(dataMap)

	return companyName, organizedCompanyProfile, organizeSecurities, nil
}

func parseHTML(htmlString string) (map[string]string, string) {
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return nil, ""
	}

	dataMap := make(map[string]string)
	var companyName string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "tr" {
				parseRow(n, dataMap)
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.ElementNode && c.Data == "td" {
						for cc := c.FirstChild; cc != nil; cc = cc.NextSibling {
							if cc.Type == html.TextNode && (strings.Contains(cc.Data, "Public Company Limited") || strings.Contains(cc.Data, "บริษัท")) {
								companyName = strings.TrimSpace(cc.Data)
							}
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "pre" {
				preContent := getTextContent(c)
				dataMap["Detail of Security / Information Memorandum"] = preContent
			}
			f(c)
		}
	}
	f(doc)

	if companyName == "" {
		companyName = dataMap["Name (Name Change)"]
	}

	return dataMap, companyName
}

// func parseFilename(n *html.Node) string {
// 	var filename string
// 	for c := n.FirstChild; c != nil; c = c.NextSibling {
// 		if c.Type == html.ElementNode && c.Data == "td" {
// 			for _, attr := range c.Attr {
// 				if attr.Key == "class" && attr.Val == "table-bold" {
// 					filename = strings.TrimSpace(getTextContent(c))
// 					return filename
// 				}
// 			}
// 		}
// 	}
// 	return filename
// }

func parseRow(n *html.Node, dataMap map[string]string) {
	tdCount := 0
	var label, value string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "td" {
			tdCount++
			if tdCount == 1 {
				label = strings.TrimSpace(getTextContent(c))
			} else if tdCount == 2 {
				value = strings.TrimSpace(getTextContent(c))
				dataMap[label] = value
				label, value = "", ""
			}
		}
	}
}

func getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var buf strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		buf.WriteString(getTextContent(c))
	}
	return buf.String()
}

func organizeData(dataMap map[string]string) CompanyProfile {
	return CompanyProfile{
		Name:                             dataMap["Name (Name Change)"],
		Address:                          dataMap["Address"],
		Telephone:                        dataMap["Telephone"],
		Fax:                              dataMap["Fax"],
		URL:                              dataMap["URL"],
		EstablishmentDate:                dataMap["Establishment Date"],
		JuristicPersonRegistrationNumber: dataMap["Juristic Person Registration Number"],
		CompanyType:                      dataMap["Company Type"],
		CommonStock: CommonStock{
			AuthorizedCapital: dataMap["Authorized Capital (Common Stock)"],
			PaidUpStock:       dataMap["Paid-up Stock (Common Stock)"],
			TreasuryStock:     dataMap["Treasury Stock (Common Stock)"],
			VotingStockMinusTreasuryStock: map[string]string{
				"As of 15/07/2024": dataMap["As of 15/07/2024"],
				"As of 30/06/2024": dataMap["As of 30/06/2024"],
			},
		},
		PreferredStock: PreferredStock{
			AuthorizedCapital:             dataMap["Authorized Capital (Preferred Stock)"],
			PaidUpCapital:                 dataMap["Paid-up Capital (Preferred Stock)"],
			PaidUpStock:                   dataMap["Paid-up Stock (Preferred Stock)"],
			TreasuryStock:                 dataMap["Treasury Stock (Preferred Stock)"],
			VotingStockMinusTreasuryStock: dataMap["Voting Stock minus Treasury Stock (Preferred Stock)"],
		},
		Form56OneReportEng:        dataMap["Form56-1 One Report (Eng)"],
		Form56OneReportThai:       dataMap["Form56-1 One Report (Thai)"],
		ListedCompanySnapshotEng:  dataMap["Listed Company Snapshot (Eng)"],
		ListedCompanySnapshotThai: dataMap["Listed Company Snapshot (Thai)"],
		DividendPolicy:            dataMap["Dividend Policy"],
		AuditorAuditCompany:       dataMap["Auditor/Audit company\n    \n      (Effective Until 31/12/2024)"],
		FinanceResponsibility:     dataMap["The person taking the highest responsibility in finance and accounting"],
		AccountSupervision:        dataMap["The person supervising accounting"],
		ListingCondition:          dataMap["Listing Condition"],
	}
}

func organizeSecurities(dataMap map[string]string) Securities {
	return Securities{
		Securities:              dataMap["Securities"],
		Name:                    dataMap["Name (Name Change)"],
		Market:                  dataMap["Market"],
		IndustrySector:          dataMap["Industry/Sector (Sector Change)"],
		SecurityType:            dataMap["Security Type"],
		Status:                  dataMap["Status"],
		ListedDate:              dataMap["Listed Date"],
		Par:                     dataMap["Par (Par Change)"],
		NoOfListedShare:         dataMap["No. of Listed Share"],
		FirstTradingDate:        dataMap["First Trading Date"],
		ISINNumber:              dataMap["ISIN Number"],
		ForeignLimit:            dataMap["Foreign Limit*"],
		ForeignAvailable:        dataMap["Foreign Available*"],
		ForeignQueue:            dataMap["Foreign Queue*"],
		ForeignLimitForExercise: dataMap["Foreign Limit for Exercise*"],
		AccountForm:             dataMap["Account Form"],
		FiscalYearEnd:           dataMap["Fiscal Year End"],
		IPOPrice:                dataMap["IPO Price"],
		IPOFinancialAdvisor:     dataMap["IPO Financial Advisor"],
		SubscriptionPeriod:      dataMap["Subscription Period"],
		IPOSilentPeriod:         dataMap["IPO Silent Period"],
		Filing:                  dataMap["Filing"],
		SalesReport:             dataMap["Sales Report"],
		DetailOfSecurity:        dataMap["Detail of Security / Information Memorandum"],
	}
}
