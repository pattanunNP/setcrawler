package main

import (
	"encoding/json"
	"fmt"
	"os"
	"setNews/handle/news"
)

func main() {
	baseURL := "https://www.settrade.com/api/cms/v1/news/all?cate=60ad3cae-ba3d-4405-af14-3ed4af1e5065&fromDate=29%2F07%2F2024&toDate=05%2F08%2F2024&orderBy=date&pageIndex="
	pageSize := "&pageSize=20"
	allItems := []news.NewsItem{}
	pageIndex := 0

	for {

		url := fmt.Sprintf("%s%d%s", baseURL, pageIndex, pageSize)
		hasNextPage, err := news.FetchNews(url, &allItems)
		if err != nil {
			fmt.Println("Error fetching news:", err)
			break
		}

		if !hasNextPage {
			break
		}
	}

	// Save all news items to a single JSON file, split by title
	if err := saveToJSON("news_result.json", allItems); err != nil {
		fmt.Println("Error saving news items to JSON:", err)
	}
}

func saveToJSON(filename string, items []news.NewsItem) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make(map[string]news.NewsItem)
	for _, item := range items {
		data[item.Title] = item
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(data)
}
