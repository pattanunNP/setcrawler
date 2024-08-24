package loanproducts

type TitleLoan struct {
	Provider                  string                    `json:"provider"`
	Product                   string                    `json:"product"`
	VehicleType               []string                  `json:"vehicleType"`
	VehicleCondition          []string                  `json:"vehicleCondition"`
	LoanType                  string                    `json:"loanType"`
	InterestRate              InterestRate              `json:"interestRate"`
	CreditLimitAndInstallment CreditLimitAndInstallment `json:"creditLimitAndInstallment"`
	BorrowerQualifications    BorrowerQualifications    `json:"borrowerQualifications"`
	GeneralFees               GeneralFees               `json:"generalFees"`
	CardFees                  CardFees                  `json:"cardFees"`
	PaymentFees               PaymentFees               `json:"paymentFees"`
	OtherFees                 OtherFees                 `json:"otherFees"`
	AdditionalInfo            *AdditionalInfo           `json:"additionalInfo"`
}

type InterestRate struct {
	AnnualInterestRate     string   `json:"annualInterestRate"`
	InterestRateConditions []string `json:"interestRateConditions"`
	PenaltyInterestRate    string   `json:"penaltyInterestRate"`
}

type CreditLimit struct {
	MinLimit int `json:"min_limit"`
	MaxLimit int `json:"max_limit"`
}

type InstallmentPeriod struct {
	MinMonth int `json:"min_month"`
	MaxMonth int `json:"max_month"`
}

type CreditLimitAndInstallment struct {
	CreditLimit           CreditLimit       `json:"creditLimit"`
	CreditLimitConditions []string          `json:"creditLimitConditions"`
	InstallmentPeriod     InstallmentPeriod `json:"installmentPeriod"`
	LoanReceivingChannel  []string          `json:"loanReceivingChannel"`
}

type BorrowerQualifications struct {
	BorrowerType     []string `json:"borrowerType"`
	AgeLimit         AgeLimit `json:"ageLimit"`
	OtherConditions  *string  `json:"otherConditions,omitempty"`
	MinimumIncome    []string `json:"minimumIncome"`
	IncomeConditions *string  `json:"incomeConditions,omitempty"`
}

type AgeLimit struct {
	MinAge int `json:"min_age"`
	MaxAge int `json:"max_age"`
}

type GeneralFees struct {
	StampDuty         []string `json:"stampDuty"`
	ReturnCheque      string   `json:"returnedCheque"`
	DebtCollectionFee []string `json:"debtCollectionFee"`
}

type CardFees struct {
	CardFees            string  `json:"cardFee"`
	CardReplacementFee  string  `json:"cardReplacementFee"`
	CreditWithdrawalFee *string `json:"creditWithdrawalFee,omitempty"`
}

type PaymentFees struct {
	FreePaymentChannel                []string `json:"freePaymentChannel"`
	DeductingFromServiceProvider      string   `json:"deductingFromServiceProvider"`
	DeductingFromOtherServiceProvider string   `json:"deductingFromOtherServiceProvider"`
	ServiceProviderCounter            string   `json:"serviceProviderCounter"`
	OtherProviderCounter              string   `json:"otherProviderCounter"`
	PaymentServicePoints              []string `json:"paymentServicePoints"`
	OnlinePayment                     string   `json:"onlinePayment"`
	CDMATMPayment                     string   `json:"cdmAtmPayment"`
	PhonePayment                      string   `json:"phonePayment"`
	ChequeMoneyOrderPayment           string   `json:"chequeMoneyOrderPayment"`
	OtherChannelPayment               *string  `json:"otherChannelPayment,omitempty"`
}

type OtherFees struct {
	LitigationLawyerFee []string `json:"litigationLawyerFee"`
	OtherFeesDetails    *string  `json:"otherFeesDetails,omitempty"`
}

type AdditionalInfo struct {
	ProductWebsite *string `json:"productWebsite,omitempty"`
	FeeWebsite     *string `json:"feewebsite,omitempty"`
}
