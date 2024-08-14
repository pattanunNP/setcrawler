package analyst

// AnalystData represents the structured data for each analyst entry.
type AnalystData struct {
	ID                       string  `json:"id"`
	Symbol                   string  `json:"symbol"`
	BrokerName               string  `json:"brokerName"`
	BrokerURL                string  `json:"brokerURL"`
	AnalystName              string  `json:"analystName"`
	CurrentYearEps           float64 `json:"currentYearEps"`
	NextYearEps              float64 `json:"nextYearEps,omitempty"`
	CurrentYearNetProfit     float64 `json:"currentYearNetProfit"`
	NextYearNetProfit        float64 `json:"nextYearNetProfit,omitempty"`
	CurrentYearPe            float64 `json:"currentYearPe"`
	NextYearPe               float64 `json:"nextYearPe,omitempty"`
	CurrentYearPbv           float64 `json:"currentYearPbv"`
	NextYearPbv              float64 `json:"nextYearPbv,omitempty"`
	CurrentYearDiv           float64 `json:"currentYearDiv"`
	NextYearDiv              float64 `json:"nextYearDiv,omitempty"`
	TargetPrice              float64 `json:"targetPrice"`
	TargetPriceChange        float64 `json:"targetPriceChange,omitempty"`
	TargetPricePercentChange float64 `json:"targetPricePercentChange,omitempty"`
	Recommend                string  `json:"recommend"`
	RecommendType            string  `json:"recommendType"`
	LastUpdateDate           string  `json:"lastUpdateDate"`
	LastResearchURL          string  `json:"lastResearchURL"`
	FullResearchURL          string  `json:"fullResearchURL,omitempty"`
	LastResearchId           string  `json:"lastResearchId"`
	FullResearchId           string  `json:"fullResearchId,omitempty"`
	ResearchText             string  `json:"researchText"`
}

// Stock defines the structure for each stock entry in the JSON
type Stock struct {
	Symbol             string        `json:"symbol"`
	LastPrice          float64       `json:"lastPrice"`
	TotalCoverage      int           `json:"totalCoverage"`
	Buy                int           `json:"buy"`
	Hold               int           `json:"hold"`
	Sell               int           `json:"sell"`
	RecommendType      string        `json:"recommendType"`
	MedianTargetPrice  float64       `json:"medianTargetPrice"`
	AverageTargetPrice float64       `json:"averageTargetPrice"`
	Bullish            float64       `json:"bullish"`
	Bearish            float64       `json:"bearish"`
	AnalystData        []AnalystData `json:"analystData"`
}

// Response defines the structure of the entire JSON file
type Response struct {
	MarketTime string  `json:"marketTime"`
	Overall    []Stock `json:"overall"`
}
