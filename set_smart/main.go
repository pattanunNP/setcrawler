package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	parchange "login_token/Parchange"
	capitalmovement "login_token/capital_movement"
	"login_token/finance"
	"login_token/fund"
	historicalnews "login_token/historical_news"
	MajorShareHolder "login_token/majorShareHolder"
	"login_token/management"
	"login_token/news"
	"login_token/nvdr"
	"login_token/tfex"
	"login_token/utils"
)

func main() {
	loginURL := "https://www.setsmart.com/api/user/login"
	username := "aowjingti@gmail.com"
	password := "Zxcvasdf1!"

	getAuthDetails := func() (string, string, error) {
		return utils.PerformLogin(loginURL, username, password)
	}

	cookieStr, tokenStr, err := getAuthDetails()
	if err != nil {
		log.Fatalf("Login error: %v", err)
	}

	stocksData, err := utils.ReadJSONFile("/Users/natpisitkao/Desktop/setcrawler/set_smart/stocks_all.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	fundsData, err := utils.ReadFundsJSONFile("/Users/natpisitkao/Desktop/setcrawler/set_smart/fund_list.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %w", err)
	}

	allResults := make(map[string]utils.OrganizedData)
	allMovements := make(map[string]map[string]capitalmovement.CapitalMovement)
	var allParChanges []parchange.ParChange
	allShareHolders := make(map[string]map[string]MajorShareHolder.CombinedShareHolderData)
	allNews := make(map[string]map[string][]news.NewsItem)
	allHistoricalNews := make(map[string]map[string][]historicalnews.NewsItem)
	allManagementData := make(map[string]map[string][]management.ManagementData)
	allHisNVDR := make(map[string]map[string]nvdr.NVDRData)
	allTFEXData := make(map[string]tfex.TFEXData)
	var allFinancialData []finance.FinancialData
	allCombinedData := make(map[string]interface{})

	locales := []string{"en_US", "th_TH"}
	// counter := 0

	for _, stock := range stocksData.Stock {
		// if counter >= 10 {
		// 	break
		// }
		organizedData := utils.OrganizedData{
			CompanyProfiles: make(map[string]utils.CompanyProfile),
			Securities:      make(map[string]utils.Securities),
		}
		movements := make(map[string]capitalmovement.CapitalMovement)
		shareHolders := make(map[string]MajorShareHolder.CombinedShareHolderData)
		stockNews := make(map[string][]news.NewsItem)
		hisNews := make(map[string][]historicalnews.NewsItem)
		stockManagementData := make(map[string][]management.ManagementData)

		for _, locale := range locales {
			localeParts := strings.Split(locale, "_")
			if len(localeParts) < 2 {
				log.Printf("Invalid locale format: %s\n", locale)
				continue
			}
			localeSuffix := localeParts[1]
			movementKey := fmt.Sprintf("%s_%s", stock.Symbol, localeSuffix)

			// 		// Capital Movement
			movement, err := capitalmovement.GetCapitalMovement(cookieStr, stock.Symbol, locale)
			if err == nil {
				movements[movementKey] = *movement
				time.Sleep(5000 * time.Millisecond)
			} else {
				fmt.Printf("Capital movement request error for symbol %s, locale %s: %v\n", stock.Symbol, locale, err)
			}

			// 		// Company Profile
			companyName, companyProfile, securities, err := utils.MakeRequestWithCookies(cookieStr, stock.Symbol, locale)
			if err == nil {
				profileKey := fmt.Sprintf("%s_%s", companyName, localeSuffix)
				organizedData.CompanyProfiles[profileKey] = companyProfile
				organizedData.Securities[profileKey] = securities
				fmt.Printf("Request successful for symbol %s, locale %s\n", stock.Symbol, locale)
				time.Sleep(5000 * time.Millisecond)
			} else {
				fmt.Printf("Request error for symbol %s, locale %s: %v\n", stock.Symbol, locale, err)
			}

			// 		// ShareHolders
			shareHoldersData, err := MajorShareHolder.GetMajorShareHoldersAndDetails(cookieStr, stock.Symbol, locale)
			if err == nil {
				shareHolders[localeSuffix] = shareHoldersData
				time.Sleep(5000 * time.Millisecond)
			} else {
				fmt.Printf("Major shareholder request error for symbol %s, locale %s: %v\n", stock.Symbol, locale, err)
			}

			// 		// News
			newsItems, err := news.FetchNews(cookieStr, stock.Symbol, locale)
			if err == nil {
				stockNews[locale] = newsItems
				time.Sleep(5000 * time.Millisecond)
			} else {
				fmt.Printf("News request error for symbol %s, locale %s: %v\n", stock.Symbol, locale, err)
			}

			// 		// Historical News
			historicalNewsItems, err := historicalnews.FetchHistoricalNews(cookieStr, stock.Symbol, locale)
			if err == nil {
				hisNews[locale] = historicalNewsItems
				time.Sleep(5000 * time.Millisecond)
			} else {
				fmt.Printf("Historical news request error for symbol %s, locale %s: %v\n", stock.Symbol, locale, err)
			}

			// // Management Data
			managementData, err := management.FetchManagementHTML(cookieStr, locale, stock.Symbol)
			if err == nil {
				stockManagementData[locale] = managementData
				time.Sleep(5000 * time.Millisecond)
			} else {
				fmt.Printf("Error fetching management data for symbol %s, locale %s: %v\n", stock.Symbol, locale, err)
			}

			// 		// NVDR Data
			fmt.Printf("Fetching NVDR data for symbol: %s, locale: %s\n", stock.Symbol, locale)
			stockNVDRData, err := nvdr.FetchStockNVDRData(cookieStr, stock.Symbol, locale)
			if err == nil {
				nvdrKey := fmt.Sprintf("%s_%s", stock.Symbol, locale)
				allHisNVDR[nvdrKey] = stockNVDRData
			} else {
				fmt.Printf("Error fetching NVDR data for symbol %s and locale %s: %v\n", stock.Symbol, locale, err)
			}

			// Fetch Financial Data
			status, err := finance.FetchFinancialData(cookieStr, tokenStr, stock.Symbol)
			if err != nil {
				log.Printf("Error fetching financial data for symbol %s: %v", stock.Symbol, err)
				continue
			}
			time.Sleep(5 * time.Second)

			if status {
				financialData, err := fetchWithRetry(cookieStr, tokenStr, stock.Symbol, getAuthDetails)
				if err != nil {
					log.Printf("Error fetching latest financial statement for symbol %s: %v", stock.Symbol, err)
					continue
				}

				allFinancialData = append(allFinancialData, financialData)
				time.Sleep(5 * time.Second)
			} else {
				log.Printf("Financial status is not true for symbol %s", stock.Symbol)
			}

		}

		for _, fundData := range fundsData {
			// if counter >= 10 {
			// 	break
			// }
			//Fetch FundData
			asOfDate, err := fund.FetchFundStatistics(cookieStr, tokenStr, fundData.Symbol)
			if err != nil {
				log.Printf("Error fetching statistcs for symbol %s: %v", fundData.Symbol, err)
				continue
			}

			performanceData, err := fund.FetchFundPerformance(cookieStr, tokenStr, fundData.Symbol, asOfDate)
			if err != nil {
				log.Printf("Error fetching performance for symbol %s: %v", fundData.Symbol, err)
				continue
			}

			combinedData := map[string]interface{}{
				"symbol":           fundData.Symbol,
				"asOfDate":         asOfDate,
				"performance_data": performanceData,
			}

			allCombinedData[fundData.Symbol] = combinedData

			// counter++
			time.Sleep(5 * time.Second)
		}

		allResults[stock.Symbol] = organizedData
		allMovements[stock.Symbol] = movements
		allShareHolders[stock.Symbol] = shareHolders
		allNews[stock.Symbol] = stockNews
		allHistoricalNews[stock.Symbol] = hisNews
		allManagementData[stock.Symbol] = stockManagementData
	}

	// Par Changes
	for _, locale := range locales {
		parChanges, err := parchange.GetParChange(cookieStr, "", "E", "", "01/01/2019", time.Now().Format("02/01/2006"), locale)
		if err == nil {
			allParChanges = append(allParChanges, parChanges.ParChange...)
		} else {
			fmt.Printf("Par change request error for locale %s: %v\n", locale, err)
		}
	}

	// NVDR data
	nvdrData, err := nvdr.FetchNVDRData(cookieStr)
	if err != nil {
		log.Fatalf("Error fetching NVDR data: %v", err)
	}

	// TFEX data
	for _, locale := range locales {
		tfexData, err := tfex.FetchTFEXData(cookieStr, locale)
		if err == nil {
			allTFEXData[locale] = tfexData
		} else {
			fmt.Printf("Error fectching TFEX data for locale %s: %v\n", locale, err)
		}
	}

	// Save all results to files
	saveJSON("company_profile.json", allResults)
	saveJSON("capital_movements.json", allMovements)
	saveJSON("par_changes.json", allParChanges)
	saveJSON("major_shareholders.json", allShareHolders)
	saveJSON("news_data.json", allNews)
	saveJSON("historical_news_data.json", allHistoricalNews)
	saveJSON("management.json", allManagementData)
	saveJSON("nvdr_trading.json", nvdrData)
	saveJSON("nvdr_stock_trading_all.json", allHisNVDR)
	saveJSON("tfex_data.json", allTFEXData)
	saveJSON("financial_data.json", allFinancialData)
	saveJSON("fund_data.json", allCombinedData)
}

