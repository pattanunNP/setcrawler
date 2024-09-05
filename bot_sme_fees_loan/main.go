package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sme_fees/models"
	"sme_fees/utils"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// Define the URL and payload
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/SMEFee/CompareProductList"
	payloadTemplate := (`{"ProductIdList":"1775,1766,1764,1822,3114,2934,2956,2923,2922,2949,2952,2944,2921,4063,4103,4138,4122,3122,3747,3867,3868,3872,3873,3883,3886,3885,3105,3107,3106,3108,1138,1105,1136,1159,1115,1089,1074,1080,1117,1161,1168,1125,1103,1077,2033,2926,2954,3116,3119,3110,3115,3120,1655,1660,2911,2927,2912,2930,2913,2931,2914,2932,2915,2916,2940,2917,2941,2933,2953,2935,2937,2957,2918,2919,2942,2947,2924,2920,2948,2929,2951,2928,2950,3381,931,662,2925,2939,2958,4117,4079,4123,4140,4066,4146,4135,4097,4083,3870,3871,1654,1796,3376,2532,2531,3124,3378,1801,2723,602,599,4142,4080,1656,1657,1658,1659,1778,2530,4110,1810,1790,1816,1767,1821,1798,1800,1797,1807,1781,1783,1784,1805,1820,1772,1770,1792,1811,1817,1776,1794,1815,1814,1780,1765,1803,1791,1812,1777,1769,1768,1763,1793,1761,1819,1813,1762,1774,1808,1786,1773,1788,1802,1804,1785,1771,1787,1789,1799,1795,1779,1806,1818,4126,3133,3134,3128,3377,1809,2727,2728,601,4089,3123,3379,2730,598,603,600,4145,4070,4085,4124,4129,4072,4098,4143,4071,4087,4114,2533,4086,4100,4137,4132,3121,3132,3125,3135,418,421,420,419,422,1096,1172,1142,1086,1099,1075,1163,1094,1166,1146,1100,1070,1076,1141,1155,1123,1124,1151,1109,1112,1114,1118,1116,1170,1165,1162,1098,1082,1084,1144,1145,1133,1093,1102,1157,1072,1073,1091,1069,1095,1104,1106,1143,1152,1088,1085,1097,1137,1130,1149,1150,1140,1160,1173,1126,1127,1128,1120,1139,1079,1167,1156,1132,1083,1068,1119,1134,1107,1154,1121,1148,1081,1101,1078,1111,1169,1164,1153,1129,1071,1135,1108,1131,1171,1110,1087,1092,1090,3380,1782,3887,3757,3751,3753,3754,3745,3755,3746,3756,3749,3758,3759,3750,165,899,905,903,900,901,904,902,3743,1602,1597,1592,3137,3129,3138,3130,3131,3126,3127,3136,2895,2896,2892,2897,2894,2899,2893,2898,1598,1969,1970","Page":%d,"Limit":3}`)
	// payloadTemplate := (`{"ProductIdList":"1775,1766,1764,1822,3114,2934,2956,2923,2922,2949,2952,2944,2921,4063,4103,4138,4122,3122,3747,3867,3868,3872,3873,3883,3886,3885,3105","Page":%d,"Limit":3}`)

	// Create a new request
	initialPayload := fmt.Sprintf(payloadTemplate, 1)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(initialPayload)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	utils.AddHeader(req)

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Error loading HTML document: %v", err)
	}

	totalPages := utils.DetermineTotalPage(doc)
	if totalPages == 0 {
		log.Fatal("Could not determine the total number of pages")
	}

	var smeProducts []models.SMEProduct
	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages)
		payload := fmt.Sprintf(payloadTemplate, page)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			log.Printf("Error creating request for page %d: %v", page, err)
			continue
		}
		utils.AddHeader(req)

		// Perform the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			log.Fatalf("Error loading HTML document: %v", err)
		}

		for i := 1; i <= 3; i++ {
			col := fmt.Sprintf("col%d", i)
			serviceProvider := utils.CleanText(doc.Find(fmt.Sprintf("th.col-s.col-s-%d span", i)).Last().Text())
			product := utils.CleanText(doc.Find(fmt.Sprintf("th.font-black.text-center.prod-%s span", col)).Text())

			// Extract fee details for each product
			debtCollectionFeeArray := utils.SplitByHyphen(utils.CleanText(doc.Find(fmt.Sprintf(".attr-DebtCollectionFee .cmpr-col.%s span", col)).Text()))
			creditCheckFee := utils.CleanText(doc.Find(fmt.Sprintf(".attr-CreditBureauFee .cmpr-col.%s span", col)).Text())
			statementReIssuingFee := utils.CleanText(doc.Find(fmt.Sprintf(".attr-StatementReIssuingFee .cmpr-col.%s span", col)).Text())

			prepaymentFeeText := utils.CleanText(doc.Find(fmt.Sprintf(".attr-PrepaymentFee .cmpr-col.%s span", col)).Text())
			minPrepayment, maxPrepayment := utils.ExtractPercentages(prepaymentFeeText)

			extensionFeeText := utils.CleanText(doc.Find(fmt.Sprintf(".attr-ExtensionFee .cmpr-col.%s span", col)).Text())
			minExtension, maxExtension := utils.ExtractPercentages(extensionFeeText)

			appraisalFeeInternalText := utils.CleanText(doc.Find(fmt.Sprintf(".attr-SurveyAndAppraisalFeeByInternal .cmpr-col.%s span", col)).Text())
			minAppraisalInternal, maxAppraisalInternal := utils.ExtractAmounts(appraisalFeeInternalText)

			var debtCollectionFeeAmount int
			for _, part := range debtCollectionFeeArray {
				debtCollectionFeeAmount = utils.ExtractFirstInt(part)
				if debtCollectionFeeAmount != 0 {
					break
				}
			}

			loanFees := models.LoanFees{
				FrontEndFee:     utils.CleanText(doc.Find(fmt.Sprintf(".attr-FrontEndFee .cmpr-col.%s span", col)).Text()),
				ManagementFee:   utils.CleanText(doc.Find(fmt.Sprintf(".attr-ManagementFee .cmpr-col.%s span", col)).Text()),
				CommitmentFee:   utils.CleanText(doc.Find(fmt.Sprintf(".attr-CommitmentFee .cmpr-col.%s span", col)).Text()),
				CancellationFee: utils.CleanText(doc.Find(fmt.Sprintf(".attr-CancellationFee .cmpr-col.%s span", col)).Text()),
				PrepaymentFee: models.FeeWithPercentage{
					Description:   prepaymentFeeText,
					MinPercentage: minPrepayment,
					MaxPercentage: maxPrepayment,
				},
				ExtensionFee: models.FeeWithPercentage{
					Description:   extensionFeeText,
					MinPercentage: minExtension,
					MaxPercentage: maxExtension,
				},
				AppraisalFeeInternal: models.FeeWithAmount{
					Description: utils.CleanText(doc.Find(fmt.Sprintf(".attr-SurveyAndAppraisalFeeByInternal .cmpr-col.%s span", col)).Text()),
					MinAmount:   minAppraisalInternal,
					MaxAmount:   maxAppraisalInternal,
				},
				AppraisalFeeExternal: models.FeeWithAmount{
					Description: utils.CleanText(doc.Find(fmt.Sprintf(".attr-SurveyAndAppraisalFeeByExternal .cmpr-col.%s span", col)).Text()),
					MinAmount:   utils.ExtractFirstInt(doc.Find(fmt.Sprintf(".attr-SurveyAndAppraisalFeeByExternal .cmpr-col.%s span", col)).Text()),
					MaxAmount:   utils.ExtractFirstInt(doc.Find(fmt.Sprintf(".attr-SurveyAndAppraisalFeeByExternal .cmpr-col.%s span", col)).Text()),
				},
				DebtCollectionFee:       debtCollectionFeeArray,
				CreditCheckFee:          creditCheckFee,
				StatementReIssuingFee:   statementReIssuingFee,
				DebtCollectionFeeAmount: debtCollectionFeeAmount,
				CreditCheckFeeAmount:    utils.ExtractFirstInt(creditCheckFee),
				StatementFeeAmount:      utils.ExtractFirstInt(statementReIssuingFee),
			}

			otherFeesText := utils.CleanText(doc.Find(fmt.Sprintf(".attr-other .cmpr-col.%s span", col)).Text())
			var otherFeesValue *string
			if otherFeesText != "" {
				otherFeesValue = &otherFeesText
			} else {
				otherFeesValue = nil
			}

			otherFees := models.OtherFees{
				OtherFees: otherFeesValue, // Assign null if empty
			}

			feeWebsiteLink := doc.Find(fmt.Sprintf(".attr-Feeurl .cmpr-col.%s a", col)).AttrOr("href", "")
			var feeWebsiteLinkValue *string
			if feeWebsiteLink != "" {
				feeWebsiteLinkValue = &feeWebsiteLink
			} else {
				feeWebsiteLinkValue = nil
			}

			additionalInfo := models.AdditionalInfo{
				FeeWebsiteLink: feeWebsiteLinkValue, // Assign null if empty
			}

			smeProduct := models.SMEProduct{
				ServiceProvider: serviceProvider,
				Product:         product,
				LoanFees:        loanFees,
				OtherFees:       otherFees,
				AdditionalInfo:  additionalInfo,
			}

			smeProducts = append(smeProducts, smeProduct)
		}

		time.Sleep(2 * time.Second)
	}

	// Create the final JSON object
	responseData := struct {
		SMEProducts []models.SMEProduct `json:"smeLoanFeeDetails"`
	}{
		SMEProducts: smeProducts,
	}

	jsonData, err := json.MarshalIndent(responseData, "", " ")
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}

	err = os.WriteFile("sme_fees.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Error wrting JSON file: %v", err)
	}

	fmt.Println("JSON data saved to sme_fees.json")
}
