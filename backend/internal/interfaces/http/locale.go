package http

import (
	"net/http"
	"strings"

	"github.com/artcodefun/detective-game/backend/internal/domain"
)

func localeFromRequest(r *http.Request) domain.Locale {
	locale := strings.ToLower(strings.TrimSpace(strings.Split(r.Header.Get("Accept-Language"), ",")[0]))
	if index := strings.IndexByte(locale, '-'); index >= 0 {
		locale = locale[:index]
	}
	if locale == string(domain.LocaleEN) {
		return domain.LocaleEN
	}
	return domain.DefaultLocale
}
