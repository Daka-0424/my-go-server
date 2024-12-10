package entity

type RewardContentType string

const (
	RewardContentTypeNone    RewardContentType = "None"
	RewardContentTypeFreeGem RewardContentType = "FreeGem"
	RewardContentTypeItem    RewardContentType = "Item"
)

type RewardContent struct {
	ContentID       uint              `yaml:"contentId"`
	ContentType     RewardContentType `yaml:"contentType"`
	ContentQuantity uint              `yaml:"quantity"`
}
