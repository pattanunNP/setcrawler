package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	parchange "login_token/Parchange"
	capitalmovement "login_token/capital_movement"
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

	allResults := make(map[string]utils.OrganizedData)
	allmovements := make(map[string]map[string]capitalmovement.CapitalMovement)

	locales := []string{"en_EN", "th_TH"}
	counter := 0

	parChanges, err := parchange.GetParChange(cookieStr, "", "E", "", "1/01/2019", "17/07/2024")
	if err != nil {
		fmt.Printf("Par change request error: %v", err)
		return
	}

	for _, stock := range stocksData.Stock {
		if counter >= 5 {
			break
		}
		// Company Profile
		organizedData := utils.OrganizedData{
			CompanyProfiles: make(map[string]utils.CompanyProfile),
			Securities:      make(map[string]utils.Securities),
		}

		// Capital Movement
		movements := make(map[string]capitalmovement.CapitalMovement)

		for _, locale := range locales {

			//Capital Movement
			movement, err := capitalmovement.GetCapitalMovement(cookieStr, stock.Symbol, locale)
			if err != nil {
				fmt.Printf("Capital movement request error for symbol %s, locale %s: %v", stock.Symbol, locale, err)
				continue
			}
			localeSuffix := strings.Split(locale, "_")[1]
			movementKey := fmt.Sprintf("%s_%s", stock.Symbol, localeSuffix)
			movements[movementKey] = *movement

			// Company_Profile Reqiest
			companyName, companyProfile, securities, err := utils.MakeRequestWithCookies(cookieStr, stock.Symbol, locale)
			if err != nil {
				fmt.Printf("Request error for symbol %s, locale %s: %v\n", stock.Symbol, locale, err)
			} else {
				localeSuffix := strings.Split(locale, "_")[1]
				profileKey := fmt.Sprintf("%s_%s", companyName, localeSuffix)
				organizedData.CompanyProfiles[profileKey] = companyProfile
				organizedData.Securities[profileKey] = securities
				fmt.Printf("Request successful for symbol %s, locale %s\n", stock.Symbol, locale)
			}

			// Parchange

		}

		allResults[stock.Symbol] = organizedData
		allmovements[stock.Symbol] = movements
		counter++
	}

	// Save all results to a single file
	file, err := os.Create("company_profile.json")
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

	err = capitalmovement.SaveToFile("capital_movements.json", allmovements)
	if err != nil {
		fmt.Println("error saving capital movments JSON file:", err)
		return
	}

	err = parchange.SaveToFile("par_changes.json", parChanges)
	if err != nil {
		fmt.Println("Error saving par changes JSON file:", err)
		return
	}

	//fmt.Println("Data successfully written to all_results.json")
}
