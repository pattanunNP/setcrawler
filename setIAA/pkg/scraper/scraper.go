package scraper

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// FetchHTML fetches the HTML content of the given URL
func FetchHTML(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	htmlContent, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(htmlContent), nil
}

// ExtractScriptContentUsingGoquery extracts the content inside the <script> tag that starts with "window.__NUXT__=" using goquery
func ExtractScriptContent(htmlContent string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		fmt.Println("Error loading HTML content into goquery:", err)
		return "Script content not found"
	}

	// Find the <script> tag and look for one containing "window.__NUXT__="
	var scriptContent string
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) == "script" {
			html, _ := s.Html()
			if strings.Contains(html, "window.__NUXT__=") {
				scriptContent = html
				return
			}
		}
	})

	if scriptContent == "" {
		return "Script content not found"
	}
	return scriptContent
}
