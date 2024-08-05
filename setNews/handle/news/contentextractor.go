package news

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FetchHTMLContent(url string) string {
	var result map[string]string

	if strings.Contains(url, "prachachat") {
		htmlContent := fetchPrachachatContent(url)
		result = map[string]string{
			"source":          "prachachat",
			"raw_html":        htmlContent,
			"article_content": ExtractContent(htmlContent),
		}
	} else if strings.Contains(url, "thunhoon") {
		htmlContent := fetchThunhoonContent(url)
		result = map[string]string{
			"source":          "thunhoon",
			"raw_html":        htmlContent,
			"article_content": ExtractContent(htmlContent),
		}
	} else {
		fmt.Println("Fetching HTML content from:", url)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching HTML:", err)
			return ""
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading HTML body:", err)
			return ""
		}

		htmlContent := string(body)
		result = map[string]string{
			"source":          "other",
			"raw_html":        htmlContent,
			"article_content": ExtractContent(htmlContent),
		}
	}

	return result["raw_html"]
}

func fetchPrachachatContent(url string) string {
	tokenData, err := getTokenData()
	if err != nil {
		fmt.Println("Error fetching token data:", err)
		return ""
	}

	bundle := tokenData["bundle"].(string)
	bidId := tokenData["bidId"].(string)
	cValue := `{"c":1}`

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.Header.Set("Cookie", fmt.Sprintf("khaos=%s; bidId=%s", bundle, bidId))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("c", cValue)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching Prachachat HTML:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading Prachachat HTML body:", err)
		return ""
	}

	htmlContent := string(body)
	fmt.Println("Fetched Prachachat HTML content: ", htmlContent)

	if strings.Contains(htmlContent, "Enable JavaScript and cookies to continue") {
		fmt.Println("JavaScript challenge detected. Content may not be fully loaded.")
		return ""
	}

	return htmlContent
}

func getTokenData() (map[string]interface{}, error) {
	tokenURL := "https://gum.criteo.com/sid/json?origin=prebid&topUrl=https%3A%2F%2Fwww.prachachat.net%2F&domain=www.prachachat.net&bundle=Er-__F8lMkZ6JTJGNlllVlFYbnhvaVRKWVE0cVk0RWt5SUsyNmdvRFlHQjdkenRKd2J6RlR3VUloblBRRzZIcHdYUUozV3RBakFpJTJGSDhaNWJkUFZHZTd1RHYzaFJHSyUyRlVOOW00N3J5YTJxSEZtSjRUTnk4VENUYWNhTEpjSlVKMnZ0UVJibUx1R2MzeDVoSks4OVNCRmZ6VlQ5Y2FvQSUzRCUzRA&cw=1&pbt=1&lsw=1"
	req, err := http.NewRequest("GET", tokenURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Usert-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching token data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading token data: %v", err)
	}

	var tokenData map[string]interface{}
	err = json.Unmarshal(body, &tokenData)
	if err != nil {
		return nil, fmt.Errorf("error parsing token data: %v", err)
	}

	return tokenData, nil
}

func fetchThunhoonContent(url string) string {
	violationURL := "https://player.gliacloud.com/violations/thunhoon.com"

	req, err := http.NewRequest("GET", violationURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("If-Modified-Since", "Mon, 05 Aug 2024 08:26:30 GMT")
	req.Header.Set("Origin", "https://thunhoon.com")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://thunhoon.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", "\"Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\", \"Google Chrome\";v=\"126\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making initial request:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	fmt.Println("Initial request response from Thunhoon:", string(body))

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error fetching Thunhoon HTML:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading Thunhoon HTML body:", err)
		return ""
	}

	htmlContent := string(body)
	fmt.Println("Fetched Thunhoon HTML content:", htmlContent)

	if strings.Contains(htmlContent, "Enable JavaScript and cookies to continue") {
		fmt.Println("JavaScript challenge detected. Content may not be fully loaded.")
		return ""
	}
	return htmlContent
}

func extractArticleContent(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}

	articleSelection := doc.Find("article.post")
	if articleSelection.Length() == 0 {
		return "", fmt.Errorf("no article content found")
	}

	articleHTML, err := articleSelection.Html()
	if err != nil {
		return "", err
	}

	cleanedArticleHTML := CleanHTMLContent(articleHTML)
	return cleanedArticleHTML, nil
}

func ExtractContent(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return ""
	}

	patterns := []string{
		"#article .entry-content",
		"#the-post .entry-content",
		"div.article-content",
		"div.post-body",
		"article",
	}

	var content string
	for _, pattern := range patterns {
		content = doc.Find(pattern).Text()
		if content != "" {
			return CleanHTMLContent(content)
		}
	}

	content = doc.Find("body").Text()
	return CleanHTMLContent(content)
}

func CleanHTMLContent(content string) string {
	content = html.UnescapeString(content)

	re := regexp.MustCompile(`(?i)<(script|style)[^>]*>.*?</\\1>`)
	content = re.ReplaceAllString(content, "")

	re = regexp.MustCompile(`(?i)<[^>]*(on\w+|style)="[^"]*"[^>]*>`)
	content = re.ReplaceAllString(content, "")

	unwantedPatterns := []string{
		`(?i)\(function \(w, d, s, l, i\) \{.*?\}\)\(window, document, 'script', 'dataLayer', 'GTM-[A-Z0-9]+'\);`,
		`(?i)\.icon-social \{.*?\}`,
		`(?i)/\*.*?\*/`,
	}

	for _, pattern := range unwantedPatterns {
		re = regexp.MustCompile(pattern)
		content = re.ReplaceAllString(content, "")
	}

	re = regexp.MustCompile(`(?i)<!--.*?-->`)
	content = re.ReplaceAllString(content, "")

	re = regexp.MustCompile(`<[^>]+>`)
	content = re.ReplaceAllString(content, "")

	content = regexp.MustCompile(`\s+`).ReplaceAllString(content, " ")

	content = strings.TrimSpace(content)

	return content
}

func ExtractOnlyText(input string) string {
	jsonPattern := regexp.MustCompile(`"([^"]+)"\s*:\s*"([^"]+)"`)
	htmlPattern := regexp.MustCompile(`<[^>]+>`)
	specialCharsPattern := regexp.MustCompile(`[^\w\s]`)
	whitespacePattern := regexp.MustCompile(`\s+`)

	cleaned := jsonPattern.ReplaceAllString(input, " ")
	cleaned = htmlPattern.ReplaceAllString(cleaned, " ")
	cleaned = specialCharsPattern.ReplaceAllString(cleaned, " ")
	cleaned = whitespacePattern.ReplaceAllString(cleaned, " ")
	cleaned = strings.TrimSpace(cleaned)

	return cleaned
}
