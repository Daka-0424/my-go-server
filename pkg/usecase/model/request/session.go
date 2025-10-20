package request

type Session struct {
	UserId uint   `json:"user_id"`
	Uuid   string `json:"uuid"`
}
