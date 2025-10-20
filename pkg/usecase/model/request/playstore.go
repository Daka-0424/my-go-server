package request

type PlayStoreBilling struct {
	PurchaseItemID uint   `json:"purchaseItemId"`
	Receipt        string `json:"receipt"`
	Signature      string `json:"signature"`
}
