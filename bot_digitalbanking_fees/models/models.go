package models

type DigitalBanking struct {
	Provider   string            `json:"provider"`
	Product    string            `json:"product"`
	Service    ServiceDetails    `json:"service"`
	Fees       FeeDetails        `json:"fees"`
	Additional AdditionalDetails `json:"additional_information"`
}

type ServiceDetails struct {
	ServiceType string   `json:"service_type"`
	MainFeature string   `json:"main_feature"`
	CustomGroup []string `json:"customer_group"`
}

type FeeDetails struct {
	PromptPayTransferFee  string   `json:"promptpay_transfer_fee"`
	InterbankTransferFee  []string `json:"interbank_transfer_fee"`
	IntrabankTransferFee  string   `json:"intrabank_transfer_fee"`
	CardlessWithdrawalFee []string `json:"cardless_withdrawal_fee"`
	EntranceFee           string   `json:"entrance_fee"`
	AnnualFee             string   `json:"annual_fee"`
	OtherFees             []string `json:"other_fees"`
}

type AdditionalDetails struct {
	ServiceWebsite *string `json:"service_website"`
	FeeWebsite     *string `json:"fee_website"`
}
