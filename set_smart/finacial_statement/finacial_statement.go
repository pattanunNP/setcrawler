package financial

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// FetchFinancialStatements makes an HTTP GET request to the given URL and returns the response body as a string.
func FetchFinancialStatements(cookieStr string, symbol string) (string, error) {
	url := fmt.Sprintf("https://www.setsmart.com/ssm/financialStatement?symbol=%s", symbol)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Referer", "https://www.setsmart.com/ism/nvdrTrading.html")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("cookie", cookieStr)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(body), nil
}

func DownloadFinancialStatements(cookieStr string) (string, error) {
	url := "https://www.setsmart.com/ssm/financialStatement"

	// Create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Create a temporary directory to store the downloaded file
	tmpDir, err := ioutil.TempDir("", "download")
	if err != nil {
		return "", fmt.Errorf("error creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir) // Clean up

	// Configure Chrome to use the temporary directory for downloads
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // Set to true for headless mode
		chromedp.UserDataDir(tmpDir),     // Set the download directory
	)

	// Create a new allocator context with the options
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Create a new Chrome context
	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	// Add cookies to the context
	if err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		for _, cookie := range parseCookies(cookieStr) {
			chromedp.SetCookie(cookie.Name, cookie.Value).
				WithDomain(cookie.Domain).
				WithPath(cookie.Path).
				WithExpires(cookie.Expires).
				WithHTTPOnly(cookie.HTTPOnly).
				Do(ctx)
		}
		return nil
	})); err != nil {
		return "", fmt.Errorf("error setting cookies: %v", err)
	}

	// Run tasks to navigate to the URL and click the export link
	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#export-link`, chromedp.ByID),
		chromedp.Click(`#export-link`, chromedp.ByID),
		chromedp.Sleep(10 * time.Second), // Wait for the download to complete
	}); err != nil {
		return "", fmt.Errorf("error running tasks: %v", err)
	}

	// Find the downloaded file
	files, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		return "", fmt.Errorf("error reading temporary directory: %v", err)
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no files downloaded")
	}

	downloadedFile := files[0]
	downloadedFilePath := tmpDir + "/" + downloadedFile.Name()

	// Move the downloaded file to the current directory
	destPath := "./" + downloadedFile.Name()
	err = os.Rename(downloadedFilePath, destPath)
	if err != nil {
		return "", fmt.Errorf("error moving downloaded file: %v", err)
	}

	return destPath, nil
}

func parseCookies(cookieStr string) []*http.Cookie {
	cookies := []*http.Cookie{}
	cookiePairs := strings.Split(cookieStr, ";")
	for _, pair := range cookiePairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}
		cookies = append(cookies, &http.Cookie{
			Name:  parts[0],
			Value: parts[1],
		})
	}
	return cookies
}
