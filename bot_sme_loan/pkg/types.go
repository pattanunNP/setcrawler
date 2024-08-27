package pkg

type ProductCondition struct {
	CreditLineType        string   `json:"credit_line_type"`
	Collateral            []string `json:"collateral"`
	ProductConditions     []string `json:"product_conditions"`
	BorrowerAge           string   `json:"borrower_age"`
	ApplicationConditions []string `json:"application_conditions"`
}

type CreditAndLoanTerms struct {
	CreditLimit               string      `json:"credit_limit"`
	CreditLimitConditions     []string    `json:"credit_limit_conditions"`
	BorrowingPeriod           []string    `json:"borrowing_period"`
	BorrowingPeriodConditions interface{} `json:"borrowing_period_conditions"`
}

type InterestRates struct {
	InterestRatesPerYear []string `json:"interest_rate_per_year"`
	DefaultInterest      string   `json:"defaultInterest"`
}

type Fees struct {
	FrontEndFee           string      `json:"front_end_fee"`
	ManagementFee         string      `json:"management_fee"`
	CommitmentFee         string      `json:"commitment_fee"`
	CancellationFee       string      `json:"cancellation_fee"`
	PrepaymentFee         []string    `json:"prepayment_fee"`
	ExtensionFee          string      `json:"extension_fee"`
	AnnualFee             string      `json:"annual_fee"`
	InternalAppraisalFee  string      `json:"internal_appraisal_fee"`
	ExternalAppraisalFee  []string    `json:"external_appraisal_fee"`
	DebtCollectionFee     string      `json:"debt_collection_fee"`
	CreditCheckFee        string      `json:"credit_check_fee"`
	StatementReissuingFee string      `json:"statement_reissuing_fee"`
	OtherFees             interface{} `json:"other_fees"`
}

type AdditionalInfo struct {
	ProductWebsite *string `json:"product_website"`
	FeeWebsite     *string `json:"fee_website"`
}

type Product struct {
	ServiceProvider string             `json:"service_provider"`
	Product         string             `json:"product"`
	InterestRates   InterestRates      `json:"interest_rates"`
	ProductDetails  ProductCondition   `json:"product_details"`
	CreditTerms     CreditAndLoanTerms `json:"credit_and_loan_terms"`
	Fees            Fees               `json:"fees"`
	AdditionalInfo  AdditionalInfo     `json:"additional_info"`
}
