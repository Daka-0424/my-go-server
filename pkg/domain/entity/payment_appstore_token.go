package entity

import "gorm.io/gorm"

const (
	// appstore_receiptのstatus
	StatusSuccess                      = 0
	StatusUnreadableJson               = 21000
	StatusInvalidReceiptData           = 21002
	StatusReceiptNotAuthenticated      = 21003
	StatusSharedSecretMismatched       = 21004
	StatusServerUnavailable            = 21005
	StatusSubscriptionExpired          = 21006
	StatusSandboxReceiptGoesProduction = 21007
	StatusProductionReceiptGoesSandbox = 21008
)

type PaymentAppstoreToken struct {
	gorm.Model
	TransactionID     string `json:"transaction_id"`
	AppAccountToken   string `json:"app_account_token"`
	BuindleID         string `json:"bundle_id"`
	Currency          string `json:"currency"`
	Enviroment        string `json:"enviroment"`
	ProductID         string `json:"product_id"`
	Price             uint   `json:"price"`
	PurchaseDate      uint   `json:"purchase_date"`
	Quantity          uint   `json:"quantity"`
	RevocationDate    uint   `json:"revocation_date"`
	UserID            uint   `json:"user_id"`
	EarnedPointID     uint   `json:"earned_point_id"`
	PlatformProductID uint   `json:"platform_product_id"`
	PlatformProduct   PlatformProduct
}

func NewPaymentAppstoreToken(
	transactionID string,
	appAccountToken string,
	buindleID string,
	currency string,
	enviroment string,
	productID string,
	price uint,
	purchaseDate uint,
	quantity uint,
	revocationDate uint,
	userID uint,
	platformProduct *PlatformProduct,
) *PaymentAppstoreToken {
	return &PaymentAppstoreToken{
		TransactionID:     transactionID,
		AppAccountToken:   appAccountToken,
		BuindleID:         buindleID,
		Currency:          currency,
		Enviroment:        enviroment,
		ProductID:         productID,
		Price:             price,
		PurchaseDate:      purchaseDate,
		Quantity:          quantity,
		RevocationDate:    revocationDate,
		UserID:            userID,
		PlatformProductID: platformProduct.ID,
	}
}
