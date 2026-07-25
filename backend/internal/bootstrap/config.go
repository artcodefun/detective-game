package bootstrap

import "os"

type Config struct {
	Port            string
	OpenRouterKey   string
	OpenRouterModel string
	MongoURI        string
	MongoDatabase   string
}

func LoadConfig() Config {
	return Config{
		Port:            envOrDefault("PORT", "8080"),
		OpenRouterKey:   os.Getenv("OPENROUTER_API_KEY"),
		OpenRouterModel: envOrDefault("OPENROUTER_MODEL", "deepseek/deepseek-v4-flash"),
		MongoURI:        envOrDefault("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase:   envOrDefault("MONGO_DATABASE", "detective_game"),
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
