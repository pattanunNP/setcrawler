package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Stock struct {
	Symbol string `json:"symbol"`
}

type StockData struct {
	Stock []Stock `json:"stocks"`
}

type Fund struct {
	Symbol   string `json:"symbol"`
	NameEN   string `json:"nameEN"`
	NameTH   string `json:"nameTH"`
	ID       string `json:"id"`
	AimcType string `json:"aimcType"`
}

func ReadJSONFile(filename string) (*StockData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	var stocksData StockData
	if err := json.Unmarshal(data, &stocksData); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return &stocksData, nil
}

func ReadFundsJSONFile(filepath string) ([]Fund, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening JSON file: %w", err)
	}
	defer file.Close()

	var funds []Fund
	err = json.NewDecoder(file).Decode(&funds)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return funds, nil
}
