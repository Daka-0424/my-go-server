package request

type AppStoreBilling struct {
	TransactionID  string `json:"transaction_id"`
	PurchaseItemID uint   `json:"purchase_item_id"`
}
