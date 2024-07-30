package nvdr

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type NVDRTurnover struct {
	Volume string `json:"Volume"`
	Value  string `json:"Value"`
}

type NVDRData struct {
	Buy        NVDRTurnover `json:"NVDR Turnover (Buy)"`
	Sell       NVDRTurnover `json:"NVDR Turnover (Sell)"`
	BuySell    NVDRTurnover `json:"NVDR Turnover (Buy+Sell)"`
	Net        NVDRTurnover `json:"NVDR Turnover (Net)"`
	Underlying NVDRTurnover `json:"Turnover of Underlying (Buy=Sell+Mkt.TO)"`
	BuyPct     NVDRTurnover `json:"% of NVDR-Buy"`
	SellPct    NVDRTurnover `json:"% of NVDR-Sell"`
	NVDRPct    NVDRTurnover `json:"% of NVDR to its Underlying Securities Turnover"`
	SecCount   string       `json:"NO. of Sec"`
}

type Data struct {
	SET   map[string]NVDRData `json:"SET"`
	MAI   map[string]NVDRData `json:"mai"`
	Total NVDRData            `json:"Grand Total"`
}

func FetchNVDRData(cookieStr string) (Data, error) {
	url := "https://www.setsmart.com/ism/nvdrTrading.html"

	req, err := http.NewRequest("POST", url, strings.NewReader(`action=rankByMarket&tradingMethod=AOM&securityType=&sector=A_0_99_0_M&lstFavorite=0&quickPeriod=&lstPeriod=D&showBeginDate=27%2F07%2F2023&beginDate=27%2F07%2F2023&showEndDate=26%2F07%2F2024&endDate=26%2F07%2F2024&lstDisplay=all`))
	if err != nil {
		return Data{}, err
	}

	// Set headers
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookieStr)
	req.Header.Set("Origin", "https://www.setsmart.com")
	req.Header.Set("Referer", "https://www.setsmart.com/ism/nvdrTrading.html")
	req.Header.Set("Sec-CH-UA", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Data{}, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return Data{}, err
	}

	data := Data{
		SET: make(map[string]NVDRData),
		MAI: make(map[string]NVDRData),
	}

	var section string

	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		if s.HasClass("indextab") {
			sectionText := s.Find("td").Text()
			section = strings.TrimSpace(sectionText)
			return
		}

		if section == "" {
			return
		}

		var nvdrData NVDRData
		parseRow(s, &nvdrData)

		if section == "SET" {
			currentType := strings.TrimSpace(s.Find("td").First().Text())
			data.SET[currentType] = nvdrData
		} else if section == "mai" {
			currentType := strings.TrimSpace(s.Find("td").First().Text())
			data.MAI[currentType] = nvdrData
		} else if section == "Total" || section == "Grand Total" {
			data.Total = nvdrData
		}
	})

	return data, nil
}

func parseRow(s *goquery.Selection, nvdrData *NVDRData) {
	cells := s.Find("td")
	nvdrData.Buy.Volume = cells.Eq(1).Text()
	nvdrData.Buy.Value = cells.Eq(2).Text()
	nvdrData.Sell.Volume = cells.Eq(3).Text()
	nvdrData.Sell.Value = cells.Eq(4).Text()
	nvdrData.BuySell.Volume = cells.Eq(5).Text()
	nvdrData.BuySell.Value = cells.Eq(6).Text()
	nvdrData.Net.Volume = cells.Eq(7).Text()
	nvdrData.Net.Value = cells.Eq(8).Text()
	nvdrData.Underlying.Volume = cells.Eq(9).Text()
	nvdrData.Underlying.Value = cells.Eq(10).Text()
	nvdrData.BuyPct.Volume = cells.Eq(11).Text()
	nvdrData.BuyPct.Value = cells.Eq(12).Text()
	nvdrData.SellPct.Volume = cells.Eq(13).Text()
	nvdrData.SellPct.Value = cells.Eq(14).Text()
	nvdrData.NVDRPct.Volume = cells.Eq(15).Text()
	nvdrData.NVDRPct.Value = cells.Eq(16).Text()
	nvdrData.SecCount = cells.Eq(17).Text()
}

func FetchStockNVDRData(cookieStr, symbol, locale string) (map[string]NVDRData, error) {
	url := fmt.Sprintf("https://www.setsmart.com/ism/nvdrHistorical.html?symbol=%s&locale=%s&beginDate=27%%2F07%%2F2021&quickPeriod=&showEndDate=26%%2F07%%2F2024&endDate=26%%2F07%%2F2024&showBeginDate=27%%2F07%%2F2021&lstDisplay=all&lstPeriod=D", symbol, locale)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookieStr)
	req.Header.Set("Origin", "https://www.setsmart.com")
	req.Header.Set("Referer", "https://www.setsmart.com/ism/nvdrHistorical.html")
	req.Header.Set("Sec-CH-UA", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var links []string
	doc.Find("span.pagelinks a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			links = append(links, "https://www.setsmart.com/ism/"+link)
		}
	})
	data := make(map[string]NVDRData)

	for _, link := range links {
		pageData, err := fetchPageData(link, cookieStr)
		if err != nil {
			fmt.Println("Error fetching page data:", err)
			continue
		}

		for date, nvdrData := range pageData {
			data[date] = nvdrData
		}
		time.Sleep(500 * time.Millisecond)
	}
	return data, nil
}

func fetchPageData(url, cookieStr string) (map[string]NVDRData, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Cookie", cookieStr)
	req.Header.Set("Origin", "https://www.setsmart.com")
	req.Header.Set("Referer", "https://www.setsmart.com/ism/nvdrHistorical.html")
	req.Header.Set("Sec-CH-UA", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	pageData := make(map[string]NVDRData)
	doc.Find("#stock tbody tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		date := strings.TrimSpace(s.Find("td").First().Text())
		isoDate, err := convertToISO8601(date)
		if err != nil {
			fmt.Println("Error converting date:", err)
			return
		}
		var nvdrData NVDRData
		parseStockRow(s, &nvdrData)
		pageData[isoDate] = nvdrData
	})

	return pageData, nil
}

func parseStockRow(s *goquery.Selection, nvdrData *NVDRData) {
	cells := s.Find("td")
	nvdrData.Buy.Volume = cells.Eq(1).Text()
	nvdrData.Buy.Value = cells.Eq(2).Text()
	nvdrData.Sell.Volume = cells.Eq(3).Text()
	nvdrData.Sell.Value = cells.Eq(4).Text()
	nvdrData.BuySell.Volume = cells.Eq(5).Text()
	nvdrData.BuySell.Value = cells.Eq(6).Text()
	nvdrData.Net.Volume = cells.Eq(7).Text()
	nvdrData.Net.Value = cells.Eq(8).Text()
	nvdrData.Underlying.Volume = cells.Eq(9).Text()
	nvdrData.Underlying.Value = cells.Eq(10).Text()
	nvdrData.BuyPct.Volume = cells.Eq(11).Text()
	nvdrData.BuyPct.Value = cells.Eq(12).Text()
	nvdrData.SellPct.Volume = cells.Eq(13).Text()
	nvdrData.SellPct.Value = cells.Eq(14).Text()
	nvdrData.NVDRPct.Volume = cells.Eq(15).Text()
	nvdrData.NVDRPct.Value = cells.Eq(16).Text()
	nvdrData.SecCount = cells.Eq(17).Text()
}

func convertToISO8601(date string) (string, error) {
	parseDate, err := time.Parse("02/01/2006", date)
	if err != nil {
		return "", err
	}
	return parseDate.Format("2006-01-02"), nil
}
