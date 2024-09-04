package models

// Struct to hold detailed fee information
type Fee struct {
	EntranceFee                      string           `json:"entrance_fee"`
	EntranceFeeAmount                int              `json:"entrance_fee_amount,omitempty"`
	AnnualFee                        *AnnualFeeDetail `json:"annual_fee,omitempty"`
	CardReplacementFee               *[]string        `json:"card_replacement_fee"`
	CardReplacementFeeAmount         int              `json:"card_replacement_fee_amount,omitempty"`
	PINReplacementFee                *string          `json:"pin_replacement_fee"`
	PINReplacementFeeAmount          int              `json:"pin_replacement_fee_amount,omitempty"`
	StatementRequestFee              *[]string        `json:"statement_request_fee"`
	StatementRequestFeeAmount        int              `json:"statement_request_fee_amount,omitempty"`
	TransactionSlipRequestFee        *string          `json:"transaction_slip_request_fee"`
	TransactionSlipRequestFeeAmount  int              `json:"transaction_slip_request_fee_amount,omitempty"`
	TransactionVerificationFee       *string          `json:"transaction_verification_fee"`
	TransactionVerificationFeeAmount int              `json:"transaction_verification_fee_amount,omitempty"`
}

// Struct to hold details about domestic transaction fees
type DomesticTransaction struct {
	FreeTransactionCount       int       `json:"free_transaction_count"`
	FreeTransactionConditions  *[]string `json:"free_transaction_conditions,omitempty"`
	BalanceInquiryFeeOut       *string   `json:"balance_inquiry_fee_out"`
	BalanceInquiryFeeOutAmount int       `json:"balance_inquiry_fee_out_amount,omitempty"`
	WithdrawFeeOut             *string   `json:"withdraw_fee_out"`
	WithdrawFeeOutAmount       int       `json:"withdraw_fee_out_amount,omitempty"`
	TransferFeeOut             *string   `json:"transfer_fee_out"`
	TransferFeeOutAmount       int       `json:"transfer_fee_out_amount,omitempty"`
	BalanceInquiryFeeIn        *string   `json:"balance_inquiry_fee_in"`
	BalanceInquiryFeeInAmount  int       `json:"balance_inquiry_fee_in_amount,omitempty"`
	WithdrawFeeIn              *string   `json:"withdraw_fee_in"`
	WithdrawFeeInAmount        int       `json:"withdraw_fee_in_amount,omitempty"`
	TransferFeeIn              *string   `json:"transfer_fee_in"`
	TransferFeeInAmount        int       `json:"transfer_fee_in_amount,omitempty"`
	TransferLimit10k           *string   `json:"transfer_limit_10k"`
	TransferLimit10kAmount     int       `json:"transfer_limit_10k_amount,omitempty"`
	TransferLimit50k           *string   `json:"transfer_limit_50k"`
	TransferLimit50kAmount     int       `json:"transfer_limit_50k_amount,omitempty"`
	AdditionalFee              *string   `json:"additional_fee"`
	AdditionalFeeAmount        int       `json:"additional_fee_amount,omitempty"`
}

// Struct to hold annual fee details
type AnnualFeeDetail struct {
	Amount     int     `json:"amount"`
	Conditions *string `json:"conditions"`
}

// Struct to hold international transaction fee details
type InternationalTransaction struct {
	WithdrawalFee              *string `json:"withdrawal_fee"`
	WithdrawalFeeAmount        int     `json:"withdrawal_fee_amount,omitempty"`
	BalanceInquiryFee          *string `json:"balance_inquiry_fee"`
	BalanceInquiryFeeAmount    int     `json:"balance_inquiry_fee_amount,omitempty"`
	CurrencyExchangeFee        *string `json:"currency_exchange_fee"`
	CurrencyExchangeFeePercent float64 `json:"currency_exchange_fee_percent,omitempty"`
}

// Struct to hold other fees
type OtherFees struct {
	OtherFees       *[]string `json:"other_fees"`
	OtherFeesAmount int       `json:"other_fees_amount,omitempty"`
}

// Struct to hold additional information
type AdditionalInfo struct {
	FeeWebsite *string `json:"fee_website_link"`
}

// Struct to hold debit fee information, integrating with Fee and DomesticTransaction
type DebitFee struct {
	Provider          string                    `json:"provider"`
	Product           string                    `json:"product"`
	GeneralFees       Fee                       `json:"general_fees"`
	DomesticFees      DomesticTransaction       `json:"domestic_transaction_fees"`
	InternationalFees *InternationalTransaction `json:"international_fees,omitempty"`
	OtherFees         *OtherFees                `json:"other_fees,omitempty"`
	AdditionalInfo    *AdditionalInfo           `json:"additional_info"`
}