func fetchWithRetry(cookiStr, tokenStr, symbol string, getAuthDetails func() (string, string, error)) (finance.FinancialData, error) {
	const maxRetries = 3

	for i := 0; i < maxRetries; i++ {
		finacialData, err := finance.FetchLatestFinancialStatement(cookiStr, tokenStr, symbol)
		if err == nil {
			return finacialData, nil
		}

		if i < maxRetries-1 {
			if isUnauthorizedError(err) {
				log.Printf("Received 401 error, re-authenticating... (attempt %d/%d)", i+1, maxRetries)
				var authErr error
				cookiStr, tokenStr, authErr = getAuthDetails()
				if authErr != nil {
					return finance.FinancialData{}, fmt.Errorf("re-authentication failed: %v", authErr)
				}
			} else {
				log.Printf("Retrying after errorL %v (attempt %d/%d)", err, i+1, maxRetries)
				time.Sleep(5 * time.Second)
			}
		} else {
			return finance.FinancialData{}, err
		}
	}
	return finance.FinancialData{}, fmt.Errorf("max retries reached for symbol %s", symbol)
}

func isUnauthorizedError(err error) bool {
	return fmt.Sprintf("%v", err) == "received non-200 response code: 401"
}

func saveJSON(filename string, data interface{}) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating JSON file %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(data)
	if err != nil {
		fmt.Printf("Error encoding JSON data to file %s: %v\n", filename, err)
	}
}
