package model

type GoogleLoginURL struct {
	URL string `json:"url"`
}

func NewGoogleLoginURL(url string) GoogleLoginURL {
	return GoogleLoginURL{
		URL: url,
	}
}
