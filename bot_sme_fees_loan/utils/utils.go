package utils

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func AddHeader(req *http.Request) {
	// Set required headers
	req.Header.Set("accept", "text/plain, */*; q=0.01")
	req.Header.Set("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Header.Set("cookie", `verify=test; verify=test; verify=test; mycookie=\u0021IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0"; _uid6672=16B5DEBD.57; _ctout6672=1; visit_time=28; _ga_NLQFGWVNXN=GS1.1.1725016859.64.1.1725016906.13.0.0`)
	req.Header.Set("origin", "https://app.bot.or.th")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://app.bot.or.th/1213/MCPD/FeeApp/SMEFee/CompareProduct")
	req.Header.Set("sec-ch-ua", `"Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome";v="128"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	req.Header.Set("verificationtoken", "aloS91rw5leV0ZlTTHbe5EibwBFT2NZCWCpUhF1f0dpc4OnaesrsJRhMk7mzeZazmpNmGWOVOglAxByYweMdOka6G4saHdGvmNsCgckht441,7oWd47HUtlFoZIB4DZ7dcKHP8XDrwvARi3SS1PXPdFZ-EOD5UF8XkqRdmq-cT5k99J0wqYw8M-2WKJdOqPf1V4y26tI9QjkCRdf7Yy5F71A1")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
}

// CleanText removes unwanted characters and trims whitespace
func CleanText(text string) string {
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, ",", "")
	text = strings.Join(strings.Fields(text), " ")
	return text
}

func CleanAndSplit(text string) []string {
	cleaned := CleanText(text)
	if cleaned == "" {
		return nil
	}
	return strings.Split(cleaned, "-")
}

func ExtractFirstInt(text string) int {
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(text)
	if match == "" {
		return 0
	}
	number, err := strconv.Atoi(match)
	if err != nil {
		return 0
	}
	return number
}

func DetermineTotalPage(doc *goquery.Document) int {
	totalPages := 1

	doc.Find("ul.pagination li a").Each(func(i int, s *goquery.Selection) {
		pagenum, exists := s.Attr("data-page")
		if exists {
			page, err := strconv.Atoi(pagenum)
			if err == nil && page > totalPages {
				totalPages = page
			}
		}
	})
	return totalPages
}

func ExtractPercentages(text string) (minPercentage, maxPercentage float64) {
	// Regex to capture percentages, both individual and ranges with decimals
	re := regexp.MustCompile(`(\d+(\.\d+)?) ?%? ?-? ?(\d+(\.\d+)?)? ?%`)

	// Find all matches in the text
	matches := re.FindStringSubmatch(text)

	// Check if there are matches
	if len(matches) > 0 {
		// Extract the first percentage as minPercentage
		minPercentage, err1 := strconv.ParseFloat(matches[1], 64)
		if err1 != nil {
			minPercentage = 0
		}

		// If there's a second percentage, extract it as maxPercentage
		if len(matches) > 3 && matches[3] != "" {
			maxPercentage, err2 := strconv.ParseFloat(matches[3], 64)
			if err2 != nil {
				maxPercentage = minPercentage
			}
			return minPercentage, maxPercentage
		}

		// If no second percentage, return the first one as both min and max
		return minPercentage, minPercentage
	}

	// Default return if no percentages are found
	return 0, 0
}

func ExtractAmounts(text string) (minAmount, maxAmount int) {
	// Adjust regex to focus on numeric ranges, ignoring extra words
	re := regexp.MustCompile(`(\d+(?:,\d{3})*)`)
	matches := re.FindAllString(text, -1)

	// Handle cases where we have two numbers (min and max) or just one
	if len(matches) >= 2 {
		minAmount, err := strconv.Atoi(strings.ReplaceAll(matches[0], ",", ""))
		if err != nil {
			return 0, 0
		}
		maxAmount, err := strconv.Atoi(strings.ReplaceAll(matches[1], ",", ""))
		if err != nil {
			return minAmount, minAmount // If second value fails, fallback to first
		}
		return minAmount, maxAmount
	} else if len(matches) == 1 {
		// Only one value, use as both min and max
		amount, err := strconv.Atoi(strings.ReplaceAll(matches[0], ",", ""))
		if err == nil {
			return amount, amount
		}
	}

	// Default return in case of no matches or error
	return 0, 0
}

func SplitByHyphen(text string) []string {
	// Clean up spaces and split by hyphen
	parts := strings.Split(text, "-")

	// Trim whitespace from each part
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}

	// Return the resulting array
	return parts
}
