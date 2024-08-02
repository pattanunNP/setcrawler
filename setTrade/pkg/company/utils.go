package company

import (
	"encoding/json"
	"fmt"
	"os"
)

func SaveToJSON(filename string, data interface{}) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating JSON file %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	err = encoder.Encode(data)
	if err != nil {
		fmt.Printf("Error encoding JSON data to file %s: %v\n", filename, err)
	}
}
