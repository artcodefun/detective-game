package bootstrap

import "os"

type Config struct {
	Port              string
	OpenRouterKey     string
	OpenRouterModel   string
	MongoURI          string
	MongoDatabase     string
	IOSMinVersion     string
	AndroidMinVersion string
	IOSUpdateURL      string
	AndroidUpdateURL  string
}

func LoadConfig() Config {
	return Config{
		Port:              envOrDefault("PORT", "8080"),
		OpenRouterKey:     os.Getenv("OPENROUTER_API_KEY"),
		OpenRouterModel:   envOrDefault("OPENROUTER_MODEL", "deepseek/deepseek-v4-flash"),
		MongoURI:          envOrDefault("MONGO_URI", "mongodb://localhost:27017/?replicaSet=rs0"),
		MongoDatabase:     envOrDefault("MONGO_DATABASE", "detective_game"),
		IOSMinVersion:     envOrDefault("IOS_MIN_SUPPORTED_VERSION", "0.0.0"),
		AndroidMinVersion: envOrDefault("ANDROID_MIN_SUPPORTED_VERSION", "0.0.0"),
		IOSUpdateURL:      os.Getenv("IOS_UPDATE_URL"),
		AndroidUpdateURL:  os.Getenv("ANDROID_UPDATE_URL"),
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
