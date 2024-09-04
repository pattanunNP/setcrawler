package models

type ChequeFee struct {
	Providers      string         `json:"providers"`
	FeesTypes      FeesTypes      `json:"fees_types"`
	OthersFees     OthersFees     `json:"others_fees"`
	AdditionalInfo AdditionalInfo `json:"additional_info"`
}

type FeesTypes struct {
	ChequeBookPurchase                []FeeDetail `json:"cheque_book_purchase"`
	ChequeDepositAcross               []FeeDetail `json:"cheque_deposit_across"`
	ChequeDepositInbranch             []FeeDetail `json:"cheque_deposit_inbranch"`
	ChequeReturnFromInstrument        []FeeDetail `json:"cheque_return_from_insufficient_funds"`
	ChequeFeeReturned                 []FeeDetail `json:"cheque_fee_returned"`
	ChequeGiftPurchase                FeeDetail   `json:"cheque_gift_purchase"`
	ChequeCashWithdrawAcross          []FeeDetail `json:"cheque_cash_withdraw_across"`
	ChequeCashWithdrawInbranch        FeeDetail   `json:"cheque_cash_withdraw_inbranch"`
	CashierChequePurchase             []FeeDetail `json:"cashier_cheque_purchase"`
	CashierChequeCashWithdrawAcross   []FeeDetail `json:"cashier_cheque_cash_withdraw_across"`
	CashierChequeCashWithdrawInbranch []FeeDetail `json:"cashier_cheque_cash_withdraw_inbranch"`
	DraftPurchaseFee                  []FeeDetail `json:"draft_purchase_fee"`
	PublicationFee                    FeeDetail   `json:"publication_fee"`
	ChequeCancellationFee             FeeDetail   `json:"cheque_cancellation_fee"`
	ChequeAdvanceDepositFee           FeeDetail   `json:"cheque_advance_deposit_fee"`
}

type FeeDetail struct {
	Text          string   `json:"text"`
	MinFee        *float64 `json:"min_fee,omitempty"`
	MaxFee        *float64 `json:"max_fee,omitempty"`
	PercentageFee *float64 `json:"percentage_fee,omitempty"`
	FeeUnit       string   `json:"fee_unit,omitempty"`
	Condition     string   `json:"condition,omitempty"`
	Service       string   `json:"service,omitempty"`
}

type OthersFees struct {
	OtherFees FeeDetail `json:"other_fees"`
}

type AdditionalInfo struct {
	WebsiteFeeLink string `json:"website_fee_link"`
}
