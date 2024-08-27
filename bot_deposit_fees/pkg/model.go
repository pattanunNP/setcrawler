package pkg

type DepositFee struct {
	Provider       string         `json:"provider"`
	Product        string         `json:"product"`
	Fees           ProductFees    `json:"product_fees"`
	GeneralFees    GeneralFees    `json:"general_fees"`
	OtherFees      OtherFees      `json:"other_fees"`
	AdditionalInfo AdditionalInfo `json:"additional_info"`
}

type ProductFees struct {
	AccountMaintenanceFee             string  `json:"account_maintenance_fee"`
	StatementRequireFee               *string `json:"statement_require_fee"`
	StatementRequireSixMonth          string  `json:"statement_require_six_month"`
	StatementRequireSixMonthToTwoYear string  `json:"statement_require_six_month_to_two_year"`
	StatementRequireTwoYear           string  `json:"statement_require_two_year"`
	ShortMessageService               string  `json:"short_message_service"`
	ShortMessageServiceFeeMonthly     string  `json:"short_message_service_fee_monthly"`
	ShortMessageServiceAnnualFee      string  `json:"short_message_service_annual_fee"`
	LostPassBookFee                   string  `json:"lost_passbook_fee"`
	AccountCloseFee                   string  `json:"account_close_fee"`
}

type GeneralFees struct {
	CoinCollectFee                         *string `json:"coin_collect_fee"`
	BranchFee                              *string `json:"brance_fee"`
	KioskOtherFee                          *string `json:"kiosk_other_fee"`
	KioskFee                               *string `json:"kiosk_fee"`
	AgentFee                               *string `json:"agent_fee"`
	ShopAgentFee                           *string `json:"shop_agent_fee"`
	PostAgentFee                           *string `json:"post_agent_fee"`
	TopupAgentFee                          *string `json:"topup_agent_fee"`
	OtherAgentFee                          *string `json:"other_agent_fee"`
	TransferBetweenSavingCurrentAccountFee *string `json:"transfer_between_saving_current_account_fee"`
	TransferBetweenBankingFee              *string `json:"transfer_between_banking_fee"`
}

type OtherFees struct {
	OtherFees *string `json:"other_fee"`
}

type AdditionalInfo struct {
	FeeURL *string `json:"fee_url"`
}
