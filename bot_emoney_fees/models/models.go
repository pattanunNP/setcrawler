package models

type EmoneyFee struct {
	Provider              string          `json:"provider"`
	Product               string          `json:"product"`
	TopUp                 TopUpDetails    `json:"top_up"`
	GeneralFees           GeneralFees     `json:"general_fees"`
	SpendingFees          SpendingFees    `json:"spending_fees"`
	TerminationFees       TerminationFees `json:"termination_fees"`
	OtherFees             OtherFees       `json:"other_fees"`
	AdditionalInformation AdditionalInfo  `json:"additional_info"`
}

type TopUpDetails struct {
	NoFeeChannels []string `json:"no_fee_channels"`
	FeeChannels   string   `json:"fee_channels"`
}

type GeneralFees struct {
	EntranceFee           string  `json:"entrance_fee"`
	AnnualFee             string  `json:"annual_fee"`
	CardReplacementFee    string  `json:"card_replacement_fee"`
	CardReplacementAmount float64 `json:"card_replacement_amount"`
	CardReplacementCond   string  `json:"card_replacement_conditions"`
	MaintenanceFee        string  `json:"maintenance_fee"`
}

type SpendingFees struct {
	SpendingFee              string  `json:"spending_fee"`
	SpendingAlertFee         string  `json:"spending_alert_fee"`
	OverseasWithdrawalFee    string  `json:"overseas_withdrawal_fee"`
	OverseasWithdrawalAmount float64 `json:"overseas_withdrawal_amount"`
	OverseasWithdrawalCond   string  `json:"overseas_withdrawal_conditions"`
	CurrencyConversionFee    string  `json:"currency_conversion_fee"`
	CurrencyConversionRate   float64 `json:"currency_conversion_rate"`
	CurrencyConversionCond   string  `json:"currency_conversion_conditions"`
}

type TerminationFees struct {
	CashRefundFee   string `json:"cash_refund_fee"`
	TerminationFees int    `json:"termination_fee"`
}

type OtherFees struct {
	OtherFeesDetails []string `json:"other_fee_details"`
}

type AdditionalInfo struct {
	FeeURL string `json:"fee_url"`
}
