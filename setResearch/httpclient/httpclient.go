package httpclient

import "net/http"

func NewClient() *http.Client {
	return &http.Client{}
}

func SetRequestHeaders() http.Header {
	headers := http.Header{}
	headers.Set("accept", "application/json, text/plain, */*")
	headers.Set("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	headers.Set("priority", "u=1, i")
	headers.Set("referer", "https://www.settrade.com/th/research/analyst-research/main")
	headers.Set("sec-ch-ua", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
	headers.Set("sec-ch-ua-mobile", "?0")
	headers.Set("sec-ch-ua-platform", "macOS")
	headers.Set("sec-fetch-dest", "empty")
	headers.Set("sec-fetch-mode", "cors")
	headers.Set("sec-fetch-site", "same-origin")
	headers.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	headers.Set("x-channel", "WEB_SETTRADE")
	headers.Set("x-client-uuid", "7b3565cc-e330-4025-b664-c81d166840f3")
	return headers
}
