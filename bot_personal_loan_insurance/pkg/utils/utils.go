package utils

import (
	"bot_personal_insurance/pkg/models"
	"encoding/json"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func CleanString(input string) string {
	trimmed := strings.TrimSpace(strings.ReplaceAll(input, "\n", " "))
	return strings.Join(strings.Fields(trimmed), " ")
}

func FilterEmptyStrings(input []string) []string {
	var result []string
	for _, str := range input {
		trimmed := strings.TrimSpace(str)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func ParseAgeRange(ageRange string) (*int, *int) {
	ageRange = strings.TrimSpace(ageRange)
	if ageRange == "ไม่มีกำหนด" || ageRange == "" {
		return nil, nil
	}

	agePattern := regexp.MustCompile(`(\d+)-(\d+)`)
	matches := agePattern.FindStringSubmatch(ageRange)

	if len(matches) == 3 {
		minAge, _ := strconv.Atoi(matches[1])
		maxAge, _ := strconv.Atoi(matches[2])
		return &minAge, &maxAge
	}
	return nil, nil
}

func ParseTextIntoArray(text, delimeter string) []string {
	parts := strings.Split(text, delimeter)
	parts = FilterEmptyStrings(parts)
	if len(parts) == 0 {
		return nil
	}
	return parts
}

func WriteJSON(data []models.PersonalLoan, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, jsonData, 0644)
}

func NullIfEmpty(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func ParseFloatFromText(text, keyword string) *float64 {
	re := regexp.MustCompile(`\b(\d+(\.\d+)?)\b`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 0 {
		value, err := strconv.ParseFloat(matches[1], 64)
		if err == nil {
			return &value
		}
	}
	return nil
}

func ParseStampDuty(text string) (float64, int) {
	re := regexp.MustCompile(`(\d+\.\d+)% ของวงเงินสินเชื่อ.*ไม่เกิน (\d+\,?\d*) บาท`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 2 {
		percentage, _ := strconv.ParseFloat(matches[1], 64)
		maxAmount, _ := strconv.Atoi(strings.ReplaceAll(matches[2], ",", ""))
		return percentage, maxAmount
	}
	return 0, 0
}

func ParseCreditCheckFees(text string) (int, int) {
	individualFee, corporateFee := 0, 0
	if strings.Contains(text, "บุคคลธรรมดา") {
		individualFee = parseFeeValueFromText(text, "บุคคลธรรมดา")
	}
	if strings.Contains(text, "นิติบุคคลธรรมดา") {
		corporateFee = parseFeeValueFromText(text, "นิติบุคคลธรรมดา")
	}
	return individualFee, corporateFee
}

func ParseCounterServiceFee(text string) []models.ExtractedCounterServiceFee {
	re := regexp.MustCompile(`(.+?) (\d+) บาท`)
	matches := re.FindAllStringSubmatch(text, -1)
	var fees []models.ExtractedCounterServiceFee
	for _, match := range matches {
		if len(match) == 3 {
			fee, _ := strconv.Atoi(match[2])
			fees = append(fees, models.ExtractedCounterServiceFee{
				Location: match[1],
				Fee:      fee,
			})
		}
	}
	return fees
}

func parseFeeValueFromText(text, keyword string) int {
	re := regexp.MustCompile(keyword + ` (\d+) บาท`)
	match := re.FindStringSubmatch(text)
	if len(match) > 1 {
		fee, _ := strconv.Atoi(match[1])
		return fee
	}
	return 0
}

func ParseBaseMRR(text string) *float64 {
	re := regexp.MustCompile(`MRR.*?(\d+\.\d+)%`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		value, err := strconv.ParseFloat(matches[1], 64)
		if err == nil {
			return &value
		}
	}
	return nil
}

func ParseDefaultInterestRate(text string) *float64 {
	re := regexp.MustCompile(`(\d+)%`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		value, err := strconv.ParseFloat(matches[1], 64)
		if err == nil {
			return &value
		}
	}
	return nil
}

func ParseDebtCollectionFees(text string) []models.ExtractedDebtCollectionFee {
	var fees []models.ExtractedDebtCollectionFee

	// Regular expressions to match different fee conditions
	reGeneral := regexp.MustCompile(`(\d+)\s*บาท/งวด`)
	reOneInstallment := regexp.MustCompile(`ไม่เกิน\s*(\d+)\s*บาท/บัญชี/รอบการทวงถามหนี้\s*\(กรณีมีหนี้ค้างชำระ 1 งวด\)`)
	reMoreInstallments := regexp.MustCompile(`ไม่เกิน\s*(\d+)\s*บาท/บัญชี/รอบการทวงถามหนี้\s*\(กรณีมีหนี้ค้างชำระมากกว่า 1 งวด\)`)
	reNoCharge := regexp.MustCompile(`หนี้ที่ถึงกำหนดชำระสะสมไม่เกิน\s*(\d+)\s*บาท`)

	// Extract values based on different conditions
	if match := reGeneral.FindStringSubmatch(text); len(match) > 1 {
		fee, _ := strconv.Atoi(match[1])
		fees = append(fees, models.ExtractedDebtCollectionFee{
			Condition:         "general",
			FeePerInstallment: fee,
		})
	}

	if match := reOneInstallment.FindStringSubmatch(text); len(match) > 1 {
		fee, _ := strconv.Atoi(match[1])
		fees = append(fees, models.ExtractedDebtCollectionFee{
			Condition:        "one installment overdue",
			MaxFeePerAccount: fee,
		})
	}

	if match := reMoreInstallments.FindStringSubmatch(text); len(match) > 1 {
		fee, _ := strconv.Atoi(match[1])
		fees = append(fees, models.ExtractedDebtCollectionFee{
			Condition:        "more than one installment overdue",
			MaxFeePerAccount: fee,
		})
	}

	if match := reNoCharge.FindStringSubmatch(text); len(match) > 1 {
		fees = append(fees, models.ExtractedDebtCollectionFee{
			Condition: "no charge if debt <= 1,000 THB",
			Fee:       0,
		})
	}

	return fees
}

func ParseOtherFees(text string) []models.ExtractedOtherFees {
	var fees []models.ExtractedOtherFees
	re := regexp.MustCompile(`วงเงิน(.+?)บาท ค่าธรรมเนียม (\d+) บาท`)

	matches := re.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) == 3 {
			loanRange := match[1]
			fee, _ := strconv.Atoi(match[2])
			fees = append(fees, models.ExtractedOtherFees{
				LoanAmountRange: strings.TrimSpace(loanRange),
				Fee:             fee,
			})
		}
	}

	return fees
}
