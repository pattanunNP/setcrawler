package models

type FeeDetail struct {
	Text      []string `json:"text"`
	MinAmount *int     `json:"minAmount"`
	MaxAmount *int     `json:"maxAmount"`
}

type GeneralFees struct {
	NewVehicleRegistrationFee                  FeeDetail `json:"newVehicleRegistrationFee"`
	OwnershipTransferFeeOneStep                FeeDetail `json:"ownershipTransferFeeOneStep"`
	OwnershipTransferFeeTwoStep                FeeDetail `json:"ownershipTransferFeeTwoStep"`
	VehicleInspectionFee                       FeeDetail `json:"vehicleInspectionFee"`
	ServiceProviderOwnershipTransferFeeOneStep FeeDetail `json:"serviceProviderOwnershipTransferFeeOneStep"`
	ServiceProviderOwnershipTransferFeeTwoStep FeeDetail `json:"serviceProviderOwnershipTransferFeeTwoStep"`
	LeaseTransferFee                           FeeDetail `json:"leaseTransferFee"`
	ContractTerminationFee                     FeeDetail `json:"contractTerminationFee"`
	LatePaymentPenalty                         FeeDetail `json:"latePaymentPenalty"`
	DebtCollectionFee                          FeeDetail `json:"debtCollectionFee"`
	TaxRenewalFee                              FeeDetail `json:"taxRenewalFee"`
	LicensePlateProcessingFee                  FeeDetail `json:"licensePlateProcessingFee"`
	RegistrationAddressChangeFee               FeeDetail `json:"registrationAddressChangeFee"`
	DocumentCopyServiceFee                     FeeDetail `json:"documentCopyServiceFee"`
}

type PaymentFees struct {
	DirectDebitFromProviderAccount      FeeDetail `json:"directDebitFromProviderAccount"`
	DirectDebitFromOtherProviderAccount FeeDetail `json:"directDebitFromOtherProviderAccount"`
	ProviderBranchPayment               FeeDetail `json:"providerBranchPayment"`
	OtherBranchPayment                  FeeDetail `json:"otherBranchPayment"`
	PaymentServicePoints                FeeDetail `json:"paymentServicePoints"`
	OnlinePayment                       FeeDetail `json:"onlinePayment"`
	CDMAtmPayment                       FeeDetail `json:"cdmAtmPayment"`
	PhonePayment                        FeeDetail `json:"phonePayment"`
	ChequeMoneyOrderPayment             FeeDetail `json:"chequeMoneyOrderPayment"`
	OtherChannelsPayment                FeeDetail `json:"otherChannelsPayment"`
}

type OtherFees struct {
	OtherFeesAndCharges FeeDetail `json:"otherFeesAndCharges"`
}

type AdditionalInfo struct {
	FeeWebsiteLinks *string `json:"feeWebsiteLinks"`
}

type CarLoanFeeDetail struct {
	Provider              string         `json:"provider"`
	ContractEffectiveDate *string        `json:"contranctEffectiveDate"`
	GeneralFees           GeneralFees    `json:"generalFees"`
	PaymentFees           PaymentFees    `json:"paymentFees"`
	OtherFees             OtherFees      `json:"otherFees"`
	AdditionalInfo        AdditionalInfo `json:"additional_info"`
}
