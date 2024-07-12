package main

import (
	"fmt"
	MajorShareHolder "login_token/majorShareHolder"
	"login_token/utils"
)

func main() {
	loginURL := "https://www.setsmart.com/api/user/login"
	username := "aowjingti@gmail.com"
	password := "Zxcvasdf1!"

	cookieStr, _, err := utils.PerformLogin(loginURL, username, password)
	if err != nil {
		fmt.Println("Error during login:", err)
		return
	}

	postURL := "https://www.setsmart.com/ism/majorshareholder.html"
	formData := "radChoice=1&txtSymbol=scb&radShow=2&submit.x=21&submit.y=11&hidAction=go&hidLastContentType="
	postResponse, err := MajorShareHolder.FetchDetailedShareholderData(postURL, cookieStr, formData)
	if err != nil {
		fmt.Println("Error fetching detailed shareholder data: ", err)
		return
	}

	fmt.Println("Detailed shareholder data response:", postResponse)
}
