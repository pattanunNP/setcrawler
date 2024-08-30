package models

// GeneralFee represents the general fee details
type GeneralFee struct {
	InternalEvaluation    FeeDetail `json:"internal_evaluation"`
	ExternalEvaluation    FeeDetail `json:"external_evaluation"`
	StampDuty             FeeDetail `json:"stamp_duty"`
	MortgageFee           FeeDetail `json:"mortgage_fee"`
	CreditCheck           FeeDetail `json:"credit_check"`
	ReturnedChequeFee     FeeDetail `json:"returned_cheque_fee"`
	InsufficientFundsFee  FeeDetail `json:"insufficient_funds_fee"`
	StatementReIssuingFee []string  `json:"statement_reissuing_fee"`
	DebtCollectionFee     []string  `json:"debt_collection_fee"`
}

// PaymentFee represents the payment fee details
type PaymentFee struct {
	DebitFromAccount         FeeDetail `json:"debit_from_account"`
	DebitFromOtherAccount    FeeDetail `json:"debit_from_other_account"`
	PayAtProviderBranch      FeeDetail `json:"pay_at_provider_branch"`
	PayAtOtherBranch         FeeDetail `json:"pay_at_other_branch"`
	PayAtServicePoint        FeeDetail `json:"pay_at_service_point"`
	PayOnline                FeeDetail `json:"pay_online"`
	PayViaCDMATM             FeeDetail `json:"pay_via_cdm_atm"`
	PayViaPhone              FeeDetail `json:"pay_via_phone"`
	PayViaChequeOrMoneyOrder FeeDetail `json:"pay_via_cheque_or_money_order"`
	PayViaOtherChannels      FeeDetail `json:"pay_via_other_channels"`
}

// OtherFee represents other fees
type OtherFee struct {
	OtherFee string `json:"other_fee"`
}

// AdditionalInfo represents additional information details
type AdditionalInfo struct {
	FeeWebsite string `json:"fee_website"`
}

// FeeDetail represents a fee with original text and extracted numerical details
type FeeDetail struct {
	OriginalText string   `json:"original_text"`
	Percentage   *float64 `json:"percentage,omitempty"`
	FeeAmount    *int     `json:"fee_amount,omitempty"`
}

// PersonalFeeDetails represents the complete details of personal fees
type PersonalFeeDetails struct {
	ServiceProvider string         `json:"service_provider"`
	Product         string         `json:"product"`
	GeneralFees     GeneralFee     `json:"general_fees"`
	PaymentFees     PaymentFee     `json:"payment_fees"`
	OtherFees       OtherFee       `json:"other_fees"`
	AdditionalInfo  AdditionalInfo `json:"additional_info"`
}
