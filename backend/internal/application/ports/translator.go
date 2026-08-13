package ports

import "github.com/artcodefun/detective-game/backend/internal/domain"

type Translator interface {
	Translate(locale domain.Locale, value domain.Translation) string
}
