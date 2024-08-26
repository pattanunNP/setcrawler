package model

type InterestRateDetails struct {
	InterestWithServiceFee          string   `json:"interest_with_service"`
	InterestWithServiceFeeCondition []string `json:"interest_with_servicefee_condition"`
	DefaultInterestRate             string   `json:"default_interest_rate"`
	DefaultInterestRateCondition    []string `json:"default_interestrate_condition"`
}

type LoanAmount struct {
	MinLoanAmount int `json:"min_loan_amount"`
	MaxLoanAmount int `json:"max_loan_amount"`
}

type LoanDuration struct {
	MinMonth *int `json:"min_loan_month"`
	MaxMonth *int `json:"max_loan_amount"`
}

type RepaymentConditions struct {
	Conditions []string `json:"conditions_repayment"`
}

type CreditApprovalConditions struct {
	LoanAmount          *LoanAmount         `json:"loan_amount"`
	ApprovalConditions  *string             `json:"approval_conditions"`
	LoanDuration        LoanDuration        `json:"loan_duration"`
	RepaymentConditions RepaymentConditions `json:"repayment_conditions"`
}

// ApplicantConditions holds conditions related to the applicant
type ApplicantConditions struct {
	MinAge                  *int     `json:"min_age"`
	MaxAge                  *int     `json:"max_age"`
	ApplicantQualifications []string `json:"applicant_qualifications"`
}

// ProductDetails holds details about the product
type ProductDetails struct {
	Details    []string `json:"details"`
	CreditLine []string `json:"credit_line"`
}

type PaymentFees struct {
	FreeChannels              []string `json:"free_payment_channels"`
	DeductFromServiceProvider string   `json:"deduct_from_provider_account"`
	DeductFromOtherProvider   string   `json:"deduct_from_other_provider_account"`
	ServiceProviderBranch     string   `json:"pay_at_provider_branch"`
	OtherProviderBranch       string   `json:"pay_at_other_branch"`
	PaymentCounters           string   `json:"pay_at_payment_counters"`
	OnlineChannels            []string `json:"online_payment_channels"`
	ATMCDMChannels            []string `json:"atm_cdm_payment_channels"`
	TelephoneChannels         string   `json:"phone_payment_channels"`
	ChequeMoneyOrderChannels  string   `json:"cheque_money_order_channels"`
	OtherChannels             string   `json:"other_payment_channels"`
}

type AdditionalInfo struct {
	ProductWebste string `json:"product_website_link"`
}

// ProductInfo holds information about the product
type ProductInfo struct {
	ServiceProvider       string                   `json:"service_provider"`
	Product               string                   `json:"product"`
	ProductDetails        ProductDetails           `json:"productDetails"`
	ApplicantConditions   ApplicantConditions      `json:"applicant_conditions"`
	CreditApprovalDetails CreditApprovalConditions `json:"credit_approval_conditions"`
	InterestRateDetails   InterestRateDetails      `json:"interest_rate_details"`
	PaymentFees           PaymentFees              `json:"payment_fees"`
	AdditionalInfo        AdditionalInfo           `json:"additional_info"`
}
