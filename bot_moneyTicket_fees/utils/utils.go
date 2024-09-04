package utils

import (
	"moneyticket_fees/models"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func AddHeader(req *http.Request) {
	// Set headers as per the curl command
	req.Header.Set("Accept", "text/plain, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Cookie", `verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0"; _uid6672=16B5DEBD.48; _ga_NLQFGWVNXN=GS1.1.1724995102.59.0.1724995102.60.0.0; visit_time=2040`)
	req.Header.Set("Origin", "https://app.bot.or.th")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://app.bot.or.th/1213/MCPD/FeeApp/OtherFee/AvalAndAcceptanceServiceFee/CompareProduct")
	req.Header.Set("Sec-CH-UA", `"Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome";v="128"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	req.Header.Set("VerificationToken", `L5XnbAMWiq5G0jv50yXvsGD00QP9Z0FekrB5YZEFVhTITrb8_ZMbQGCWGA6pWpbRfIu8Kba6PvszFBGDo3bf2uUj6uWlt3PM-NO0IF6hWEE1,In5_ggox11IW_m3u6owT5WuvRMCpSn2_H3Vx303rhE7368CVjoIoucIl9psmzKY9oE1sl28jSVmKYRdJN5Lp5GAnKqDndYQEQvmVCx9aQrE1`)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
}

func CleanText(text string) string {
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.TrimSpace(text)
	text = strings.Join(strings.Fields(text), " ")
	return text
}

func ParseFeeDetailsAsArray(doc *goquery.Document, selector string) []string {
	var details []string
	doc.Find(selector).Each(func(index int, item *goquery.Selection) {
		text := CleanText(item.Text())
		if text != "" {
			// Split by numbered items "1.", "2.", etc.
			splitDetails := splitByNumberedPattern(text)
			details = append(details, splitDetails...)
		}
	})

	return details
}

func ExtractFeeDetails(doc *goquery.Document, selector string) []string {
	var details []string
	doc.Find(selector).Each(func(index int, item *goquery.Selection) {
		// Combine text from all <span> elements within the selector
		text := CleanText(item.Text())
		// Process text and split by line breaks
		parts := SplitAndCleanText(text)
		details = append(details, parts...)
	})

	// Filter out empty strings
	var filteredDetails []string
	for _, detail := range details {
		if detail != "" {
			filteredDetails = append(filteredDetails, detail)
		}
	}

	return filteredDetails
}

func splitByNumberedPattern(text string) []string {
	// Regular expression to match numbered patterns
	re := regexp.MustCompile(`(\d+\.)`)
	parts := re.Split(text, -1)
	var result []string

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}

	return result
}

func SplitAndCleanText(text string) []string {
	var result []string
	// Split by numbers "1.", "2.", etc.
	parts := strings.Split(text, " ")
	var currentText string
	for _, part := range parts {
		if strings.HasPrefix(part, "1.") || strings.HasPrefix(part, "2.") {
			if currentText != "" {
				result = append(result, strings.TrimSpace(strings.ReplaceAll(currentText, ",", "")))
			}
			currentText = part
		} else {
			currentText += " " + part
		}
	}
	if currentText != "" {
		result = append(result, strings.TrimSpace(strings.ReplaceAll(currentText, ",", "")))
	}
	return result
}

func ExtractNumericInfo(acceptanceFee, avalFee []string) models.ExtractedInfo {
	var extractedInfo models.ExtractedInfo

	// Define regex patterns for extracting numbers
	percentagePattern := regexp.MustCompile(`(\d+(\.\d+)?)%`)
	minFeePattern := regexp.MustCompile(`ขั้นต่ำ (\d+) บาท`)
	cancellationFeePattern := regexp.MustCompile(`ค่าธรรมเนียมฉบับละ (\d+) บาท`)

	// Extract values for acceptance fees
	for _, fee := range acceptanceFee {
		if percentagePattern.MatchString(fee) {
			percentage, _ := strconv.ParseFloat(percentagePattern.FindStringSubmatch(fee)[1], 64)
			extractedInfo.MaxAcceptanceFeePercentage = &percentage
		}
		if minFeePattern.MatchString(fee) {
			minFee, _ := strconv.Atoi(minFeePattern.FindStringSubmatch(fee)[1])
			extractedInfo.MinAcceptanceFeeBaht = &minFee
		}
		if cancellationFeePattern.MatchString(fee) {
			cancellationFee, _ := strconv.Atoi(cancellationFeePattern.FindStringSubmatch(fee)[1])
			extractedInfo.CancellationFeeBaht = &cancellationFee
		}
	}

	// Extract values for aval fees
	for _, fee := range avalFee {
		if percentagePattern.MatchString(fee) {
			percentage, _ := strconv.ParseFloat(percentagePattern.FindStringSubmatch(fee)[1], 64)
			extractedInfo.MaxAvalFeePercentage = &percentage
		}
		if minFeePattern.MatchString(fee) {
			minFee, _ := strconv.Atoi(minFeePattern.FindStringSubmatch(fee)[1])
			extractedInfo.MinAvalFeeBaht = &minFee
		}
	}

	return extractedInfo
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
