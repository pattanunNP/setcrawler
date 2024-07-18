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

	mappedData := mapDataKeys(dataMap, locale)

	// Organize data into desired structure
	organizedCompanyProfile := organizeData(mappedData)
	organizeSecurities := organizeSecurities(mappedData)

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
		if n.Type == html.ElementNode && n.Data == "tr" {
			labelNode, valueNode := extractLabelValueNodes(n)
			if labelNode != nil && valueNode != nil {
				label := strings.TrimSpace(getTextContent(labelNode))
				value := strings.TrimSpace(getTextContent(valueNode))
				dataMap[label] = value
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	companyName = dataMap["Name"]

	return dataMap, companyName
}

func extractLabelValueNodes(n *html.Node) (*html.Node, *html.Node) {
	var labelNode, valueNode *html.Node
	tdCount := 0
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "td" {
			tdCount++
			if tdCount == 1 {
				labelNode = c
			} else if tdCount == 2 {
				valueNode = c
			}
		}
	}
	return labelNode, valueNode
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

// func parseRow(n *html.Node, dataMap map[string]string) {
// 	tdCount := 0
// 	var label, value string
// 	for c := n.FirstChild; c != nil; c = c.NextSibling {
// 		if c.Type == html.ElementNode && c.Data == "td" {
// 			tdCount++
// 			if tdCount == 1 {
// 				label = strings.TrimSpace(getTextContent(c))
// 			} else if tdCount == 2 {
// 				value = strings.TrimSpace(getTextContent(c))
// 				dataMap[label] = value
// 				label, value = "", ""
// 			}
// 		}
// 	}
// }

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

func mapDataKeys(dataMap map[string]string, locale string) map[string]string {
	mappedData := make(map[string]string)
	for key, value := range dataMap {
		switch key {
		case "Name (Name Change)", "ชื่อบริษัท (การเปลี่ยนชื่อ)":
			mappedData["Name (Name Change)"] = value
		case "Address", "ที่อยู่":
			mappedData["Address"] = value
		case "Telephone", "เบอร์โทรศัพท์":
			mappedData["Telephone"] = value
		case "Fax", "เบอร์โทรสาร":
			mappedData["Fax"] = value
		case "URL":
			mappedData["URL"] = value
		case "Establishment Date", "วันที่ก่อตั้งบริษัท":
			mappedData["Establishment Date"] = value
		case "Juristic Person Registration Number", "เลขทะเบียนนิติบุคคล":
			mappedData["Juristic Person Registration Number"] = value
		case "Company Type", "ประเภทบริษัท":
			mappedData["Company Type"] = value
		case "Authorized Capital (Common Stock)", "ทุนจดทะเบียน (หุ้นสามัญ)":
			mappedData["Authorized Capital (Common Stock)"] = value
		case "Paid-up Capital (Common Stock)", "ทุนจดทะเบียนชำระแล้ว (หุ้นสามัญ)":
			mappedData["Paid-up Capital (Common Stock)"] = value
		case "Paid-up Stock (Common Stock)", "จำนวนหุ้นชำระแล้ว (หุ้นสามัญ)":
			mappedData["Paid-up Stock (Common Stock)"] = value
		case "Treasury Stock (Common Stock)", "จำนวนหุ้นซื้อคืน (หุ้นสามัญ)":
			mappedData["Treasury Stock (Common Stock)"] = value
		case "Voting Stock minus Treasury Stock (Common Stock)", "จำนวนหุ้นที่มีสิทธิออกเสียง หัก หุ้นซื้อคืน (หุ้นสามัญ)":
			mappedData["Voting Stock minus Treasury Stock (Common Stock)"] = value
		case "Form56-1 One Report (Thai)", "แบบฟอร์ม 56-1 One Report (ไทย)":
			mappedData["Form56-1 One Report (Thai)"] = value
		case "Form56-1 One Report (Eng)", "แบบฟอร์ม 56-1 One Report (อังกฤษ)":
			mappedData["Form56-1 One Report (Eng)"] = value
		case "Listed Company Snapshot (Thai)", "Listed Company Snapshot (ไทย)":
			mappedData["Listed Company Snapshot (Thai)"] = value
		case "Listed Company Snapshot (Eng)", "Listed Company Snapshot (อังกฤษ)":
			mappedData["Listed Company Snapshot (Eng)"] = value
		case "Dividend Policy", "นโยบายเงินปันผล":
			mappedData["Dividend Policy"] = value
		case "Auditor/Audit company", "ผู้ตรวจสอบบัญชี/สำนักงานตรวจสอบบัญชี":
			mappedData["Auditor/Audit company"] = value
		case "The person taking the highest responsibility in finance and accounting", "ผู้รับผิดชอบสูงสุดในสายงานบัญชีและการเงิน":
			mappedData["The person taking the highest responsibility in finance and accounting"] = value
		case "The person supervising accounting", "ผู้ควบคุมดูแลการทำบัญชี":
			mappedData["The person supervising accounting"] = value
		case "Listing Condition", "เงื่อนไขการจดทะเบียนในตลท.":
			mappedData["Listing Condition"] = value
		case "Securities", "หลักทรัพย์":
			mappedData["Securities"] = value
		case "Market", "ตลาด":
			mappedData["Market"] = value
		case "Industry/Sector", "กลุ่มอุตสาหากรร/หมวดอุตสาหกรรม":
			mappedData["Industry/Sector"] = value
		case "Security Type", "ประเภทหลักทรัพย์":
			mappedData["Security Type"] = value
		case "Status", "สถานะ":
			mappedData["Status"] = value
		case "Listed Date", "วันที่จดทะเบียนกับตลท":
			mappedData["Listed Date"] = value
		case "Par", "ราคาพาร์":
			mappedData["Par"] = value
		case "No. of Listed Share", "จำนวนหุ้นจดทะเบียนกับตลท.":
			mappedData["No. of Listed Share"] = value
		case "First Trading Date", "วันที่เริ่มต้นซื้อขาย":
			mappedData["First Trading Date"] = value
		case "ISIN Number", "เลขรหัสหลักทรัพย์สากล":
			mappedData["ISIN Number"] = value
		case "Foreign Limit", "ข้อจำกัดหุ้นต่างด้าว":
			mappedData["Foreign Limit"] = value
		case "Foreign Available", "จำนวนหุ้นคงเหลือเพื่อโอน":
			mappedData["Foreign Available"] = value
		case "Foreign Queue", "จำนวนหุ้นต่างด้าวที่รอการโอน":
			mappedData["Foreign Queue"] = value
		case "Foreign Limit for Exercise", "ข้อจำกัดหุ้นต่างด้าวสำหรับการแปลงสภาพ":
			mappedData["Foreign Limit for Exercise"] = value
		case "Account Form", "รูปแบบงบการเงิน":
			mappedData["Account Form"] = value
		case "Fiscal Year End", "วันที่สิ้นสุดรอบระยะเวลาบัญชี":
			mappedData["Fiscal Year End"] = value
		case "IPO Price", "ราคาเสนอขายหุ้นแก่ประชาชนครั้งแรก":
			mappedData["IPO Price"] = value
		case "IPO Financial Advisor", "ที่ปรึกษาทางการเงิน IPO":
			mappedData["IPO Financial Advisor"] = value
		case "Subscription Period", "ช่วงเวลาการจองซื้อหุ้น":
			mappedData["Subscription Period"] = value
		case "IPO Silent Period", "ระยะเวลาห้ามซื้อขายหุ้น IPO":
			mappedData["IPO Silent Period"] = value
		case "Filing", "แบบ Filing":
			mappedData["Filing"] = value
		case "Sales Report", "รายงานผลการขาย":
			mappedData["Sales Report"] = value
		case "Detail of Security / Information Memorandum", "รายละเอียดหลักทรัพย์ / สรุปข้อสนเทศหลักทรัพย์":
			mappedData["Detail of Security / Information Memorandum"] = value
		default:
			mappedData[key] = value
		}

	}
	return mappedData
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
