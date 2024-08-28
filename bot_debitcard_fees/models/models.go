package models

// Struct to hold detailed fee information
type Fee struct {
	EntranceFee                string           `json:"entrance_fee"`
	AnnualFee                  *AnnualFeeDetail `json:"annual_fee,omitempty"`
	CardReplacementFee         *[]string        `json:"card_replacement_fee"`
	PINReplacementFee          *string          `json:"pin_replacement_fee"`
	StatementRequestFee        *[]string        `json:"statement_request_fee"`
	TransactionSlipRequestFee  *string          `json:"transaction_slip_request_fee"`
	TransactionVerificationFee *string          `json:"transaction_verification_fee"`
}

// Struct to hold details about domestic transaction fees
type DomesticTransaction struct {
	FreeTransactionCount     *[]string `json:"free_transaction_count"`
	OwnProviderTransaction   *string   `json:"own_provider_transaction,omitempty"`
	BalanceInquiryFeeOut     *string   `json:"balance_inquiry_fee_out"`
	WithdrawFeeOut           *string   `json:"withdraw_fee_out"`
	TransferFeeOut           *string   `json:"transfer_fee_out"`
	OtherProviderTransaction *string   `json:"other_provider_transaction,omitempty"`
	BalanceInquiryFeeIn      *string   `json:"balance_inquiry_fee_in"`
	WithdrawFeeIn            *string   `json:"withdraw_fee_in"`
	TransferFeeIn            *string   `json:"transfer_fee_in"`
	BalanceInquiryFeeOutAlt  *string   `json:"balance_inquiry_fee_out_alt"`
	WithdrawFeeOutAlt        *string   `json:"withdraw_fee_out_alt"`
	TransferFeeOutAlt        *string   `json:"transfer_fee_out_alt"`
	CrossProviderTransferFee *string   `json:"cross_provider_transfer_fee,omitempty"`
	TransferLimit10k         *string   `json:"transfer_limit_10k"`
	TransferLimit50k         *string   `json:"transfer_limit_50k"`
	AdditionalFee            *string   `json:"additional_fee"`
}

// Struct to hold annual fee details
type AnnualFeeDetail struct {
	Amount     int     `json:"amount"`
	Conditions *string `json:"conditions"`
}

// Struct to hold overall fee data, including general fees and domestic transaction fees
type FeesData struct {
	GeneralFees  Fee                 `json:"general_fees"`
	DomesticFees DomesticTransaction `json:"domestic_transaction_fees"`
}

type InternationalTransaction struct {
	WithdrawalFee       *string `json:"withdrawal_fee"`
	BalanceInquiryFee   *string `json:"balance_inquiry_fee"`
	CurrencyExchangeFee *string `json:"currency_exchange_fee"`
}

type OtherFees struct {
	OtherFees *[]string `json:"other_fees"`
}

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
