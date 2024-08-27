package main

import (
	"bath_net/pkg"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/BAHTNETFee/CompareProductList"
	payloadTemplate := `{"ProductIdList":"162152,2,5,17,4,157479,23,26,15,449176,27,6,150920,194031,9,162151,240,35,16,162568,13,237,28,32,241,163222,33,34,37,155024,30,24","Page":%d,"Limit":3}`

	initialPayload := fmt.Sprintf(payloadTemplate, 1)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(initialPayload)))
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	pkg.AddHeaders(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	totalPages := pkg.DetermineTotalPage(doc)
	if totalPages == 0 {
		fmt.Println("Could not determine the total number of pages")
		return
	}

	var bathnetFees []pkg.BathnetFee

	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages)
		payload := fmt.Sprintf(payloadTemplate, page)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			fmt.Println("Error creating request for page", page, ":", err)
			continue
		}

		// Set headers
		pkg.AddHeaders(req)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request for page", page, ":", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Request for page %d failed with status: %d\n", page, resp.StatusCode)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body for page", page, ":", err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			fmt.Println("Error parsing HTML for page", page, ":", err)
			continue
		}

		for i := 1; i <= 3; i++ {
			col := "col" + strconv.Itoa(i)
			provider := doc.Find("th.attr-header.attr-prod.font-black.text-center.cmpr-col." + col + " span").Last().Text()

			feeDetails := pkg.FeeDetails{
				TransferWithinBangkokAndVicinity:             pkg.ExtractFeeData(doc, "attr-SenderRecipientInBkk", col),
				TransferFromBangkokToRegion:                  pkg.ExtractFeeData(doc, "attr-SenderInBkk", col),
				TransferFromRegionToBangkok:                  pkg.ExtractFeeData(doc, "attr-SenderInRegion", col),
				TransferWithinRegion:                         pkg.ExtractFeeData(doc, "attr-SenderReceiverInRegion", col),
				TransferFromBangkokToOtherBankAccount:        pkg.ExtractFeeData(doc, "attr-SenderInBkkToOtherAcc", col),
				TransferFromRegionToOtherBankAccount:         pkg.ExtractFeeData(doc, "attr-SenderInRegionToOtherAcc", col),
				ReceiveTransferInBangkokFromOtherBankAccount: pkg.ExtractFeeData(doc, "attr-SenderInBkkFromOtherAcc", col),
				ReceiveTransferInRegionFromOtherBankAccount:  pkg.ExtractFeeData(doc, "attr-SenderInRegionFromOtherAcc", col),
			}

			additionInfo := pkg.AdditionalInfo{
				FeeWebiteLink: doc.Find("tr.attr-header.attr-Feeurl td.cmpr-col."+col+" a.prod-url").AttrOr("href", ""),
			}

			bathnetFees = append(bathnetFees, pkg.BathnetFee{
				Provider:       provider,
				Fees:           feeDetails,
				AdditionalInfo: additionInfo,
			})
		}
	}

	jsonData, err := json.MarshalIndent(bathnetFees, "", " ")
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		return
	}

	err = os.WriteFile("bahtnet_fees.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Failed to write JSON to file:", err)
		return
	}

	fmt.Println("Data successfully saved to bahtnet_fees.json")
}
