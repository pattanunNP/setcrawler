package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	historicalnews "login_token/historical_news"
	"login_token/management"
	"login_token/news"
	"login_token/utils"
)

func main() {
	loginURL := "https://www.setsmart.com/api/user/login"
	username := "aowjingti@gmail.com"
	password := "Zxcvasdf1!"

	cookieStr, _, err := utils.PerformLogin(loginURL, username, password)
	if err != nil {
		fmt.Println("Login error:", err)
		return
	}

	stocksData, err := utils.ReadJSONFile("/Users/natpisitkao/Desktop/setcrawler/set_smart/stocks_all.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	// allResults := make(map[string]utils.OrganizedData)
	// allmovements := make(map[string]map[string]capitalmovement.CapitalMovement)
	// var allParChanges []parchange.ParChange
	// allShareHolders := make(map[string]map[string]MajorShareHolder.CombinedShareHolderData)
	allNews := make(map[string]map[string][]news.NewsItem)
	allHistoricalNews := make(map[string]map[string][]historicalnews.NewsItem)
	allManagementData := make(map[string]map[string][]management.ManagementData)

	locales := []string{"en_US", "th_TH"}
	counter := 0

	// for _, stock := range stocksData.Stock {
	// 	if counter >= 10 {
	// 		break
	// 	}
	// 	// Company Profile
	// 	organizedData := utils.OrganizedData{
	// 		CompanyProfiles: make(map[string]utils.CompanyProfile),
	// 		Securities:      make(map[string]utils.Securities),
	// 	}

	// 	// Capital Movement
	// 	movements := make(map[string]capitalmovement.CapitalMovement)
	// 	shareHolders := make(map[string]MajorShareHolder.CombinedShareHolderData)
	stockNews := make(map[string][]news.NewsItem)

	// 	for _, locale := range locales {

	// 		//Capital Movement
	// 		movement, err := capitalmovement.GetCapitalMovement(cookieStr, stock.Symbol, locale)
	// 		if err != nil {
	// 			fmt.Printf("Capital movement request error for symbol %s, locale %s: %v", stock.Symbol, locale, err)
	// 			continue
	// 		}
	// 		localeSuffix := strings.Split(locale, "_")[1]
	// 		movementKey := fmt.Sprintf("%s_%s", stock.Symbol, localeSuffix)
	// 		movements[movementKey] = *movement

	// 		// Company_Profile Reqiest
	// 		companyName, companyProfile, securities, err := utils.MakeRequestWithCookies(cookieStr, stock.Symbol, locale)
	// 		if err != nil {
	// 			fmt.Printf("Request error for symbol %s, locale %s: %v\n", stock.Symbol, locale, err)
	// 		} else {
	// 			localeSuffix := strings.Split(locale, "_")[1]
	// 			profileKey := fmt.Sprintf("%s_%s", companyName, localeSuffix)
	// 			organizedData.CompanyProfiles[profileKey] = companyProfile
	// 			organizedData.Securities[profileKey] = securities
	// 			fmt.Printf("Request successful for symbol %s, locale %s\n", stock.Symbol, locale)
	// 		}

	// 		// ShareHolders
	// 		shareHoldersData, err := MajorShareHolder.GetMajorShareHoldersAndDetails(cookieStr, stock.Symbol, locale)
	// 		if err != nil {
	// 			fmt.Printf("Major shareholder request error for symbol %s, locale %s: %v", stock.Symbol, locale, err)
	// 			continue
	// 		}

	// 		shareHolders[localeSuffix] = shareHoldersData

	// 	}

	// 	allResults[stock.Symbol] = organizedData
	// 	allmovements[stock.Symbol] = movements
	// 	allShareHolders[stock.Symbol] = shareHolders

	// 	counter++
	// }

	for _, stock := range stocksData.Stock {
		if counter >= 1 {
			break
		}

		HisNews := make(map[string][]historicalnews.NewsItem)
		stockManagementData := make(map[string][]management.ManagementData)

		for _, locale := range locales {

			// News
			newsItems, err := news.FetchNews(cookieStr, stock.Symbol, locale)
			if err != nil {
				fmt.Printf("News request error for symbol %s, locale %s: %v", stock.Symbol, locale, err)
				continue
			}
			stockNews[locale] = newsItems

			fmt.Printf("Fetching historical news for symbol: %s, locale: %s\n", stock.Symbol, locale)
			historicalNews, err := historicalnews.FetchHistoricalNews(cookieStr, stock.Symbol, locale)
			if err != nil {
				fmt.Printf("Historical news request error for symbol %s, locale %s: %v\n", stock.Symbol, locale, err)
				continue
			}
			HisNews[locale] = historicalNews
			time.Sleep(500 * time.Millisecond)

			fmt.Printf("Fetching management data for symbol: %s, locale: %s\n", stock.Symbol, locale)
			managementData, err := management.FetchManagementHTML(cookieStr, locale, stock.Symbol)
			if err != nil {
				fmt.Printf("Error fetching management data for symbol %s, locale %s: %v\n", stock.Symbol, locale, err)
				continue
			}
			stockManagementData[locale] = managementData
			time.Sleep(500 * time.Millisecond)
		}

		allNews[stock.Symbol] = stockNews
		allHistoricalNews[stock.Symbol] = HisNews
		allManagementData[stock.Symbol] = stockManagementData
		counter++
	}

	// for _, locale := range locales {
	// 	parChanges, err := parchange.GetParChange(cookieStr, "", "E", "", "01/01/2019", "17/07/2024", locale)
	// 	if err != nil {
	// 		fmt.Printf("Par change request error for locale %s: %v", locale, err)
	// 		continue
	// 	}
	// 	allParChanges = append(allParChanges, parChanges.ParChange...)
	// }

	// Save all results to a single file
	// saveJSON("company_profile.json", allResults)
	// saveJSON("capital_movements.json", allmovements)
	// saveJSON("par_changes.json", allParChanges)
	// saveJSON("major_shareholders.json", allShareHolders)
	saveJSON("news_data.json", allNews)
	saveJSON("historical_news_data.json", allHistoricalNews)
	saveJSON("management.json", allManagementData)
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
