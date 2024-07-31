package finance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type FinancialStatusResponse struct {
	Status       bool   `json:"status"`
	ErrorMessage string `json:"errorMessage"`
}

type Period struct {
	Quarter                string  `json:"quarter"`
	Year                   int     `json:"year"`
	BeginDate              string  `json:"beginDate"`
	EndDate                string  `json:"endDate"`
	FinancialStatementLink string  `json:"financialStatementLink"`
	Fscomp                 bool    `json:"fscomp"`
	FsType                 string  `json:"fsType"`
	RestatementDate        *string `json:"restatementDate"`
	StatementType          string  `json:"statementType"`
}

type Value struct {
	Amount        *float64 `json:"amount"`
	PercentChange *float64 `json:"percentChange"`
	Adjusted      *bool    `json:"adjusted"`
}

type Account struct {
	AccountName string  `json:"accnountName"`
	Level       int     `json:"level"`
	Divider     int     `json:"divider"`
	DefaultItem bool    `json:"defaultItem"`
	Values      []Value `json:"values"`
	Format      string  `json:"format"`
}

type FinancialData struct {
	Company       string `json:"company"`
	FinancialData struct {
		Periods  []Period  `json:"periods"`
		Accounts []Account `json:"accounts"`
	} `json:"financialData"`
}

func FetchFinancialData(cookieStr, tokenStr, symbol string) (bool, error) {
	url := fmt.Sprintf("https://www.setsmart.com/ism/api/savesymbol.api?symbol=%s", symbol)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return false, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", tokenStr))
	req.Header.Add("cookie", cookieStr)
	req.Header.Add("referer", "https://www.setsmart.com/ssm/financialStatement;symbol=SCB")
	req.Header.Add("sec-ch-ua", "\"Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\", \"Google Chrome\";v=\"126\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("error reading response body: %w", err)
	}

	var FinancialStatusResponse FinancialStatusResponse
	if err := json.Unmarshal(body, &FinancialStatusResponse); err != nil {
		return false, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	if FinancialStatusResponse.Status {
		return true, nil
	}

	return false, fmt.Errorf("financial status error: %s", FinancialStatusResponse.ErrorMessage)
}

func FetchLatestFinancialStatement(cookieStr, tokenStr, symbol string) (FinancialData, error) {
	url := fmt.Sprintf("https://www.setsmart.com/api/stock/%s/latest-financialstatement?statement=balance_sheet&type=consolidate&compare=qoq&amountType=quarter&lang=en", symbol)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return FinancialData{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", tokenStr))
	req.Header.Add("cookie", cookieStr)
	req.Header.Add("referer", fmt.Sprintf("https://www.setsmart.com/ssm/financialStatement;symbol=%s", symbol))
	req.Header.Add("sec-ch-ua", "\"Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\", \"Google Chrome\";v=\"126\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return FinancialData{}, fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return FinancialData{}, fmt.Errorf("error reading response body: %w", err)
	}

	if res.StatusCode != 200 {
		return FinancialData{}, fmt.Errorf("received non-200 response code: %d", res.StatusCode)
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return FinancialData{}, fmt.Errorf("error unmarshalling financial data: %w", err)
	}

	financialData := FinancialData{
		Company: symbol,
		FinancialData: struct {
			Periods  []Period  `json:"periods"`
			Accounts []Account `json:"accounts"`
		}{
			Periods:  parsePeriods(responseData["periods"]),
			Accounts: parseAccount(responseData["accounts"]),
		},
	}
	return financialData, nil
}

func parsePeriods(data interface{}) []Period {
	var periods []Period
	if data == nil {
		return periods
	}

	periodsData, ok := data.([]interface{})
	if !ok {
		return periods
	}

	for _, item := range periodsData {
		periodMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		period := Period{
			Quarter:                periodMap["quarter"].(string),
			Year:                   int(periodMap["year"].(float64)),
			BeginDate:              periodMap["beginDate"].(string),
			EndDate:                periodMap["endDate"].(string),
			FinancialStatementLink: periodMap["financialStatementLink"].(string),
			Fscomp:                 periodMap["fscomp"].(bool),
			FsType:                 periodMap["fsType"].(string),
			RestatementDate:        parseStringPointer(periodMap["restatementDate"]),
			StatementType:          periodMap["statementType"].(string),
		}
		periods = append(periods, period)
	}
	return periods
}

func parseAccount(data interface{}) []Account {
	var accounts []Account
	if data == nil {
		return accounts
	}

	accountsData, ok := data.([]interface{})
	if !ok {
		return accounts
	}

	for _, item := range accountsData {
		accountMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		account := Account{
			AccountName: accountMap["accountName"].(string),
			Level:       int(accountMap["level"].(float64)),
			Divider:     int(accountMap["divider"].(float64)),
			DefaultItem: accountMap["defaultItem"].(bool),
			Values:      parseValues(accountMap["values"]),
			Format:      accountMap["format"].(string),
		}

		accounts = append(accounts, account)
	}
	return accounts
}

func parseValues(data interface{}) []Value {
	var values []Value
	if data == nil {
		return values
	}

	valuesData, ok := data.([]interface{})
	if !ok {
		return values
	}

	for _, item := range valuesData {
		valueMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		value := Value{
			Amount:        parseFloatPointer(valueMap["amount"]),
			PercentChange: parseFloatPointer(valueMap["percentChange"]),
			Adjusted:      parseBoolPointer(valueMap["adjusted"]),
		}

		values = append(values, value)
	}

	return values
}

func parseStringPointer(data interface{}) *string {
	if data == nil {
		return nil
	}
	str, ok := data.(string)
	if !ok {
		return nil
	}
	return &str
}

func parseFloatPointer(data interface{}) *float64 {
	if data == nil {
		return nil
	}
	floatVal, ok := data.(float64)
	if !ok {
		return nil
	}
	return &floatVal
}

func parseBoolPointer(data interface{}) *bool {
	if data == nil {
		return nil
	}
	boolVal, ok := data.(bool)
	if !ok {
		return nil
	}
	return &boolVal
}

func SaveFinancialDataToJSON(data FinancialData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating JSON file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	err = encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("error encoding JSON data: %w", err)
	}
	return nil
}
