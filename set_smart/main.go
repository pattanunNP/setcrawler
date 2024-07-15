package main

import (
	"encoding/json"
	"fmt"
	capitalmovement "login_token/capital_movement"
	"login_token/utils"
	"os"
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

	allResults := make(map[string]utils.OrganizedData)
	locales := []string{"en_EN", "th_TH"}
	counter := 0

	for _, stock := range stocksData.Stock {
		if counter >= 1 {
			break
		}
		// organizedData := utils.OrganizedData{
		// 	CompanyProfiles: make(map[string]utils.CompanyProfile),
		// 	Securities:      make(map[string]utils.Securities),
		// }
		for _, locale := range locales {
			response, err := capitalmovement.MakeCapitalMovementRequest(cookieStr, stock.Symbol, locale)
			if err != nil {
				fmt.Printf("Request error for symbol %s, locale %s: %v\n", stock.Symbol, locale, err)
				continue
			}
			fmt.Println(response)
			// 	companyName, companyProfile, securities, err := utils.MakeRequestWithCookies(cookieStr, stock.Symbol, locale)
			// 	if err != nil {
			// 		fmt.Printf("Request error for symbol %s, locale %s: %v\n", stock.Symbol, locale, err)
			// 	} else {
			// 		localeSuffix := strings.Split(locale, "_")[1]
			// 		profileKey := fmt.Sprintf("%s_%s", companyName, localeSuffix)
			// 		organizedData.CompanyProfiles[profileKey] = companyProfile
			// 		organizedData.Securities[profileKey] = securities
			// 		fmt.Printf("Request successful for symbol %s, locale %s\n", stock.Symbol, locale)
			// 	}
			// }
			// allResults[stock.Symbol] = organizedData
			counter++
		}
	}

	// Save all results to a single file
	file, err := os.Create("all_results.json")
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
		return
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(allResults, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON: ", err)
		return
	}

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing to JSON file: ", err)
		return
	}

}
