package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"setNews/handle/news"
	"time"
)

func main() {
	baseURL := "https://www.settrade.com/api/cms/v1/news/all?cate=60ad3cae-ba3d-4405-af14-3ed4af1e5065&fromDate=29%2F07%2F2024&toDate=05%2F08%2F2024&orderBy=date&pageIndex="
	pageSize := "&pageSize=20"
	allItems := []news.NewsItem{}
	pageIndex := 0

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Timeout reached, stopping the fetching process.")
				done <- struct{}{}
				return
			default:
				url := fmt.Sprintf("%s%d%s", baseURL, pageIndex, pageSize)
				hasNextPage, err := news.FetchNews(url, &allItems)
				if err != nil {
					fmt.Println("Error fetching news:", err)
					done <- struct{}{}
					return
				}

				if !hasNextPage {
					done <- struct{}{}
					return
				}

				pageIndex++
			}
		}
	}()
	<-done

	saveToJSON("news_results.json", allItems)
}

func saveToJSON(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(data)
}
