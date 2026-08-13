package i18n

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/artcodefun/detective-game/backend/internal/domain"
)

//go:embed locales/ru.json
var russianCatalog []byte

//go:embed locales/en.json
var englishCatalog []byte

type Translator struct {
	catalogs map[domain.Locale]map[string]string
}

func NewTranslator() *Translator {
	return &Translator{catalogs: map[domain.Locale]map[string]string{
		domain.LocaleRU: parseCatalog(russianCatalog),
		domain.LocaleEN: parseCatalog(englishCatalog),
	}}
}

func (t *Translator) Translate(locale domain.Locale, value domain.Translation) string {
	template := t.catalogs[locale][value.Key]
	if template == "" {
		template = t.catalogs[domain.DefaultLocale][value.Key]
	}
	if template == "" {
		return string(value.Key)
	}
	for name, param := range value.Params {
		template = strings.ReplaceAll(template, "{"+name+"}", fmt.Sprint(param))
	}
	return template
}

func parseCatalog(data []byte) map[string]string {
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		panic(fmt.Sprintf("parse translation catalog: %v", err))
	}
	catalog := make(map[string]string, len(raw))
	for key, value := range raw {
		catalog[key] = value
	}
	return catalog
}
