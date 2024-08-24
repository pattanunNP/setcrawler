package model

// LoanProduct struct to hold loan product information
type LoanProduct struct {
	ServiceProvider         string                 `json:"serviceProvider"`
	Product                 string                 `json:"product"`
	InterestRate            InterestRate           `json:"interestRate"`
	InterestConditions      *string                `json:"interestConditions"`
	InterestPromotions      []string               `json:"interestPromotions"`
	PenaltyRates            []PenaltyRate          `json:"penaltyRates"`
	MinimumPayment          *string                `json:"minimumPayment"`
	CreditLimit             []string               `json:"creditLimit"`
	LoanAmount              LoanAmount             `json:"loanAmount"`
	LoanDurationMonths      LoanDurationMonths     `json:"loanDurationMonths"`
	MoneyTransferMethod     *string                `json:"moneyTransferMethod"`
	MoneyTransferConditions []string               `json:"moneyTransferConditions"`
	ApplicantRequirements   ApplicantRequirements  `json:"applicantRequirements"`
	ApplicationConditions   []string               `json:"applicationConditions"`
	RevolvingCreditFeeInfo  RevolvingCreditFeeInfo `json:"revolvingCreditFeeInfo"`
	ServiceFee              ServiceFee             `json:"serviceFee"`
	PaymentFeeInfo          PaymentFeeInfo         `json:"paymentFeeInfo"`
}

// InterestRate struct to hold interest rate details
type InterestRate struct {
	SalaryEmployee *string `json:"salaryEmployee"`
	BusinessOwner  *string `json:"businessOwner"`
}

// PenaltyRate struct to hold penalty rate details
type PenaltyRate struct {
	Rate      *string  `json:"rate"`
	Condition []string `json:"condition"`
}

type LoanAmount struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type LoanDurationMonths struct {
	MinMonth int `json:"min_month"`
	MaxMonth int `json:"max_month"`
}

type ApplicantRequirements struct {
	SalaryEmployee            SalaryEmployeeRequirements `json:"salaryEmployee"`
	BusinessOwnerRequirements BusinessOwnerRequirements  `json:"businessOwner"`
}

type SalaryEmployeeRequirements struct {
	Age            int    `json:"age"`
	MinimumIncome  int    `json:"minimumIncome"`
	WorkExperience string `json:"workExperience"`
}

type BusinessOwnerRequirements struct {
	Age              int    `json:"age"`
	MinimumIncome    int    `json:"minimumIncome"`
	BusinessDuration string `json:"businessDuration"`
}

type ServiceFee struct {
	CreditCheckFee                *string `json:"creditCheckFee"`
	StampDuty                     *string `json:"stampDuty"`
	EarlyRepaymentFee             *string `json:"earlyRepaymentFee"`
	ChequeReturnedFee             *string `json:"chequeReturnedFee"`
	InsufficientDirectDebitCharge *string `json:"insufficientDirectDebitCharge"`
	StatementReissuingFee         *string `json:"statementReissuingFee"`
	TransactionVerificationFee    *string `json:"transactionVerificationFee"`
	CollectionFee                 *string `json:"collectionFee"`
	OtherFees                     *string `json:"otherFees"`
}

type RevolvingCreditFeeInfo struct {
	InstallmentService  *string `json:"installmentService"`
	Benefits            *string `json:"benefits"`
	CardFee             *string `json:"cardFee"`
	CardReplacementFee  *string `json:"cardReplacementFee"`
	PinReissuingFee     *string `json:"pinReissuingFee"`
	FxRiskConversionFee *string `json:"fxRiskConversionFee"`
}

type PaymentMethods struct {
	NoFee                   []string `json:"noFee"`
	BranchService           *string  `json:"branchService"`
	DeductingFromOtherBank  *string  `json:"deductingFromOtherBank"`
	ProviderBranchService   *string  `json:"providerBranchService"`
	OtherProviderBranch     []string `json:"otherProviderBranch"`
	CounterService          []string `json:"counterService"`
	OnlinePayment           []string `json:"onlinePayment"`
	CDMATMPayment           []string `json:"cDMATMPayment"`
	PhonePayment            *string  `json:"phonePayment"`
	ChequeMoneyOrderPayment *string  `json:"chequeMoneyOrderPayment"`
	OtherChannelPayment     *string  `json:"otherChannelPayment"`
}

type AdditionalInfo struct {
	ProductWebsite *string `json:"productWebsite"`
	FeeWebsite     *string `json:"feeWebsite"`
}

type PaymentFeeInfo struct {
	PaymentMethods PaymentMethods `json:"paymentMethods"`
	AdditionalInfo AdditionalInfo `json:"additionalInfo"`
}
