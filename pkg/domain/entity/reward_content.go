package entity

type RewardContentType string

const (
	RewardContentTypeNone    RewardContentType = "None"
	RewardContentTypeFreeGem RewardContentType = "FreeGem"
	RewardContentTypeItem    RewardContentType = "Item"
)

type RewardContent struct {
	ContentID       uint              `yaml:"contentId" json:"content_id"`
	ContentType     RewardContentType `yaml:"contentType" json:"content_type"`
	ContentQuantity uint              `yaml:"quantity" json:"content_quantity"`
}
