package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/net/html"
)

// FetchCSRFToken fetches the CSRF token from the login page
func FetchCSRFToken() (string, []*http.Cookie, error) {
	loginPageURL := "https://www.setsmart.com/ssm/login"
	resp, err := http.Get(loginPageURL)
	if err != nil {
		return "", nil, fmt.Errorf("error fetching login page: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("error reading login page body: %w", err)
	}

	// Parse the HTML to extract the CSRF token
	token := extractCSRFToken(body)
	cookies := resp.Cookies()

	return token, cookies, nil
}

// extractCSRFToken extracts the CSRF token from the HTML body
func extractCSRFToken(body []byte) string {
	token := ""
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return token
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			for _, a := range n.Attr {
				if a.Key == "name" && a.Val == "_csrf" {
					for _, a := range n.Attr {
						if a.Key == "value" {
							token = a.Val
							return
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return token
}

// FetchAccessToken logs in and returns the cookies
func FetchAccessToken() ([]*http.Cookie, error) {
	envFilePath := filepath.Join(".", ".env")
	err := godotenv.Load(envFilePath)
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	username := os.Getenv("MYID")
	password := os.Getenv("PASSWORD")

	if username == "" || password == "" {
		return nil, fmt.Errorf("USERNAME or PASSWORD environment variable is missing")
	}

	csrfToken, initialCookies, err := FetchCSRFToken()
	if err != nil {
		return nil, fmt.Errorf("error fetching CSRF token: %w", err)
	}

	apiURL := "https://www.setsmart.com/api/user/login"

	loginPayload := map[string]string{
		"username": username,
		"password": password,
		"_csrf":    csrfToken,
	}
	payloadBytes, err := json.Marshal(loginPayload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling login payload: %w", err)
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("error creating cookie jar: %w", err)
	}

	client := &http.Client{Jar: jar}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	for _, cookie := range initialCookies {
		req.AddCookie(cookie)
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://www.setsmart.com")
	req.Header.Set("Referer", "https://www.setsmart.com/ssm/login")
	req.Header.Set("Sec-CH-UA", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-OK response: %d - %s", resp.StatusCode, string(body))
	}

	cookies := resp.Cookies()

	return cookies, nil
}

// SetupCookieAndToken sets up an HTTP client with the access token and cookies
func SetupCookieAndToken() (*http.Client, []*http.Cookie, error) {
	cookies, err := FetchAccessToken()
	if err != nil {
		return nil, nil, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating cookie jar: %w", err)
	}
	client := &http.Client{Jar: jar}

	// Add cookies to jar
	jar.SetCookies(&url.URL{Scheme: "https", Host: "www.setsmart.com"}, cookies)

	return client, cookies, nil
}
