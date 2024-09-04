package models

import "time"

type FeeDetails struct {
	InwardRemittance  InwardRemittance  `json:"inward_remittance"`
	OutwardRemittance OutwardRemittance `json:"outward_remittance"`
}

type InwardRemittance struct {
	Fee                     []string    `json:"fee"`
	FeeNumeric              interface{} `json:"fee_numeric"` // Can be float64, []float64, or nil
	ExchangeCompensationFee string      `json:"exchange_compensation_fee"`
}

type OutwardRemittance struct {
	FeeType                 string            `json:"fee_type"`
	Conditions              []string          `json:"conditions"`
	ConditionsNumeric       ConditionsNumeric `json:"conditions_numeric"`
	ExchangeCompensationFee string            `json:"exchange_compensation_fee"`
}

type ConditionsNumeric struct {
	PromotionStartDate *time.Time `json:"promotion_start_date,omitempty"`
	PromotionEndDate   *time.Time `json:"promotion_end_date,omitempty"`
	FeePerTransaction  *float64   `json:"fee_per_transaction,omitempty"`
	TransactionRange   []float64  `json:"transaction_range,omitempty"`
	CancellationFee    *float64   `json:"cancellation_fee,omitempty"`
}

type CheckAndDraftFees struct {
	TravelerChequeBuyingFee        []string                `json:"traveler_cheque_buying_fee"`
	TravelerChequeBuyingFeeNumeric CheckAndDraftFeeNumeric `json:"traveler_cheque_buying_fee_numeric"`
	TravelerChequeSellingFee       []string                `json:"traveler_cheque_selling_fee"`
	DraftBuyingFee                 []string                `json:"draft_buying_fee"`
	DraftBuyingFeeNumeric          CheckAndDraftFeeNumeric `json:"draft_buying_fee_numeric"`
	DraftSellingFee                []string                `json:"draft_selling_fee"`
	DraftSellingFeeNumeric         CheckAndDraftFeeNumeric `json:"draft_selling_fee_numeric"`
	ForeignBillBuyingFee           []string                `json:"foreign_bill_buying_fee"`
	ForeignBillBuyingFeeNumeric    CheckAndDraftFeeNumeric `json:"foreign_bill_buying_fee_numeric"`
	ForeignBillSellingFee          []string                `json:"foreign_bill_selling_fee"`
	ExchangeCompensationFee        []string                `json:"exchange_compensation_fee"`
}

type CheckAndDraftFeeNumeric struct {
	BaseFee        *float64  `json:"base_fee,omitempty"`
	StampDuty      *float64  `json:"stamp_duty,omitempty"`
	ForeignCharge  *float64  `json:"foreign_charge,omitempty"`
	ReturnFee      *float64  `json:"return_fee,omitempty"`
	StopPaymentFee *float64  `json:"stop_payment_fee,omitempty"`
	OtherFees      []float64 `json:"other_fees,omitempty"`
}

type LetterOfCreditFees struct {
	ForeignLC         string        `json:"foreign_lc"`
	ForeignLCNumeric  LCFeesNumeric `json:"foreign_lc_numeric"`
	DomesticLC        string        `json:"domestic_lc"`
	DomesticLCNumeric LCFeesNumeric `json:"domestic_lc_numeric"`
}

type LCFeesNumeric struct {
	PercentFee *float64  `json:"percent_fee,omitempty"`
	MinFee     *float64  `json:"min_fee,omitempty"`
	AmendFee   *float64  `json:"amend_fee,omitempty"`
	OtherFees  []float64 `json:"other_fees,omitempty"`
}

type BillCollectionFees struct {
	InwardBillFee          string      `json:"inward_bill_fee"`
	OutwardBillFeeExporter string      `json:"outward_bill_fee_exporter"`
	OutwardBillFeeImporter string      `json:"outward_bill_fee_importer"`
	ImportBillFee          string      `json:"import_bill_fee"`
	ExportBillFeeSeller    InvoiceFees `json:"export_bill_fee_seller"`
	ExportBillFeeBuyer     string      `json:"export_bill_fee_buyer"`
}

type InvoiceFees struct {
	FirstInvoice       string `json:"first_invoice"`
	SubsequentInvoices string `json:"subsequent_invoices"`
}

type OtherFees struct {
	OtherFee []string `json:"other_fee"`
}

type AdditionalInfo struct {
	FeeURL *string `json:"fee_url"`
}

type Fees struct {
	Provider                  string             `json:"international_fees"`
	InternationalTransferFees FeeDetails         `json:"international_transfer_fees"`
	CheckAndDraftFees         CheckAndDraftFees  `json:"check_and_draft_fees"`
	LetterOfCreditFees        LetterOfCreditFees `json:"letter_of_credit_fees"`
	BillCollectionFees        BillCollectionFees `json:"bill_collection_fees"`
	OtherFees                 OtherFees          `json:"other_fees"`
	AdditionalInformation     AdditionalInfo     `json:"additional_information"`
}
