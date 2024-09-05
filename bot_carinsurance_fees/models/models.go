package models

import (
	"carinsurance/utils"
	"sort"
	"strings"
)

type CarInsuranceDetails struct {
	ServiceProvider       string                `json:"service_provider"`
	Product               string                `json:"product"`
	GeneralFees           GeneralFees           `json:"general_fees"`
	CardFees              CardFees              `json:"card_fees"`
	PaymentFees           PaymentFees           `json:"payment_fees"`
	OtherFees             OtherFees             `json:"other_fees"`
	AdditionalInformation AdditionalInformation `json:"additional_information"`
}

type GeneralFees struct {
	LatePaymentInterest      string    `json:"late_payment_interest"`
	LatePaymentInterestValue float64   `json:"late_payment_interest_value,omitempty"`
	DebtCollectionFee        []string  `json:"debt_collection_fee"`
	DebtCollectionFeeValues  []int     `json:"debt_collection_fee_values,omitempty"`
	StampDutyFee             []string  `json:"stamp_duty_fee"`
	StampDutyFeeValues       []float64 `json:"stamp_duty_fee_values,omitempty"`
	ChequeReturnFee          string    `json:"cheque_return_fee"`
	ChequeReturnFeeValue     float64   `json:"cheque_return_fee_value,omitempty"`
}

type CardFees struct {
	CardFee             string  `json:"card_fee"`
	CardReplacementFee  string  `json:"card_replacement_fee"`
	CreditWithdrawalFee *string `json:"credit_withdrawal_fee,omitempty"`
}

type PaymentFees struct {
	FreePaymentChannels          []string `json:"free_payment_channels"`
	ProviderAccountDeductionFee  string   `json:"provider_account_deduction_fee"`
	OtherProviderAccountFee      string   `json:"other_provider_account_fee"`
	OtherProviderAccountFeeValue float64  `json:"other_provider_account_fee_value,omitempty"`
	ServiceProviderBranchFee     string   `json:"service_provider_branch_fee"`
	OtherBranchFee               string   `json:"other_branch_fee"`
	ServiceCounterFee            string   `json:"service_counter_fee"`
	OnlinePaymentFee             []string `json:"online_payment_fee"`
	CDMATMPaymentFee             string   `json:"cdm_atm_payment_fee"`
	TelephonePaymentFee          string   `json:"telephone_payment_fee"`
	ChequeOrMoneyOrderPaymentFee string   `json:"cheque_or_money_order_payment_fee"`
	OtherPaymentChannelsFee      string   `json:"other_payment_channels_fee"`
}

type OtherFees struct {
	LawyerFeeLitigation []string  `json:"lawyer_fee_litigation"`
	LawyerFeeValues     []float64 `json:"lawyer_fee_values,omitempty"`
	OtherFees           *string   `json:"other_fees,omitempty"`
}

type AdditionalInformation struct {
	FeeWebsiteLink *string `json:"fee_website_link,omitempty"`
}

func (c *CarInsuranceDetails) PopulateComparableFields() {
	// Populate comparable fields based on text values
	c.GeneralFees.LatePaymentInterestValue = utils.ConvertTextToFloat(c.GeneralFees.LatePaymentInterest)
	c.GeneralFees.ChequeReturnFeeValue = utils.ConvertTextToFloat(c.GeneralFees.ChequeReturnFee)
	c.GeneralFees.DebtCollectionFeeValues = utils.ExtractNumbersFromText(strings.Join(c.GeneralFees.DebtCollectionFee, " "))
	c.GeneralFees.StampDutyFeeValues = utils.ExtractFloatNumbersFromText(strings.Join(c.GeneralFees.StampDutyFee, " "))
}

func SortByServiceProvider(details []CarInsuranceDetails) {
	sort.Slice(details, func(i, j int) bool {
		return details[i].ServiceProvider < details[j].ServiceProvider
	})
}

func FilterByMaxLatePaymentInterest(details []CarInsuranceDetails, maxInterest float64) []CarInsuranceDetails {
	var filtered []CarInsuranceDetails
	for _, detail := range details {
		if detail.GeneralFees.LatePaymentInterestValue <= maxInterest {
			filtered = append(filtered, detail)
		}
	}
	return filtered
}
