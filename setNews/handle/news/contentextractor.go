package news

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FetchHTMLContent(url string) NewsItem {

	newsItem := NewsItem{URL: url}

	switch {
	case strings.Contains(url, "prachachat"):
		doc, err := fetchPrachachatContent(url)
		if err != nil {
			return newsItem
		}
		newsItem.HTMLContent = ExtractContent(doc)
		newsItem.Title = ExtractTitle(doc)
	case strings.Contains(url, "thunhoon"):
		htmlContent := fetchThunhoonContent(url)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
		if err != nil {
			fmt.Println("Error parsing Thunhoon HTML:", err)
			return newsItem
		}
		newsItem.HTMLContent = ExtractThunhoonContent(htmlContent)
		newsItem.Title = ExtractTitle(doc)
	case strings.Contains(url, "businesstoday"):
		htmlContent := fetchBusinessTodayContent(url)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
		if err != nil {
			fmt.Println("Error parsing Business Tody HTML:", err)
			return newsItem
		}
		newsItem.HTMLContent = ExtractBusinessTodayContent(doc)
		newsItem.Title = ExtractTitle(doc)
	case strings.Contains(url, "moneyandbanking"):
		htmlContent := fetchMoneyAndBankingContent(url)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
		if err != nil {
			fmt.Println("Error parsing Money and Banking HTML:", err)
			return newsItem
		}
		newsItem.HTMLContent = ExtractContent(doc)
		newsItem.Title = ExtractTitle(doc)
	default:
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching HTML:", err)
			return newsItem
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading HTML body:", err)
			return newsItem
		}

		htmlContent := string(body)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
		if err != nil {
			fmt.Println("Error parsing HTML:", err)
			return newsItem
		}
		newsItem.HTMLContent = ExtractContent(doc)
		newsItem.Title = ExtractTitle(doc)
	}
	return newsItem
}

func fetchMoneyAndBankingContent(url string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Cookie", "121995=1; verify=test; _cbclose=1; _cbclose12598=1; _uid12598=48B82399.1; _ctout12598=1; _ga_HNFW6HCLEH=GS1.1.1723022831.1.0.1723022831.60.0.0; _ga=GA1.3.197662334.1723022832; _gid=GA1.3.1202061379.1723022832; _gat_UA-23886484-2=1; _cc_id=b1b3e9c5a08d9757c99ac9a189c1e275; panoramaId_expiry=1723109232332; panoramaId=ba7f945564e00215dd81fe01f3cba9fb927a15eda3c9cbd36b893ee58ea0f8ab; panoramaIdType=panoDevice; __gads=ID=c4168dffc840fc9d:T=1723022832:RT=1723022832:S=ALNI_MYyNfcujXzDb4OvsRCRWj-Dxv1NLg; __gpi=UID=00000ec7b7e808fb:T=1723022832:RT=1723022832:S=ALNI_Ma9fcuhJNsO0Td7pPoLlApmP5xrPw; __eoi=ID=d380205082b71cbf:T=1723022832:RT=1723022832:S=AA-Afjbayk_Dx0m0v8E4or4Llt_L; cto_bundle=tyglWl9jRGhjeFQyYmNBMUVybjc5Wm5yOUdZTGtZRnc5JTJCcGx4Mko4NjlEaU54SnY0OTNnZWQxQ2hGNXdwZ1RKaGxxeWh0ZVFqd2txZDAlMkJMaWtDSyUyRmg1RVZ0bXFSenF6QmdpNWRmdldNeVFhVEFxQ1NESGE3aVhlTnlXdGI1R3E2aWZQRm81UVZ1MUxYVVJMNlJmNDN0OGc2ODRqYW16OHY1UnEycyUyRmJNVEFKdjhIN2JPMjdKUXhNUGN5Uzd6Q3JRenFaTUo3Q2RNZSUyRjhLbE5iRGxzem5nY1N2ZyUzRCUzRA; FCNEC=%5B%5B%22AKsRol9FCnmdNn0moJhM622B9n_tJqAf5mITEU1BKPFiZ_UgAI35zPLHVZJmVR2p8_FckoIKVyxGoJ1sb9elVuELM-N8pBndV1WiFepg48NOPX1ihGOTM9UPKJTwC1GwJ3pCCtGTfz_qqJp1iFAqdKQSLexc0jx02Q%3D%3D%22%5D%5D; cookieyes-consent=consentid:d1NTMFN3bmZDbjd3ZFBiWVAzWDhYNzd3VXptN1IxeUQ,consent:yes,action:yes,necessary:yes,functional:yes,analytics:yes,performance:yes,advertisement:yes")
	req.Header.Set("Priority", "u=0, i")
	req.Header.Set("Referer", url)
	req.Header.Set("Sec-Ch-Ua", "\"Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\", \"Google Chrome\";v=\"126\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching Money and Banking HTML:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading Money and Banking HTML body:", err)
		return ""
	}

	return string(body)
}

