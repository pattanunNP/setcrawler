package main

import (
	"bytes"
	"emoney_fees/models"
	"emoney_fees/parser"
	"emoney_fees/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// URL and data to be sent in the POST request
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/EMoneyFee/CompareProductList"
	// payloadTemplate := `{"ProductIdList":"10658,10669,10670,10664,10675,10677,10644,10653,10690,10689,10684,10683,10687,10688,10260,93,104,10523,10239,10150,10151,10485,10397,10392,10474,29,10107,10106,10105,10108,85,57,54,56,10584,92,10220,10294,10292,10228,10227,10219,10153,10152,28,31,10674,10660,10662,10234,10235,30,10398,10484,10486,55,10668,60,10685,10686,53,10291,10632,10609,10619,10643,10290,10296,10637,10616,58","Page":%d,"Limit":3}`
	payloadTemplate := `{"ProductIdList":"10658,10669,10670,10664,10675,10677,10644,10653,10690,10689,10684,10683,10687,10688,10260,93,104,10523","Page":%d,"Limit":3}`

	initialpayload := fmt.Sprintf(payloadTemplate, 1)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(initialpayload)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	utils.AddHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	totalPages := utils.DetermineTotalPage(doc)
	if totalPages == 0 {
		fmt.Println("Could not determine the total number of pages")
		return
	}
	fmt.Printf("Total pages to process: %d\n", totalPages)

	var fees []models.EmoneyFee

	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages)

		payload := fmt.Sprintf(payloadTemplate, page)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			fmt.Printf("Error creating request for page %d: %v\n", page, err)
			continue
		}

		utils.AddHeader(req)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request for page %d: %v\n", page, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Request for page %d failed with status: %d\n", page, resp.StatusCode)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response for page %d: %v\n", page, err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(body))
		if err != nil {
			fmt.Printf("Error parsing HTML for page %d: %v\n", page, err)
			continue
		}
		for i := 1; i <= 3; i++ {
			col := "col" + strconv.Itoa(i)
			provider := utils.CleanText(doc.Find(fmt.Sprintf("th.%s span", col)).Text())
			product := utils.CleanText(doc.Find(fmt.Sprintf("th.prod-col%d span", i)).Text())

			topUpDetails := parser.ParseTopUpDetails(doc, i)
			generalFees := parser.ParseGeneralFees(doc, i)
			spendingFes := parser.ParseSpendingFees(doc, i)
			terminationFees := parser.ParseTerminationFees(doc, i)
			otherFees := parser.ParseOtherFes(doc, i)
			additionalInfo := parser.AdditionalInfo(doc, i)

			fees = append(fees, models.EmoneyFee{
				Provider:              provider,
				Product:               product,
				TopUp:                 topUpDetails,
				GeneralFees:           generalFees,
				SpendingFees:          spendingFes,
				TerminationFees:       terminationFees,
				OtherFees:             otherFees,
				AdditionalInformation: additionalInfo,
			})
		}

		time.Sleep(3 * time.Second)

	}

	jsonData, err := json.MarshalIndent(fees, "", " ")
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return
	}

	err = os.WriteFile("emoney_fees.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Data sucessfully written to emoney_fees.json")
}
