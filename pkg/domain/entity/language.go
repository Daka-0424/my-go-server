package entity

type LanguageType string

const (
	LanguageTypeJapanese LanguageType = "ja"
	LanguageTypeEnglish  LanguageType = "en"
)

var LanguageTypeList = []LanguageType{
	LanguageTypeJapanese,
	LanguageTypeEnglish,
}

type LanguageIndex uint

const (
	LanguageIndexJapanese LanguageIndex = iota
	LanguageIndexEnglish
	LanguageIndexMax
)

type LanguageBase struct {
	Japanese string `yaml:"japanese"`
	English  string `yaml:"english"`
}

func (l *LanguageBase) update(texts []string) {
	l.Japanese = texts[LanguageIndexJapanese]
	l.English = texts[LanguageIndexEnglish]
}

func (l *LanguageBase) GetLanguage(lang LanguageType) string {
	switch lang {
	case LanguageTypeJapanese:
		return l.Japanese
	case LanguageTypeEnglish:
		return l.English
	default:
		return l.English
	}
}

func (l *LanguageBase) GetTexts() []string {
	return []string{
		l.Japanese,
		l.English,
	}
}

func ConvertLanguageType(lang string) LanguageType {
	switch lang {
	case "ja", "jp":
		return LanguageTypeJapanese
	case "en":
		return LanguageTypeEnglish
	default:
		return LanguageTypeEnglish
	}
}

func GetIndexLanguageType(idx uint) LanguageType {
	switch LanguageIndex(idx) {
	case LanguageIndexJapanese:
		return LanguageTypeJapanese
	case LanguageIndexEnglish:
		return LanguageTypeEnglish
	default:
		return LanguageTypeEnglish
	}
}
