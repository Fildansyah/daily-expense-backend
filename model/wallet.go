package model

import "expense.com/m/model/concern"

type WalletTypeEnum string

const (
	WalletTypeWallet  WalletTypeEnum = "WALLET"
	WalletTypeDebit   WalletTypeEnum = "DEBIT"
	WalletTypeCredit  WalletTypeEnum = "CREDIT"
	WalletTypeEwallet WalletTypeEnum = "E_WALLET"
	WalletTypeSaving  WalletTypeEnum = "SAVING"
)

var WalletTypeValues = []WalletTypeEnum{WalletTypeWallet, WalletTypeDebit, WalletTypeCredit, WalletTypeEwallet}

type Wallet struct {
	concern.BaseFields
	UserID       string
	User         User `gorm:"-"`
	Name         string
	Balance      float64
	Currency     string
	WalletTypeId string
	WalletType   WalletTypeEnum `gorm:"-"`
	Color        string
	Icon         string
	IsExclude    bool
}
