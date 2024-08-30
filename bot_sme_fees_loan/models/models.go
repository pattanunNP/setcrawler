package models

type SMEProduct struct {
	ServiceProvider string         `json:"serviceProvider"`
	Product         string         `json:"product"`
	FeeDetails      FeeDetails     `json:"fee_details"`
	LoanFees        LoanFees       `json:"loan_fees"`
	OtherFees       OtherFees      `json:"other_fees"`
	AdditionalInfo  AdditionalInfo `json:"additional_info"`
}

// Define the structures for storing different fee details
type FeeDetails struct {
	LoanFees       LoanFees       `json:"loanFees"`
	OtherFees      OtherFees      `json:"otherFees"`
	AdditionalInfo AdditionalInfo `json:"additionalInfo"`
}

type LoanFees struct {
	FrontEndFee           string `json:"frontEndFee"`
	ManagementFee         string `json:"managementFee"`
	CommitmentFee         string `json:"commitmentFee"`
	CancellationFee       string `json:"cancellationFee"`
	PrepaymentFee         string `json:"prepaymentFee"`
	ExtensionFee          string `json:"extensionFee"`
	AnnualFee             string `json:"annualFee"`
	AppraisalFeeInternal  string `json:"appraisalFeeInternal"`
	AppraisalFeeExternal  string `json:"appraisalFeeExternal"`
	DebtCollectionFee     string `json:"debtCollectionFee"`
	CreditCheckFee        string `json:"creditCheckFee"`
	StatementReIssuingFee string `json:"statementReIssuingFee"`

	// New fields to store extracted numerical values
	DebtCollectionFeeAmount int `json:"debtCollectionFeeAmount"`
	CreditCheckFeeAmount    int `json:"creditCheckFeeAmount"`
	StatementFeeAmount      int `json:"statementReIssuingFeeAmount"`
}

type OtherFees struct {
	OtherFees *string `json:"otherFees"` // Pointer to allow null value
}

type AdditionalInfo struct {
	FeeWebsiteLink *string `json:"feeWebsiteLink"` // Pointer to allow null value
}
