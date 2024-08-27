package pkg

import (
	"encoding/json"
	"fmt"
	"os"
)

func SaveProductToFile(products []Product, fileName string) error {
	jsonData, err := json.MarshalIndent(products, "", " ")
	if err != nil {
		return fmt.Errorf("failed to convert to JSON: %v", err)
	}

	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to save JSON to file: %v", err)
	}

	return nil
}
