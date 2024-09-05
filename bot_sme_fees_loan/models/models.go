package models

type SMEProduct struct {
	ServiceProvider string         `json:"serviceProvider"`
	Product         string         `json:"product"`
	LoanFees        LoanFees       `json:"loan_fees"`
	OtherFees       OtherFees      `json:"other_fees"`
	AdditionalInfo  AdditionalInfo `json:"additional_info"`
}

type LoanFees struct {
	FrontEndFee             string            `json:"frontEndFee"`
	ManagementFee           string            `json:"managementFee"`
	CommitmentFee           string            `json:"commitmentFee"`
	CancellationFee         string            `json:"cancellationFee"`
	PrepaymentFee           FeeWithPercentage `json:"prepaymentFee"`
	ExtensionFee            FeeWithPercentage `json:"extensionFee"`
	AppraisalFeeInternal    FeeWithAmount     `json:"appraisalFeeInternal"`
	AppraisalFeeExternal    FeeWithAmount     `json:"appraisalFeeExternal"`
	DebtCollectionFee       []string          `json:"debtCollectionFee"`
	CreditCheckFee          string            `json:"creditCheckFee"`
	StatementReIssuingFee   string            `json:"statementReIssuingFee"`
	DebtCollectionFeeAmount int               `json:"debtCollectionFeeAmount"`
	CreditCheckFeeAmount    int               `json:"creditCheckFeeAmount"`
	StatementFeeAmount      int               `json:"statementReIssuingFeeAmount"`
}

type FeeWithPercentage struct {
	Description   string  `json:"description"`
	MinPercentage float64 `json:"minPercentage"`
	MaxPercentage float64 `json:"maxPercentage"`
}

type FeeWithAmount struct {
	Description string `json:"description"`
	MinAmount   int    `json:"minAmount"`
	MaxAmount   int    `json:"maxAmount"`
}

type OtherFees struct {
	OtherFees *string `json:"otherFees"`
}

type AdditionalInfo struct {
	FeeWebsiteLink *string `json:"feeWebsiteLink"`
}
