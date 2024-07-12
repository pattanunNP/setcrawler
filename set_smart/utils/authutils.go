package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func PerformLogin(loginURL, username, password string) (string, string, error) {
	// Create a cookie jar to store cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", "", fmt.Errorf("error creating cookie jar: %w", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	// Login credentials
	credentials := map[string]string{
		"username": username,
		"password": password,
	}

	// Convert credentials to JSON
	jsonData, err := json.Marshal(credentials)
	if err != nil {
		return "", "", fmt.Errorf("error marshalling JSON: %w", err)
	}

	// Create a POST request for login
	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("error creating login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Origin", "https://www.setsmart.com")
	req.Header.Set("Referer", "https://www.setsmart.com/ssm/login")
	req.Header.Set("Sec-CH-UA", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", "macOS")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	// Perform the login request
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("error making login request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading login response body: %w", err)
	}

	// Extract cookies from the response
	u, _ := url.Parse(loginURL)
	cookies := jar.Cookies(u)

	// Create a cookie string from the extracted cookies
	var cookieStr strings.Builder
	for _, cookie := range cookies {
		cookieStr.WriteString(fmt.Sprintf("%s=%s; ", cookie.Name, cookie.Value))
	}

	// Extract access token from response body (assuming it's in JSON format)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", "", fmt.Errorf("access token not found in response")
	}

	// Add the access token to the cookies
	cookieStr.WriteString(fmt.Sprintf("access_grant=%s; ", accessToken))

	return cookieStr.String(), accessToken, nil
}
