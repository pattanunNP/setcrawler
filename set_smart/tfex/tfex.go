package tfex

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type VolumeData struct {
	Total    string `json:"total"`
	Average  string `json:"average"`
	High     string `json:"high"`
	HighDate string `json:"high_date"`
	Low      string `json:"low"`
	LowDate  string `json:"low_date"`
}

type InstrumentData struct {
	Instrument  string     `json:"instrument"`
	Volume      VolumeData `json:"volume"`
	OI          string     `json:"oi"`
	Transaction string     `json:"transaction"`
}

type TFEXData struct {
	Futures   []InstrumentData `json:"futures"`
	Options   []InstrumentData `json:"options"`
	Total     InstrumentData   `json:"total_market"`
	DailyData []DailyData      `json:"daily_data"`
}

type DailyData struct {
	Date    string `json:"date"`
	Futures struct {
		Volume       string `json:"volume"`
		OI           string `json:"oi"`
		Transactions string `json:"transactions"`
	} `json:"futures"`
	Options struct {
		Call struct {
			Volume       string `json:"volume"`
			OI           string `json:"oi"`
			Transactions string `json:"transactions"`
		} `json:"call"`
		Put struct {
			Volume       string `json:"volume"`
			OI           string `json:"oi"`
			Transactions string `json:"transactions"`
		} `json:"put"`
		Total struct {
			Volume       string `json:"volume"`
			OI           string `json:"oi"`
			Transactions string `json:"transactions"`
		} `json:"total"`
	} `json:"options"`
	Total struct {
		Volume       string `json:"volume"`
		OI           string `json:"oi"`
		Transactions string `json:"transactions"`
	} `json:"total"`
}

func FetchTFEXData(cookieStr, locale string) (TFEXData, error) {
	url := "https://www.setsmart.com/ism/tfexTrading.html"

	beginDate, endDate := "27/07/2021", "26/07/2024"
	if locale == "th_TH" {
		beginDate = convertToThaiDate(beginDate)
		endDate = convertToThaiDate(endDate)
	}

	data := fmt.Sprintf("underlyMarket=&lstDisplay=T&tradingMethod=&quickPeriod=&lstPeriod=D&showBeginDate=%s&beginDate=%s&showEndDate=%s&endDate=%s&locale=%s&submit.x=9&submit.y=12", beginDate, beginDate, endDate, endDate, locale)

	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return TFEXData{}, fmt.Errorf("error creating request: %w", err)
	}
	// Set headers
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Add("cache-control", "max-age=0")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cookie", cookieStr)
	req.Header.Add("origin", "https://www.setsmart.com")
	req.Header.Add("priority", "u=0, i")
	req.Header.Add("referer", "https://www.setsmart.com/ism/tfexTrading.html")
	req.Header.Add("sec-ch-ua", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", `"macOS"`)
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TFEXData{}, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return TFEXData{}, err
	}

	tfexData := TFEXData{}
	parseTFEXData(doc, &tfexData)
	parseDailyData(doc, &tfexData.DailyData)

	paginationLinks := extractPaginationLinks(doc)
	fmt.Println("Pagiantion links extracted:", paginationLinks)

	for _, link := range paginationLinks {
		fmt.Println("Pagination link:", link)
		if err := fetchAndParsePage(link, cookieStr, &tfexData); err != nil {
			return tfexData, err
		}
	}

	return tfexData, nil
}

func fetchAndParsePage(link, cookieStr string, tfexData *TFEXData) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return err
	}
	req.Header.Add("cookie", cookieStr)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	parseDailyData(doc, &tfexData.DailyData)
	return nil
}

func extractPaginationLinks(doc *goquery.Document) []string {
	var links []string
	baseURL := "https://www.setsmart.com"
	doc.Find(".pagelinks a.olink").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			link := baseURL + href
			links = append(links, link)
		}
	})
	return links
}

func parseTFEXData(doc *goquery.Document, tfexData *TFEXData) {
	var currentSection string

	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		if s.HasClass("industry") {
			currentSection = CleanText(s.Find("td").Text())
			return
		}

		if currentSection == "Futures" || currentSection == "Options" || currentSection == "ฟิวเจอร์ส" || currentSection == "ออปชั่น" {
			var instrumentData InstrumentData
			parseRow(s, &instrumentData)
			if currentSection == "Futures" || currentSection == "ฟิวเจอร์ส" {
				tfexData.Futures = append(tfexData.Futures, instrumentData)
			} else if currentSection == "Options" || currentSection == "ออปชั่น" {
				tfexData.Options = append(tfexData.Options, instrumentData)
			}
		} else if currentSection == "Total Market" || currentSection == "มูลค่ารวมตลาด" {
			var totalData InstrumentData
			parseRow(s, &totalData)
			tfexData.Total = totalData
		}
	})
}

