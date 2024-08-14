package analyst

import (
	"fmt"
	"html"
	"regexp"
	"setIAA/pkg/pdftools"
	"strconv"
	"strings"
)

// parseAnalystData processes each line and assigns the parsed data to the appropriate field in the AnalystData struct
func parseAnalystData(lines []string) (AnalystData, error) {
	data := AnalystData{}
	re := regexp.MustCompile(`(?P<key>[\w]+):\s?(?P<value>.+)`)

	urlRegex := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)

	var buffer string
	inMultiLineValue := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if inMultiLineValue {
			if strings.HasSuffix(line, `"`) {
				buffer += " " + line
				inMultiLineValue = false
			} else {
				buffer += " " + line
				continue
			}
		}

		if strings.HasPrefix(line, `"`) && !strings.HasSuffix(line, `"`) {
			buffer = line
			inMultiLineValue = true
			continue
		}

		if buffer != "" && !inMultiLineValue {
			line = buffer + " " + line
			buffer = ""
		}

		line = html.UnescapeString(line)

		match := re.FindStringSubmatch(line)
		if len(match) == 0 {
			continue
		}

		key := strings.TrimSpace(match[1])
		value := strings.TrimSpace(match[2])
		value = strings.ReplaceAll(value, `\"`, "")

		switch key {
		case "id":
			data.ID = value
		case "symbol":
			data.Symbol = value
		case "brokerName":
			data.BrokerName = value
		case "brokerURL":
			urlMatch := urlRegex.FindString(value)
			if urlMatch != "" {
				data.BrokerURL = urlMatch
			}
		case "analystName":
			data.AnalystName = value
		case "currentYearEps":
			data.CurrentYearEps, _ = strconv.ParseFloat(value, 64)
		case "nextYearEps":
			data.NextYearEps, _ = strconv.ParseFloat(value, 64)
		case "currentYearNetProfit":
			data.CurrentYearNetProfit, _ = strconv.ParseFloat(value, 64)
		case "nextYearNetProfit":
			data.NextYearNetProfit, _ = strconv.ParseFloat(value, 64)
		case "currentYearPe":
			data.CurrentYearPe, _ = strconv.ParseFloat(value, 64)
		case "nextYearPe":
			data.NextYearPe, _ = strconv.ParseFloat(value, 64)
		case "currentYearPbv":
			data.CurrentYearPbv, _ = strconv.ParseFloat(value, 64)
		case "nextYearPbv":
			data.NextYearPbv, _ = strconv.ParseFloat(value, 64)
		case "currentYearDiv":
			data.CurrentYearDiv, _ = strconv.ParseFloat(value, 64)
		case "nextYearDiv":
			data.NextYearDiv, _ = strconv.ParseFloat(value, 64)
		case "targetPrice":
			data.TargetPrice, _ = strconv.ParseFloat(value, 64)
		case "targetPriceChange":
			data.TargetPriceChange, _ = strconv.ParseFloat(value, 64)
		case "targetPricePercentChange":
			data.TargetPricePercentChange, _ = strconv.ParseFloat(value, 64)
		case "recommend":
			data.Recommend = value
		case "recommendType":
			data.RecommendType = value
		case "lastUpdateDate":
			data.LastUpdateDate = value
		case "lastResearchURL":
			urlMatch := urlRegex.FindString(value)
			if urlMatch != "" {
				data.LastResearchURL = urlMatch

				// Download and extract text from PDF
				researchText, err := pdftools.ProcessPDF(data.LastResearchURL)
				if err != nil {
					fmt.Printf("Error extracting text from PDF at URL %s: %v\n", data.LastResearchURL, err)
				} else {
					data.ResearchText = cleanExtractedText(researchText)
				}
			}
		case "fullResearchURL":
			data.FullResearchURL = value
		case "lastResearchId":
			data.LastResearchId = value
		case "fullResearchId":
			data.FullResearchId = value
		default:
			fmt.Printf("Unhandled key: %s\n", key)
		}
	}

	return data, nil
}

// ExtractAnalystData extracts relevant analyst data from the script content
func ExtractAnalystData(scriptContent string) []AnalystData {
	decodedContent := decodeUnicodeEscape(scriptContent)

	analystDataParts := strings.Split(decodedContent, "},{")
	var result []AnalystData

	unwantedKeywords := []string{
		"datetime", "localDatetime", "price", "volume", "value", "highlightData", "quotationChartAccumulated", "historicalTrading", "esg", "sectorComparison",
	}

	for _, part := range analystDataParts {
		if strings.Contains(part, "analystName") {
			include := true
			for _, keyword := range unwantedKeywords {
				if strings.Contains(part, keyword) {
					include = false
					break
				}
			}
			if include {
				cleanPart := strings.Trim(part, "{} ")
				lines := strings.Split(cleanPart, ",")

				analystData, err := parseAnalystData(lines)
				if err == nil {
					result = append(result, analystData)
				} else {
					fmt.Printf("Error parsing analyst data: %v\n", err)
				}
			}
		}
	}

	return result
}

func decodeUnicodeEscape(s string) string {
	replacer := strings.NewReplacer(
		"\\u003c", "<",
		"\\u003e", ">",
		"\\u0026", "&",
		"\\u0022", `"`,
		"\\u0027", "'",
		"\\u002F", "/",
	)
	return replacer.Replace(s)
}

func cleanExtractedText(text string) string {
	cleanedText := strings.ReplaceAll(text, "The encoding parameter is ignored when writing to the console.\n", "")
	cleanedText = strings.TrimSpace(cleanedText)
	return cleanedText
}
