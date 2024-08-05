package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type NewsItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type ApiResponse struct {
	PageIndex  int        `json:"pageIbdex"`
	PageSize   int        `json:"pageSize"`
	TotalCount int        `json:"totalCount"`
	TotalPages int        `json:"totalPages"`
	Items      []NewsItem `json:"items"`
}

func main() {
	baseURL := "https://www.settrade.com/api/cms/v1/research-settrade/popular-research?frequency=Daily&language=TH&pageIndex=%d&pageSize=20"
	var allItems []NewsItem

	for pageIndex := 0; pageIndex < 3; pageIndex++ {
		url := fmt.Sprintf(baseURL, pageIndex)
		apiResponse, err := fetcNews(url)
		if err != nil {
			fmt.Println("Error fetching news:", err)
			return
		}

		allItems = append(allItems, apiResponse.Items...)
	}

	file, err := json.MarshalIndent(allItems, "", " ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	err = os.WriteFile("analytics.json", file, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Analytics Items Saved to analytics.json")
}

func fetcNews(url string) (*ApiResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Add("cookie", "your_cookie_here")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("referer", "https://www.settrade.com/th/news-and-articles/news/main")
	req.Header.Add("sec-ch-ua", "\"Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\", \"Google Chrome\";v=\"126\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Add("x-channel", "WEB_SETTRADE")
	req.Header.Add("x-client-uuid", "0e1ad0c3-1bc5-4316-b5fd-81f00388d3fe")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}
	return &apiResponse, nil
}
