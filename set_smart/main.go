package main

import (
	"fmt"
	"login_token/functions"
)

func main() {
	// allAPIURL := "https://www.setsmart.com/api/search/stock/all"
	// activeAPIURL := "https://www.setsmart.com/api/search/stock/active"
	outputFileName := "stocks_all.json"

	// err := functions.UpdateAndSaveSymbolsStatus(allAPIURL, activeAPIURL, outputFileName)
	// if err != nil {
	// 	fmt.Println("error updating symbols status:", err)
	// 	return
	// }

	locales := []string{"en_US", "th_TH"}

	symbols, err := functions.ReadSymbolsFromFile(outputFileName)
	if err != nil {
		fmt.Println("error reading symbols from file:", err)
		return
	}

	err = functions.MakePostRequests(symbols, locales)
	if err != nil {
		fmt.Println("error making POST request:", err)
	}

	// accessToken, cookies, err := utils.FetchAccessToken()
	// if err != nil {
	// 	fmt.Println("error: %w", err)
	// 	return
	// }

	// fmt.Println("Access Token: ", accessToken)
	// fmt.Println("Cookies:")
	// for _, cookie := range cookies {
	// 	fmt.Printf("%s: %s\n", cookie.Name, cookie.Value)
	// }
}
