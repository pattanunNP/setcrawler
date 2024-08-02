package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type NewResponse struct {
	Description string `json:"description"`
}

type ApiResponse struct {
	HasNextPage     bool          `json:"hasNextPage"`
	HasPreviousPage bool          `json:"hasPreviousPage"`
	IndexFrom       int           `json:"indexFrom"`
	Items           []NewResponse `json:"items"`
}

func main() {
	url := "https://www.settrade.com/api/cms/v1/news/all?cate=60ad3cae-ba3d-4405-af14-3ed4af1e5065&fromDate=26%2F07%2F2024&toDate=02%2F08%2F2024&orderBy=date&pageIndex=0&pageSize=20"
	descriptions := []string{}

	fetchNews(url, &descriptions)

	prettyJSON, err := json.MarshalIndent(descriptions, "", "  ")
	if err != nil {
		fmt.Println("Error generating pretty JSON:", err)
		return
	}

	fmt.Println(string(prettyJSON))

	time.Sleep(10 * time.Second)
	os.Exit(0)
}

func fetchNews(url string, description *[]string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("sec-ch-ua", "\"Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\", \"Google Chrome\";v=\"126\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Add("x-channel", "WEB_SETTRADE")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Referer", "https://www.settrade.com/th/news-and-articles/news/main")
	req.Header.Add("x-client-uuid", "1806db1b-1f7c-4a20-8160-4cddd1ba1b2b")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("Cookie", "incap_ses_1010_2685215=gldAH9lFEHsHc1QEyT0EDoXErGYAAAAA6yeNVlfxA1t7BWjuLLa4OQ==; nlbi_2685215=595mfnC/EVq7CjE2wZdY4QAAAABNotSbhY8OVk7VD2KZ7QaD; visid_incap_2685215=alApkGF1RzG9vFke8yL3bv3/oGYAAAAAQkIPAAAAAACfpUbl+owb2MEPz6nvfNqV; charlot=02043556-9608-4ef8-88b8-0065bf0bffeb")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, item := range apiResponse.Items {
		*description = append(*description, item.Description)
	}

	if apiResponse.HasNextPage {
		nextPageIndex := apiResponse.IndexFrom + len(apiResponse.Items)
		nextUrl := fmt.Sprintf("https://www.settrade.com/api/cms/v1/news/all?cate=60ad3cae-ba3d-4405-af14-3ed4af1e5065&fromDate=26%%2F07%%2F2024&toDate=02%%2F08%%2F2024&orderBy=date&pageIndex=%d&pageSize=20", nextPageIndex)
		fetchNews(nextUrl, description)
	}
}
