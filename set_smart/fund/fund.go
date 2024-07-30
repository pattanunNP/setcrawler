package fund

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

type FundData struct {
	Symbol             string `json:"symbol"`
	FundName           string `json:"fund_name"`
	AMC                string `json:"amc"`
	FundClassification string `json:"fund_classification"`
	ManagementStyle    string `json:"management_style"`
	DividendPolicy     string `json:"dividend_policy"`
	Risk               string `json:"risk"`
	LTFRMF             string `json:"ltf_rmf"`
	IPOStartDate       string `json:"ipo_start_date"`
	IPOEndDate         string `json:"ipo_end_date"`
	Factsheet          string `json:"factsheet"`
}

func FetchFundTableData(cookieStr string) ([]FundData, error) {
	// Create context with timeout
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()

	// Increase timeout to 120 seconds
	ctx, cancel = context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	// Create a slice to hold the fund data
	var fundData []FundData

	// Debug information to check if page content is loaded
	var tableContent string

	// Run chromedp tasks
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.setsmart.com/ssm/fundIPO"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Set the cookie using JavaScript
			err := chromedp.Evaluate(fmt.Sprintf(`document.cookie = "%s";`, cookieStr), nil).Do(ctx)
			if err != nil {
				return err
			}
			fmt.Println("Cookie set successfully")
			return nil
		}),
		chromedp.Sleep(10*time.Second),                                   // Wait for the page to load completely
		chromedp.WaitVisible("tr.ng-star-inserted", chromedp.ByQueryAll), // Wait for table rows to be visible
		chromedp.OuterHTML("body", &tableContent, chromedp.ByQuery),      // Capture the table content for debugging
		chromedp.Evaluate(`(() => {
			const rows = Array.from(document.querySelectorAll('tr.ng-star-inserted'));
			if (rows.length === 0) {
				console.log('No rows found');
				return [];
			}
			return rows.map(row => {
				const cells = row.querySelectorAll('td');
				return {
					symbol: cells[1]?.innerText.trim(),
					fundName: cells[2]?.innerText.trim(),
					amc: cells[3]?.innerText.trim(),
					fundClassification: cells[4]?.innerText.trim(),
					managementStyle: cells[5]?.innerText.trim(),
					dividendPolicy: cells[6]?.innerText.trim(),
					risk: cells[7]?.innerText.trim(),
					ltfRmf: cells[8]?.innerText.trim(),
					ipoStartDate: cells[9]?.innerText.trim(),
					ipoEndDate: cells[10]?.innerText.trim(),
					factsheet: cells[11]?.querySelector('a')?.href.trim()
				};
			});
		})()`, &fundData),
	)
	if err != nil {
		return nil, err
	}
	sdsegsdsgs
	// Print table content for debugging purposes
	fmt.Println("Table Content:", tableContent)

	return fundData, nil
}
