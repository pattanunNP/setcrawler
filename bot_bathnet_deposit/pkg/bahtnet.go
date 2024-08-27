package pkg

type FeeDetails struct {
	TransferWithinBangkokAndVicinity             []string `json:"transfer_within_bangkok_and_vicinity"`
	TransferFromBangkokToRegion                  []string `json:"transfer_from_bangkok_to_region"`
	TransferFromRegionToBangkok                  []string `json:"transfer_from_region_to_bangkok"`
	TransferWithinRegion                         []string `json:"transfer_within_region"`
	TransferFromBangkokToOtherBankAccount        []string `json:"transfer_from_bangkok_to_other_bank_account"`
	TransferFromRegionToOtherBankAccount         []string `json:"transfer_from_region_to_other_bank_account"`
	ReceiveTransferInBangkokFromOtherBankAccount []string `json:"receive_transfer_in_bangkok_from_other_bank_account"`
	ReceiveTransferInRegionFromOtherBankAccount  []string `json:"receive_transfer_in_region_from_other_bank_account"`
}

type AdditionalInfo struct {
	FeeWebiteLink string `json:"fee_website_link"`
}

type BathnetFee struct {
	Provider       string         `json:"provider"`
	Fees           FeeDetails     `json:"fees_details"`
	AdditionalInfo AdditionalInfo `json:"additional_info"`
}
