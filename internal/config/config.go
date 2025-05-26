package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SupabaseURL        string
	SupabaseAnonKey    string
	SupabaseServiceKey string
	DatabaseURL        string
	JWTSecret          string
	Port               string
}

func Load() *Config {
	paths := []string{".env", "../.env", "../../.env"}

	loaded := false
	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			log.Printf("Loaded .env from: %s", path)
			loaded = true
			break
		}
	}

	if !loaded {
		log.Println("Warning: No .env file found, using system environment variables")
	}

	config := &Config{
		SupabaseURL:        os.Getenv("SUPABASE_URL"),
		SupabaseAnonKey:    os.Getenv("SUPABASE_ANON_KEY"),
		SupabaseServiceKey: os.Getenv("SUPABASE_SERVICE_KEY"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		Port:               os.Getenv("PORT"),
	}

	// Validate required fields
	if config.SupabaseURL == "" {
		log.Fatal("SUPABASE_URL is required")
	}
	if config.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	return config
}
