package models

type HouseLoanFeesDetails struct {
	Provider       string         `json:"provider"`
	Product        string         `json:"product"`
	InterestRates  InterestRates  `json:"interest_rates"`
	PaymentsFees   PaymentFees    `json:"payments_fees"`
	OtherFees      OtherFees      `json:"other_fees"`
	AdditionalInfo AdditionalInfo `json:"additional_info"`
}

type InterestRates struct {
	DefaultInterestRate        FeeDetail   `json:"default_interest_rate"`
	SurveyAndAppraisalFee      []FeeDetail `json:"survey_and_appraisal_fee"`
	MRTA                       FeeDetail   `json:"mrta"`
	StampDuty                  FeeDetail   `json:"stamp_duty"`
	MortgageFee                FeeDetail   `json:"mortgage_fee"`
	TransferOwnershipFee       FeeDetail   `json:"transfer_ownership_fee"`
	CreditBureauFee            FeeDetail   `json:"credit_bureau_fee"`
	FireInsurancePremium       FeeDetail   `json:"fire_insurance_premium"`
	OtherChequeReturnedFee     FeeDetail   `json:"other_cheque_returned_fee"`
	InsufficientDirectDebitFee FeeDetail   `json:"insufficient_direct_debit_fee"`
	CopyStatementReissuingFee  FeeDetail   `json:"copy_statement_reissuing_fee"`
	ChequeReturnedFee          FeeDetail   `json:"cheque_returned_fee"`
	DebtCollectionFee          []FeeDetail `json:"debt_collection_fee"`
	ChangingInterestRateFee    FeeDetail   `json:"changing_interest_rate_fee"`
	RefinanceFee               FeeDetail   `json:"refinance_fee"`
}

type PaymentFees struct {
	DirectDebitFromProvider      FeeDetail   `json:"direct_debit_from_provider"`
	DirectDebitFromOtherProvider []FeeDetail `json:"direct_debit_from_other_provider"`
	AtProviderBranch             FeeDetail   `json:"at_provider_branch"`
	AtOtherProviderBranch        FeeDetail   `json:"at_other_provider_branch"`
	AtPaymentServicePoint        []FeeDetail `json:"at_payment_service_point"`
	OnlinePayment                []FeeDetail `json:"online_payment"`
	CDMOrATM                     FeeDetail   `json:"cdm_or_atm"`
	PhonePayment                 FeeDetail   `json:"phone_payment"`
	ChequeOrMoneyOrderPayment    FeeDetail   `json:"cheque_or_money_order_payment"`
	OtherPaymentChannels         FeeDetail   `json:"other_payment_channels"`
}

type FeeDetail struct {
	Original []string `json:"original_text"`
	Numeric  *float64 `json:"numeric,omitempty"`
	MinFee   float64  `json:"min_fee,omitempty"`
	MaxFee   float64  `json:"max_fee,omitempty"`
}

type OtherFees struct {
	OtherFees *string `json:"other_fees"`
}

type AdditionalInfo struct {
	FeewebsiteLink string `json:"fee_website_link"`
}
