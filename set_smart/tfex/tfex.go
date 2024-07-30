package tfex

import (
	"fmt"
	"net/http"
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

	data := fmt.Sprintf("underlyMarket=&lstDisplay=T&tradingMethod=&quickPeriod=&lstPeriod=D&showBeginDate=27%%2F07%%2F2021&beginDate=27%%2F07%%2F2021&showEndDate=26%%2F07%%2F2024&endDate=26%%2F07%%2F2024&locale=%s&submit.x=9&submit.y=12", locale)

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

	return tfexData, nil
}

func parseTFEXData(doc *goquery.Document, tfexData *TFEXData) {
	var currentSection string

	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		if s.HasClass("industry") {
			currentSection = s.Find("td").Text()
			return
		}

		if currentSection == "Futures" || currentSection == "Options" {
			var instrumentData InstrumentData
			parseRow(s, &instrumentData)
			if currentSection == "Futures" {
				tfexData.Futures = append(tfexData.Futures, instrumentData)
			} else if currentSection == "Options" {
				tfexData.Options = append(tfexData.Options, instrumentData)
			}
		} else if currentSection == "Total Market" {
			var totalData InstrumentData
			parseRow(s, &totalData)
			tfexData.Total = totalData
		}
	})
}

func parseRow(s *goquery.Selection, instrumentData *InstrumentData) {
	cells := s.Find("td")
	instrumentData.Instrument = cells.Eq(0).Text()
	instrumentData.Volume.Total = cells.Eq(1).Text()
	instrumentData.Volume.Average = cells.Eq(2).Text()
	instrumentData.Volume.High = cells.Eq(3).Text()
	instrumentData.Volume.HighDate, _ = convertToISO8601(cells.Eq(4).Text())
	instrumentData.Volume.Low = cells.Eq(5).Text()
	instrumentData.Volume.LowDate, _ = convertToISO8601(cells.Eq(6).Text())
	instrumentData.OI = cells.Eq(7).Text()
	instrumentData.Transaction = cells.Eq(8).Text()
}

func parseDailyData(doc *goquery.Document, dailyData *[]DailyData) {
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		if s.Find("td").First().Text() == "Total" || s.Find("td").First().Text() == "Avg/Day" {
			return
		}
		var data DailyData
		parseDailyRow(s, &data)
		*dailyData = append(*dailyData, data)
	})
}

func parseDailyRow(s *goquery.Selection, dailyData *DailyData) {
	cells := s.Find("td")
	dailyData.Date, _ = convertToISO8601(cells.Eq(0).Text())
	dailyData.Futures.Volume = cells.Eq(1).Text()
	dailyData.Futures.OI = cells.Eq(2).Text()
	dailyData.Futures.Transactions = cells.Eq(3).Text()
	dailyData.Options.Call.Volume = cells.Eq(4).Text()
	dailyData.Options.Call.OI = cells.Eq(5).Text()
	dailyData.Options.Call.Transactions = cells.Eq(6).Text()
	dailyData.Options.Put.Volume = cells.Eq(7).Text()
	dailyData.Options.Put.OI = cells.Eq(8).Text()
	dailyData.Options.Put.Transactions = cells.Eq(9).Text()
	dailyData.Options.Total.Volume = cells.Eq(10).Text()
	dailyData.Options.Total.OI = cells.Eq(11).Text()
	dailyData.Options.Total.Transactions = cells.Eq(12).Text()
	dailyData.Total.Volume = cells.Eq(13).Text()
	dailyData.Total.OI = cells.Eq(14).Text()
	dailyData.Total.Transactions = cells.Eq(15).Text()
}

func cleanTex(text string) string {
	return strings.TrimSpace(strings.ReplaceAll(text, "\n", ""))
}

func convertToISO8601(date string) (string, error) {
	parsedDate, err := time.Parse("02/01/2006", date)
	if err != nil {
		return "", err
	}
	return parsedDate.Format("2006-01-02"), nil
}
