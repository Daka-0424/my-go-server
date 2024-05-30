package entity

import (
	"time"

	"gorm.io/gorm"
)

type GooglePlayPurchaseType int

const (
	GooglePlayPurchaseType__TEST        GooglePlayPurchaseType = iota // テスト購入 (purchaseType: 0)
	GooglePlayPurchaseType__PROMO                                     // プロモーション (purchaseType: 1)
	GooglePlayPurchaseType__REWARDED                                  // 動画視聴などの報酬 (purchaseType: 2)
	GooglePlayPurchaseType__GENERAL     GooglePlayPurchaseType = 100  // 一般購入 (purchaseType: null)
	GooglePlayPurchaseType__UNSPECIFIED GooglePlayPurchaseType = -1   // 不明 (purchaseType: -1)
)

type PaymentPlaystoreToken struct {
	gorm.Model
	OrderID                     string    `json:"order_id"`
	PackageName                 string    `json:"package_name"`
	ProductID                   string    `json:"product_id"`
	PurchaseState               int64     `json:"purchase_state"`
	Purchased                   time.Time `json:"purchased"`
	PurchaseToken               string    `json:"purchase_token"`
	PurchaseTimeMillis          int64     `json:"purchase_time_millis"`
	Quantity                    int64     `json:"quantity"`
	RegionCode                  string    `json:"region_code"`
	ConsumeState                int64     `json:"consume_state"`
	Kind                        string    `json:"kind"`
	DeveloperPayload            string    `json:"developer_payload"`
	AcknowledgementState        int64     `json:"acknowledgement_state"`
	ObfuscatedExternalAccountID string    `json:"obfuscated_external_account_id"`
	ObfuscatedExternalProfileID string    `json:"obfuscated_external_profile_id"`
	PlatformProductID           uint      `json:"platform_product_id"`
	Signature                   string    `json:"signature"`
	PurchaseType                int       `json:"purchase_type"`
	UserID                      uint      `json:"user_id"`
	EarnedPointID               uint
	PlatformProduct             PlatformProduct
}

func NewPaymentPlaystoreToken(
	orderID string,
	packageName string,
	productID string,
	purchaseState int64,
	purchaseToken string,
	purchaseTimeMillis int64,
	quantity int64,
	regionCode string,
	consumeState int64,
	kind string,
	developerPayload string,
	acknowledgementState int64,
	obfuscatedExternalAccountID string,
	obfuscatedExternalProfileID string,
	platformProductID uint,
	signature string,
	purchaseType GooglePlayPurchaseType,
	userID uint,
) *PaymentPlaystoreToken {
	return &PaymentPlaystoreToken{
		OrderID:                     orderID,
		PackageName:                 packageName,
		ProductID:                   productID,
		PurchaseState:               purchaseState,
		Purchased:                   time.Unix(int64(purchaseTimeMillis)/1000, 0),
		PurchaseToken:               purchaseToken,
		PurchaseTimeMillis:          purchaseTimeMillis,
		Quantity:                    quantity,
		RegionCode:                  regionCode,
		ConsumeState:                consumeState,
		Kind:                        kind,
		DeveloperPayload:            developerPayload,
		AcknowledgementState:        acknowledgementState,
		ObfuscatedExternalAccountID: obfuscatedExternalAccountID,
		ObfuscatedExternalProfileID: obfuscatedExternalProfileID,
		PlatformProductID:           platformProductID,
		Signature:                   signature,
		PurchaseType:                int(purchaseType),
		UserID:                      userID,
	}
}

func GetPurchaseType(purchaseType *int64) GooglePlayPurchaseType {
	if purchaseType == nil {
		return GooglePlayPurchaseType__GENERAL
	}

	switch *purchaseType {
	case 0:
		return GooglePlayPurchaseType__TEST
	case 1:
		return GooglePlayPurchaseType__PROMO
	case 2:
		return GooglePlayPurchaseType__REWARDED
	default:
		return GooglePlayPurchaseType__UNSPECIFIED
	}
}
