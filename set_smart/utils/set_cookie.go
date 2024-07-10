package authutils

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

// FetchAccessToken logs in and returns the access token and cookies
// func FetchAccessToken(username, password string) ([]*http.Cookie, string, error) {
// 	apiURL := "https://www.setsmart.com/api/user/login"

// 	loginPayload := map[string]string{
// 		"username": username,
// 		"password": password,
// 	}
// 	payloadBytes, err := json.Marshal(loginPayload)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("error marshalling login payload: %w", err)
// 	}

// 	jar, err := cookiejar.New(nil)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("error creating cookie jar: %w", err)
// 	}

// 	client := &http.Client{Jar: jar}

// 	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
// 	if err != nil {
// 		return nil, "", fmt.Errorf("error creating request: %w", err)
// 	}

// 	req.Header.Set("Accept", "application/json, text/plain, */*")
// 	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Origin", "https://www.setsmart.com")
// 	req.Header.Set("Referer", "https://www.setsmart.com/ssm/login")
// 	req.Header.Set("Sec-CH-UA", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
// 	req.Header.Set("Sec-CH-UA-Mobile", "?0")
// 	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
// 	req.Header.Set("Sec-Fetch-Dest", "empty")
// 	req.Header.Set("Sec-Fetch-Mode", "cors")
// 	req.Header.Set("Sec-Fetch-Site", "same-origin")
// 	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("error sending request: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		body, _ := io.ReadAll(resp.Body)
// 		return nil, "", fmt.Errorf("received non-OK response: %d - %s", resp.StatusCode, string(body))
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("error reading response body: %w", err)
// 	}

// 	var loginResp LoginResponse
// 	err = json.Unmarshal(body, &loginResp)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("error unmarshalling response: %w", err)
// 	}

// 	cookies := resp.Cookies()

// 	return cookies, loginResp.AccessToken, nil
// }

// SetupClientWithToken sets up an HTTP client with the access token and cookies
func SetupClientWithToken(cookies []*http.Cookie) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("error creating cookie jar: %w", err)
	}

	client := &http.Client{Jar: jar}

	jar.SetCookies(&url.URL{Scheme: "https", Host: "wwww.setsmart.com"}, cookies)

	return client, nil
}

// MakeAuthenticatedRequest makes an authenticated request using the provided client, access token, and cookies
func MakeAuthenticatedRequest(client *http.Client, requestURL string, additionalHeaders map[string]string, cookies []*http.Cookie) (string, error) {
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Set additional headers
	for key, value := range additionalHeaders {
		req.Header.Set(key, value)
	}

	// Construct the Cookie header manually
	var cookieHeader strings.Builder
	for _, cookie := range cookies {
		cookieHeader.WriteString(cookie.Name + "=" + cookie.Value + "; ")
	}
	req.Header.Set("Cookie", strings.TrimRight(cookieHeader.String(), "; "))

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	return string(body), nil
}

// PrintCookies prints all details of each cookie
func PrintCookies(cookies []*http.Cookie) {
	for _, cookie := range cookies {
		fmt.Printf("Name: %s\n", cookie.Name)
		fmt.Printf("Value: %s\n", cookie.Value)
		fmt.Printf("Domain: %s\n", cookie.Domain)
		fmt.Printf("Path: %s\n", cookie.Path)
		fmt.Printf("Expires: %s\n", cookie.Expires)
		fmt.Printf("Secure: %t\n", cookie.Secure)
		fmt.Printf("HttpOnly: %t\n", cookie.HttpOnly)
		fmt.Printf("SameSite: %s\n", cookie.SameSite)
		fmt.Println()
	}
}
