package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Product represents the extracted product information.
type Product struct {
	Provider     string `json:"provider"`
	ProductName  string `json:"product_name"`
	CardType     string `json:"card_type"`
	BenefitType  string `json:"benefit_type"`
	Feature      string `json:"feature"`
	MinAge       string `json:"min_age"`
	MinIncome    string `json:"min_income"`
	InterestFree string `json:"interest_free"`
	CreditLimit  string `json:"credit_limit"`
	EntranceFee  string `json:"entrance_fee"`
	AnnualFee    string `json:"annual_fee"`
	OtherFees    string `json:"other_fees"`
	CashAdvance  string `json:"cash_advance"`
	ProductURL   string `json:"product_url"`
	FeeURL       string `json:"fee_url"`
}

// NewProduct represents the desired JSON structure.
type NewProduct struct {
	Provider          string   `json:"provider"`
	ProductName       string   `json:"product_name"`
	CardType          string   `json:"card_type"`
	MainBenefit       string   `json:"main_benefit"`
	ProductFeatures   []string `json:"product_features"`
	MaximumCreditLine string   `json:"maximum_credit_line"`
	MinimumAge        string   `json:"minimum_age"`
	IncomeCondition   struct {
		Income    string `json:"income"`
		Condition string `json:"condition"`
	} `json:"income_condition"`
	InterestFreePeriod string `json:"interest_free_period"`
	Fees               struct {
		EntranceFee string `json:"entrance_fee"`
		AnnualFee   struct {
			FirstYear       string `json:"first_year"`
			SubsequentYears string `json:"subsequent_years"`
			Conditions      string `json:"conditions"`
		} `json:"annual_fee"`
		FxRiskFee      string `json:"fx_risk_fee"`
		CashAdvanceFee struct {
			Amount    string `json:"amount"`
			Condition string `json:"conditions"`
		} `json:"cash_advance_fee"`
		AdditionalInfo struct {
			ProductWebsite string `json:"product_website"`
			FeeWebsite     string `json:"fee_website"`
		} `json:"additional_info"`
	} `json:"fees"`
}

