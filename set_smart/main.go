package main

import (
	"encoding/json"
	"fmt"
	"io"
	"login_token/utils"
	"net/http"
)

func main() {
	// Setup client with cookies and access token
	client, err := utils.SetupCookieAndToken()
	if err != nil {
		fmt.Println("Error setting up cookies and token:", err)
		return
	}

	// Create a GET request to the protected endpoint
	apiURL := "https://www.setsmart.com/api/search/stock/all"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Unmarshal the JSON response for pretty printing and saving to a file
	var jsonResponse interface{}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return
	}

	// Pretty print the JSON response
	prettyJSON, err := json.MarshalIndent(jsonResponse, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	fmt.Println(string(prettyJSON))

	// Save the JSON response to a file
	// err = os.WriteFile("response.json", prettyJSON, 0644)
	// if err != nil {
	// 	fmt.Println("Error writing JSON to file:", err)
	// 	return
	// }

	// fmt.Println("Response saved to response.json")
}
