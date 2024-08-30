package models

type MoneyTicketFees struct {
	Provider       string         `json:"provider"`
	AcceptanceFee  []string       `json:"acceptance_fee"`
	AvalFee        []string       `json:"aval_fee"`
	OtherFees      *string        `json:"other_fees,omitempty"`
	AdditionalInfo AdditionalInfo `json:"additional_info"`
}

type AdditionalInfo struct {
	FeeLinks []string `json:"fee_links"`
}
