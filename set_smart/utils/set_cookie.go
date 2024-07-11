package authutils

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// SetupClientWithCookies sets up an HTTP client with the provided cookies
func SetupClientWithCookies(cookies []*http.Cookie) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("error creating cookie jar: %w", err)
	}

	client := &http.Client{Jar: jar}

	// Add cookies to jar
	jar.SetCookies(&url.URL{Scheme: "https", Host: "www.setsmart.com"}, cookies)

	return client, nil
}

// MakeAuthenticatedRequest makes an authenticated request using the provided client and cookies
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