func fetchPrachachatContent(url string) (*goquery.Document, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching Prachachat HTML:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading Prachachat HTML body:", err)
		return nil, err
	}

	htmlContent := string(body)
	if strings.Contains(htmlContent, "Enable JavaScript and Cookies") {
		fmt.Println("JavaScript challenge detected. Content may not be fully loaded")
		return nil, fmt.Errorf("JavaScript challenge detected")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		fmt.Println("Error parsing Prachachat HTML:", err)
		return nil, err
	}
	return doc, nil
}

func fetchThunhoonContent(url string) string {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", url, nil)
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

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching Thunhoon HTML:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading Thunhoon HTML body:", err)
		return ""
	}

	htmlContent := string(body)
	if strings.Contains(htmlContent, "Enable JavaScript and Cookies to continue") {
		fmt.Println("JavaScript challenge detected. Content may not be fully loaded")
		return ""
	}
	return htmlContent
}

func ExtractThunhoonContent(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return ""
	}

	content := doc.Find("div.col-lg-8.col-sm-12").Text()
	if content != "" {
		return CleanHTMLContent(content)
	}
	return ""
}

func fetchBusinessTodayContent(url string) string {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Cookie", "cookielawinfo-checkbox-necessary=yes; cookielawinfo-checkbox-non-necessary=yes; _gid=GA1.2.414608578.1722998883; _fbp=fb.1.1722998883662.812601603845803557; CookieLawInfoConsent=eyJuZWNlc3NhcnkiOnRydWUsIm5vbi1uZWNlc3NhcnkiOnRydWV9; viewed_cookie_policy=yes; __gads=ID=cf070b2a4cab58bb:T=1722998883:RT=1723017303:S=ALNI_Mb255-ufB7OuDvUW6HwBZKSB6cOtw; __gpi=UID=00000eb7c58092a0:T=1722998883:RT=1723017303:S=ALNI_MbWwS9_EKDoT4nQQb-7cliJp1ZPVw; __eoi=ID=9ee17185c854a7ab:T=1722998883:RT=1723017303:S=AA-Afja787KsM2-_6lx_8NcHoCk1; _ga=GA1.2.817401315.1722998883; _gat_UA-144169061-1=1; _ga_K4XGKBQ5CR=GS1.1.1723017303.3.1.1723017389.58.0.0")
	req.Header.Set("Priority", "u=0, i")
	req.Header.Set("Sec-Ch-Ua", "\"Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\", \"Google Chrome\";v=\"126\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching Business Today HTML:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading Business Today HTML body:", err)
		return ""
	}

	htmlContent := string(body)
	if strings.Contains(htmlContent, "403 - Forbidden") {
		fmt.Println("Access forbidden detected. Content may not be fully loaded")
		return ""
	}
	return htmlContent
}

func ExtractBusinessTodayContent(doc *goquery.Document) string {
	content := doc.Find("div.tdb-block-inner.td-fix-index").Text()
	if content != "" {
		return CleanHTMLContent(content)
	}
	return ""
}

func ExtractContent(doc *goquery.Document) string {

	patterns := []string{
		"#article .entry-content",
		"#the-post .entry-content",
		"div.article-content",
		"div.post-body",
		"article",
	}

	for _, pattern := range patterns {
		content := doc.Find(pattern).Text()
		if content != "" {
			return CleanHTMLContent(content)
		}
	}

	return CleanHTMLContent(doc.Find("body").Text())
}

func ExtractTitle(doc *goquery.Document) string {
	title := doc.Find("title").Text()
	if title == "" {
		title = doc.Find("h1").Text()
	}
	return title
}

func CleanHTMLContent(content string) string {
	content = html.UnescapeString(content)

	re := regexp.MustCompile(`(?i)(function\s*\(w, d, s, l, i\)\s*{[\s\S]*?event:\s*'gtm\.'\s*})`)
	content = re.ReplaceAllString(content, "")

	re = regexp.MustCompile(`(?i)<(script|style)[^>]*>.*?</\\1>`)
	content = re.ReplaceAllString(content, "")

	re = regexp.MustCompile(`(?i)<[^>]+>`)
	content = re.ReplaceAllString(content, "")

	re = regexp.MustCompile(`\s+`)
	content = re.ReplaceAllString(content, " ")

	content = strings.TrimSpace(content)

	return content
}
