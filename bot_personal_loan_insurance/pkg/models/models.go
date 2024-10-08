package models

type LoanDetails struct {
	CreditLimit                string   `json:"credit_limit"`
	CreditLimitCondition       []string `json:"credit_limit_condition,omitempty"`
	InstallmentPeriod          string   `json:"installment_period"`
	InstallmentPeriodCondition string   `json:"installment_period_condition,omitempty"`
	InstallmentPlanDetail      string   `json:"installment_plan_detail,omitempty"`
}

type BorrowerDetails struct {
	MinAge             *int     `json:"min_age,omitempty"`
	MaxAge             *int     `json:"max_age,omitempty"`
	BorrowerConditions []string `json:"borrower_conditions,omitempty"`
}

type FeeDetails struct {
	InternalAppraisalFee string                 `json:"internal_appraisal_fee"`
	ExternalAppraisalFee string                 `json:"external_appraisal_fee"`
	StampDutyFee         []string               `json:"stamp_duty_fee,omitempty"`
	MortgageFee          string                 `json:"mortgage_fee"`
	CreditCheckFee       []string               `json:"credit_check_fee,omitempty"`
	ReturnedChequeFee    string                 `json:"returned_cheque_fee"`
	InsufficientFundsFee string                 `json:"insufficient_funds_fee"`
	StatementReIssueFee  string                 `json:"statement_reissue_fee"`
	DebtCollectionFee    []string               `json:"debt_collection_fee,omitempty"`
	OtherFees            []string               `json:"other_fees,omitempty"`
	ExtractedFees        ExtractedFeesContainer `json:"extracted_fees,omitempty"`
}

type PaymentFees struct {
	DirectDebitProvider        string                       `json:"direct_debit_provider"`
	DirectDebitOtherProvider   string                       `json:"direct_debit_other_provider"`
	BankCounterService         string                       `json:"bank_counter_service"`
	BankCounterOtherService    string                       `json:"bank_counter_other_service"`
	CounterServiceFee          []string                     `json:"counter_service_fee,omitempty"`
	OnlinePaymentFee           string                       `json:"online_payment_fee"`
	CDMATMPaymentFee           string                       `json:"cdm_atm_payment_fee"`
	PhonePaymentFee            string                       `json:"phone_payment_fee"`
	ChequePaymentFee           string                       `json:"cheque_payment_fee"`
	OtherPaymentChannels       []string                     `json:"other_payment_channels,omitempty"`
	CounterServiceFeeExtracted []ExtractedCounterServiceFee `json:"counter_service_fee_extracted,omitempty"`
}

type InterestDetails struct {
	InterestRatePerYear    string                      `json:"interest_rate_per_year"`
	InterestRateConditions string                      `json:"interest_rate_condition"`
	DefaultInterestRate    string                      `json:"default_interest_rate"`
	ExtractedValues        ExtractedInterestRateValues `json:"extracted_interest_rate_values,omitempty"`
}

type ProductDetails struct {
	ProductType     []string `json:"product_type"`
	CreditCharacter string   `json:"credit_character"`
	Collateral      []string `json:"collateral"`
	CreditLineType  string   `json:"credit_line_type"`
	LifeInsurance   string   `json:"life_insurance"`
}

type AdditionalInfo struct {
	Website    *string `json:"website"`
	FeeWebsite *string `json:"fee_website"`
}

type PersonalLoan struct {
	ServiceProvider string          `json:"service_provider"`
	Product         string          `json:"product"`
	ProductDetails  ProductDetails  `json:"product_details"`
	InterestDetails InterestDetails `json:"interest_rate"`
	LoanDetails     LoanDetails     `json:"loan_details"`
	BorrowerDetails BorrowerDetails `json:"borrower_details"`
	FeeDetails      FeeDetails      `json:"fee_details"`
	PaymentFees     PaymentFees     `json:"payments_fees"`
	AdditionalInfo  AdditionalInfo  `json:"additional_info"`
}

type ExtractedInterestRateValues struct {
	BaseMRR                       *float64 `json:"base_mrr,omitempty"`
	EffectiveDate                 *string  `json:"effective_date,omitempty"`
	DefaultInterestRatePercentage *float64 `json:"default_interest_rate_percentage,omitempty"`
}

type ExtractedFees struct {
	StampDutyFee   []ExtractedStampDutyFee   `json:"stamp_duty_fee_extracted,omitempty"`
	CreditCheckFee []ExtractedCreditCheckFee `json:"credit_check_fee_extracted,omitempty"`
	OtherFees      []ExtractedOtherFees      `json:"other_fees_extracted,omitempty"`
}

type ExtractedStampDutyFee struct {
	PercentageOfLoan float64 `json:"percentage_of_loan"`
	MaximumAmount    int     `json:"maximum_amount"`
}

type ExtractedCreditCheckFee struct {
	Type string `json:"type"`
	Fee  int    `json:"fee"`
}

type ExtractedOtherFees struct {
	LoanAmountRange string `json:"loan_amount_range"`
	Fee             int    `json:"fee"`
}

type ExtractedCounterServiceFee struct {
	Location string `json:"location"`
	Fee      int    `json:"fee"`
}

type ExtractedFeesContainer struct {
	StampDutyFeeExtracted      []ExtractedStampDutyFee      `json:"stamp_duty_fee_extracted,omitempty"`
	CreditCheckFeeExtracted    []ExtractedCreditCheckFee    `json:"credit_check_fee_extracted,omitempty"`
	OtherFeesExtracted         []ExtractedOtherFees         `json:"other_fees_extracted,omitempty"`
	CounterServiceFeeExtracted []ExtractedCounterServiceFee `json:"counter_service_fee_extracted,omitempty"`
	DebtCollectionFeeExtracted []ExtractedDebtCollectionFee `json:"debt_collection_fee_extracted,omitempty"`
}

type ExtractedDebtCollectionFee struct {
	Condition         string `json:"condition"`
	FeePerInstallment int    `json:"fee_per_installment,omitempty"`
	MaxFeePerAccount  int    `json:"max_fee_per_account,omitempty"`
	Fee               int    `json:"fee,omitempty"`
}
