package models

type MoneyTicketFees struct {
	Provider       string         `json:"provider"`
	AcceptanceFee  []string       `json:"acceptance_fee"`
	AvalFee        []string       `json:"aval_fee"`
	ExtractedInfo  ExtractedInfo  `json:"extracted_info"`
	OtherFees      *string        `json:"other_fees,omitempty"`
	AdditionalInfo AdditionalInfo `json:"additional_info"`
}

type AdditionalInfo struct {
	FeeLinks []string `json:"fee_links"`
}

type ExtractedInfo struct {
	MaxAcceptanceFeePercentage *float64 `json:"max_acceptance_fee_percentage,omitempty"`
	MinAcceptanceFeeBaht       *int     `json:"min_acceptance_fee_baht,omitempty"`
	CancellationFeeBaht        *int     `json:"cancellation_fee_baht,omitempty"`
	MaxAvalFeePercentage       *float64 `json:"max_aval_fee_percentage,omitempty"`
	MinAvalFeeBaht             *int     `json:"min_aval_fee_baht,omitempty"`
}