func parseRow(s *goquery.Selection, instrumentData *InstrumentData) {
	cells := s.Find("td")
	instrumentData.Instrument = CleanText(cells.Eq(0).Text())
	instrumentData.Volume.Total = CleanText(cells.Eq(1).Text())
	instrumentData.Volume.Average = CleanText(cells.Eq(2).Text())
	instrumentData.Volume.High = CleanText(cells.Eq(3).Text())
	instrumentData.Volume.HighDate, _ = convertToISO8601(CleanText(cells.Eq(4).Text()))
	instrumentData.Volume.Low = CleanText(cells.Eq(5).Text())
	instrumentData.Volume.LowDate, _ = convertToISO8601(CleanText(cells.Eq(6).Text()))
	instrumentData.OI = CleanText(cells.Eq(7).Text())
	instrumentData.Transaction = CleanText(cells.Eq(8).Text())
}

func parseDailyData(doc *goquery.Document, dailyData *[]DailyData) {
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		if s.Find("td").First().Text() == "Total" || s.Find("td").First().Text() == "Avg/Day" || s.Find("td").First().Text() == "รวม" || s.Find("td").First().Text() == "เฉลี่ย/วัน" {
			return
		}
		var data DailyData
		if parseDailyRow(s, &data) {
			*dailyData = append(*dailyData, data)
		}
	})
}

func parseDailyRow(s *goquery.Selection, dailyData *DailyData) bool {
	cells := s.Find("td")
	for i := 0; i < 16; i++ {
		text := CleanText(cells.Eq(i).Text())
		if text == "" || containsUnwantedPattern(text) {
			return false
		}
		switch i {
		case 0:
			dailyData.Date, _ = convertToISO8601(text)
		case 1:
			dailyData.Futures.Volume = text
		case 2:
			dailyData.Futures.OI = text
		case 3:
			dailyData.Futures.Transactions = text
		case 4:
			dailyData.Options.Call.Volume = text
		case 5:
			dailyData.Options.Call.OI = text
		case 6:
			dailyData.Options.Call.Transactions = text
		case 7:
			dailyData.Options.Put.Volume = text
		case 8:
			dailyData.Options.Put.OI = text
		case 9:
			dailyData.Options.Put.Transactions = text
		case 10:
			dailyData.Options.Total.Volume = text
		case 11:
			dailyData.Options.Total.OI = text
		case 12:
			dailyData.Options.Total.Transactions = text
		case 13:
			dailyData.Total.Volume = text
		case 14:
			dailyData.Total.OI = text
		case 15:
			dailyData.Total.Transactions = text
		}
	}
	return true
}

func containsUnwantedPattern(text string) bool {
	unwantedPatterns := []string{
		"function openFavoriteQuery", "g_Calendar.localeSensitiveShow", "MM_jumpMenu", "Underlying :",
		"Display :", "Trading Method :", "All(AOM+Block Trade)", "Period :",
		"Volume (Contracts)", "OI(Contracts)", "Transaction(Deals)", "All              Equity Index              Single Stock              Metal              Deferred Contract              Currency              Energy              Interest Rate              Agriculture",
		"หน้า :", "ปริมาณ (สัญญา)", "สถานะคงค้าง", "จำนวนรายการ", "มูลค่ารวมตลาด", "Contract size of SET50 Index Futures was changed from 1,000 Baht per index point to 200 Baht per index point from May 6, 2014 onward.",
	}

	for _, pattern := range unwantedPatterns {
		if strings.Contains(text, pattern) {
			return true
		}
	}
	return false
}

func CleanText(text string) string {
	cleaned := strings.TrimSpace(strings.ReplaceAll(text, "\n", ""))
	return cleaned
}

func convertToISO8601(date string) (string, error) {
	if strings.Contains(date, "/256") {
		parts := strings.Split(date, "/")
		if len(parts) == 3 {
			year, err := strconv.Atoi(parts[2])
			if err != nil {
				return "", err
			}
			year -= 543
			date = fmt.Sprintf("%s/%s/%d", parts[0], parts[1], year)
		}
	}
	parsedDate, err := time.Parse("02/01/2006", date)
	if err != nil {
		return "", err
	}
	return parsedDate.Format("2006-01-02"), nil
}

func convertToThaiDate(date string) string {
	parts := strings.Split(date, "/")
	if len(parts) == 3 {
		year, err := strconv.Atoi(parts[2])
		if err == nil {
			year += 543
			return fmt.Sprintf("%s/%s/%d", parts[0], parts[1], year)
		}
	}
	return date
}
