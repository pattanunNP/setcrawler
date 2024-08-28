package models

type GeneralFees struct {
	EntranceFeeMainCard          string   `json:"main_card_entrance_fee"`
	AnnualFeeMainCard            []string `json:"main_card_annual_fee"`
	CurrencyConversionRisk       []string `json:"currency_conversion_risk"`
	CashAdvanceFee               string   `json:"cash_advance_fee"`
	ReplacementCardFee           string   `json:"replacement_card_fee"`
	EntranceFeeSupplementaryCard string   `json:"supplementary_card_entrance_fee"`
	AnnualFeeSupplementaryCard   []string `json:"supplementary_card_annual_fee"`
	NewPINRequestFee             string   `json:"new_pin_request_fee"`
	StatementCopyFee             string   `json:"statement_copy_fee"`
	TransactionVerificationFee   string   `json:"transaction_verification_fee"`
	SalesSlipCopyFee             string   `json:"sales_slip_copy_fee"`
	ReturnedChequeFee            string   `json:"returned_cheque_fee"`
	TaxPaymentFee                string   `json:"tax_payment_fee"`
	DebtCollectionFee            []string `json:"debt_collection_fee"`
}

type PaymentFees struct {
	FeeFreeChannels        []string `json:"fee_free_channels"`
	DirectDebitServiceFee  string   `json:"direct_debit_service_fee"`
	DirectDebitOtherFee    string   `json:"direct_debit_other_fee"`
	BankCounterFee         string   `json:"bank_counter_fee"`
	OtherBankCounterFee    string   `json:"other_bank_counter_fee"`
	PaymentServicePointFee []string `json:"payment_service_point_fee"`
	OnlinePaymentFee       string   `json:"online_payment_fee"`
	ATMPaymentFee          string   `json:"atm_payment_fee"`
	PhonePaymentFee        string   `json:"phone_payment_fee"`
	ChequeOrMoneyOrderFee  string   `json:"cheque_or_money_order_fee"`
	OtherPaymentChannels   []string `json:"other_payment_channels"`
}

type OthersFees struct {
	OtherFees *string `json:"other_fees"`
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