func main() {
	url := "https://app.bot.or.th/1213/MCPD/ProductApp/Credit/CompareProductList"
	payloadTemplate := `{"ProductIdList":"3629,3622,3630,4633,4632,4634,4655,4656,1604,1607,2259,5573,5161,2251,5534,5537,5177,5491,5522,5523,3627,3631,3600,3664,3624,3601,3625,3603,3633,3626,3604,3661,3620,3634,3636,3635,3606,3638,3637,3640,3639,3605,3672,3668,3628,3662,3621,3648,3607,3642,3649,3643,3651,3613,3644,3652,3609,3645,3653,3646,3647,3615,3656,3610,3650,3671,3658,3655,3616,3670,3612,4653,4657,4658,4659,4660,5568,5528,2256,5570,2245,5531,4471,4472,4467,5479,5494,5484,5497,5483,5503,1603,4760,4765,5473,3804,3800,4445,4482,4483,4479,4480,4452,4453,4476,4477,4454,4456,4457,4458,4460,4461,4462,4463,4464,4465,4444,4446,4447,4449,4450,4473,4470,5486,5498,5501,5540,5555,2242,5560,5496,1605,1600,1601,1602,4475,2246,5563,5565,5539,5567,5525,5556,2244,5562,2255,5536,2260,5574,2252,5546,5548,5538,5577,5518,5514,5521,5517,5527,3802,5529,5542,2258,5572,2249,5533,5516,4459,4474,4469,5499,5492,5481,5148,5435,3608,3641,3669,3611,3617,3660,5114,5287,5282,5285,5193,5137,5211,5173,5162,5156,5222,5126,5202,5158,5147,5240,5294,5296,5323,5306,5300,5329,749,753,752,751,750,754,5489,5500,5480,5495,4484,4481,4478,4455,4466,4451,3808,3797,3798,5502,5482,5566,2257,5571,2247,5532,5575,5576,5541,5554,2248,5544,2240,5558,5429,3996,3993,5520,3994,3995,3602,3632,3654,3665,3663,3618,3657,3623,1606,5526,5553,5552,3806,3799,3801,2253,5535,5564,3614,3666,3667,3619,3659,2250,5545,2241,5559,5549,2254,5547,2243,5561,5550,4448,4468,5519,5493,5266,5305,5304,5246,4654,5462,5213,5233,5478,5467,3997,5180,5557,5543,5551,5524,5335","Page":%d,"Limit":3}`

	// First, get the total number of pages from the initial request
	initialPage := 1
	payload := fmt.Sprintf(payloadTemplate, initialPage)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set headers (same as your previous implementation)
	setHeaders(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to retrieve data: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Error loading HTML: %v", err)
	}

	// Extract total number of pages
	totalPages := 1
	doc.Find("ul.pagination li.MoveLast a").Each(func(i int, s *goquery.Selection) {
		dataPage, exists := s.Attr("data-page")
		if exists {
			totalPages, err = strconv.Atoi(dataPage)
			if err != nil {
				log.Fatalf("Error converting data-page to int: %v", err)
			}
		}
	})

	var allNewProducts []NewProduct

	// Loop through each page
	for page := 1; page <= totalPages; page++ {
		payload := fmt.Sprintf(payloadTemplate, page)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}

		setHeaders(req)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error sending request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Failed to retrieve data: %v", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response: %v", err)
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			log.Fatalf("Error loading HTML: %v", err)
		}
		// Extract data from the table and map directly to NewProduct
		doc.Find("thead tr.attr-header").Each(func(i int, s *goquery.Selection) {
			product := Product{}
			product.Provider = cleanString(s.Find("th.col-s-0").Text())
			product.ProductName = cleanString(s.Find(".text-bold").Text())
			product.CardType = cleanString(s.Find("th.font-black").Contents().Not("span").Text())

			doc.Find("tbody tr.attr-header").Each(func(i int, s *goquery.Selection) {
				header := cleanString(s.Find("td.text-center.frst-col span").Text())
				switch header {
				case "ประเภทสิทธิประโยชน์เด่น":
					product.BenefitType = cleanString(s.Find("td.cmpr-col.col1 span.text-bold").Text())
				case "ลักษณะเด่น":
					product.Feature = cleanString(s.Find("td.cmpr-col.col1 span").Text())
				case "อายุผู้สมัครบัตรหลัก":
					product.MinAge = cleanString(s.Find("td.cmpr-col.col1 span").Text())
				case "รายได้ขั้นต่ำ และเงื่อนไขในการสมัคร":
					product.MinIncome = cleanString(s.Find("td.cmpr-col.col1 span.text-primary").Text())
				case "ระยะเวลาสูงสุดที่ปลอดดอกเบี้ย":
					product.InterestFree = cleanString(s.Find("td.cmpr-col.col1 span").Text())
				case "วงเงินสูงสุด":
					product.CreditLimit = cleanString(s.Find("td.cmpr-col.col1 span").Text())
				case "ค่าธรรมเนียมแรกเข้าบัตรหลัก":
					product.EntranceFee = cleanString(s.Find("td.cmpr-col.col1 span").Text())
				case "ค่าธรรมเนียมรายปีบัตรหลัก":
					product.AnnualFee = cleanString(s.Find("td.cmpr-col.col1 span").Text())
				case "ค่าความเสี่ยงจากการแปลงสกุลเงิน":
					product.OtherFees = cleanString(s.Find("td.cmpr-col.col1 span").Text())
				case "ค่าธรรมเนียมเบิกถอนเงินสด":
					product.CashAdvance = cleanString(s.Find("td.cmpr-col.col1 span").Text())
				}
			})

			product.ProductURL = cleanString(s.Find("tbody tr.attr-header.attr-url td.cmpr-col.col1 a").AttrOr("href", ""))
			product.FeeURL = cleanString(s.Find("tbody tr.attr-header.attr-feeurl td.cmpr-col.col1 a").AttrOr("href", ""))

			newProducts := mapToNewProduct(product)
			allNewProducts = append(allNewProducts, newProducts...)
		})

		// Stop for 5 seconds before making the next request
		time.Sleep(5 * time.Second)
	}

	// Convert the combined new products to JSON and save to a file
	newJSON, err := json.MarshalIndent(allNewProducts, "", "  ")
	if err != nil {
		log.Fatalf("Error converting to JSON: %v", err)
	}

	file, err := os.Create("new_products.json")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(newJSON)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	log.Println("Product data saved to new_products.json")
}

