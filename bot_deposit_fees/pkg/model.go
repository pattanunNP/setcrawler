package pkg

type DepositFee struct {
	Provider           string             `json:"provider"`
	Product            string             `json:"product"`
	Fees               ProductFees        `json:"product_fees"`
	NumericProductFees NumericProductFees `json:"numeric_product_fees"`
	GeneralFees        GeneralFees        `json:"general_fees"`
	NumericGeneralFees NumericGeneralFees `json:"numeric_general_fees"`
	OtherFees          OtherFees          `json:"other_fees"`
	AdditionalInfo     AdditionalInfo     `json:"additional_info"`
}

type ProductFees struct {
	AccountMaintenanceFee             []string `json:"account_maintenance_fee"`
	StatementRequireFee               []string `json:"statement_require_fee"`
	StatementRequireSixMonth          []string `json:"statement_require_six_month"`
	StatementRequireSixMonthToTwoYear []string `json:"statement_require_six_month_to_two_year"`
	StatementRequireTwoYear           []string `json:"statement_require_two_year"`
	ShortMessageService               []string `json:"short_message_service"`
	ShortMessageServiceFeeMonthly     []string `json:"short_message_service_fee_monthly"`
	ShortMessageServiceAnnualFee      []string `json:"short_message_service_annual_fee"`
	LostPassBookFee                   []string `json:"lost_passbook_fee"`
	AccountCloseFee                   []string `json:"account_close_fee"`
}

type NumericProductFees struct {
	AccountMaintenanceFee             *float64 `json:"account_maintenance_fee"`
	StatementRequireFee               *float64 `json:"statement_require_fee"`
	StatementRequireSixMonth          *float64 `json:"statement_require_six_month"`
	StatementRequireSixMonthToTwoYear *float64 `json:"statement_require_six_month_to_two_year"`
	StatementRequireTwoYear           *float64 `json:"statement_require_two_year"`
	ShortMessageService               *float64 `json:"short_message_service"`
	ShortMessageServiceFeeMonthly     *float64 `json:"short_message_service_fee_monthly"`
	ShortMessageServiceAnnualFee      *float64 `json:"short_message_service_annual_fee"`
	LostPassBookFee                   *float64 `json:"lost_passbook_fee"`
	AccountCloseFee                   *float64 `json:"account_close_fee"`
}

type GeneralFees struct {
	CoinCollectFee                         []string `json:"coin_collect_fee"`
	BranchFee                              []string `json:"branch_fee"`
	KioskOtherFee                          []string `json:"kiosk_other_fee"`
	KioskFee                               []string `json:"kiosk_fee"`
	AgentFee                               []string `json:"agent_fee"`
	ShopAgentFee                           []string `json:"shop_agent_fee"`
	PostAgentFee                           []string `json:"post_agent_fee"`
	TopupAgentFee                          []string `json:"topup_agent_fee"`
	OtherAgentFee                          []string `json:"other_agent_fee"`
	TransferBetweenSavingCurrentAccountFee []string `json:"transfer_between_saving_current_account_fee"`
	TransferBetweenBankingFee              []string `json:"transfer_between_banking_fee"`
}

type NumericGeneralFees struct {
	CoinCollectFee                         *float64 `json:"coin_collect_fee"`
	BranchFee                              *float64 `json:"branch_fee"`
	KioskOtherFee                          *float64 `json:"kiosk_other_fee"`
	KioskFee                               *float64 `json:"kiosk_fee"`
	AgentFee                               *float64 `json:"agent_fee"`
	ShopAgentFee                           *float64 `json:"shop_agent_fee"`
	PostAgentFee                           *float64 `json:"post_agent_fee"`
	TopupAgentFee                          *float64 `json:"topup_agent_fee"`
	OtherAgentFee                          *float64 `json:"other_agent_fee"`
	TransferBetweenSavingCurrentAccountFee *float64 `json:"transfer_between_saving_current_account_fee"`
	TransferBetweenBankingFee              *float64 `json:"transfer_between_banking_fee"`
}

type OtherFees struct {
	OtherFees []string `json:"other_fee"`
}

type AdditionalInfo struct {
	FeeURL *string `json:"fee_url"`
}
