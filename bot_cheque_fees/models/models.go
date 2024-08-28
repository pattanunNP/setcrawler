package models

type ChequeFee struct {
	Providers      string         `json:"providers"`
	FeesTypes      FeesTypes      `json:"fees_types"`
	OthersFees     OthersFees     `json:"others_fees"`
	AdditionalInfo AdditionalInfo `json:"additional_info"`
}

type FeesTypes struct {
	ChequeBookPurchase                []string `json:"cheque_book_purchase"`
	ChequeDepositAcross               []string `json:"cheque_deposit_across"`
	ChequeDepositInbranch             []string `json:"cheque_deposit_inbranch"`
	ChequeReturnFromInstrument        []string `json:"cheque_return_from_insufficient_funds"`
	ChequeFeeReturned                 []string `json:"cheque_fee_returned"`
	ChequeGiftPurchase                string   `json:"cheque_gift_purchase"`
	ChequeCashWithdrawAcross          []string `json:"cheque_cash_withdraw_across"`
	ChequeCashWithdrawInbranch        string   `json:"cheque_cash_withdraw_inbranch"`
	CashierChequePurchase             []string `json:"cashier_cheque_purchase"`
	CashierChequeCashWithdrawAcross   []string `json:"cashier_cheque_cash_withdraw_across"`
	CashierChequeCashWithdrawInbranch []string `json:"cashier_cheque_cash_withdraw_inbranch"`
	DraftPurchaseFee                  []string `json:"draft_purchase_fee"`
	PublicationFee                    string   `json:"publication_fee"`
	ChequeCancellationFee             string   `json:"cheque_cancellation_fee"`
	ChequeAdvanceDepositFee           string   `json:"cheque_advance_deposit_fee"`
}

type OthersFees struct {
	OtherFees string `json:"other_fees"`
}

type AdditionalInfo struct {
	WebsiteFeeLink string `json:"website_fee_link"`
}
