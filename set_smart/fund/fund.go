package fund

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type FundStatusResponse struct {
	Status       bool   `json:"status"`
	ErrorMessage string `json:"errorMessage"`
}

type Fund struct {
	Symbol   string `json:"symbol"`
	NameEN   string `json:"nameEN"`
	NameTH   string `json:"nameTH"`
	ID       string `json:"id"`
	AmicTypr string `json:"aimcType"`
}

type FundResponse struct {
	Funds []Fund `json:"funds"`
}

func FetchFundData(cookieStr, tokenStr string, symbols []string) error {
	allData := make(map[string]interface{})

	for _, symbol := range symbols {
		asOfDate, err := FetchFundStatistics(cookieStr, tokenStr, symbol)
		if err != nil {
			return fmt.Errorf("error fetching fund statistics: %w", err)
		}

		performanceData, err := FetchFundPerformance(cookieStr, tokenStr, symbol, asOfDate)
		if err != nil {
			fmt.Printf("Error fetching performance for fund SCB2576: %v\n", err)
		}

		combinedData := map[string]interface{}{
			"performance_data": performanceData,
		}
		allData[symbol] = combinedData

		time.Sleep(5 * time.Second)
	}

	saveCombineDataToFile("fund_data.json", allData)

	return nil
}

func FetchFundList(cookieStr, tokenStr, date string) ([]Fund, error) {
	url := "https://www.setsmart.com/api/fund/list"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", tokenStr))
	req.Header.Add("cookie", cookieStr)
	req.Header.Add("referer", "https://www.setsmart.com/ssm/fundInformation")
	req.Header.Add("sec-ch-ua", "\"Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\", \"Google Chrome\";v=\"126\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("received non-200 response code: %d", res.StatusCode)
	}

	var fundResponse FundResponse
	err = json.Unmarshal(body, &fundResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling fund list: %w", err)
	}
	return fundResponse.Funds, nil
}

func FetchFundPerformance(cookieStr, tokenStr, symbol, date string) (map[string]interface{}, error) {
	parsedDate, err := time.Parse(time.RFC3339, date)

	formattedDate := parsedDate.Format("02/01/2006")
	encodeDate := url.QueryEscape(formattedDate)
	url := fmt.Sprintf("https://www.setsmart.com/api/fund/%s/performance?date=%s", symbol, encodeDate)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", tokenStr))
	req.Header.Add("cookie", cookieStr)
	req.Header.Add("referer", fmt.Sprintf("https://www.setsmart.com/ssm/fundInformation;fundSymbol=%s", symbol))
	req.Header.Add("sec-ch-ua", "\"Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\", \"Google Chrome\";v=\"126\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("received non-200 response code: %d", res.StatusCode)
	}

	var performanceData map[string]interface{}
	err = json.Unmarshal(body, &performanceData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling fund performance data: %w", err)
	}

	return performanceData, nil
}

func FetchFundHistoricalPerformance(cookieStr, tokenStr, symbol, date string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://www.setsmart.com/api/fund/%s/historical-performance?date=%s", symbol, date)
	fmt.Printf("Requesting URL: %s\n", url) // Log the request URL
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating requset: %w", err)
	}

	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", tokenStr))
	req.Header.Add("cookie", cookieStr)
	req.Header.Add("referer", fmt.Sprintf("https://www.setsmart.com/ssm/fundInformation;fundSymbol=%s", symbol))
	req.Header.Add("sec-ch-ua", "\"Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\", \"Google Chrome\";v=\"126\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("received non-200 response code: %d", res.StatusCode)
	}

	var historicalPerformanceData map[string]interface{}
	err = json.Unmarshal(body, &historicalPerformanceData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling fund historical performance data: %w", err)
	}
	return historicalPerformanceData, nil
}

func FetchFundStatistics(cookieStr, tokenStr, symbol string) (string, error) {
	url := fmt.Sprintf("https://www.setsmart.com/api/fund/%s/statistic", symbol)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", tokenStr))
	req.Header.Add("cookie", cookieStr)
	req.Header.Add("referer", fmt.Sprintf("https://www.setsmart.com/ssm/fundInformation;fundSymbol=%s", symbol))
	req.Header.Add("sec-ch-ua", "\"Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\", \"Google Chrome\";v=\"126\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("received non-200 response code: %d", res.StatusCode)
	}

	var statisticData map[string]interface{}
	err = json.Unmarshal(body, &statisticData)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling fund statistics data: %w", err)
	}

	asOfDate, ok := statisticData["asOfDate"].(string)
	if !ok {
		return "", fmt.Errorf("asOfDate not found in response")
	}

	fmt.Println(asOfDate)
	return asOfDate, nil
}

func saveCombineDataToFile(filename string, data map[string]interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	err = encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("error encoding JSON data to file: %w", err)
	}
	return nil
}
