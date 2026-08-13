package domain

type Locale string

const (
	LocaleRU      Locale = "ru"
	LocaleEN      Locale = "en"
	DefaultLocale Locale = LocaleRU
)

type Translation struct {
	Key    string         `bson:"key"`
	Params map[string]any `bson:"params,omitempty"`
}

func T(key string) Translation {
	return Translation{Key: key}
}

func TWith(key string, params map[string]any) Translation {
	return Translation{Key: key, Params: params}
}
