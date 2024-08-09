package research

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"setResearch/httpclient"
	"setResearch/utils"
)

type ResearchItem struct {
	UUID         string            `json:"uuid"`
	RenderType   string            `json:"renderType"`
	Title        string            `json:"title"`
	URL          string            `json:"url"`
	IsSuggestTag bool              `json:"isSuggestTag"`
	IsTodayTag   bool              `json:"isTodayTag"`
	Symbol       string            `json:"symbol"`
	SymbolId     string            `json:"symbolId"`
	Market       string            `json:"market"`
	CateUuid     string            `json:"cateUuid"`
	SubCateUuid  string            `json:"subCateUuid"`
	CateCode     string            `json:"cateCode"`
	SubCateCode  string            `json:"subCateCode"`
	CateName     string            `json:"cateName"`
	SubCateName  string            `json:"subCateName"`
	StartDate    string            `json:"startDate"`
	Source       string            `json:"source"`
	Views        int               `json:"views"`
	HTMLFile     string            `json:"htmlContent"`
	FileURL      string            `json:"fileurl"`
	PDFContent   map[string]string `json:"pdfcontent"`
}

type ResearchItems struct {
	IndexFrom  int            `json:"indexFrom"`
	PageIndex  int            `json:"pageIndex"`
	PageSize   int            `json:"pageSize"`
	TotalCount int            `json:"totalcount"`
	TotalPages int            `json:"totalPages"`
	Items      []ResearchItem `json:"items"`
}

type ResponseData struct {
	ResearchItems ResearchItems `json:"researchItems"`
}

func FetchResearchItems(client *http.Client, pageIndex int) (*ResponseData, error) {
	url := fmt.Sprintf("https://www.settrade.com/api/cms/v1/research-settrade/search?startDate=02%%2F08%%2F2024&endDate=09%%2F08%%2F2024&pageSize=20&pageIndex=%d", pageIndex)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header = httpclient.SetRequestHeaders()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var responseData ResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return &responseData, nil
}

func ProcessResearchItem(client *http.Client, item *ResearchItem) error {
	htmlContent, err := utils.FetchHTMLContent(client, item.URL)
	if err != nil {
		return fmt.Errorf("error fetching HTML content: %v", err)
	}

	htmlContent = utils.DecodeUnicodeEscapeSequences(htmlContent)
	fileURL := utils.ExtractFileURLFromHTML(htmlContent)
	item.FileURL = fileURL

	if fileURL != "" {
		pdfContent, err := ProcessPDF(fileURL)
		if err != nil {
			return fmt.Errorf("error processign PDF: %v", err)
		}
		item.PDFContent = map[string]string{"Full PDF": pdfContent}
	}
	return nil
}
