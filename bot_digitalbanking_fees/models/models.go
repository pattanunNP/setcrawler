package models

type DigitalBanking struct {
	Provider   string            `json:"provider"`
	Product    string            `json:"product"`
	Service    ServiceDetails    `json:"service"`
	Fees       FeeDetails        `json:"fees"`
	Additional AdditionalDetails `json:"additional_information"`
}

type ServiceDetails struct {
	Type           string          `json:"type"`
	MainFeature    string          `json:"main_feature"`
	CustomerGroups []CustomerGroup `json:"customer_groups"`
}

type CustomerGroup struct {
	Description         string   `json:"description"`
	AgeRequirement      string   `json:"age_requirement"`
	AccountRequirements []string `json:"account_requirements"`
}

type FeeDetails struct {
	PromptPayTransfer  Fee                 `json:"promptpay_transfer"`
	InterbankTransfer  []TransferCondition `json:"interbank_transfer"`
	IntrabankTransfer  Fee                 `json:"intrabank_transfer"`
	CardlessWithdrawal Fee                 `json:"cardless_withdrawal"`
	EntranceFee        Fee                 `json:"entrance_fee"`
	AnnualFee          Fee                 `json:"annual_fee"`
	OtherFees          []OtherFee          `json:"other_fees"`
}

type Fee struct {
	FeeText     string `json:"fee_text"`
	FeeAmount   int    `json:"fee_amount,omitempty"`
	Description string `json:"description,omitempty"`
}

type TransferCondition struct {
	Description    string `json:"description"`
	ConditionText  string `json:"condition_text"`
	ConditionRange Range  `json:"condition_range"`
	FeeText        string `json:"fee_text"`
	FeeAmount      int    `json:"fee_amount"`
}

type Range struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type OtherFee struct {
	Description string              `json:"description"`
	Conditions  []OtherFeeCondition `json:"conditions"`
}

type OtherFeeCondition struct {
	ConditionText string                 `json:"condition_text"`
	Currency      map[string]CurrencyFee `json:"currency"`
}

type CurrencyFee struct {
	FeeText   string `json:"fee_text"`
	FeeAmount int    `json:"fee_amount"`
}

type AdditionalDetails struct {
	ServiceWebsite *string `json:"service_website"`
	FeeWebsite     *string `json:"fee_website"`
}
