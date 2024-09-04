package models

type FeeDetail struct {
	Text       string   `json:"text"`
	Amount     *int     `json:"amount,omitempty"`
	Percentage *float64 `json:"percentage,omitempty"`
}

type AnnualFeeDetail struct {
	Text          string `json:"text"`
	InitialAmount *int   `json:"initial_amount,omitempty"`
}

type GeneralFees struct {
	EntranceFeeMainCard          FeeDetail         `json:"main_card_entrance_fee"`
	AnnualFeeMainCard            []AnnualFeeDetail `json:"main_card_annual_fee"`
	CurrencyConversionRisk       FeeDetail         `json:"currency_conversion_risk"`
	CashAdvanceFee               FeeDetail         `json:"cash_advance_fee"`
	ReplacementCardFee           FeeDetail         `json:"replacement_card_fee"`
	EntranceFeeSupplementaryCard FeeDetail         `json:"supplementary_card_entrance_fee"`
	AnnualFeeSupplementaryCard   []AnnualFeeDetail `json:"supplementary_card_annual_fee"`
	NewPINRequestFee             FeeDetail         `json:"new_pin_request_fee"`
	StatementCopyFee             FeeDetail         `json:"statement_copy_fee"`
	TransactionVerificationFee   FeeDetail         `json:"transaction_verification_fee"`
	SalesSlipCopyFee             FeeDetail         `json:"sales_slip_copy_fee"`
	ReturnedChequeFee            FeeDetail         `json:"returned_cheque_fee"`
	TaxPaymentFee                FeeDetail         `json:"tax_payment_fee"`
	DebtCollectionFee            []FeeDetail       `json:"debt_collection_fee"`
}

type PaymentFees struct {
	FeeFreeChannels       []string  `json:"fee_free_channels"`
	DirectDebitServiceFee FeeDetail `json:"direct_debit_service_fee"`
	BankCounterFee        FeeDetail `json:"bank_counter_fee"`
	OnlinePaymentFee      FeeDetail `json:"online_payment_fee"`
	ATMPaymentFee         FeeDetail `json:"atm_payment_fee"`
	PhonePaymentFee       FeeDetail `json:"phone_payment_fee"`
	OtherPaymentChannels  []string  `json:"other_payment_channels"`
}

type OthersFees struct {
	OtherFees *FeeDetail `json:"other_fees"`
}

type AdditionalInfo struct {
	WebsiteFeeLink string `json:"website_fee_link"`
}

type CreditCardFee struct {
	Provider       string         `json:"provider"`
	Product        string         `json:"product"`
	GeneralFees    GeneralFees    `json:"general_fees"`
	PaymentFees    PaymentFees    `json:"payment_fees"`
	OtherFees      OthersFees     `json:"other_fees"`
	AdditionalInfo AdditionalInfo `json:"additional_info"`
}
