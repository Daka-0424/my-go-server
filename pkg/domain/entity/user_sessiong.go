package entity

type UserSetting struct {
	EntityBase
	UserID      uint `gorm:"<-:create;not null;index:idx_user_id,priority:1"`
	IsSoundPlay bool
	Language    LanguageType

	gormAuxiliary
}

func NewUserSetting(userID uint, languageCode string) *UserSetting {
	return &UserSetting{
		UserID:      userID,
		IsSoundPlay: true,
		Language:    ConvertLanguageType(languageCode),
	}
}

func (u *UserSetting) SetSoundPlay(isSoundPlay bool) {
	u.IsSoundPlay = isSoundPlay
	u.updateColumn("is_sound_play")
}

func (u *UserSetting) SetLanguage(languageCode LanguageType) {
	u.Language = languageCode
	u.updateColumn("language")
}
