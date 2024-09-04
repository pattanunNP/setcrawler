package pkg

type ExtractedDetails struct {
	Payer     *string `json:"payer,omitempty"`
	Fee       *int    `json:"fee,omitempty"`
	Currency  *string `json:"currency,omitempty"`
	Condition *string `json:"condition,omitempty"`
}

type FeeEntry struct {
	Description string           `json:"description"`
	Extracted   ExtractedDetails `json:"extracted"`
}

type FeeDetails struct {
	TransferWithinBangkokAndVicinity             []FeeEntry `json:"transfer_within_bangkok_and_vicinity"`
	TransferFromBangkokToRegion                  []FeeEntry `json:"transfer_from_bangkok_to_region"`
	TransferFromRegionToBangkok                  []FeeEntry `json:"transfer_from_region_to_bangkok"`
	TransferWithinRegion                         []FeeEntry `json:"transfer_within_region"`
	TransferFromBangkokToOtherBankAccount        []FeeEntry `json:"transfer_from_bangkok_to_other_bank_account"`
	TransferFromRegionToOtherBankAccount         []FeeEntry `json:"transfer_from_region_to_other_bank_account"`
	ReceiveTransferInBangkokFromOtherBankAccount []FeeEntry `json:"receive_transfer_in_bangkok_from_other_bank_account"`
	ReceiveTransferInRegionFromOtherBankAccount  []FeeEntry `json:"receive_transfer_in_region_from_other_bank_account"`
}

type AdditionalInfo struct {
	FeeWebiteLink string `json:"fee_website_link"`
}

type BathnetFee struct {
	Provider       string         `json:"provider"`
	Fees           FeeDetails     `json:"fees_details"`
	AdditionalInfo AdditionalInfo `json:"additional_info"`
}
