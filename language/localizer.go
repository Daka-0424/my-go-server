package language

import (
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const (
	LanguageJapanese = "ja"
	LanguageEnglish  = "en"
)

const (
	defaultLanguage = LanguageJapanese
)

type Localizer struct {
	localizers map[string]*i18n.Localizer
	mu         sync.RWMutex
}

func NewLocalizer(langFiles []string) (*Localizer, error) {
	bundle := i18n.NewBundle(language.Japanese)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	for _, file := range langFiles {
		if _, err := bundle.LoadMessageFile(file); err != nil {
			return nil, err
		}
	}

	return &Localizer{
		localizers: map[string]*i18n.Localizer{
			"ja": i18n.NewLocalizer(bundle, "ja"),
		},
	}, nil
}

func (l *Localizer) MustLocalize(errorCode, languageType string, templateData map[string]interface{}) string {
	lang := strings.Split(languageType, ",")

	// 動的な言語サポートの実装を視野に入れて、mutexを使用
	// ただし、現時点では言語はjaとenのみ
	l.mu.RLock()
	localizer, ok := l.localizers[lang[0]]
	l.mu.RUnlock()
	if !ok {
		localizer = l.localizers[defaultLanguage]
	}

	cfg := &i18n.LocalizeConfig{
		MessageID:    errorCode,
		TemplateData: templateData,
	}
	return localizer.MustLocalize(cfg)
}
