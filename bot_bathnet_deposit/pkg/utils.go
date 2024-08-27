package pkg

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ExtractFeeData(doc *goquery.Document, attrClass, col string) []string {
	var lines []string
	doc.Find("tr." + attrClass + " td.cmpr-col." + col + " span").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		cleanedText := CleanText(text)
		if cleanedText != "" {
			lines = append(lines, cleanedText)
		}
	})
	if len(lines) == 0 {
		return nil
	}

	combinedText := strings.Join(lines, " ")
	cleanedText := CleanText(combinedText)
	feeLines := strings.Split(cleanedText, "/")
	for i := range feeLines {
		feeLines[i] = strings.TrimSpace(feeLines[i])
	}
	return feeLines
}

func CleanText(text string) string {
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, "\t", "")
	text = strings.ReplaceAll(text, "\u00a0", " ")
	text = strings.Join(strings.Fields(text), " ")
	return strings.TrimSpace(text)
}

func DetermineTotalPage(doc *goquery.Document) int {
	totalPages := 1

	doc.Find("ul.pagination li a").Each(func(i int, s *goquery.Selection) {
		pageNum, exists := s.Attr("data-page")
		if exists {
			page, err := strconv.Atoi(pageNum)
			if err == nil && page > totalPages {
				totalPages = page
			}
		}
	})
	return totalPages
}

func AddHeaders(req *http.Request) {
	// Adding headers
	req.Header.Add("accept", "text/plain, */*; q=0.01")
	req.Header.Add("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("cookie", `verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95@AdobeOrg=1; AMCV_F915091E62ED182D0A495F95@AdobeOrg=179643557|MCIDTS|19951|MCMID|53550622918316951353729640026118558196|MCAAMLH-1724305541|3|MCAAMB-1724305541|RKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y|MCOPTOUT-1723707941s|NONE|vVersion|5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0"; _uid6672=16B5DEBD.34; _ctout6672=1; _ga_NLQFGWVNXN=GS1.1.1724750238.42.1.1724751414.60.0.0`)
	req.Header.Add("origin", "https://app.bot.or.th")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("referer", "https://app.bot.or.th/1213/MCPD/FeeApp/BAHTNETFee/CompareProduct")
	req.Header.Add("sec-ch-ua", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", `"macOS"`)
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	req.Header.Add("verificationtoken", "y-9dmHuzupWbGPX11HGC1PKTjK4vNlIdlwZnB2Q-2OrJOoPnX0ZURIqqIaqMpQhrMir1Sgyd-4F7XkqYxr02e8h5rXq88jVJBJOYIjqU0EU1,zOoZLFiN5-W3LghtPZSIDJSOOwAGWgU2ve3oi46kqXYaJCCSX9ehTjfO580uQ2yY9kD2vloSFekcfgNy6d1TVctAZp8OR2ebO7r1rc5I7Ds1")
	req.Header.Add("x-requested-with", "XMLHttpRequest")
}
