package MajorShareHolder

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func FetchDetailedShareholderData(postURL, cookieStr, formData string) (string, error) {
	req, err := http.NewRequest("POST", postURL, bytes.NewBufferString(formData))
	if err != nil {
		return "", fmt.Errorf("error creating POST request: %w", err)
	}

	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Add("cache-control", "max-age=0")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cookie", cookieStr)
	req.Header.Add("origin", "https://www.setsmart.com")
	req.Header.Add("priority", "u=0, i")
	req.Header.Add("referer", "https://www.setsmart.com/ism/majorshareholder.html?locale=en_US")
	req.Header.Add("sec-ch-ua", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "macOS")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error maing POST requesy: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}
