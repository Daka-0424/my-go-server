package entity

type SocialGoogle struct {
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Provider      string `json:"provider"`
	ProviderID    string `json:"provider_id"`
}
