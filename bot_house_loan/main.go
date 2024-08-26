package main

import (
	"fmt"
	"house_loan/pkg/models"
	"house_loan/pkg/parser"
	"house_loan/pkg/request"
	"house_loan/pkg/writer"
	"log"
)

func main() {
	url := "https://app.bot.or.th/1213/MCPD/ProductApp/HomeLoan/CompareProductList"
	initialPayload := `{"ProductIdList":"21430,21389,21396,21425,21452,21423,21404,17220,17216,700,703,710,21443,21440,21419,21403,21411,21392,21414,21451,17218,17217,17226,1739,21431,17222,21444,21390,21418,21397,21441,20822,708,704,21433,21004,21006,21432,21391,17224,21413,17223,21408,21005,21007,21439,21427,15220,17219,17221,20229,20854,21417,21442,20833,21412,21401,21426,21436,21400,20789,20791,20776,20261,20306,21012,15219,17225,20852,20281,20328,21020,15222,21437,21098,21060,20241,20313,21013,21128,21113,21008,21010,17227,20283,20262,20330,20307,21021,20770,20771,20768,20783,21428,21252,21263,20255,20300,20263,21009,21011,15221,20308,20321,20269,21445,20252,20267,20331,20246,20272,20260,21429,21036,21242,20924,21225,21210,21016,21245,21192,21095,21415,20843,20250,20265,20333,20259,21233,709,705,20790,20794,20796,20793,20782,20798,20795,21253,21261,21254,21262,21193,21059,20316,20257,21014,21030,21118,21037,21133,21395,21017,16243,20298,20277,20305,20304,20244,20324,20856,20287,20311,21135,20319,20268,20777,20785,20266,20279,20293,21022,16213,16234,21249,21260,20276,21239,21115,21015,20942,20336,20256,20291,20299,14980,21218,21214,20253,20294,21185,21049,21099,21244,21023,16224,20775,21221,21215,21212,21227,20289,21407,21448,20302,20309,21224,21142,21053,20318,20264,20258,21438,16204,16254,16249,20314,21196,20290,20312,20278,21124,21197,21201,21189,21018,16194,16244,20315,16223,21258,21259,21250,21251,20280,20326,21194,21038,21101,21134,21019,20922,706,20270,21047,21226,21114,20774,14981,21035,20325,21084,21172,21039,16233,20271,20245,21173,21032,21131,21119,20274,21116,21209,20929,16214,21182,21087,21103,20934,20931,16183,21085,21092,21238,21424,21398,21420,21170,21130,21220,20923,20937,21195,21127,21041,21246,21257,21086,16176,21175,20920,16171,16193,21223,20928,20941,21126,16203,21093,21255,21256,21247,21248,20921,16184,20787,20926,1738,20219,20221,20223,20225,20227,701,20940,21409,20296,17711,14982,20301,20284,20780,20781,20784,20792,20788,20786,20933,20249,20329,20286,20251,21406,20936,20919,20779,20778,20797,20772,20773,20769,20254,20220,20222,20224,20226,20228,20335,20292,707,20273,20927,20282,20327,20248,21449,20930,20297,20334,14965,20935,21421,21199,20925,702,20932,12618,21089,20938,20939,20322,20317,20320,20303,21045,21456,21405,20323,20275,20285,20243,20242,20247,20295,21184,21112,21216,21040,21100,21455,1740,20310,20332,20288,21051,21110,21206,21139,21048,21111,21241,21061,21141,21208,21217,21056,21240,21121,21198,21132,21138,21213,21043,21219,21058,21243,21123,21091,21140,21211,21222,21454,21187,21096,21453,21129,17710,17709,17712,16469,16461,15456,15455","Page":1,"Limit":3}`

	firstPageBody, err := request.FetchHTML(url, initialPayload)
	if err != nil {
		log.Fatalf("Error fetching first page: %v", err)
	}

	totalPages := parser.DetermineTotalPage(firstPageBody)
	if totalPages == 0 {
		log.Fatalf("Could not determine the total number of pages")
	}

	var allHouseLoans []models.HouseLoan

	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages)
		payload := fmt.Sprintf(`{"ProductIdList":"21430,21389,21396,21425,21452,21423,21404,17220,17216,700,703,710,21443,21440,21419,21403,21411,21392,21414,21451,17218,17217,17226,1739,21431,17222,21444,21390,21418,21397,21441,20822,708,704,21433,21004,21006,21432,21391,17224,21413,17223,21408,21005,21007,21439,21427,15220,17219,17221,20229,20854,21417,21442,20833,21412,21401,21426,21436,21400,20789,20791,20776,20261,20306,21012,15219,17225,20852,20281,20328,21020,15222,21437,21098,21060,20241,20313,21013,21128,21113,21008,21010,17227,20283,20262,20330,20307,21021,20770,20771,20768,20783,21428,21252,21263,20255,20300,20263,21009,21011,15221,20308,20321,20269,21445,20252,20267,20331,20246,20272,20260,21429,21036,21242,20924,21225,21210,21016,21245,21192,21095,21415,20843,20250,20265,20333,20259,21233,709,705,20790,20794,20796,20793,20782,20798,20795,21253,21261,21254,21262,21193,21059,20316,20257,21014,21030,21118,21037,21133,21395,21017,16243,20298,20277,20305,20304,20244,20324,20856,20287,20311,21135,20319,20268,20777,20785,20266,20279,20293,21022,16213,16234,21249,21260,20276,21239,21115,21015,20942,20336,20256,20291,20299,14980,21218,21214,20253,20294,21185,21049,21099,21244,21023,16224,20775,21221,21215,21212,21227,20289,21407,21448,20302,20309,21224,21142,21053,20318,20264,20258,21438,16204,16254,16249,20314,21196,20290,20312,20278,21124,21197,21201,21189,21018,16194,16244,20315,16223,21258,21259,21250,21251,20280,20326,21194,21038,21101,21134,21019,20922,706,20270,21047,21226,21114,20774,14981,21035,20325,21084,21172,21039,16233,20271,20245,21173,21032,21131,21119,20274,21116,21209,20929,16214,21182,21087,21103,20934,20931,16183,21085,21092,21238,21424,21398,21420,21170,21130,21220,20923,20937,21195,21127,21041,21246,21257,21086,16176,21175,20920,16171,16193,21223,20928,20941,21126,16203,21093,21255,21256,21247,21248,20921,16184,20787,20926,1738,20219,20221,20223,20225,20227,701,20940,21409,20296,17711,14982,20301,20284,20780,20781,20784,20792,20788,20786,20933,20249,20329,20286,20251,21406,20936,20919,20779,20778,20797,20772,20773,20769,20254,20220,20222,20224,20226,20228,20335,20292,707,20273,20927,20282,20327,20248,21449,20930,20297,20334,14965,20935,21421,21199,20925,702,20932,12618,21089,20938,20939,20322,20317,20320,20303,21045,21456,21405,20323,20275,20285,20243,20242,20247,20295,21184,21112,21216,21040,21100,21455,1740,20310,20332,20288,21051,21110,21206,21139,21048,21111,21241,21061,21141,21208,21217,21056,21240,21121,21198,21132,21138,21213,21043,21219,21058,21243,21123,21091,21140,21211,21222,21454,21187,21096,21453,21129,17710,17709,17712,16469,16461,15456,15455","Page":%d,"Limit":3}`, page)
		// payload := fmt.Sprintf(`{"ProductIdList":"21430,21389,21396,21425,21452,21423,21404,17220,17216,700,703,710,21443,21440,21419,"Page":%d,"Limit":3}`, page)
		pageBody, err := request.FetchHTML(url, payload)
		if err != nil {
			log.Printf("error fetching data for page %d: %v", page, err)
			continue
		}

		loans, err := parser.ParseHouseLoanDetails(pageBody)
		if err != nil {
			log.Printf("Error parsing loan details for page %d %v", page, err)
			continue
		}
		allHouseLoans = append(allHouseLoans, loans...)
	}

	err = writer.WriteJSON(allHouseLoans, "house_loan.json")
	if err != nil {
		log.Fatalf("Error writing JSON file: %v", err)
	}

	fmt.Println("House loan data has been saved to hous_loan.json")
}
