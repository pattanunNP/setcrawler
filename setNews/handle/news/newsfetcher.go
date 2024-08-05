package news

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type NewsItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	HTMLContent string `json:"htmlContent,omitempty"`
}

type ApiResponse struct {
	PageIndex       int        `json:"pageIndex"`
	PageSize        int        `json:"pageSize"`
	TotalCount      int        `json:"totalcount"`
	TotalPages      int        `json:"totalPages"`
	IndexFrom       int        `json:"indexFrom"`
	Items           []NewsItem `json:"items"`
	HasNextPage     bool       `json:"hasNextPage"`
	HasPreviousPage bool       `json:"hasPreviousPage"`
}

func FetchNews(url string, allItems *[]NewsItem) (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return false, err
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
		fmt.Println(err)
		return false, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return false, err
	}

	*allItems = append(*allItems, apiResponse.Items...)
	return apiResponse.HasNextPage, nil
}
