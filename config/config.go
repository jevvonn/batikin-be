package config

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string `env:"APP_ENV,required"`
	AppPort string `env:"APP_PORT,required"`

	DbHost     string `env:"DB_HOST,required"`
	DbPort     string `env:"DB_PORT,required"`
	DbUser     string `env:"DB_USER,required"`
	DbPassword string `env:"DB_PASSWORD,required"`
	DbName     string `env:"DB_NAME,required"`

	JWTSecret string `env:"JWT_SECRET,required"`

	OPENAIAPIKey string `env:"OPENAI_API_KEY,required"`
	GeminiAPIKey string `env:"GEMINI_API_KEY,required"`

	SupabaseURL  string `env:"SUPABASE_PROJECT_URL,required"`
	SupbaseToken string `env:"SUPABASE_PROJECT_TOKEN,required"`

	MidtransServerKey string `env:"MIDTRANS_SERVER_KEY,required"`
	RoboflowAPIKey    string `env:"ROBOFLOW_API_KEY,required"`
}

var cfg Config

func Load() Config {
	if cfg.AppEnv == "" {
		New()
	}

	return cfg
}

func New() Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// parse
	err = env.Parse(&cfg)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	return cfg
}
