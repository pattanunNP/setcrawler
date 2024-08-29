package models

// Define the struct to hold the JSON structure
type FeeDetails struct {
	InwardRemittance  string `json:"inward_remittance"`
	OutwardRemittance string `json:"outward_remittance"`
}

type CheckAndDraftFees struct {
	TravelerChequeBuyingFee  string `json:"traveler_cheque_buying_fee"`
	TravelerChequeSellingFee string `json:"traveler_cheque_selling_fee"`
	DraftBuyingFee           string `json:"draft_buying_fee"`
	DraftSellingFee          string `json:"draft_selling_fee"`
	ForeignBillBuyingFee     string `json:"foreign_bill_buying_fee"`
	ForeignBillSellingFee    string `json:"foreign_bill_selling_fee"`
}

type LetterOfCreditFees struct {
	ForeignLC  string `json:"foreign_lc"`
	DomesticLC string `json:"domestic_lc"`
}

type BillCollectionFees struct {
	InwardBillFee          string `json:"inward_bill_fee"`
	OutwardBillFeeExporter string `json:"outward_bill_fee_exporter"`
	OutwardBillFeeImporter string `json:"outward_bill_fee_importer"`
	ImportBillFee          string `json:"import_bill_fee"`
	ExportBillFeeSeller    string `json:"export_bill_fee_seller"`
	ExportBillFeeBuyer     string `json:"export_bill_fee_buyer"`
}

type OtherFees struct {
	OtherFee string `json:"other_fee"`
}

type AdditionalInfo struct {
	FeeURL string `json:"fee_url"`
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
