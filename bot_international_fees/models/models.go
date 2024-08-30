package models

type FeeDetails struct {
	InwardRemittance  InwardRemittance  `json:"inward_remittance"`
	OutwardRemittance OutwardRemittance `json:"outward_remittance"`
}

type InwardRemittance struct {
	Fee                     string `json:"fee"`
	ExchangeCompensationFee string `json:"exchange_compensation_fee"`
}

type OutwardRemittance struct {
	FeeType                 string   `json:"fee_type"`
	Conditions              []string `json:"conditions"`
	ExchangeCompensationFee string   `json:"exchange_compensation_fee"`
}

type CheckAndDraftFees struct {
	TravelerChequeBuyingFee  string `json:"traveler_cheque_buying_fee"`
	TravelerChequeSellingFee string `json:"traveler_cheque_selling_fee"`
	DraftBuyingFee           string `json:"draft_buying_fee"`
	DraftSellingFee          string `json:"draft_selling_fee"`
	ForeignBillBuyingFee     string `json:"foreign_bill_buying_fee"`
	ForeignBillSellingFee    string `json:"foreign_bill_selling_fee"`
	ExchangeCompensationFee  string `json:"exchange_compensation_fee"`
}

type LetterOfCreditFees struct {
	ForeignLC  string `json:"foreign_lc"`
	DomesticLC string `json:"domestic_lc"`
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
