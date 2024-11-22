package entity

type RewardContentType string

const (
	RewardContentTypeNone    RewardContentType = "None"
	RewardContentTypeFreeGem RewardContentType = "FreeGem"
	RewardContentTypeItem    RewardContentType = "Item"
)

type RewardContent struct {
	ContentID       uint              `yaml:"contentId" gorm:"content_id"`
	ContentType     RewardContentType `yaml:"contentType" gorm:"content_type"`
	ContentQuantity uint              `yaml:"quantity" gorm:"content_quantity"`
}
