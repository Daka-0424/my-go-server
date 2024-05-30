package model

type RequirementVersion struct {
	Version    string `json:"version"`
	NeedUpdate bool   `json:"need_update"`
	IsReview   bool   `json:"is_review"`
}

func NewRequirementVersion(version string, needUpdate bool, isReview bool) *RequirementVersion {
	return &RequirementVersion{
		Version:    version,
		NeedUpdate: needUpdate,
		IsReview:   isReview,
	}
}
