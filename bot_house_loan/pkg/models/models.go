package models

type InterestRate struct {
	AverageInterestRateThreeYears float64  `json:"average_interest_rate_three_years"`
	InterestRateConditions        []string `json:"interest_rate_conditions"`
	EffectiveInterestRate         string   `json:"effective_interest_rate"`
	MaximumNormalInterestRate     float64  `json:"maximum_normal_interest_rate"`
	DefaultInterestRate           string   `json:"default_interest_rate"`
}

type ProductDetails struct {
	LoanType                  string  `json:"loan_type"`
	CollateralType            string  `json:"collateral_type"`
	BorrowerQualifications    string  `json:"borrower_qualifications"`
	LoanConditions            string  `json:"loan_conditions"`
	CombinedLoanConditions    *string `json:"combined_loan_conditions,omitempty"`
	ProductSpecificConditions *string `json:"product_specific_conditions,omitempty"`
	BorrowerAge               *int    `json:"borrower_age"`
	MinimumIncome             string  `json:"minimum_income"`
	ApplicationConditions     string  `json:"application_conditions"`
}

type LoanCreditRepayment struct {
	CreditLimitRange      string   `json:"credit_limit_range"`
	LTVRatio              []string `json:"ltv_ratio"`
	CreditLimitConditions *string  `json:"credit_limit_conditions,omitempty"`
	LoanTerm              string   `json:"loan_term"`
	RepaymentConditions   string   `json:"repayment_conditions"`
}

type InsuranceDetails struct {
	MRTAConditions      string `json:"mrta_conditions"`
	MRTACancellationFee string `json:"mrta_cancellation_fee"`
}

type PaymentFees struct {
	DeductingFromBankAccount  string   `json:"deducting_from_bank_account"`
	DeductingFromOtherBankAC  []string `json:"deducting_from_other_bank_account"`
	BankCounterService        string   `json:"bank_counter_service"`
	OtherBankCounterService   string   `json:"other_bank_counter_service"`
	OtherCounterService       []string `json:"other_counter_service"`
	OnlinePayment             []string `json:"online_payment"`
	CDMATMPayment             string   `json:"cdm_atm_payment"`
	PhonePayment              string   `json:"phone_payment"`
	ChequeOrMoneyOrderPayment string   `json:"cheque_or_money_order_payment"`
	OtherChannelPayment       string   `json:"other_channel_payment"`
}

type GeneralFees struct {
	SurveyAndAppraisalFee     []string `json:"survey_and_appraisal_fee"`
	StampDuty                 string   `json:"stamp_duty"`
	MortgageFee               string   `json:"mortgage_fee"`
	TransferFee               string   `json:"transfer_fee"`
	CreditInfoVerificationFee string   `json:"credit_info_verification_fee"`
	FireInsurancePremium      string   `json:"fire_insurance_premium"`
	ChequeReturnFee           string   `json:"cheque_return_fee"`
	DeficiencyBalanceFee      string   `json:"deficiency_balance_fee"`
	StatementCopyFee          string   `json:"statement_copy_fee"`
	ChequeReturnFine          []string `json:"cheque_return_fine"`
	DebtCollectionFee         []string `json:"debt_collection_fee"`
	InterestRateChangeFee     string   `json:"interest_rate_change_fee"`
	RefinanceFee              []string `json:"refinance_fee"`
	OtherFees                 *string  `json:"other_fees,omitempty"`
}

type HouseLoan struct {
	ServiceProvider     string              `json:"service_provider"`
	Product             string              `json:"product"`
	InterestRate        InterestRate        `json:"interest_rate"`
	ProductDetails      ProductDetails      `json:"product_details"`
	LoanCreditRepayment LoanCreditRepayment `json:"loan_credit_repayment"`
	InsuranceDetails    InsuranceDetails    `json:"insurance_details"`
	GeneralFees         GeneralFees         `json:"general_fees"`
	PaymentFees         PaymentFees         `json:"payment_fees"`
	ProductWebsite      string              `json:"product_website"`
	FeeWebsite          string              `json:"fee_website"`
}