func setHeaders(req *http.Request) {
	req.Header.Set("Accept", "text/plain, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Cookie", `verify=test; verify=test; verify=test; mycookie=\u0021IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; _uid6672=16B5DEBD.6; visit_time=2171; _ga_NLQFGWVNXN=GS1.1.1723798115.10.1.1723798115.60.0.0`)
	req.Header.Set("Origin", "https://app.bot.or.th")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://app.bot.or.th/1213/MCPD/ProductApp/Credit/CompareProduct")
	req.Header.Set("Sec-Ch-Ua", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	req.Header.Set("Verificationtoken", "jZKLX6IRz_o2TRzkK3Muo-MQJJyTKCCqYH_nqRqFPSlMqzPukrPg--T1qwu7lBF6ikqALmGzObfK1bEHefY_iHAfKTD0-PqTh6CzSbnPS4M1,qZsSwtss8Ueiv1fJBVwqoyWusJv_BnjoVmkYsXcNn3E6JAWC5FBaJ1jlFqgXFD_9nIEjon7NPJ-AfhnaK11irhNoXh9jfWaN2j8YIaZp-3g1")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
}

func mapToNewProduct(product Product) []NewProduct {
	// Initialize a list to hold the new product information
	var newProducts []NewProduct

	// Directly use the product name as extracted from the scraped data
	productName := cleanString(product.ProductName)
	cardType := cleanString(product.CardType)

	// Create a single NewProduct instance
	newProduct := NewProduct{
		Provider:    cleanString(product.Provider),
		ProductName: productName,
		CardType:    cardType,
		MainBenefit: cleanString(product.BenefitType),
		ProductFeatures: []string{
			cleanString(product.Feature),
		},
		MaximumCreditLine: cleanString(product.CreditLimit),
		MinimumAge:        cleanString(product.MinAge),
	}
	newProduct.IncomeCondition.Income = ""
	newProduct.IncomeCondition.Condition = cleanString(product.MinIncome)
	newProduct.InterestFreePeriod = cleanString(product.InterestFree)

	// Populate the Fees struct.
	newProduct.Fees.EntranceFee = cleanString(product.EntranceFee)
	annualFeeParts := strings.Split(product.AnnualFee, "ปีถัดไป:")
	if len(annualFeeParts) == 2 {
		newProduct.Fees.AnnualFee.FirstYear = strings.TrimSpace(strings.Split(annualFeeParts[0], "ปีแรก:")[1])
		newProduct.Fees.AnnualFee.SubsequentYears = strings.TrimSpace(annualFeeParts[1])
	}
	newProduct.Fees.FxRiskFee = cleanString(product.OtherFees)
	newProduct.Fees.CashAdvanceFee.Amount = cleanString(product.CashAdvance)
	newProduct.Fees.AdditionalInfo.ProductWebsite = cleanString(product.ProductURL)
	newProduct.Fees.AdditionalInfo.FeeWebsite = cleanString(product.FeeURL)

	newProducts = append(newProducts, newProduct)

	return newProducts
}

// cleanString is a helper function to trim and clean strings.
func cleanString(str string) string {
	return strings.TrimSpace(str)
}

// cleanSpaces removes extra spaces and line breaks within a string.
func cleanSpaces(str string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(str, " ")
}
