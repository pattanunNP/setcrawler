package main

import (
	"bytes"
	"deposit_fee/pkg"
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
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/DepositFee/CompareProductList"
	payloadTemplate := `{"ProductIdList":"64178,64190,64179,64191,64180,64192,64181,64193,64182,64194,64183,64184,64176,64188,64177,64189,64157,64158,64159,64160,64168,64169,64161,64162,64170,64171,64163,64164,64172,64165,64166,64167,64147,64148,64151,64152,64153,64154,64155,64156,64324,64230,64244,64221,64368,64349,64306,64331,64267,64376,64308,64316,64318,64353,64356,64354,64242,64286,64382,64384,64361,64256,64377,62848,62851,62850,62853,62852,62855,62854,62857,62856,62858,62842,63610,63554,63708,63658,63545,63613,63695,63617,63570,63606,63680,63656,63518,63585,63688,63706,63644,63530,63659,63523,63687,63635,63704,63541,63643,63707,63528,63698,63962,63916,63964,63966,63951,63917,63926,63918,63927,63910,63976,63953,63937,63954,63919,63920,63975,63930,63978,63967,63984,63905,63935,63936,63907,63981,63922,63977,37616,37615,37623,63898,63887,63899,63889,63900,63890,63891,63902,63893,63903,63895,63904,63897,63888,63896,63892,63883,63886,61129,61101,61148,61102,61118,61116,61120,61104,61139,61109,61126,61096,61144,61098,61128,61103,61119,61137,61124,61141,61150,61095,61138,61107,61142,61097,61123,61147,61112,61127,61151,61105,61131,61135,61117,61136,61110,61125,61145,61149,61133,61115,61146,61130,61132,61106,61108,61113,63872,63858,63873,63859,63874,63860,63875,63861,63876,63862,63877,63863,63878,63864,63879,63865,63880,63866,63881,63867,63870,63857,64020,64001,64021,64002,63993,64022,64013,63994,64003,64014,63995,64023,64015,63996,64004,64016,63997,64024,64017,63998,64018,64019,64000,63985,64005,22391,22386,22392,22388,22393,22390,22394,22398,22395,22399,22396,22400,22397,22383,22384,23262,23259,23263,23260,23264,62209,62247,62198,62204,62266,62190,62186,62264,62187,62265,62210,62236,62188,62226,62227,62194,62212,62243,62193,62195,62228,62184,62201,62213,62205,62237,62234,63724,63722,63725,63723,63726,62860,62859,62862,62861,62863,63513,22222,22216,22217,22218,22219,22220,22221,62511,62512,62513,62514,62515,62516,62517,62501,62500,62502,62503,62504,62505,62506,62507,62508,62509,62510,63484,63485,63486,63487,63488,63489,63490,63491,64084,64104,64121,64075,64062,64082,64096,64067,64038,64126,64118,64033,64053,64089,64092,64060,64097,64068,64025,64045,64091,64111,64113,64133,64135,64063,64039,64070,64046,64115,64041,64125,64134,64035,64129,64105,64123,64049,64127,64094,64137,64142,64086,64108,64057,64055,64031,64100,64098,64028,64026,64119,64071,64076,64116,64114,64090,64034,64042,64085,64061,64122,64048,64058,64141,64093,64130,64136,64037,64112,64064,64107,64083,64056,64078,64054,64027,64079,64032,64050,64047,64124,64095,64066,64087,64139,64073,64120,64128,64101,64072,64138,64109,64421,64487,64484,64415,64497,64449,64446,64392,64476,64411,64512,64447,64496,64490,64493,64475,64406,64503,64395,64430,64455,64429,64459,64492,64474,64391,64412,64510,64456,64464,64494,64504,64404,64393,64501,64394,64502,64451,64427,64444,64457,64436,64445,64505,64499,64478,64498,64433,64489,64479,64509,64511,64413,64441,64485,64488,64434,64467,64463,64402,64428,63445,63446,63447,63448,63449,63450,63442,63443,63437,63452,63454,63455,63458,63459,52351,52354,52353,52356,52355,52358,52357,59183,59174,59184,59175,59185,59166,59170,59180,59171,59181,59172,59182,59173,63327,63337,63346,63356,63366,63374,63378,63271,63385,63278,63392,63288,63294,63300,63217,63221,63314,63322,63417,63340,63426,63359,63435,63405,22380,22385,22381,22387,22382,22389,63438,63439,63440,63441,63444,64187,64173,64174,64185,64186,64175,64143,64144,64149,64150,64145,64146,64365,64246,64277,64371,64337,64328,62843,62845,62847,63642,63646,63615,63709,63652,63650,63602,63974,63909,63982,63929,63925,63908,63894,63882,63885,63884,61140,61121,61100,61099,61114,61134,61122,63853,63856,63868,63869,63855,63854,64007,63988,64008,63989,64009,63986,64006,63987,60673,62251,62225,63483,63492,59176,59167,59177,59168,59178,59169,59179,63287,63297,63227,63230,64290,63983,63928,64012,63990,64010,63991,64011,63992,23258,23261,62254,64051,64117,64080,64131,64106,64077,64040,64069,64044,64036,64065,64132,64074,64030,64110,64103,64140,64081,64088,64052,64059,64491,64448,52352,63231,63317,37605,37602,23265,23266,23272,23273,23274,23267,23275,23268,23269,23276,23270,23271,62208,62231,62202,62258,62259,62239,62192,62241,62199,62260,62242,62222,62257,62252,62235","Page":%d,"Limit":3}`
	// payloadTemplate := `{"ProductIdList":"64178,64190,64179,64191,64180,64192,64181,64193,64182,64194,64183,64184,64176,64188,64177,64189,64157,64158,64159,64160,64168,64169,64161,64162,64170,64171","Page":%d,"Limit":3}`

	initialPayload := fmt.Sprintf(payloadTemplate, 1)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(initialPayload)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Accept", "text/plain, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Cookie", `verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0"; _uid6672=16B5DEBD.33; _ctout6672=1; visit_time=7; _ga_NLQFGWVNXN=GS1.1.1724745088.41.1.1724745119.29.0.0`)
	req.Header.Set("Origin", "https://app.bot.or.th")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://app.bot.or.th/1213/MCPD/FeeApp/DepositFee/CompareProduct")
	req.Header.Set("Sec-CH-UA", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	req.Header.Set("VerificationToken", "74S1-mNrhwaDuQUWwscFjUYtJxQ5s99O3DVEe2QhG6FrtkAOjfZLWJu7TeeFW_RbnEb8EtcH9HUPGUI_tkDzmFfI3_BOj1Wq1pm_WfymWlw1,KBGT8pCk39B8HzWaGPhLk-7YNCBN094abFf_jWruSaJs8C8GXoTiDjhKhyncG0bhofWy3But0_8W8DO13KVGo5iCkTM2JthdeaHmDFZR2_A1")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	totalPages := pkg.DetermineTotalPage(doc)
	if totalPages == 0 {
		log.Fatalf("Could not determine the total number of pages")
	}

	var depositFess []pkg.DepositFee

	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages)
		payload := fmt.Sprintf(payloadTemplate, page)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			log.Fatalf("Error creating request for page %d: %v", page, err)
		}

		// Set headers
		req.Header.Set("Accept", "text/plain, */*; q=0.01")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		req.Header.Set("Cookie", `verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0"; _uid6672=16B5DEBD.33; _ctout6672=1; visit_time=7; _ga_NLQFGWVNXN=GS1.1.1724745088.41.1.1724745119.29.0.0`)
		req.Header.Set("Origin", "https://app.bot.or.th")
		req.Header.Set("Priority", "u=1, i")
		req.Header.Set("Referer", "https://app.bot.or.th/1213/MCPD/FeeApp/DepositFee/CompareProduct")
		req.Header.Set("Sec-CH-UA", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
		req.Header.Set("Sec-CH-UA-Mobile", "?0")
		req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
		req.Header.Set("VerificationToken", "74S1-mNrhwaDuQUWwscFjUYtJxQ5s99O3DVEe2QhG6FrtkAOjfZLWJu7TeeFW_RbnEb8EtcH9HUPGUI_tkDzmFfI3_BOj1Wq1pm_WfymWlw1,KBGT8pCk39B8HzWaGPhLk-7YNCBN094abFf_jWruSaJs8C8GXoTiDjhKhyncG0bhofWy3But0_8W8DO13KVGo5iCkTM2JthdeaHmDFZR2_A1")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error making request for page %d: %v", page, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Request for page %d failed with status: %d", page, resp.StatusCode)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body for page %d: %v", page, err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			log.Printf("Error parsing HTML for page %d: %v", page, err)
			continue
		}

		for i := 1; i <= 3; i++ {
			col := "col" + strconv.Itoa(i)

			provider := pkg.CleanString(doc.Find("th.col-s.col-s-" + strconv.Itoa(i) + " span").Last().Text())
			product := pkg.CleanString(doc.Find("th.font-black.text-center.prod-" + col + " span").Text())

			fees := pkg.ProductFees{
				AccountMaintenanceFee:             pkg.ExtractFee(doc, "attr-AccountMaintenanceFee", col),
				StatementRequireFee:               pkg.ExtractFeePtr(doc, "attr-StatementRequireFee", col),
				StatementRequireSixMonth:          pkg.ExtractFee(doc, "attr-StatementRequireSixMonth", col),
				StatementRequireSixMonthToTwoYear: pkg.ExtractFee(doc, "attr-StatementRequiresixMonthToTwoYear", col),
				StatementRequireTwoYear:           pkg.ExtractFee(doc, "attr-StatementRequireTwoYear", col),
				ShortMessageService:               pkg.ExtractFee(doc, "attr-ShortMessageService", col),
				ShortMessageServiceFeeMonthly:     pkg.ExtractFee(doc, "attr-ShortMessageServiceFeeMonthy", col),
				ShortMessageServiceAnnualFee:      pkg.ExtractFee(doc, "attr-ShortMessageServiceAnnaulFee", col),
				LostPassBookFee:                   pkg.ExtractFee(doc, "attr-LostPassBookFee", col),
				AccountCloseFee:                   pkg.ExtractFee(doc, "attr-AccountCloseFee", col),
			}

			general := pkg.GeneralFees{
				CoinCollectFee:                         pkg.ExtractFeePtr(doc, "attr-CoinCollectFee", col),
				BranchFee:                              pkg.ExtractFeePtr(doc, "attr-BRFee", col),
				KioskOtherFee:                          pkg.ExtractFeePtr(doc, "attr-KioskOtherFee", col),
				KioskFee:                               pkg.ExtractFeePtr(doc, "attr-KioskFee", col),
				AgentFee:                               pkg.ExtractFeePtr(doc, "attr-AgentFee", col),
				ShopAgentFee:                           pkg.ExtractFeePtr(doc, "attr-ShopAgentFee", col),
				PostAgentFee:                           pkg.ExtractFeePtr(doc, "attr-PostAgentFee", col),
				TopupAgentFee:                          pkg.ExtractFeePtr(doc, "attr-TopupAgentFee", col),
				OtherAgentFee:                          pkg.ExtractFeePtr(doc, "attr-OtherAgentFee", col),
				TransferBetweenSavingCurrentAccountFee: pkg.ExtractFeePtr(doc, "attr-TransferBetweenSavingCurrentAccoutnFee", col),
				TransferBetweenBankingFee:              pkg.ExtractFeePtr(doc, "attr-TransferBetweenBankingFee", col),
			}

			otherFees := pkg.OtherFees{
				OtherFees: pkg.ExtractFeePtr(doc, "attr-OtherFee", col),
			}

			additionalInfo := pkg.AdditionalInfo{
				FeeURL: pkg.ExtractURL(doc, "attr_feeurl", col),
			}

			depositFess = append(depositFess, pkg.DepositFee{
				Provider:       provider,
				Product:        product,
				Fees:           fees,
				GeneralFees:    general,
				OtherFees:      otherFees,
				AdditionalInfo: additionalInfo,
			})
		}

		file, err := os.Create("deposit_fees.json")
		if err != nil {
			log.Fatalf("Error creating JSON file: %v", err)
		}
		defer file.Close()

		jsonData, err := json.MarshalIndent(depositFess, "", " ")
		if err != nil {
			log.Fatalf("Error marshaling JSON data: %v", err)
		}

		_, err = file.Write(jsonData)
		if err != nil {
			log.Fatalf("Error writing to JSON file: %v", err)
		}

		fmt.Println("Data successfully saved to depit_fees.json")
	}
}
