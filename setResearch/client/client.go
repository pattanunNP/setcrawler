package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ResearchItem struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	CateName    string `json:"cateName"`
	SubCateName string `json:"subCateName"`
	StartDate   string `json:"startDate"`
	Source      string `json:"source"`
	FileURL     string `json:"fileUrl,omitempty"`
}

type Response struct {
	ResearchItems struct {
		IndexFrom  int            `json:"indexFrom"`
		PageIndex  int            `json:"pageIndex"`
		PageSize   int            `json:"pageSize"`
		TotalCount int            `json:"totalCount"`
		TotalPages int            `json:"totalPages"`
		Items      []ResearchItem `json:"items"`
	} `json:"researchItems"`
}

func FetchPage(url string, pageIndex int) (*Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s&pageIndex=%d", url, pageIndex), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	setHeaders(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func FetchHTMLContent(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

func setHeaders(req *http.Request) {
	req.Header.Set("sec-ch-ua", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Set("x-channel", "WEB_SETTRADE")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Referer", "https://www.settrade.com/th/research/analyst-research/main")
	req.Header.Set("x-client-uuid", "fce4b3b8-c539-42f6-81c1-5028d592f745")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
}
