package main

import (
	"bytes"
	"debitcard_fees/models"
	"debitcard_fees/parser"
	"debitcard_fees/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/DebitFee/CompareProductList"
	// Define the request payload
	payloadTemplate := `{"ProductIdList":"1474,1606,1634,1234,1237,1492,1502,1377,1378,1379,1380,1381,1382,1365,1366,1256,976,954,950,961,13,14,67,1459,1466,1460,1467,1468,1463,731,733,720,721,722,723,719,718,726,727,725,724,1490,1473,1476,1475,1478,1477,1482,1481,1484,1480,1483,1488,1485,1487,1585,1641,1587,1593,1618,1627,1642,1592,1612,1594,1235,1239,1236,1240,1504,1491,1496,1503,1501,1494,1499,473,474,475,477,472,246,946,958,962,972,964,960,16,15,17,11,12,682,1457,1471,1465,1472,1461,1462,1469,1456,1470,1464,1458,1,2,744,746,750,751,748,730,732,749,752,138,1601,1649,1500,1489,1497,1493,1498,139,1479,1486,1238,1369,1367,728,729,1495","Page":%d,"Limit":3}`

	initialPayload := fmt.Sprintf(payloadTemplate, 1)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(initialPayload)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	utils.AddHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	totalPages := utils.DetermineTotalPage(doc)
	if totalPages == 0 {
		log.Fatal("Could not determine the total number of pages")
	}

	var debitFees []models.DebitFee

	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages)
		payload := fmt.Sprintf(payloadTemplate, page)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			log.Printf("Error creating request for page %d: %v", page, err)
			continue
		}

		utils.AddHeader(req)

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error making request for page %d: %v", page, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Request for page %d failed with status: %d\n", page, resp.StatusCode)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response for page %d: %v", page, err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			log.Printf("Error parsingg HTML for page %d: %v", page, err)
			continue
		}

		for i := 1; i <= 3; i++ {
			col := "col" + strconv.Itoa(i)
			provider := utils.CleanText(doc.Find(fmt.Sprintf("th.%s span", col)).Text())
			product := utils.CleanText(doc.Find(fmt.Sprintf("th.prod-col%d span", i)).Text())

			generalFees := parser.ParseGeneralFees(doc, i)
			domesticFees := parser.ParseDomesticFees(doc, i)
			internationalFees := parser.ParseInternationalFees(doc, i)
			otherFees := parser.ParseOtherFees(doc, i)
			additionalInfo := parser.ParseAdditionalInfo(doc, i)

			debitFees = append(debitFees, models.DebitFee{
				Provider:          provider,
				Product:           product,
				GeneralFees:       generalFees,
				DomesticFees:      domesticFees,
				InternationalFees: internationalFees,
				OtherFees:         otherFees,
				AdditionalInfo:    additionalInfo,
			})
		}
	}

	jsonData, err := json.MarshalIndent(debitFees, "", " ")
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}

	err = os.WriteFile("debit_fees.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing JSON to file: %v", err)
	}

	fmt.Println("Data successfully written to debit_fees.json")
}
